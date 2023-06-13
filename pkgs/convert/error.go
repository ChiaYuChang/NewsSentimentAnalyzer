package convert

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func ToPgError(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr, true
	} else {
		return nil, false
	}
}
