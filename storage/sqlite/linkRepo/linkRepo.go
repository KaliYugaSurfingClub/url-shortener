package linkRepo

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

func New(db *sqlx.DB) *LinkRepo {
	return &LinkRepo{db: transaction.NewQueries(db)}
}

func (r *LinkRepo) GetActualByAlias(ctx context.Context, alias string) (*model.Link, error) {
	const op = "storage.sqlite.LinkRepo.GetActualByAlias"

	query := actualOnly(`SELECT * FROM link WHERE alias=?`)
	link := &entity.Link{}

	err := r.db.QueryRowContext(ctx, query, alias).StructScan(link)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, core.ErrLinkNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return link.ToModel(), nil
}

func (r *LinkRepo) GetTotalCountByUserId(ctx context.Context, userId string) (int64, error) {
	const op = "storage.sqlite.LinkRepo.GetTotalCountByUserId"

	query := `SELECT COUNT(*) FROM link WHERE created_by = ?`
	var totalCount int64

	err := r.db.GetContext(ctx, &totalCount, query, userId)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCount, nil
}

func (r *LinkRepo) GetByUserId(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, error) {
	const op = "storage.sqlite.LinkRepo.GetByUserId"

	query := withGetParams(`SELECT * FROM link WHERE created_by = ?`, params)
	entities := make([]entity.Link, 0)

	err := r.db.SelectContext(ctx, &entities, query, userId)

	if errors.Is(err, sql.ErrNoRows) {
		return make([]*model.Link, 0), nil
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	links := make([]*model.Link, 0, len(entities))
	for _, ent := range entities {
		links = append(links, ent.ToModel())
	}

	return links, nil
}

func (r *LinkRepo) Save(ctx context.Context, link model.Link) (int64, error) {
	const op = "storage.sqlite.LinkRepo.Save"

	query := `
		INSERT INTO link(created_by, original, alias, expiration_date, max_clicks) 
		VALUES (?, ?, ?, ?, ?)
	`

	res, err := r.db.ExecContext(ctx, query,
		link.CreatedBy,
		link.Original,
		link.Alias,
		entity.SqlExpirationDate(link.ExpirationDate),
		entity.SqlMaxClicks(link.MaxClicks),
	)

	if err != nil && errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintUnique) {
		return -1, fmt.Errorf("%s: %w", op, core.ErrAliasExists)
	}
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func (r *LinkRepo) UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error {
	op := "storage.sqlite.LinkRepo.UpdateLastAccess"

	query := actualOnly(`UPDATE link SET last_access=?, clicks_count=clicks_count+1 WHERE id=?`)

	_, err := r.db.ExecContext(ctx, query, timestamp, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
