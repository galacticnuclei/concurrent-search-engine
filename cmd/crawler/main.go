package main

import (
	"time"

	"github.com/galacticnuclei/concurrent-search-engine/internal/crawler"
	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
)

func main() {
	f := frontier.New(100)

	for i := 1; i <= 5; i++ {
		go crawler.StartWorker(i, f)
	}

	f.Push("https://google.com")
	f.Push("https://github.com")
	f.Push("https://golang.org")
	f.Push("https://example.com")

	time.Sleep(5 * time.Second)
}
