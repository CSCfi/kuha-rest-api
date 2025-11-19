package fis

import (
	"context"
	"database/sql"

	fissqlc "github.com/DeRuina/KUHA-REST-API/internal/db/fis"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type RaceNKStore struct {
	db *sql.DB
}

func (s *RaceNKStore) GetNordicCombinedSeasons(ctx context.Context) ([]int32, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	rows, err := q.GetNordicCombinedSeasons(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]int32, 0, len(rows))
	for _, v := range rows {
		if v.Valid {
			out = append(out, v.Int32)
		}
	}
	return out, nil
}

func (s *RaceNKStore) GetNordicCombinedDisciplines(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	rows, err := q.GetNordicCombinedDisciplines(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(rows))
	for _, v := range rows {
		if v.Valid {
			out = append(out, v.String)
		}
	}
	return out, nil
}

func (s *RaceNKStore) GetNordicCombinedCategories(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	rows, err := q.GetNordicCombinedCategories(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(rows))
	for _, v := range rows {
		if v.Valid {
			out = append(out, v.String)
		}
	}
	return out, nil
}

func (s *RaceNKStore) GetRacesNK(ctx context.Context, seasons []int32, disciplines, cats []string) ([]fissqlc.ARacenk, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	params := fissqlc.GetRacesNKParams{
		Column1: seasons,
		Column2: disciplines,
		Column3: cats,
	}
	return q.GetRacesNK(ctx, params)
}

func (s *RaceNKStore) GetLastRowRaceNK(ctx context.Context) (fissqlc.ARacenk, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	return q.GetLastRowRaceNK(ctx)
}

func (s *RaceNKStore) InsertRaceNK(ctx context.Context, in InsertRaceNKClean) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	return q.InsertRaceNK(ctx, mapInsertRaceNKToParams(in))
}

func (s *RaceNKStore) UpdateRaceNKByID(ctx context.Context, in UpdateRaceNKClean) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	_, err := q.UpdateRaceNKByID(ctx, mapUpdateRaceNKToParams(in))
	return err
}

func (s *RaceNKStore) DeleteRaceNKByID(ctx context.Context, raceID int32) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := fissqlc.New(s.db)
	_, err := q.DeleteRaceNKByID(ctx, raceID)
	return err
}
