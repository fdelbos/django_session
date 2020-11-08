package django_session

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	PGXSession struct {
		Pool *pgxpool.Pool
	}
)

const query = `
	select
		session_data
	from
		django_session
	where
		session_key = $1
		and expire_date > now()
`

func (pgxs *PGXSession) Fetch(ctx context.Context, key string, dest interface{}) error {
	encodedValue := ""
	err := pgxs.Pool.QueryRow(ctx, query, key).Scan(&encodedValue)
	if err == sql.ErrNoRows {
		return ErrSessionInvalid
	} else if err != nil {
		return err
	}

	res, err := decodeString(encodedValue)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, dest)
}
