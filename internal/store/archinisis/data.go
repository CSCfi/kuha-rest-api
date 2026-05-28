package archinisis

import (
	"context"
	"database/sql"

	archsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/archinisis"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type DataStore struct {
	db *sql.DB
}

type ArchAthleteResponse struct {
	NationalID  string   `json:"national_id"`
	FirstName   *string  `json:"first_name,omitempty"`
	LastName    *string  `json:"last_name,omitempty"`
	Initials    *string  `json:"initials,omitempty"`
	DateOfBirth *string  `json:"date_of_birth,omitempty"`
	Height      *float64 `json:"height,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
}

type ArchMeasurementResponse struct {
	MeasurementGroupID int32   `json:"measurement_group_id"`
	MeasurementID      *int32  `json:"measurement_id,omitempty"`
	Discipline         *string `json:"discipline,omitempty"`
	SessionName        *string `json:"session_name,omitempty"`
	Place              *string `json:"place,omitempty"`
	RaceID             *int32  `json:"race_id,omitempty"`
	StartTime          *string `json:"start_time,omitempty"`
	StopTime           *string `json:"stop_time,omitempty"`
	NbSegments         *int32  `json:"nb_segments,omitempty"`
	Comment            *string `json:"comment,omitempty"`
}

// Race report methods

func (s *DataStore) GetRaceReportSessions(ctx context.Context, sporttiID string) ([]int32, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	rows, err := q.GetRaceReportSessionIDsBySporttiID(ctx, sporttiID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *DataStore) GetRaceReport(ctx context.Context, sporttiID string, sessionID int32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	return q.GetRaceReport(ctx, archsqlc.GetRaceReportParams{
		SporttiID: sporttiID,
		SessionID: sessionID,
	})
}

func (s *DataStore) UpsertRaceReport(ctx context.Context, sporttiID string, sessionID int32, raceReport string) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := archsqlc.New(tx)

	if err := q.UpsertReport(ctx, archsqlc.UpsertReportParams{
		SessionID:  sessionID,
		RaceReport: raceReport,
	}); err != nil {
		return err
	}

	if err := q.UpsertReportUser(ctx, archsqlc.UpsertReportUserParams{
		SessionID: sessionID,
		SporttiID: sporttiID,
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DataStore) GetSporttiIDsBySessionID(ctx context.Context, sessionID int32) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	return q.GetSporttiIDsBySessionID(ctx, sessionID)
}

// Athlete methods

func (s *DataStore) UpsertAthleteOnly(ctx context.Context, athlete archsqlc.UpsertAthleteParams) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	return q.UpsertAthlete(ctx, athlete)
}

func (s *DataStore) GetAthleteByID(ctx context.Context, sporttiID string) (*ArchAthleteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	a, err := q.GetAthleteBySporttiID(ctx, sporttiID)
	if err != nil {
		return nil, err
	}

	return &ArchAthleteResponse{
		NationalID:  a.NationalID,
		FirstName:   utils.StringPtrOrNil(a.FirstName),
		LastName:    utils.StringPtrOrNil(a.LastName),
		Initials:    utils.StringPtrOrNil(a.Initials),
		DateOfBirth: utils.FormatDatePtr(a.DateOfBirth),
		Height:      utils.NullNumericToFloatPtr(a.Height),
		Weight:      utils.NullNumericToFloatPtr(a.Weight),
	}, nil
}

// Measurement methods

func (s *DataStore) UpsertMeasurements(ctx context.Context, sporttiID string, measurements []archsqlc.UpsertMeasurementParams) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := archsqlc.New(tx)

	for i := range measurements {
		if !measurements[i].NationalID.Valid || measurements[i].NationalID.String == "" {
			measurements[i].NationalID = utils.NullString(sporttiID)
		}
		if err := q.UpsertMeasurement(ctx, measurements[i]); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DataStore) GetMeasurementsByID(ctx context.Context, sporttiID string) ([]ArchMeasurementResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	ms, err := q.GetMeasurementsBySporttiID(ctx, utils.NullString(sporttiID))
	if err != nil {
		return nil, err
	}

	result := make([]ArchMeasurementResponse, 0, len(ms))
	for _, m := range ms {
		result = append(result, toMeasurementResponse(m))
	}
	return result, nil
}

func (s *DataStore) GetMeasurementByMeasurementID(ctx context.Context, measurementID int32) (*ArchMeasurementResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	const query = `
SELECT measurement_group_id, measurement_id, national_id, discipline, session_name,
       place, race_id, start_time, stop_time, nb_segments, comment
FROM measurement WHERE measurement_id = $1`

	var m archsqlc.Measurement
	row := s.db.QueryRowContext(ctx, query, measurementID)
	if err := row.Scan(
		&m.MeasurementGroupID,
		&m.MeasurementID,
		&m.NationalID,
		&m.Discipline,
		&m.SessionName,
		&m.Place,
		&m.RaceID,
		&m.StartTime,
		&m.StopTime,
		&m.NbSegments,
		&m.Comment,
	); err != nil {
		return nil, err
	}

	resp := toMeasurementResponse(m)
	return &resp, nil
}

func (s *DataStore) DeleteMeasurementByMeasurementID(ctx context.Context, measurementID int32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	const query = `DELETE FROM measurement WHERE measurement_id = $1 RETURNING national_id`

	var nationalID sql.NullString
	row := s.db.QueryRowContext(ctx, query, measurementID)
	if err := row.Scan(&nationalID); err != nil {
		return "", err
	}
	return nationalID.String, nil
}

func toMeasurementResponse(m archsqlc.Measurement) ArchMeasurementResponse {
	return ArchMeasurementResponse{
		MeasurementGroupID: m.MeasurementGroupID,
		MeasurementID:      utils.Int32PtrOrNil(m.MeasurementID),
		Discipline:         utils.StringPtrOrNil(m.Discipline),
		SessionName:        utils.StringPtrOrNil(m.SessionName),
		Place:              utils.StringPtrOrNil(m.Place),
		RaceID:             utils.Int32PtrOrNil(m.RaceID),
		StartTime:          utils.FormatTimestampPtr(m.StartTime),
		StopTime:           utils.FormatTimestampPtr(m.StopTime),
		NbSegments:         utils.Int32PtrOrNil(m.NbSegments),
		Comment:            utils.StringPtrOrNil(m.Comment),
	}
}
