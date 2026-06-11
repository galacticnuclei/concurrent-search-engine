package crawler

import (
	"fmt"
	"time"

	"github.com/galacticnuclei/concurrent-search-engine/internal/frontier"
)

func StartWorker(id int, f *frontier.Frontier) {
	for url := range f.Pop() {
		fmt.Printf("Worker %d crawling %s\n", id, url)

		time.Sleep(time.Second)
	}
}
