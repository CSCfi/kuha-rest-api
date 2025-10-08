package kamk

import (
	"context"
	"database/sql"
)

type Injuries interface {
	AddInjury(ctx context.Context, userID string, in InjuryInput) error
	MarkInjuryRecovered(ctx context.Context, userID string, injuryID int32) (int64, error)
	GetActiveInjuries(ctx context.Context, userID string) ([]Injury, error)
	GetMaxInjuryID(ctx context.Context, userID string) (int32, error)
}

// KAMKStorage
type KAMKStorage struct {
	db       *sql.DB
	injuries Injuries
}

// Methods
func (s *KAMKStorage) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *KAMKStorage) Injuries() Injuries {
	return s.injuries
}

// NewKAMKStorage creates a new KAMKStorage instance
func NewKAMKStorage(db *sql.DB) *KAMKStorage {
	return &KAMKStorage{
		db:       db,
		injuries: &InjuriesStore{db: db},
	}
}
