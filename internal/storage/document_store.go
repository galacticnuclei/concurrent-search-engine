package storage

import (
	"database/sql"

	"github.com/galacticnuclei/concurrent-search-engine/internal/models"
)

type DocumentStore struct {
	db *sql.DB
}

func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore{
		db: db,
	}
}

func (s *DocumentStore) Save(doc *models.Document) error {
	_, err := s.db.Exec(
		`
		INSERT INTO documents (
			url,
			title,
			content
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (url) DO NOTHING
		`,
		doc.URL,
		doc.Title,
		doc.Content,
	)

	return err
}
