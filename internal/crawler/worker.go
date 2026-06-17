package crawler

import (
	"fmt"
	"strings"

	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
	"github.com/galacticnuclei/concurrent-search-engine/internal/robots"
	"github.com/galacticnuclei/concurrent-search-engine/internal/storage"
)

func StartWorker(
	id int,
	f *frontier.Frontier,
	visited *frontier.Visited,
	robotsCache *robots.Cache,
	store *storage.DocumentStore,
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

		doc.Content = strings.ToValidUTF8(
			doc.Content,
			"",
		)

		const MaxContentSize = 100000

		if len(doc.Content) > MaxContentSize {
			doc.Content =
				doc.Content[:MaxContentSize]
		}

		changed, err :=
			store.ContentChanged(
				doc.URL,
				doc.Content,
			)

		if err != nil {
			fmt.Printf(
				"Worker %d failed to check content: %v\n",
				id,
				err,
			)
			continue
		}

		if !changed {
			fmt.Printf(
				"Worker %d skipped %s (unchanged)\n",
				id,
				url,
			)
			continue
		}
		if err := store.Save(doc); err != nil {
			fmt.Printf(
				"Worker %d failed to save document: %v\n",
				id,
				err,
			)
		}

		if err := store.SaveLinks(doc); err != nil {
			fmt.Printf(
				"Worker %d failed to save links: %v\n",
				id,
				err,
			)
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
