package parser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExtractTitle(doc *goquery.Document) string {
	return strings.TrimSpace(
		doc.Find("title").First().Text(),
	)
}

func ExtractContent(doc *goquery.Document) string {
	return strings.TrimSpace(
		doc.Find("body").Text(),
	)
}

func ExtractLinks(doc *goquery.Document) []string {
	var links []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})

	return links
}
