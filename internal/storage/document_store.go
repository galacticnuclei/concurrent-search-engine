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
	contentHash := HashContent(doc.Content)

	_, err := s.db.Exec(
		`
		INSERT INTO documents (
			url,
			title,
			content,
			content_hash
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (content_hash) DO NOTHING
		`,
		doc.URL,
		doc.Title,
		doc.Content,
		contentHash,
	)

	return err
}

func (s *DocumentStore) SaveLinks(
	doc *models.Document,
) error {
	for _, link := range doc.Links {
		_, err := s.db.Exec(
			`
			INSERT INTO links (
				from_url,
				to_url
			)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
			`,
			doc.URL,
			link,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
