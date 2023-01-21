package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/html"
)

type Scraper struct {
	Pending   chan string // pending urls to scrap
	Indexed   int64
	Entries   sync.Map
	Whitelist map[string]bool
}

func (s *Scraper) Start() {
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.worker()
		}()
	}
	wg.Wait()
}

func (s *Scraper) worker() {
	for entry := range s.Pending {
		s.scrapOne(entry)
	}
}

func (s *Scraper) scrapOne(scrapUrl string) {

	// log.Println("Scraping:", scrapUrl)

	refenceUrl, _ := url.Parse(scrapUrl)

	resp, err := http.DefaultClient.Get(scrapUrl)
	if err != nil {
		log.Println("scrapOne:", err.Error())
		return
	}
	defer resp.Body.Close()

	// todo: check status code
	if resp.StatusCode == http.StatusInternalServerError {
		newWhitelist := map[string]bool{}
		for host := range s.Whitelist {
			if refenceUrl.Host == host {
				continue
			}
			newWhitelist[host] = true
		}
		s.Whitelist = newWhitelist
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("ERROR: status code", resp.StatusCode, scrapUrl)
	}

	// todo: check headers (content-type, content-length, etc)

	urls := GetUrls(resp.Body)
	s.Entries.Store(scrapUrl, &Entry{
		When: time.Now(),
	})
	atomic.AddInt64(&s.Indexed, 1)

	for _, uu := range urls {
		relativeUrl, err := url.Parse(uu)
		if err != nil {
			log.Println("ERR:", uu, err.Error())
			continue
		}
		absoluteUrl := refenceUrl.ResolveReference(relativeUrl)

		if !s.Whitelist[absoluteUrl.Host] {
			continue
		}

		u := absoluteUrl.String()
		if _, exist := s.Entries.Load(u); exist {
			continue
		}

		s.Entries.Store(u, nil)
		s.Pending <- u // add new url to scrap
	}

}

func NewScraper() *Scraper {
	return &Scraper{
		Pending:   make(chan string, 100000),
		Entries:   sync.Map{},
		Whitelist: nil, // if nil, not whitelist is applied :D
	}
}

type Entry struct {
	When time.Time
	// todo: whatever...
	// Links []string
}

// GetUrls retrieve all urls found on the stream
func GetUrls(r io.Reader) (result []string) {

	d := html.NewTokenizer(r)

	for {
		d.Next()
		t := d.Token()

		switch t.Type {
		case html.ErrorToken:
			return // todo: reason?
		case html.StartTagToken:

			name := strings.ToLower(t.Data)

			if name == "a" {
				for _, attribute := range t.Attr {
					if attribute.Key == "href" {
						result = append(result, attribute.Val)
					}
				}
			}

			if name == "img" {
				for _, attribute := range t.Attr {
					if attribute.Key == "src" {
						result = append(result, attribute.Val)
					}
				}
			}

			// todo: the rest of tags with references, script, link, ...

		default:
			// something else...
		}

	}

	return
}
