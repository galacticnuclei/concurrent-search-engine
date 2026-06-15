package indexer

import "sort"

type InvertedIndex struct {
	Terms map[string]map[int]int
}

func New() *InvertedIndex {
	return &InvertedIndex{
		Terms: make(map[string]map[int]int),
	}
}
func (idx *InvertedIndex) AddDocument(
	docID int,
	text string,
) {
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

type Result struct {
	DocID int
	Score int
}

func (idx *InvertedIndex) SearchRanked(
	term string,
) []Result {
	postings := idx.Terms[term]

	var results []Result

	for docID, freq := range postings {
		results = append(
			results,
			Result{
				DocID: docID,
				Score: freq,
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
