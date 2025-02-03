package utv

import (
	"database/sql"
)

// Define UTV interface
type UTV interface {
	// Define methods here (e.g., GetAthletes, GetResults, etc.)
}

// Implement UTV storage
type UTVStorage struct {
	db *sql.DB
}

// NewUTVStorage initializes storage for UTV database tables
func NewUTVStorage(db *sql.DB) *UTVStorage {
	return &UTVStorage{db: db}
}
