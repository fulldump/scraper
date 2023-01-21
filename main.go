package main

import (
	"fmt"
	"time"
)

func main() {

	myScraper := NewScraper()
	myScraper.Pending <- "https://es.wikipedia.org/wiki/Wikipedia:Portada" // like a seed :D

	myScraper.Whitelist = map[string]bool{
		"www.treeweb.es":   true,
		"goose.blue":       true,
		"es.wikipedia.org": true,
	}

	go func() {
		for {
			fmt.Println("Pending:", len(myScraper.Pending))
			fmt.Printf("Entries: %v\n", len(myScraper.Entries))
			time.Sleep(3 * time.Second)
		}
	}()

	myScraper.Start()
}
