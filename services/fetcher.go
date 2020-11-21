package services

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/romiras/url-meta-scraper/pkg/events"
)

type (
	Bytes []byte

	IFetcher interface {
		Fetch(url string) (Bytes, error)
	}

	Fetcher struct {
		Client *http.Client
	}
)

func NewHTTPClient() *http.Client {
	err := loadUserAgents()
	if err != nil {
		log.Fatalln(err)
	}

	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: time.Millisecond * 3000,
			}).Dial,
			TLSHandshakeTimeout: time.Millisecond * 5000,
		},
		Timeout: time.Millisecond * 1500,
	}
}

const MaxFetchBytes = 4096

func (f Fetcher) Fetch(url string) (*events.UrlScrapedEvent, int, error) {
	noStatusCode := int(-1)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, noStatusCode, err
	}

	req.Header.Add("Range", "bytes=0-"+strconv.FormatInt(MaxFetchBytes, 10))

	ua := getRandomUA()
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, noStatusCode, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return nil, resp.StatusCode, nil
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, noStatusCode, err
	}
	sz := len(buf)
	if sz > MaxFetchBytes {
		sz = MaxFetchBytes
	}

	headers := make(map[string]string)
	headers["Content-Encoding"] = getContentType(resp)

	urlScraped := events.UrlScrapedEvent{
		URL:       url,
		UpdatedAt: time.Now().Unix(),
		Headers:   headers,
		Body:      string(buf[:sz]),
	}

	return &urlScraped, resp.StatusCode, nil
}

func getContentType(r *http.Response) string {
	contentType := r.Header.Get("Content-type")
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}

// var out io.Writer
// _, err = io.Copy(out, rd)
// if err != nil {
// 	log.Fatal(err)
// }

// buf, err := ioutil.ReadAll(rd)
// if err != nil {
// 	log.Fatal(err)
// }
