package main

import (
	"github.com/romiras/url-meta-scraper/registries"
)

func main() {
	reg := registries.NewRegistry()
	defer reg.Close()

	err := reg.URLSubscriber.Consume(reg.ConsumerHandler, reg.Logger)
	if err != nil {
		reg.Logger.Fatal(err)
		return
	}

	reg.Logger.Info(" [*] Waiting for messages. To exit press CTRL+C")
	select {}

	reg.Logger.Info("shutting down")
}
