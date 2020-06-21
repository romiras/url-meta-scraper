package services

import (
	"errors"
	"syscall"
)

type FetchHelper struct {
}

const MaxAttempts = 3

func (hlp FetchHelper) Try(fetcher Fetcher, url string, attempts uint) ([]byte, error) {
	buf, statusCode, err := fetcher.Fetch(url)
	attempts++
	if err != nil {
		if CanRetry(statusCode, attempts, err) {
			Retry(url, attempts)
		} else {
			return nil, err
		}
	}

	return buf, nil
}

func CanRetry(statusCode int, attempts uint, err error) bool {
	if errors.Is(err, syscall.ECONNREFUSED) {
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

func Retry(url string, attempts uint) {
	// TODO: (URL, attempt nr.) -> queue 'failed-urls'
}
