package utv

import (
	"context"
	"database/sql"
	"encoding/json"

	utvsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/utv"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type SourceCacheStore struct {
	db *sql.DB
}

func (s *SourceCacheStore) GetAll(ctx context.Context) ([]utvsqlc.SourceCache, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := utvsqlc.New(s.db)
	return q.GetSourceCacheAll(ctx)
}

func (s *SourceCacheStore) GetBySource(ctx context.Context, source string) (json.RawMessage, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := utvsqlc.New(s.db)
	return q.GetSourceCacheDataBySource(ctx, source)
}

func (s *SourceCacheStore) Upsert(ctx context.Context, source string, data json.RawMessage) error {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	q := utvsqlc.New(s.db)

	return q.UpsertSourceCache(ctx, utvsqlc.UpsertSourceCacheParams{
		Source:  source,
		Column2: data,
	})
}
