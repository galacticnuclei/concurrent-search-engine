package indexer

func CombineScore(
	tfidf float64,
	pageRank float64,
) float64 {
	return 0.8*tfidf +
		0.2*pageRank
}
