package services

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/romiras/url-meta-scraper/pkg"
)

// Dispatcher ?

type FetchHelper struct {
}

const MaxAttempts = 3

func (hlp FetchHelper) Try(fetcher *Fetcher, url string, attempts uint) (*pkg.UrlScraped, uint, error) {
	urlScraped, statusCode, err := fetcher.Fetch(url)
	attempts++
	// fmt.Println("\t", attempts, statusCode, url)
	if err != nil {
		if CanRetry(statusCode, attempts, err) {
			// Retry(url, attempts)
			return nil, attempts, err
		} else {
			return nil, 0, err // give up
		}
	}

	if urlScraped == nil {
		return nil, 0, errors.New("UrlScraped is nil")
	}

	return urlScraped, 0, nil
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
