package tietoevry

import (
	"context"
	"database/sql"

	tietoevrysqlc "github.com/DeRuina/KUHA-REST-API/internal/db/tietoevry"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type SymptomsStore struct {
	db *sql.DB
}

func NewSymptomsStore(db *sql.DB) *SymptomsStore {
	return &SymptomsStore{db: db}
}

func (s *SymptomsStore) InsertSymptom(ctx context.Context, arg tietoevrysqlc.InsertSymptomParams) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := tietoevrysqlc.New(s.db)
	return q.InsertSymptom(ctx, arg)
}
