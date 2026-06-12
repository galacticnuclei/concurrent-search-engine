package crawler

import (
	"fmt"

	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
)

func StartWorker(id int, f *frontier.Frontier) {
	for url := range f.Pop() {
		fmt.Printf("Worker %d received %s\n", id, url)

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
	}
}
