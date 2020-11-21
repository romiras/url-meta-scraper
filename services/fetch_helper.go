package services

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"syscall"

	"github.com/romiras/url-meta-scraper/pkg"
	"github.com/romiras/url-meta-scraper/pkg/events"
)

// Dispatcher ?

type FetchHelper struct{}

const MaxAttempts = 3

func (hlp FetchHelper) Try(fetcher *Fetcher, url string, attempts uint) (*events.UrlScrapedEvent, uint, error) {
	urlScraped, statusCode, err := fetcher.Fetch(url)
	if err != nil || isSuccess(statusCode) {
		attempts++
		if err != nil {
			log.Print("\tError:", err.Error())
		} else {
			log.Print("\tFetch failed")
		}
		if CanRetry(statusCode, attempts, err) {
			return nil, attempts, err
		}
		return nil, pkg.NoRetry, err // give up
	}
	return urlScraped, pkg.NoRetry, nil
}

func CanRetry(statusCode int, attempts uint, err error) bool {
	if errors.Is(err, syscall.ECONNREFUSED) {
		fmt.Println("ECONNREFUSED")
		return false
	}
	switch {
	case statusCode >= 300 && statusCode < 500:
		return false
	case statusCode >= 500 && statusCode != 502 && statusCode != 503:
		return attempts < MaxAttempts
	}
	return true
}

func isSuccess(statusCode int) bool {
	return statusCode != http.StatusOK && statusCode != http.StatusPartialContent
}
