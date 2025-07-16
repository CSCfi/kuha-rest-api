package utv

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	utvsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/utv"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
	"github.com/google/uuid"
)

type CoachtechDataStore struct {
	db *sql.DB
}

func (s *CoachtechDataStore) GetStatus(ctx context.Context, userID uuid.UUID) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	queries := utvsqlc.New(s.db)
	return queries.GetCoachtechStatus(ctx, userID)
}

func (s *CoachtechDataStore) GetData(ctx context.Context, userID string, afterStr, beforeStr *string) ([]json.RawMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	uid, err := utils.ParseUUID(userID)
	if err != nil {
		return nil, err
	}

	after, err := utils.ParseDatePtr(afterStr)
	if err != nil {
		return nil, err
	}

	before, err := utils.ParseDatePtr(beforeStr)
	if err != nil {
		return nil, err
	}

	arg := utvsqlc.GetCoachtechDataParams{
		UserID:     uid,
		AfterDate:  utils.NullTimeIfEmpty(after),
		BeforeDate: utils.NullTimeIfEmpty(before),
	}

	queries := utvsqlc.New(s.db)

	data, err := queries.GetCoachtechData(ctx, arg)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, utils.ErrQueryTimeOut
		}
		return nil, err
	}

	return data, nil
}
