package crawler

import (
	"fmt"

	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
	"github.com/galacticnuclei/concurrent-search-engine/internal/robots"
)

func StartWorker(
	id int,
	f *frontier.Frontier,
	visited *frontier.Visited,
	robotsCache *robots.Cache,
) {
	for url := range f.Pop() {
		fmt.Printf("Worker %d received %s\n", id, url)

		if !robotsCache.Allowed(url) {
			fmt.Printf(
				"Worker %d skipped %s (robots.txt)\n",
				id,
				url,
			)
			continue
		}

		doc, err := Crawl(url)
		if err != nil {
			fmt.Printf("Worker %d error: %v\n", id, err)
			continue
		}

		fmt.Printf(
			"Document: %s (%d chars)\n",
			doc.Title,
			len(doc.Content),
		)

		fmt.Printf(
			"Links Found: %d\n",
			len(doc.Links),
		)

		newURLs := 0

		for _, link := range doc.Links {
			if visited.AddIfNotSeen(link) {
				f.Push(link)
				newURLs++
			}
		}

		fmt.Printf(
			"Discovered %d new URLs\n",
			newURLs,
		)

		fmt.Printf(
			"Visited Count: %d\n",
			visited.Count(),
		)
	}
}
