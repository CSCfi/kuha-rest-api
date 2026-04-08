package archinisis

import (
	"context"
	"database/sql"
)

// Interfaces
type Users interface {
	DeleteUserBySporttiID(ctx context.Context, sporttiID string) (string, error)
}

type Data interface {
	GetRaceReportSessions(ctx context.Context, sporttiID string) ([]int32, error)
	GetRaceReport(ctx context.Context, sporttiID string, sessionID int32) (string, error)
	UpsertRaceReport(ctx context.Context, sporttiID string, sessionID int32, raceReport string) error
	UpsertData(ctx context.Context, payload ArchDataPayload) error
	GetDataBySporttiID(ctx context.Context, sporttiID string) (*ArchDataResponse, error)
	GetSporttiIDsBySessionID(ctx context.Context, sessionID int32) ([]string, error)
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
