package klab

import (
	"context"
	"database/sql"
	"time"
)

const DataTimeout = 30 * time.Second

// Interfaces
type Users interface {
	GetAllSporttiIDs(ctx context.Context) ([]string, error)
}

// kLABStorage
type KLABStorage struct {
	db    *sql.DB
	users Users
}

// Methods
func (s *KLABStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *KLABStorage) Users() Users {
	return s.users
}

// NewKLABStorage creates a new KLABStorage instance
func NewKLABStorage(db *sql.DB) *KLABStorage {
	return &KLABStorage{
		db:    db,
		users: &UsersStore{db: db},
	}
}
