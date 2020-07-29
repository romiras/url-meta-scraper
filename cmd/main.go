package main

import (
	"fmt"
	"os"

	"github.com/romiras/url-meta-scraper/consumers/drivers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func newLogger() *logrus.Logger {
	logger := logrus.New()
	if os.Getenv("GO_ENV") == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}

const DefaultAmqpURI = "amqp://guest:guest@localhost:5672"

func handler(deliveries <-chan amqp.Delivery, logger log.Logger) {
	multiAck := false
	for msg := range deliveries {
		logger.Info("-> msg")
		fmt.Println(string(msg.Body))
		err := msg.Ack(multiAck)
		if err != nil {
			logger.Error(err)
		}
	}
	logger.Info("handle: deliveries channel closed / no new messages")
}

func main() {
	logger := newLogger()

	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	queue := "urls"
	cons, err := drivers.NewAmqpConsumer(amqpURI, queue, logger)
	if err != nil {
		logger.Fatal(err)
	}

	err = cons.Consume(handler, logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	select {}

	logger.Info("shutting down")

	if err := cons.Close(); err != nil {
		logger.Fatalf("error during shutdown: %s", err)
	}
}
