package registries

import (
	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/initializers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/pkg"
	"github.com/romiras/url-meta-scraper/producers"
	"github.com/romiras/url-meta-scraper/services"
	"github.com/streadway/amqp"
)

type Registry struct {
	Fetcher             *services.Fetcher
	FetchHelper         *services.FetchHelper
	URLSubscriber       consumers.IConsumer
	ScrapedURLPublisher producers.TaskProducer
	Logger              log.Logger
}

func NewRegistry() *Registry {
	logger := initializers.NewLogger()
	urlSubscriber := initializers.NewURLSubscriber(logger)
	scrapedURLPublisher := initializers.NewScrapedURLPublisher(logger)

	return &Registry{
		&services.Fetcher{
			Client: services.NewHTTPClient(),
		},
		&services.FetchHelper{},
		urlSubscriber,
		scrapedURLPublisher,
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

func (reg *Registry) handleURL(url string) bool {
	urlScraped, attempts, err := reg.FetchHelper.Try(reg.Fetcher, url, 0)
	if err != nil {
		if attempts == 0 {
			reg.Logger.Info("Giving up", err.Error())
		} else {
			reg.Retry(url, attempts)
		}
		return false
	}

	err = reg.produceScrapedURLTask(urlScraped)
	if err != nil {
		reg.Logger.Error(err)
		return false
	}

	return true
}

func (reg *Registry) produceScrapedURLTask(urlScraped *pkg.UrlScraped) error {
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

func (reg *Registry) Retry(url string, attempts uint) {
	reg.Logger.Info("\tRetry", attempts, url, " -> failed-urls")
	// TODO: (URL, attempt nr.) -> queue 'failed-urls'
	// err := reg.RedisPublisher.Publish("failed-urls", payload)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Consider Exponential Backoff algorithm  https://blog.miguelgrinberg.com/post/how-to-retry-with-class
}

func (reg *Registry) ConsumerHandler(deliveries <-chan amqp.Delivery, logger log.Logger) {
	multiAck := false
	for msg := range deliveries {
		logger.Info("-> msg")
		reg.handleURL(string(msg.Body))

		err := msg.Ack(multiAck)
		if err != nil {
			logger.Error(err)
		}
	}
	logger.Info("handle: deliveries channel closed / no new messages")
}

/*
func (reg *Registry) scanURLs(urlsRcv chan string, payloadsTx chan []byte) {
	for url := range urlsRcv {
		urlScraped, attempts, err := reg.FetchHelper.Try(reg.Fetcher, url, 0)
		if err != nil {
			if attempts == 0 {
				reg.Logger.Info("Giving up", err.Error())
			} else {
				// Retry(reg, url, attempts)
				reg.Logger.Info("Retry", attempts, url)
				// urlsRcv <- url
			}
			continue
		}

		payload, err := json.Marshal(urlScraped)
		if err != nil {
			reg.Logger.Fatal(err)
		}

		payloadsTx <- payload
	}
}
*/
