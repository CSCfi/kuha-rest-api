package tietoevry

import (
	"context"
	"database/sql"

	tietoevrysqlc "github.com/DeRuina/KUHA-REST-API/internal/db/tietoevry"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type ExercisePayload struct {
	Exercise tietoevrysqlc.InsertExerciseParams
	HRZones  []tietoevrysqlc.InsertExerciseHRZoneParams
	Samples  []tietoevrysqlc.InsertExerciseSampleParams
	Sections []tietoevrysqlc.InsertExerciseSectionParams
}

type ExercisesStore struct {
	db *sql.DB
}

func NewExercisesStore(db *sql.DB) *ExercisesStore {
	return &ExercisesStore{
		db: db,
	}
}

func (s *ExercisesStore) InsertExerciseBundle(ctx context.Context, input ExercisePayload) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	// Start a transaction to ensure atomicity
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // This will be a no-op if the transaction is committed

	q := tietoevrysqlc.New(tx)

	// 1. Insert base exercise
	if err := q.InsertExercise(ctx, input.Exercise); err != nil {
		return err
	}

	// 2. Insert HR zones
	for _, zone := range input.HRZones {
		if err := q.InsertExerciseHRZone(ctx, zone); err != nil {
			return err
		}
	}

	// 3. Insert samples
	for _, sample := range input.Samples {
		if err := q.InsertExerciseSample(ctx, sample); err != nil {
			return err
		}
	}

	// 4. Insert sections
	for _, section := range input.Sections {
		if err := q.InsertExerciseSection(ctx, section); err != nil {
			return err
		}
	}

	// Commit the transaction if all inserts succeeded
	return tx.Commit()
}
