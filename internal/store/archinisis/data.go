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
	MeasurementGroupID *int32  `json:"measurement_group_id,omitempty"`
	MeasurementID      int32   `json:"measurement_id"`
	Discipline         *string `json:"discipline,omitempty"`
	SessionName        *string `json:"session_name,omitempty"`
	Place              *string `json:"place,omitempty"`
	RaceID             *int32  `json:"race_id,omitempty"`
	StartTime          *string `json:"start_time,omitempty"`
	StopTime           *string `json:"stop_time,omitempty"`
	NbSegments         *int32  `json:"nb_segments,omitempty"`
	Comment            *string `json:"comment,omitempty"`
}

type ArchDataResponse struct {
	ArchAthleteResponse
	Measurements []ArchMeasurementResponse `json:"measurements"`
}

// Race report methods

func (s *DataStore) GetRaceReportSessions(ctx context.Context, sporttiID string) ([]int32, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	rows, err := q.GetRaceReportSessionIDsBySporttiID(ctx, utils.NullString(sporttiID))
	if err != nil {
		return nil, err
	}

	out := make([]int32, 0, len(rows))
	for _, r := range rows {
		if r.Valid {
			out = append(out, r.Int32)
		}
	}
	return out, nil
}

func (s *DataStore) GetRaceReport(ctx context.Context, sporttiID string, sessionID int32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	res, err := q.GetRaceReport(ctx, archsqlc.GetRaceReportParams{
		SporttiID: utils.NullString(sporttiID),
		SessionID: utils.NullInt32(sessionID),
	})
	if err != nil {
		return "", err
	}
	return res.String, nil
}

// UpsertRaceReport stores one report row per (sportti_id, session_id).
// The report table has no unique constraint on that pair, so the upsert is
// an UPDATE followed by an INSERT when no row was updated, inside one transaction.
func (s *DataStore) UpsertRaceReport(ctx context.Context, sporttiID string, sessionID int32, raceReport string) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := archsqlc.New(tx)

	updated, err := q.UpdateReport(ctx, archsqlc.UpdateReportParams{
		SporttiID:  utils.NullString(sporttiID),
		SessionID:  utils.NullInt32(sessionID),
		RaceReport: utils.NullString(raceReport),
	})
	if err != nil {
		return err
	}

	if updated == 0 {
		if err := q.InsertReport(ctx, archsqlc.InsertReportParams{
			SporttiID:  utils.NullString(sporttiID),
			SessionID:  utils.NullInt32(sessionID),
			RaceReport: utils.NullString(raceReport),
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Combined data methods

func (s *DataStore) UpsertData(ctx context.Context, payload ArchDataPayload) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := archsqlc.New(tx)

	if err := q.UpsertAthlete(ctx, payload.Athlete); err != nil {
		return err
	}

	if err := upsertMeasurementsTx(ctx, q, payload.Athlete.NationalID, payload.Measurements); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *DataStore) GetDataBySporttiID(ctx context.Context, sporttiID string) (*ArchDataResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)

	// Athlete must exist
	a, err := q.GetAthleteBySporttiID(ctx, sporttiID)
	if err != nil {
		return nil, err
	}

	// Measurements may be empty
	ms, err := q.GetMeasurementsBySporttiID(ctx, utils.NullString(sporttiID))
	if err != nil {
		return nil, err
	}

	resp := ArchDataResponse{
		ArchAthleteResponse: toAthleteResponse(a),
		Measurements:        make([]ArchMeasurementResponse, 0, len(ms)),
	}
	for _, m := range ms {
		resp.Measurements = append(resp.Measurements, toMeasurementResponse(m))
	}
	return &resp, nil
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

	resp := toAthleteResponse(a)
	return &resp, nil
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

	if err := upsertMeasurementsTx(ctx, q, sporttiID, measurements); err != nil {
		return err
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

	q := archsqlc.New(s.db)
	m, err := q.GetMeasurementByMeasurementID(ctx, measurementID)
	if err != nil {
		return nil, err
	}

	resp := toMeasurementResponse(m)
	return &resp, nil
}

func (s *DataStore) DeleteMeasurementByMeasurementID(ctx context.Context, measurementID int32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := archsqlc.New(s.db)
	nationalID, err := q.DeleteMeasurementByMeasurementID(ctx, measurementID)
	if err != nil {
		return "", err
	}
	return nationalID.String, nil
}

// Helpers

// upsertMeasurementsTx makes sure every referenced measurement_group row exists,
// then upserts each measurement keyed on measurement_id. Runs inside the caller's tx.
func upsertMeasurementsTx(ctx context.Context, q *archsqlc.Queries, sporttiID string, measurements []archsqlc.UpsertMeasurementParams) error {
	seenGroups := make(map[int32]struct{}, len(measurements))

	for i := range measurements {
		m := &measurements[i]

		if !m.NationalID.Valid || m.NationalID.String == "" {
			m.NationalID = utils.NullString(sporttiID)
		}

		if m.MeasurementGroupID.Valid {
			gid := m.MeasurementGroupID.Int32
			if _, done := seenGroups[gid]; !done {
				if err := q.EnsureMeasurementGroup(ctx, gid); err != nil {
					return err
				}
				seenGroups[gid] = struct{}{}
			}
		}

		if err := q.UpsertMeasurement(ctx, *m); err != nil {
			return err
		}
	}
	return nil
}

func toAthleteResponse(a archsqlc.Athlete) ArchAthleteResponse {
	return ArchAthleteResponse{
		NationalID:  a.NationalID,
		FirstName:   utils.StringPtrOrNil(a.FirstName),
		LastName:    utils.StringPtrOrNil(a.LastName),
		Initials:    utils.StringPtrOrNil(a.Initials),
		DateOfBirth: utils.FormatDatePtr(a.DateOfBirth),
		Height:      utils.NullNumericToFloatPtr(a.Height),
		Weight:      utils.NullNumericToFloatPtr(a.Weight),
	}
}

func toMeasurementResponse(m archsqlc.Measurement) ArchMeasurementResponse {
	return ArchMeasurementResponse{
		MeasurementGroupID: utils.Int32PtrOrNil(m.MeasurementGroupID),
		MeasurementID:      m.MeasurementID,
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
