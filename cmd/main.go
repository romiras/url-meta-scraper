package main

import (
	"fmt"
	"os"

	"github.com/romiras/url-meta-scraper/consumers/drivers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var urlsChan chan string

/*
func scanURLs(reg *registries.Registry, urlsRcv chan string, payloadsTx chan []byte, logger log.Logger) {
	for url := range urlsRcv {
		urlScraped, attempts, err := reg.FetchHelper.Try(reg.Fetcher, url, 0)
		if err != nil {
			if attempts == 0 {
				GiveUp(err)
			} else {
				// Retry(reg, url, attempts)
				logger.Info("Retry", attempts, url)
				// urlsRcv <- url
			}
			continue
		}

		payload, err := json.Marshal(urlScraped)
		if err != nil {
			panic(err)
		}

		payloadsTx <- payload
	}
}

func sender(reg *registries.Registry, payloadsRcv chan []byte) {
	for payload := range payloadsRcv {
		err := reg.RedisPublisher.Publish("scraped-urls", payload)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func Do(reg *registries.Registry) {
	payloadsChan := make(chan []byte)

	go func() {
		sender(reg, payloadsChan)
	}()

	go func() {
		scanURLs(reg, urlsChan, payloadsChan)
		close(payloadsChan)
	}()
}

func GiveUp(err error) {
	log.Println(err.Error())
}

func Retry(reg *registries.Registry, url string, attempts uint) {
	log.Println("\tRetry", attempts, url, "-> failed-urls")
	// TODO: (URL, attempt nr.) -> queue 'failed-urls'
	// err := reg.RedisPublisher.Publish("failed-urls", payload)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Consider Exponential Backoff algorithm  https://blog.miguelgrinberg.com/post/how-to-retry-with-class
}
*/

func newLogger() *logrus.Logger {
	logger := logrus.New()
	if os.Getenv("GO_ENV") == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}

const DefaultAmqpURI = "amqp://guest:guest@localhost:5672"

/*
func handle(deliveries <-chan amqp.Delivery, done chan error, logger log.Logger) {
	for msg := range deliveries {
		logger.Info("-> msg")
		fmt.Println(string(msg.Body))
	}
	logger.Info("handle: deliveries channel closed")
	done <- nil
}
*/

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

func main() {
	logger := newLogger()

	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	c, err := NewConsumer(amqpURI, "urls", "", logger)
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

	logger.Info("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

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
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go func() {
		handle(deliveries, logger)
		c.done <- nil
	}()

	return c, nil
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

	/*
		reg := registries.NewRegistry()
		reg.RedisSubscriber = services.NewRedisSubscriber("urls")
		reg.RedisPublisher = services.NewRedisPublisher()

		urlsChan = make(chan string)
		go func() {
			reg.RedisSubscriber.Consume(urlsChan)
		}()

		go func() {
			for url := range urlsChan {
				log.Println(url)
			}
		}()

		Do(reg)
		time.Sleep(time.Millisecond)

		err := reg.RedisSubscriber.Close()
		if err != nil {
			log.Println(err.Error())
			return
		}
	*/
}

// func worker(sub *services.RedisSubscriber, done chan bool) {
// 	sub.Consume()
// 	fmt.Println("worker Done")
// 	done <- true
// }

// <-done
// close(done)
