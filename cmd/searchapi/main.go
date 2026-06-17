package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/galacticnuclei/concurrent-search-engine/internal/indexer"
	"github.com/galacticnuclei/concurrent-search-engine/internal/storage"
)

type SearchResponse struct {
	URL      string  `json:"url"`
	Score    float64 `json:"score"`
	TFIDF    float64 `json:"tfidf"`
	PageRank float64 `json:"pagerank"`
}

func main() {
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := storage.NewPostgres(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	docs, err := indexer.LoadDocuments(db)
	if err != nil {
		log.Fatal(err)
	}

	graph, err := indexer.LoadLinkGraph(db)
	if err != nil {
		log.Fatal(err)
	}

	ranks := indexer.ComputePageRank(
		graph,
		20,
		0.85,
	)

	docMap, err := indexer.LoadDocumentMap(db)
	if err != nil {
		log.Fatal(err)
	}

	idx := indexer.New()
	idx.DocumentURL = docMap
	idx.PageRanks = ranks

	for _, doc := range docs {
		idx.AddDocument(
			doc.ID,
			doc.Content,
		)
	}

	for url, rank := range ranks {
		log.Printf(
			"%s -> %.6f",
			url,
			rank,
		)
		break
	}

	log.Printf(
		"PageRank entries: %d",
		len(ranks),
	)

	http.HandleFunc(
		"/search",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			query :=
				r.URL.Query().Get("q")

			results :=
				idx.SearchQuery(query)

			var response []SearchResponse

			limit := 10
			if len(results) < limit {
				limit = len(results)
			}

			for i := 0; i < limit; i++ {
				result := results[i]

				response = append(
					response,
					SearchResponse{
						URL:      docMap[result.DocID],
						Score:    result.Score,
						TFIDF:    result.TFIDF,
						PageRank: result.PageRank,
					},
				)
			}

			latency :=
				time.Since(start)

			log.Printf(
				"Query=%q Latency=%v",
				query,
				latency,
			)

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			json.NewEncoder(w).
				Encode(response)
		},
	)

	log.Println(
		"Search API running on :8080",
	)

	log.Fatal(
		http.ListenAndServe(
			":8080",
			nil,
		),
	)
}
