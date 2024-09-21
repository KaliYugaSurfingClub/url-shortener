package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"time"
	"url_shortener/core"
	"url_shortener/storage/transaction"
)

type LinkRepo struct {
	db *transaction.Queries
}

func NewLinkRepo(db *sql.DB) *LinkRepo {
	return &LinkRepo{db: transaction.NewQueries(db)}
}

func (r *LinkRepo) GetOriginalByAlias(ctx context.Context, alias string) (int64, string, error) {
	const op = "storage.sqlite.LinkRepo.GetOriginalByAlias"

	var original string
	var id int64

	err := r.db.QueryRowContext(ctx, `SELECT id, original FROM link WHERE alias=?`, alias).Scan(&id, &original)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, "", fmt.Errorf("%s: %w", op, core.ErrOriginalNotFound)
	}
	if err != nil {
		return 0, "", fmt.Errorf("%s: %w", op, err)
	}

	return id, original, nil
}

func (r *LinkRepo) UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error {
	op := "storage.sqlite.LinkRepo.UpdateLastAccess"

	_, err := r.db.ExecContext(ctx, `UPDATE link SET last_access=? WHERE id=?`, timestamp, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *LinkRepo) Save(ctx context.Context, original string, alias string) (int64, error) {
	const op = "storage.sqlite.LinkRepo.Save"

	res, err := r.db.ExecContext(ctx, `INSERT INTO link (original, alias) VALUES (?, ?)`, original, alias)
	if err != nil && errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintUnique) {
		return -1, fmt.Errorf("%s: %w", op, core.ErrAliasExists)
	}
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, core.ErrLastInsertId)
	}

	return id, nil
}
