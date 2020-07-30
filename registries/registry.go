package registries

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/services"
)

type Registry struct {
	Fetcher         *services.Fetcher
	FetchHelper     *services.FetchHelper
	RedisPublisher  *services.RedisPublisher
	RedisSubscriber *services.RedisSubscriber
	AmqpSubscriber  consumers.IConsumer
	Logger          log.Logger
}

func NewRegistry() *Registry {
	return &Registry{
		&services.Fetcher{
			Client: services.NewHTTPClient(),
		},
		&services.FetchHelper{},
		nil,
		nil,
		newLogger(),
	}
}

func newLogger() *logrus.Logger {
	logger := logrus.New()
	if os.Getenv("GO_ENV") == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}
