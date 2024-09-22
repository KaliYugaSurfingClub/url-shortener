package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"time"
	"url_shortener/core"
	"url_shortener/core/model"
	"url_shortener/storage/entity"
	"url_shortener/storage/transaction"
)

type LinkRepo struct {
	db *transaction.Queries
}

func NewLinkRepo(db *sqlx.DB) *LinkRepo {
	return &LinkRepo{db: transaction.NewQueries(db)}
}

func (r *LinkRepo) GetByAlias(ctx context.Context, alias string) (*model.Link, error) {
	const op = "storage.sqlite.LinkRepo.GetByAlias"

	stmt := `
		SELECT id, created_by, original, clicks_count, last_access_time, expiration_date, max_clicks 
		FROM link WHERE alias=?
	`

	link := &entity.Link{Alias: alias}

	err := r.db.QueryRowContext(ctx, stmt, alias).StructScan(link)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, core.ErrLinkNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return link.ToModel(), nil
}

func (r *LinkRepo) Save(ctx context.Context, link model.Link) (int64, error) {
	const op = "storage.sqlite.LinkRepo.Save"

	stmt := `
		INSERT INTO link(created_by, original, alias, expiration_date, max_clicks) 
		VALUES (?, ?, ?, ?, ?)
	`

	var id int64

	err := r.db.GetContext(ctx, &id, stmt,
		link.CreatedBy,
		link.Original,
		link.Alias,
		entity.SqlExpirationDate(link.ExpirationDate),
		entity.SqlMaxClicks(link.MaxClicks),
	)

	var sqliteErr *sqlite3.Error
	if errors.As(err, &sqliteErr) && errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintUnique) {
		return -1, fmt.Errorf("%s: %w", op, core.ErrAliasExists)
	}
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *LinkRepo) UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error {
	op := "storage.sqlite.LinkRepo.UpdateLastAccess"

	stmt := `UPDATE link SET last_access=?, clicks_count=clicks_count+1 WHERE id=?`

	_, err := r.db.ExecContext(ctx, stmt, timestamp, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
