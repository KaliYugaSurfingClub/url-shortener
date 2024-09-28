package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

func ParseConstraintError(err error) (string, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr.Message, true
	}

	return "", false
}
