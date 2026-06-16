package indexer

import (
	"math"
	"sort"
)

type InvertedIndex struct {
	Terms       map[string]map[int]int
	Positions   map[string]map[int][]int
	NumDocs     int
	DocumentURL map[int]string
	PageRanks   map[string]float64
}

func New() *InvertedIndex {
	return &InvertedIndex{
		Terms:       make(map[string]map[int]int),
		Positions:   make(map[string]map[int][]int),
		NumDocs:     0,
		DocumentURL: make(map[int]string),
		PageRanks:   make(map[string]float64),
	}
}

func (idx *InvertedIndex) AddDocument(
	docID int,
	text string,
) {
	idx.NumDocs++

	tokens := Tokenize(text)

	for position, token := range tokens {
		if _, exists := idx.Terms[token]; !exists {
			idx.Terms[token] = make(map[int]int)
		}

		idx.Terms[token][docID]++

		if _, exists := idx.Positions[token]; !exists {
			idx.Positions[token] =
				make(map[int][]int)
		}

		idx.Positions[token][docID] =
			append(
				idx.Positions[token][docID],
				position,
			)
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
	DocID    int
	TFIDF    float64
	PageRank float64
	Score    float64
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
				TFIDF: score,
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
		url := idx.DocumentURL[docID]
		pageRank := idx.PageRanks[url]

		finalScore := CombineScore(
			score,
			pageRank,
		)

		results = append(
			results,
			Result{
				DocID:    docID,
				TFIDF:    score,
				PageRank: pageRank,
				Score:    finalScore,
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
				TFIDF: score,
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

func (idx *InvertedIndex) SearchPhrase(
	query string,
) []Result {
	tokens := Tokenize(query)

	if len(tokens) != 2 {
		return nil
	}

	first := tokens[0]
	second := tokens[1]

	firstDocs := idx.Positions[first]
	secondDocs := idx.Positions[second]

	var results []Result

	for docID := range firstDocs {
		_, exists := secondDocs[docID]
		if !exists {
			continue
		}

		firstPositions :=
			firstDocs[docID]

		secondPositions :=
			secondDocs[docID]

		found := false

		for _, p1 := range firstPositions {
			for _, p2 := range secondPositions {
				if p2 == p1+1 {
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		if found {
			results = append(
				results,
				Result{
					DocID: docID,
					Score: 1,
				},
			)
		}
	}

	return results
}
