package main

import (
	"fmt"

	"github.com/galacticnuclei/concurrent-search-engine/internal/indexer"
)

func main() {
	idx := indexer.New()

	idx.AddDocument(
		1,
		"golang golang golang distributed systems",
	)

	idx.AddDocument(
		2,
		"golang distributed database",
	)

	fmt.Println(idx.Terms["golang"])
	fmt.Println(idx.Terms["distributed"])
}
