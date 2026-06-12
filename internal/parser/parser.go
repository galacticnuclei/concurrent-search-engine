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
