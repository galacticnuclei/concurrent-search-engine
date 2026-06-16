package indexer

import (
	"math"
	"sort"
)

type InvertedIndex struct {
	Terms   map[string]map[int]int
	NumDocs int
}

func New() *InvertedIndex {
	return &InvertedIndex{
		Terms:   make(map[string]map[int]int),
		NumDocs: 0,
	}
}
func (idx *InvertedIndex) AddDocument(
	docID int,
	text string,
) {
	idx.NumDocs++

	tokens := Tokenize(text)

	for _, token := range tokens {
		if _, exists := idx.Terms[token]; !exists {
			idx.Terms[token] = make(map[int]int)
		}

		idx.Terms[token][docID]++
	}
}

func (idx *InvertedIndex) Search(
	term string,
) map[int]int {
	return idx.Terms[term]
}

func (idx *InvertedIndex) IDF(
	term string,
) float64 {
	postings, exists := idx.Terms[term]
	if !exists {
		return 0
	}

	docFreq := len(postings)

	return math.Log(
		float64(idx.NumDocs) /
			float64(docFreq),
	)
}

type Result struct {
	DocID int
	Score float64
}

func (idx *InvertedIndex) SearchRanked(
	term string,
) []Result {
	postings := idx.Terms[term]

	var results []Result

	idf := idx.IDF(term)

	for docID, freq := range postings {
		tf := float64(freq)
		score := tf * idf

		results = append(
			results,
			Result{
				DocID: docID,
				Score: score,
			},
		)
	}

	sort.Slice(
		results,
		func(i, j int) bool {
			return results[i].Score >
				results[j].Score
		},
	)

	return results
}

func (idx *InvertedIndex) SearchQuery(
	query string,
) []Result {
	tokens := Tokenize(query)

	scores := make(map[int]float64)

	for _, term := range tokens {
		postings := idx.Terms[term]
		idf := idx.IDF(term)

		for docID, freq := range postings {
			tf := float64(freq)
			score := tf * idf

			scores[docID] += score
		}
	}

	var results []Result

	for docID, score := range scores {
		results = append(
			results,
			Result{
				DocID: docID,
				Score: score,
			},
		)
	}

	sort.Slice(
		results,
		func(i, j int) bool {
			return results[i].Score >
				results[j].Score
		},
	)

	return results
}

func (idx *InvertedIndex) SearchAND(
	query string,
) []Result {
	tokens := Tokenize(query)

	if len(tokens) == 0 {
		return nil
	}

	counts := make(map[int]int)

	// Count how many query terms each document contains
	for _, term := range tokens {
		postings := idx.Terms[term]

		for docID := range postings {
			counts[docID]++
		}
	}

	// Compute TF-IDF scores only for documents
	// that contain ALL query terms
	scores := make(map[int]float64)

	for docID, count := range counts {
		if count == len(tokens) {
			for _, term := range tokens {
				postings := idx.Terms[term]
				idf := idx.IDF(term)

				freq := postings[docID]
				tf := float64(freq)

				scores[docID] += tf * idf
			}
		}
	}

	// Convert scores map into results slice
	var results []Result

	for docID, score := range scores {
		results = append(
			results,
			Result{
				DocID: docID,
				Score: score,
			},
		)
	}

	// Sort by descending score
	sort.Slice(
		results,
		func(i, j int) bool {
			return results[i].Score >
				results[j].Score
		},
	)

	return results
}

func (idx *InvertedIndex) SearchOR(
	query string,
) []Result {
	return idx.SearchQuery(query)
}
