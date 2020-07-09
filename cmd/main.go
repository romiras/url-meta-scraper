package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/romiras/url-meta-scraper/consumers/drivers"
	"github.com/romiras/url-meta-scraper/registries"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var urlsChan chan string

func scanURLs(reg *registries.Registry, urlsRcv chan string, payloadsTx chan []byte) {
	for url := range urlsRcv {
		urlScraped, attempts, err := reg.FetchHelper.Try(reg.Fetcher, url, 0)
		if err != nil {
			if attempts == 0 {
				GiveUp(err)
			} else {
				// Retry(reg, url, attempts)
				log.Println("Retry", attempts, url)
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

func newLogger() *logrus.Logger {
	logger := logrus.New()
	if os.Getenv("GO_ENV") == "production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}

const DefaultAmqpURI = "amqp://guest:guest@localhost:5672"

func main() {
	logger := newLogger()

	amqpURI := os.Getenv("AMQP_URI")
	if amqpURI == "" {
		amqpURI = DefaultAmqpURI
	}

	cons, err := drivers.NewAmqpConsumer(amqpURI, "urls", logger)
	if err != nil {
		logger.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		var msgs <-chan interface{}
		logger.Info("Consume")
		msgs, err := cons.Consume()
		if err != nil {
			logger.Fatal(err)
		}

		for msg := range msgs {
			logger.Info("-> msg")
			delivery, ok := msg.(amqp.Delivery)
			if ok {
				logger.Info(string(delivery.Body))
			} else {
				logger.Info("damn")
			}
		}
		done <- true
	}()

	<-done

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
