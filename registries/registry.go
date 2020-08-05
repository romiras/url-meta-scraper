package registries

import (
	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/initializers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/producers"
	"github.com/romiras/url-meta-scraper/services"
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
