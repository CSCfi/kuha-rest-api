package tietoevry

import (
	"context"
	"database/sql"
)

// TietoevryStorage
type TietoevryStorage struct {
	db *sql.DB
}

// Methods
func (s *TietoevryStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

// NewTietoevryStorage creates a new TietoevryStorage instance
func NewTietoevryStorage(db *sql.DB) *TietoevryStorage {
	return &TietoevryStorage{
		db: db,
	}
}
