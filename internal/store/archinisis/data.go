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

type ArchDataPayload struct {
	Athlete      archsqlc.UpsertAthleteParams
	Measurements []archsqlc.UpsertMeasurementParams
}

type ArchDataResponse struct {
	// Athlete
	NationalID  string   `json:"national_id"`
	FirstName   *string  `json:"first_name,omitempty"`
	LastName    *string  `json:"last_name,omitempty"`
	Initials    *string  `json:"initials,omitempty"`
	DateOfBirth *string  `json:"date_of_birth,omitempty"`
	Height      *float64 `json:"height,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`

	// Measurements
	Measurements []ArchMeasurementResponse `json:"measurements"`
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

func (s *DataStore) GetRaceReportSessions(ctx context.Context, sporttiID string) ([]int32, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)

	rows, err := q.GetRaceReportSessionIDsBySporttiID(ctx, sql.NullString{String: sporttiID, Valid: true})
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

func (s *DataStore) GetRaceReport(ctx context.Context, sporttiID string, sessionID int32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)

	res, err := q.GetRaceReport(ctx, archsqlc.GetRaceReportParams{
		SporttiID: sql.NullString{String: sporttiID, Valid: true},
		SessionID: sql.NullInt32{Int32: sessionID, Valid: true},
	})
	if err != nil {
		return "", err
	}
	if !res.Valid {
		return "", sql.ErrNoRows
	}
	return res.String, nil
}

func (s *DataStore) UpsertRaceReport(ctx context.Context, p archsqlc.UpsertRaceReportParams) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	return q.UpsertRaceReport(ctx, p)
}

func (s *DataStore) UpsertData(ctx context.Context, payload ArchDataPayload) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := archsqlc.New(tx)

	// 1) Athlete
	if err := q.UpsertAthlete(ctx, payload.Athlete); err != nil {
		return err
	}

	// 2) Measurements
	for i := range payload.Measurements {
		if !payload.Measurements[i].NationalID.Valid || payload.Measurements[i].NationalID.String == "" {
			payload.Measurements[i].NationalID = utils.NullString(payload.Athlete.NationalID)
		}
		if err := q.UpsertMeasurement(ctx, payload.Measurements[i]); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *DataStore) GetDataBySporttiID(ctx context.Context, sporttiID string) (*ArchDataResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)

	// 1) Athlete (must exist)
	a, err := q.GetAthleteBySporttiID(ctx, sporttiID)
	if err != nil {
		return nil, err
	}

	// 2) Measurements (may be empty)
	ms, err := q.GetMeasurementsBySporttiID(ctx, utils.NullString(sporttiID))
	if err != nil {
		return nil, err
	}

	resp := ArchDataResponse{
		NationalID:   a.NationalID,
		FirstName:    utils.StringPtrOrNil(a.FirstName),
		LastName:     utils.StringPtrOrNil(a.LastName),
		Initials:     utils.StringPtrOrNil(a.Initials),
		DateOfBirth:  utils.FormatDatePtr(a.DateOfBirth),
		Height:       utils.NullNumericToFloatPtr(a.Height),
		Weight:       utils.NullNumericToFloatPtr(a.Weight),
		Measurements: make([]ArchMeasurementResponse, 0, len(ms)),
	}

	for _, m := range ms {
		resp.Measurements = append(resp.Measurements, ArchMeasurementResponse{
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
		})
	}

	return &resp, nil
}
