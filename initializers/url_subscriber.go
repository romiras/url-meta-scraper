package initializers

import (
	"os"

	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/consumers/drivers"
	"github.com/romiras/url-meta-scraper/log"
)

const (
	DefaultAmqpURI  = "amqp://guest:guest@localhost:5672"
	DefaultURLQueue = "urls"
)

func NewURLSubscriber(logger log.Logger) consumers.IConsumer {
	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	queue := os.Getenv("URLS_QUEUE_NAME")
	if queue == "" {
		queue = DefaultURLQueue
	}

	cons, err := drivers.NewAmqpConsumer(amqpURI, queue, logger)
	if err != nil {
		logger.Fatal(err)
	}

	return cons
}
