package services

import "fmt"

type Scraper struct {
}

func (s Scraper) Scrape(buf []byte) {
	dumpBuf(buf)
}

func dumpBuf(buf []byte) {
	MaxBuf := 10
	len := len(buf)
	if len > MaxBuf {
		len = MaxBuf
	}
	fmt.Println(string(buf[:len]))
}
