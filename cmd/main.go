package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/romiras/url-meta-scraper/consumers/drivers"
	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/registries"
	"github.com/streadway/amqp"
)

const DefaultAmqpURI = "amqp://guest:guest@localhost:5672"

var urlsChan chan string

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
	logger.Info(err.Error())
}

func Retry(reg *registries.Registry, url string, attempts uint) {
	logger.Info("\tRetry", attempts, url, "-> failed-urls")
	// TODO: (URL, attempt nr.) -> queue 'failed-urls'
	// err := reg.RedisPublisher.Publish("failed-urls", payload)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Consider Exponential Backoff algorithm  https://blog.miguelgrinberg.com/post/how-to-retry-with-class
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

func main() {
	reg := registries.NewRegistry()
	// reg.RedisSubscriber = services.NewRedisSubscriber("urls")
	// reg.RedisPublisher = services.NewRedisPublisher()

	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	queue := "urls"
	cons, err := drivers.NewAmqpConsumer(amqpURI, queue, reg.Logger)
	if err != nil {
		reg.Logger.Fatal(err)
	}

	err = cons.Consume(handler, reg.Logger)
	if err != nil {
		reg.Logger.Fatal(err)
	}

	reg.Logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	select {}

	reg.Logger.Info("shutting down")

	if err := cons.Close(); err != nil {
		reg.Logger.Fatalf("error during shutdown: %s", err)
	}

	/// Old

	urlsChan = make(chan string)
	go func() {
		reg.RedisSubscriber.Consume(urlsChan)
	}()

	go func() {
		for url := range urlsChan {
			reg.Logger.Info(url)
		}
	}()

	Do(reg)
	time.Sleep(time.Millisecond)

	err := reg.RedisSubscriber.Close()
	if err != nil {
		reg.Logger.Info(err.Error())
		return
	}

}
