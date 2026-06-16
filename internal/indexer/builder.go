package indexer

import (
	"database/sql"

	"github.com/galacticnuclei/concurrent-search-engine/internal/models"
)

func LoadDocuments(db *sql.DB) ([]models.Document, error) {
	rows, err := db.Query(`
		SELECT id, url, title, content
		FROM documents
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.Document

	for rows.Next() {
		var doc models.Document

		err := rows.Scan(
			&doc.ID,
			&doc.URL,
			&doc.Title,
			&doc.Content,
		)
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

func LoadDocumentMap(
	db *sql.DB,
) (map[int]string, error) {
	rows, err := db.Query(`
		SELECT id, url
		FROM documents
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	docMap := make(map[int]string)

	for rows.Next() {
		var id int
		var url string

		err := rows.Scan(&id, &url)
		if err != nil {
			return nil, err
		}

		docMap[id] = url
	}

	return docMap, nil
}
