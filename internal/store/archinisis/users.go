package archinisis

import (
	"context"
	"database/sql"

	archsqlc "github.com/DeRuina/KUHA-REST-API/internal/db/archinisis"
	"github.com/DeRuina/KUHA-REST-API/internal/utils"
)

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) DeleteUserBySporttiID(ctx context.Context, sporttiID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.QueryTimeout)
	defer cancel()

	id, err := utils.ParseSporttiID(sporttiID)
	if err != nil {
		return "", err
	}

	q := archsqlc.New(s.db)
	deletedID, err := q.DeleteAthleteByNationalID(ctx, id)
	if err != nil {
		return "", err
	}
	return deletedID, nil
}
