package main

import (
	"log"
	"time"

	"github.com/galacticnuclei/concurrent-search-engine/internal/crawler"
	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
	"github.com/galacticnuclei/concurrent-search-engine/internal/robots"
	"github.com/galacticnuclei/concurrent-search-engine/internal/storage"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=psql@1428*# dbname=search_engine sslmode=disable"

	db, err := storage.NewPostgres(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := storage.NewDocumentStore(db)

	f := frontier.New(10000)
	visited := frontier.NewVisited()
	robotsCache := robots.NewCache()

	for i := 1; i <= 5; i++ {
		go crawler.StartWorker(
			i,
			f,
			visited,
			robotsCache,
			store,
		)
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
