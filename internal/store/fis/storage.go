package fis

import (
	"context"
	"database/sql"
)

// Competitors interface
type Competitors interface {
	GetAthletesBySector(ctx context.Context, sectorCode string) ([]GetBySectorResponse, error)
	GetNationsBySector(ctx context.Context, sectorCode string) ([]string, error)
}

// FISStorage struct to hold table-specific storage
type FISStorage struct {
	competitors Competitors
}

// Methods to return each table's storage interface
func (s *FISStorage) Competitors() Competitors {
	return s.competitors
}

// Storage for FIS database tables
func NewFISStorage(db *sql.DB) *FISStorage {
	return &FISStorage{
		competitors: &CompetitorsStore{db: db},
	}
}
