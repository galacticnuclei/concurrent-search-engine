package indexer

import "database/sql"

func LoadLinkGraph(
	db *sql.DB,
) (map[string][]string, error) {
	rows, err := db.Query(`
		SELECT from_url, to_url
		FROM links
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	graph := make(
		map[string][]string,
	)

	for rows.Next() {
		var from string
		var to string

		err := rows.Scan(
			&from,
			&to,
		)
		if err != nil {
			return nil, err
		}

		graph[from] =
			append(
				graph[from],
				to,
			)
	}

	return graph, nil
}

func ComputePageRank(
	graph map[string][]string,
	iterations int,
	damping float64,
) map[string]float64 {
	ranks := make(map[string]float64)

	n := len(graph)
	if n == 0 {
		return ranks
	}

	initialRank := 1.0 / float64(n)

	for page := range graph {
		ranks[page] = initialRank
	}

	for i := 0; i < iterations; i++ {
		newRanks := make(map[string]float64)

		for page := range graph {
			newRanks[page] =
				(1.0 - damping) / float64(n)
		}

		for page, links := range graph {
			if len(links) == 0 {
				continue
			}

			share :=
				ranks[page] /
					float64(len(links))

			for _, link := range links {
				newRanks[link] +=
					damping * share
			}
		}

		ranks = newRanks
	}

	return ranks
}
