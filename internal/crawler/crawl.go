package crawler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/galacticnuclei/concurrent-search-engine/internal/models"
	"github.com/galacticnuclei/concurrent-search-engine/internal/parser"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Crawl(pageURL string) (*models.Document, error) {
	resp, err := client.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("error fetching %s: %w", pageURL, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", pageURL, err)
	}

	title := parser.ExtractTitle(doc)
	content := parser.ExtractContent(doc)

	rawLinks := parser.ExtractLinks(doc)

	var links []string

	for _, link := range rawLinks {
		normalized, ok := parser.NormalizeURL(pageURL, link)
		if ok {
			links = append(links, normalized)
		}
	}

	fmt.Printf(
		"Fetched %s (status: %d)\n",
		pageURL,
		resp.StatusCode,
	)

	return &models.Document{
		URL:     pageURL,
		Title:   title,
		Content: content,
		Links:   links,
	}, nil
}
