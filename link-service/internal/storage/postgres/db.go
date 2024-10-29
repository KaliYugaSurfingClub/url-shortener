package postgres

import (
	"context"
	"fmt"
	"github.com/KaliYugaSurfingClub/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool(postgresURL string) (*pgxpool.Pool, func(), error) {
	const op errs.Op = "storage.postgres.NewPgxPool"

	cancel := func() {}

	poolCfg, err := pgxpool.ParseConfig(postgresURL)
	if err != nil {
		return nil, cancel, fmt.Errorf("%s: %w", op, err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, cancel, fmt.Errorf("%s: %w", op, err)
	}

	cancel = func() {
		db.Close()
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, cancel, fmt.Errorf("%s: %w", op, err)
	}

	return db, cancel, nil
}
