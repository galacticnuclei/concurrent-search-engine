package main

import (
	"fmt"
	"log"

	"github.com/galacticnuclei/concurrent-search-engine/internal/indexer"
	"github.com/galacticnuclei/concurrent-search-engine/internal/storage"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=psql@1428*# dbname=search_engine sslmode=disable"

	db, err := storage.NewPostgres(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	docs, err := indexer.LoadDocuments(db)
	graph, err := indexer.LoadLinkGraph(db)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"Loaded %d pages with outgoing links\n",
		len(graph),
	)
	if err != nil {
		log.Fatal(err)
	}
	ranks := indexer.ComputePageRank(
		graph,
		20,
		0.85,
	)

	fmt.Printf(
		"Computed PageRank for %d pages\n",
		len(ranks),
	)

	docMap, err := indexer.LoadDocumentMap(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d documents\n", len(docs))
	idx := indexer.New()
	idx.DocumentURL = docMap
	idx.PageRanks = ranks
	for _, doc := range docs {
		idx.AddDocument(
			doc.ID,
			doc.Content,
		)
	}

	fmt.Printf(
		"Indexed %d unique terms\n",
		len(idx.Terms),
	)

	results := idx.SearchPhrase("github copilot")

	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		result := results[i]

		fmt.Printf(
			"%d. %s\n   score: %.4f\n\n",
			i+1,
			docMap[result.DocID],
			result.Score,
		)
	}
}
