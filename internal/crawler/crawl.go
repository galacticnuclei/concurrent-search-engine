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

func Crawl(url string) (*Document, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", url, err)
	}

	title := parser.ExtractTitle(doc)
	content := parser.ExtractContent(doc)
	links := parser.ExtractLinks(doc)
	fmt.Printf(
		"Fetched %s (status: %d)\n",
		url,
		resp.StatusCode,
	)

	return &Document{
		URL:     url,
		Title:   title,
		Content: content,
		Links:   links,
	}, nil
}
