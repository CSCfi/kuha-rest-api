package archinisis

import (
	"context"
	"database/sql"

	archsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/archinisis"
)

// Interfaces
type Users interface {
	DeleteUserBySporttiID(ctx context.Context, sporttiID string) (string, error)
}

type Data interface {
	// Race reports
	GetRaceReportSessions(ctx context.Context, sporttiID string) ([]int32, error)
	GetRaceReport(ctx context.Context, sporttiID string, sessionID int32) (string, error)
	UpsertRaceReport(ctx context.Context, sporttiID string, sessionID int32, raceReport string) error
	GetSporttiIDsBySessionID(ctx context.Context, sessionID int32) ([]string, error)
	// Athlete
	UpsertAthleteOnly(ctx context.Context, athlete archsqlc.UpsertAthleteParams) error
	GetAthleteByID(ctx context.Context, sporttiID string) (*ArchAthleteResponse, error)
	// Measurements
	UpsertMeasurements(ctx context.Context, sporttiID string, measurements []archsqlc.UpsertMeasurementParams) error
	GetMeasurementsByID(ctx context.Context, sporttiID string) ([]ArchMeasurementResponse, error)
	GetMeasurementByMeasurementID(ctx context.Context, measurementID int32) (*ArchMeasurementResponse, error)
	DeleteMeasurementByMeasurementID(ctx context.Context, measurementID int32) (string, error)
}

// ArchinisisStorage
type ArchinisisStorage struct {
	db    *sql.DB
	users Users
	data  Data
}

// Methods
func (s *ArchinisisStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *ArchinisisStorage) Users() Users {
	return s.users
}

func (s *ArchinisisStorage) Data() Data {
	return s.data
}

// NewArchinisisStorage creates a new ArchinisisStorage instance
func NewArchinisisStorage(db *sql.DB) *ArchinisisStorage {
	return &ArchinisisStorage{
		db:    db,
		users: &UsersStore{db: db},
		data:  &DataStore{db: db},
	}
}
