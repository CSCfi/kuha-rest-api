package utv

import (
	"context"
	"database/sql"
	"encoding/json"
)

// OuraData interface
type OuraData interface {
	GetDates(ctx context.Context, userID string, startDate *string, endDate *string) ([]string, error)
	GetTypes(ctx context.Context, userID string, summaryDate string) ([]string, error)
	GetData(ctx context.Context, userID string, summaryDate string, key *string) (json.RawMessage, error)
}

// PolarData interface
type PolarData interface {
	GetDates(ctx context.Context, userID string, startDate *string, endDate *string) ([]string, error)
	GetTypes(ctx context.Context, userID string, summaryDate string) ([]string, error)
	GetData(ctx context.Context, userID string, summaryDate string, key *string) (json.RawMessage, error)
}

// SuuntoData interface
type SuuntoData interface {
	GetDates(ctx context.Context, userID string, startDate *string, endDate *string) ([]string, error)
	GetTypes(ctx context.Context, userID string, summaryDate string) ([]string, error)
	GetData(ctx context.Context, userID string, summaryDate string, key *string) (json.RawMessage, error)
}

// UTVStorage struct to hold table-specific storage
type UTVStorage struct {
	oura   OuraData
	polar  PolarData
	suunto SuuntoData
}

// Methods to return each table's storage interface
func (s *UTVStorage) Oura() OuraData {
	return s.oura
}

func (s *UTVStorage) Polar() PolarData {
	return s.polar
}

func (s *UTVStorage) Suunto() SuuntoData {
	return s.suunto
}

// Storage for UTV database tables
func NewUTVStorage(db *sql.DB) *UTVStorage {
	return &UTVStorage{
		oura:   &OuraDataStore{db: db},
		polar:  &PolarDataStore{db: db},
		suunto: &SuuntoDataStore{db: db},
	}
}
