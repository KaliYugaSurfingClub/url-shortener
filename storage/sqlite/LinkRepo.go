package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"time"
	"url_shortener/core"
	"url_shortener/core/model"
	"url_shortener/storage/transaction"
)

type LinkRepo struct {
	db *transaction.Queries
}

func NewLinkRepo(db *sql.DB) *LinkRepo {
	return &LinkRepo{db: transaction.NewQueries(db)}
}

func (r *LinkRepo) GetByAlias(ctx context.Context, alias string) (*model.Link, error) {
	const op = "storage.sqlite.LinkRepo.GetByAlias"

	stmt := `SELECT id, user_id, original, clicks_count, last_access, expire_date, max_clicks FROM link WHERE alias=?`

	link := &model.Link{Alias: alias, ExpireDate: model.NoExpireDate, MaxClicks: model.UnlimitedClicks}
	var expireDate sql.NullTime
	var maxClicks sql.NullInt64

	err := r.db.QueryRowContext(ctx, stmt, alias).Scan(
		&link.Id,
		&link.CreatedBy,
		&link.Original,
		&link.ClicksCount,
		&link.LastAccess,
		&expireDate,
		&maxClicks,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, core.ErrLinkNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if expireDate.Valid {
		link.ExpireDate = expireDate.Time
	}

	if maxClicks.Valid {
		link.MaxClicks = maxClicks.Int64
	}

	return link, nil
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

func (r *LinkRepo) Save(ctx context.Context, link model.Link) (int64, error) {
	const op = "storage.sqlite.LinkRepo.Save"

	stmt := `INSERT INTO link(user_id, original, alias, expire_date, max_clicks) VALUES (?, ?, ?, ?, ?)`

	expireDate := sql.NullTime{Valid: false}
	if link.ExpireDate != model.NoExpireDate {
		expireDate = sql.NullTime{Valid: true, Time: link.ExpireDate}
	}

	maxClicks := sql.NullInt64{Valid: false}
	if link.MaxClicks != model.UnlimitedClicks {
		maxClicks = sql.NullInt64{Valid: true, Int64: link.MaxClicks}
	}

	res, err := r.db.ExecContext(ctx, stmt,
		link.CreatedBy,
		link.Original,
		link.Alias,
		expireDate,
		maxClicks,
	)

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
