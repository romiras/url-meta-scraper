package services

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/romiras/url-meta-scraper/pkg"
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
		Timeout: time.Millisecond * 50,
	}
}

const MaxFetchBytes = 4096

func (f Fetcher) Fetch(url string) (*pkg.UrlScraped, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Range", "bytes=0-"+strconv.FormatInt(MaxFetchBytes, 10))

	ua := getRandomUA()
	if ua != "" {
		req.Header.Add("User-Agent", ua)
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		log.Print("Fetch failed")
		return nil, resp.StatusCode, err
	}

	urlScraped := pkg.UrlScraped{}
	urlScraped.Headers = make(map[string]string)
	urlScraped.Headers["Content-Encoding"] = getContentType(resp)
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	sz := len(buf)
	if sz > MaxFetchBytes {
		sz = MaxFetchBytes
	}
	urlScraped.Body = string(buf[:sz])

	return &urlScraped, 0, nil
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
