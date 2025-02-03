package store

import (
	"github.com/DeRuina/KUHA-REST-API/internal/db"
	"github.com/DeRuina/KUHA-REST-API/internal/store/fis"
	"github.com/DeRuina/KUHA-REST-API/internal/store/utv"
)

// Define database interfaces
type FIS interface {
	Competitors() fis.Competitors
}

type UTV interface {
	// Add UTV methods here
}

// Define Storage struct for multiple databases
type Storage struct {
	FIS FIS
	UTV UTV
}

// NewStorage initializes storage for multiple databases
func NewStorage(databases *db.Database) *Storage { // ✅ Accept *db.Database directly
	return &Storage{
		FIS: fis.NewFISStorage(databases.FIS), // ✅ Correctly assigns FIS DB
		UTV: utv.NewUTVStorage(databases.UTV), // ✅ Correctly assigns UTV DB
	}
}
