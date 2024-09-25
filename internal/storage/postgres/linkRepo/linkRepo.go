package linkRepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/storage/transaction"
	"time"
)

type LinkRepo struct {
	db transaction.Queries
}

func New(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{db: transaction.NewQueries(db)}
}

func (r *LinkRepo) GetActiveByAlias(ctx context.Context, alias string) (*model.Link, error) {
	const op = "storage.postgres.LinkRepo.GetActiveByAlias"

	query := activeOnly(`SELECT * FROM link WHERE alias=$1`)

	row := r.db.QueryRow(ctx, query, alias)
	link, err := linkFromRow(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, core.ErrLinkNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}

func (r *LinkRepo) GetCount(ctx context.Context, userId int64, params model.LinkFilter) (int64, error) {
	const op = "storage.postgres.LinkRepo.GetCount"

	query := build(`SELECT COUNT(*) FROM link WHERE created_by = $1`).Filter(params).String()

	var totalCount int64
	err := r.db.QueryRow(ctx, query, userId).Scan(&totalCount)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCount, nil
}

func (r *LinkRepo) GetByUserId(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, error) {
	const op = "storage.postgres.LinkRepo.GetByUserId"

	query := build(`SELECT * FROM link WHERE created_by = $1`).
		Filter(params.Filter).
		Sort(params.Sort).
		String()

	links := make([]*model.Link, 0)

	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		link, err := linkFromRow(rows)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		links = append(links, link)
	}

	return links, nil
}

func (r *LinkRepo) AliasExists(ctx context.Context, alias string) (bool, error) {
	const op = "storage.postgres.LinkRepo.AliasExists"

	query := `SELECT EXISTS (SELECT 1 FROM link WHERE alias=$1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, alias).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (r *LinkRepo) CustomNameExists(ctx context.Context, customName string, userId int64) (bool, error) {
	const op = "storage.postgres.LinkRepo.CustomNameExists"

	query := `SELECT EXISTS (SELECT 1 FROM link WHERE custom_name=$1 AND created_by = $2)`

	var exists bool
	err := r.db.QueryRow(ctx, query, customName, userId).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

func (r *LinkRepo) Save(ctx context.Context, link *model.Link) (int64, error) {
	const op = "storage.postgres.LinkRepo.Save"

	query := `
		INSERT INTO link(created_by, original, alias, custom_name, expiration_date, clicks_to_expiration)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var id int64

	custom := link.Alias
	if custom == "" {
		custom = link.Alias
	}

	err := r.db.QueryRow(
		ctx, query, link.CreatedBy, link.Original, link.Alias,
		custom, link.ExpirationDate, link.ClicksToExpiration,
	).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *LinkRepo) UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error {
	op := "storage.postgres.LinkRepo.UpdateLastAccess"

	query := `UPDATE link SET last_access_time=$2, clicks_count=clicks_count+1 WHERE id=$1`

	_, err := r.db.Exec(ctx, query, id, timestamp)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
