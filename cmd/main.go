package main

import (
	"fmt"

	"github.com/romiras/url-meta-scraper/log"
	"github.com/romiras/url-meta-scraper/registries"
	"github.com/streadway/amqp"
)

/*
const DefaultAmqpURI = "amqp://guest:guest@localhost:5672"

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
		task, err := producers.NewTask(payload)
		if err != nil {
			reg.Logger.Error(err)
		}

		_, err = reg.ScrapedURLPublisher.Produce(task)
		if err != nil {
			reg.Logger.Error(err)
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
*/

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
	defer reg.Close()

	err := reg.URLSubscriber.Consume(handler, reg.Logger)
	if err != nil {
		reg.Logger.Fatal(err)
	}

	reg.Logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	select {}

	reg.Logger.Info("shutting down")
}
