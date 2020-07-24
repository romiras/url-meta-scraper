package main

import (
	"fmt"
	"os"

	"github.com/romiras/url-meta-scraper/consumers"
	"github.com/romiras/url-meta-scraper/consumers/drivers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var urlsChan chan string

func newLogger() *logrus.Logger {
	logger := logrus.New()
	if os.Getenv("GO_ENV") == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}

const DefaultAmqpURI = "amqp://guest:guest@localhost:5672"

func handle(deliveries <-chan amqp.Delivery, logger log.Logger) {
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

func main0() {
	logger := newLogger()

	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	cons, err := drivers.NewAmqpConsumer(amqpURI, "urls", logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Consume")
	// var msgs <-chan amqp.Delivery
	err = cons.Consume(handle, logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	select {}

	err = cons.Close()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Finish")
}

func main() {
	logger := newLogger()

	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	queue := "urls"
	c, err := NewConsumer(amqpURI, queue, "", logger)
	if err != nil {
		logger.Fatal(err)
	}

	err = c.Cons(handler, logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("running forever")
	select {}

	logger.Info("shutting down")

	if err := c.Shutdown(); err != nil {
		logger.Fatalf("error during shutdown: %s", err)
	}
}

func (c *Consumer) Cons(handleFunc consumers.HandleFunc, logger log.Logger) error {
	var err error

	logger.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		logger.Fatal(fmt.Errorf("Channel: %s", err))
	}

	queue := "urls"

	logger.Info("Queue bound to Exchange, starting Consume (consumer tag '%s')", c.tag)
	deliveries, err := c.channel.Consume(
		queue, // name
		c.tag, // consumerTag,
		false, // noAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		logger.Fatal(fmt.Errorf("Queue Consume: %s", err))
	}

	go func() {
		handle(deliveries, logger)
		c.done <- nil
	}()

	return nil
}

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

func NewConsumer(amqpURI, queue, ctag string, logger log.Logger) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}

	var err error

	logger.Info("dialing %s", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	return c, nil
}
