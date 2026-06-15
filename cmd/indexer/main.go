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

	fmt.Println(idx.SearchRanked("github")[:10])
}
