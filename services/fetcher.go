package services

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

type (
	Bytes []byte

	FetcherIntf interface {
		Fetch(url string) (Bytes, error)
	}

	Fetcher struct {
		Client *http.Client
	}
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Millisecond * 3000,
			}).Dial,
			TLSHandshakeTimeout: time.Millisecond * 5000,
		},
		Timeout: time.Millisecond * 5000,
	}
}

func (f Fetcher) Fetch(url string) (Bytes, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Print("Fetch failed")
		return nil, resp.StatusCode, err
	}

	var out Bytes
	out, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return out, 0, nil
}
