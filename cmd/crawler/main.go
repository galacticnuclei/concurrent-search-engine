package main

import (
	"time"

	"github.com/galacticnuclei/concurrent-search-engine/internal/crawler"
	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
	"github.com/galacticnuclei/concurrent-search-engine/internal/robots"
)

func main() {
	f := frontier.New(10000)
	visited := frontier.NewVisited()
	robotsCache := robots.NewCache()
	for i := 1; i <= 5; i++ {
		go crawler.StartWorker(i, f, visited, robotsCache)
	}

	seedURLs := []string{
		"https://google.com",
		"https://github.com",
		"https://golang.org",
		"https://example.com",
	}

	for _, url := range seedURLs {
		if visited.AddIfNotSeen(url) {
			f.Push(url)
		}
	}

	time.Sleep(60 * time.Second)
}
