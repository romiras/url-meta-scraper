package initializers

import (
	"os"

	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/producers"
	"github.com/romiras/url-meta-scraper/producers/drivers"
)

const DefaultScrapedURLQueue = "scraped-urls"

func NewScrapedURLPublisher(logger log.Logger) producers.TaskProducer {
	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	queue := os.Getenv("SCRAPED_URLS_QUEUE_NAME")
	if queue == "" {
		queue = DefaultScrapedURLQueue
	}

	producer, err := drivers.NewAmqpProducer(amqpURI, queue, logger)
	if err != nil {
		logger.Fatal(err)
	}

	return producer
}
