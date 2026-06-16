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
