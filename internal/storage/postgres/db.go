package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(postgresURL string) (*pgxpool.Pool, func(), error) {
	const op = "storage.postgres.NewPgxPool"

	poolCfg, err := pgxpool.ParseConfig(postgresURL)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	cancel := func() {
		db.Close()
	}

	return db, cancel, nil
}
