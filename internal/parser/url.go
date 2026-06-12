package parser

import (
	"net/url"
	"strings"
)

func NormalizeURL(baseURL, href string) (string, bool) {
	href = strings.TrimSpace(href)

	if href == "" {
		return "", false
	}

	// Skip page fragments
	if strings.HasPrefix(href, "#") {
		return "", false
	}

	// Skip javascript links
	if strings.HasPrefix(href, "javascript:") {
		return "", false
	}

	// Skip email links
	if strings.HasPrefix(href, "mailto:") {
		return "", false
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", false
	}

	link, err := url.Parse(href)
	if err != nil {
		return "", false
	}

	resolved := base.ResolveReference(link)

	// Only crawl HTTP/S pages
	if resolved.Scheme != "http" && resolved.Scheme != "https" {
		return "", false
	}

	// Remove fragments like:
	// https://site.com/page#section
	resolved.Fragment = ""

	return resolved.String(), true
}
