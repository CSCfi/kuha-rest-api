package fis

import (
	"context"
	"database/sql"
)

// Define Competitors interface
type Competitors interface {
	GetAthletesBySector(ctx context.Context, sectorCode string) ([]GetBySectorResponse, error)
	GetNationsBySector(ctx context.Context, sectorCode string) ([]string, error)
	GetByFiscodeJP(ctx context.Context, fiscode int32) (int32, error)
	GetByFiscodeNK(ctx context.Context, fiscode int32) (int32, error)
	GetByGenderAndNationJP(ctx context.Context, gender, nation string) ([]int32, error)
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
