package main

import (
	"log"

	"github.com/romiras/url-meta-scraper/registries"
)

func Do(reg *registries.Registry) {
	urls := []string{
		"https://golang.org/robots.txt",
		"http://localhost:8123/test",
		"http://localhost/robots.txt",
	}

	for _, url := range urls {
		buf, err := reg.FetchHelper.Try(reg.Fetcher, url, 0)
		if err != nil {
			GiveUp(err)
		}
		Done(buf)
	}
}

func Done(buf []byte) {
	// TODO: (buf) -> queue 'scraped-urls'
}

func GiveUp(err error) {
	log.Println(err.Error())
}

func main() {
	reg := registries.NewRegistry()
	Do(reg)
}
