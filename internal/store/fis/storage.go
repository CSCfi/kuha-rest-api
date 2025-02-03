package fis

import (
	"context"
	"database/sql"
)

// Define Competitors interface
type Competitors interface {
	GetAthletesBySector(ctx context.Context, sectorCode string) ([]GetBySectorResponse, error)
	GetNationsBySector(ctx context.Context, sectorCode string) ([]string, error)
}

// Implement the `Competitors()` method to return the interface
func (s *FISStorage) Competitors() Competitors {
	return s.competitors
}

// Ensure `FISStorage` implements `FIS`
type FISStorage struct {
	competitors Competitors
}

// `NewFISStorage` initializes storage for FIS database tables
func NewFISStorage(db *sql.DB) *FISStorage {
	return &FISStorage{
		competitors: &CompetitorsStore{db: db},
	}
}
