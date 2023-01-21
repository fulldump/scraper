package main

import (
	"fmt"
	"time"
)

func main() {

	myScraper := NewScraper()
	// myScraper.Pending <- "https://es.wikipedia.org/wiki/Wikipedia:Portada" // like a seed :D
	myScraper.Pending <- "https://www.treeweb.es/" // like a seed :D
	// myScraper.Pending <- "https://goose.blue/" // like a seed :D

	myScraper.Whitelist = map[string]bool{
		"www.treeweb.es":   true,
		"goose.blue":       true,
		"es.wikipedia.org": true,
	}

	go func() {
		for {
			fmt.Println("Pending:", len(myScraper.Pending))
			fmt.Println("Indexed:", myScraper.Indexed)
			time.Sleep(1 * time.Second)
		}
	}()

	myScraper.Start()
}
