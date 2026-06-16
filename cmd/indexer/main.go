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
	if err != nil {
		log.Fatal(err)
	}
	docMap, err := indexer.LoadDocumentMap(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Loaded %d documents\n", len(docs))
	idx := indexer.New()

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

	results := idx.SearchRanked("github")

	limit := 10
	if len(results) < limit {
		limit = len(results)
	}

	for i := 0; i < limit; i++ {
		result := results[i]

		fmt.Printf(
			"%d. %s\n   score: %d\n\n",
			i+1,
			docMap[result.DocID],
			result.Score,
		)
	}
}
