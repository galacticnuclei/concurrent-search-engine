package crawler

import (
	"fmt"

	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
)

func StartWorker(id int, f *frontier.Frontier) {
	for url := range f.Pop() {
		fmt.Printf("Worker %d received %s\n", id, url)

		Crawl(url)
	}
}
