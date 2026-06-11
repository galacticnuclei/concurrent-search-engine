package crawler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/galacticnuclei/concurrent-search-engine/internal/parser"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Crawl(url string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Failed to parse %s\n", url)
		return
	}

	title := parser.ExtractTitle(doc)

	fmt.Printf(
		"Fetched %s (status: %d)\n",
		url,
		resp.StatusCode,
	)
	fmt.Printf("Title: %s\n", title)
}
