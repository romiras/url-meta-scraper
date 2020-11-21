package registries

import (
	"encoding/json"
	"time"

	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/initializers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/pkg"
	"github.com/romiras/url-meta-scraper/pkg/events"
	"github.com/romiras/url-meta-scraper/producers"
	"github.com/romiras/url-meta-scraper/services"
	"github.com/streadway/amqp"
)

type Registry struct {
	Fetcher             *services.Fetcher
	FetchHelper         *services.FetchHelper
	URLSubscriber       consumers.IConsumer
	ScrapedURLPublisher producers.TaskProducer
	FailedURLPublisher  producers.TaskProducer
	Logger              log.Logger
}

func NewRegistry() *Registry {
	logger := initializers.NewLogger()
	urlSubscriber := initializers.NewURLSubscriber(logger)
	scrapedURLPublisher := initializers.NewScrapedURLPublisher(logger)
	failedURLPublisher := initializers.NewFailedURLPublisher(logger)

	return &Registry{
		&services.Fetcher{
			Client: services.NewHTTPClient(),
		},
		&services.FetchHelper{},
		urlSubscriber,
		scrapedURLPublisher,
		failedURLPublisher,
		logger,
	}
}

func (reg *Registry) Close() {
	var err error
	if err = reg.URLSubscriber.Close(); err != nil {
		reg.Logger.Error(err)
	}

	if err = reg.ScrapedURLPublisher.Close(); err != nil {
		reg.Logger.Error(err)
	}
}

func (reg *Registry) handleURL(urlEvent *events.UrlEvent) bool {
	if urlEvent == nil {
		reg.Logger.Fatal("Got urlEvent==nil")
		return false
	}
	urlScraped, attempts, err := reg.FetchHelper.Try(reg.Fetcher, urlEvent.URL, urlEvent.Attempts)
	if err != nil {
		if attempts == pkg.NoRetry {
			reg.Logger.Info("Giving up", err.Error())
		} else {
			reg.Retry(urlEvent)
		}
		return false
	}
	if urlScraped == nil {
		return false
	}

	err = reg.produceScrapedURLTask(urlScraped)
	if err != nil {
		reg.Logger.Error(err)
		return false
	}

	return true
}

func (reg *Registry) produceScrapedURLTask(urlScraped *events.UrlScrapedEvent) error {
	task, err := producers.NewTask(urlScraped)
	if err != nil {
		return err
	}

	_, err = reg.ScrapedURLPublisher.Produce(task)
	if err != nil {
		return err
	}

	return nil
}

func (reg *Registry) produceFailedURLTask(urlFailed *events.UrlFailedEvent) error {
	task, err := producers.NewTask(urlFailed)
	if err != nil {
		return err
	}

	_, err = reg.FailedURLPublisher.Produce(task)
	if err != nil {
		return err
	}

	return nil
}

func (reg *Registry) Retry(urlEvent *events.UrlEvent) {
	reg.Logger.Info("\tRetry", urlEvent.Attempts, urlEvent.URL, " -> failed-urls")
	urlFailed := events.UrlFailedEvent{
		URL:      urlEvent.URL,
		FailedAt: time.Now().Unix(),
		Attempts: urlEvent.Attempts,
	}
	err := reg.produceFailedURLTask(&urlFailed)
	if err != nil {
		reg.Logger.Fatal(err)
	}
}

func (reg *Registry) ConsumerHandler(deliveries <-chan amqp.Delivery, logger log.Logger) {
	var urlEvent *events.UrlEvent
	multiAck := false
	for msg := range deliveries {
		logger.Info("-> msg")
		err := json.Unmarshal(msg.Body, &urlEvent)
		if err != nil {
			logger.Error(err)
		}
		_ = reg.handleURL(urlEvent)

		err = msg.Ack(multiAck)
		if err != nil {
			logger.Error(err)
		}
	}
	logger.Info("handle: deliveries channel closed / no new messages")
}
