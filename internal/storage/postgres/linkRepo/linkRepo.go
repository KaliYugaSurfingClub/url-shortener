package linkRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/storage/entity"
	"shortener/internal/storage/transaction"
	"time"
)

type LinkRepo struct {
	db *transaction.Queries
}

func New(db *sqlx.DB) *LinkRepo {
	return &LinkRepo{db: transaction.NewQueries(db)}
}

func (r *LinkRepo) GetActiveByAlias(ctx context.Context, alias string) (*model.Link, error) {
	const op = "storage.sqlite.LinkRepo.GetActiveByAlias"

	query := activeOnly(`SELECT * FROM link WHERE alias=?`)

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

func (r *LinkRepo) GetByUserId(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, error) {
	const op = "storage.sqlite.LinkRepo.GetByUserId"

	query := withGetParams(`SELECT * FROM link WHERE created_by = ?`, params)

	entities := make([]entity.Link, 0)
	if err := r.db.SelectContext(ctx, &entities, query, userId); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	links := make([]*model.Link, 0, len(entities))
	for _, ent := range entities {
		links = append(links, ent.ToModel())
	}

	return links, nil
}

func (r *LinkRepo) GetCount(ctx context.Context, userId int64, params model.GetLinksParams) (int64, error) {
	const op = "storage.sqlite.LinkRepo.GetCount"

	query := withGetParams(`SELECT COUNT(*) FROM link WHERE created_by = ?`, params)
	var totalCount int64

	err := r.db.GetContext(ctx, &totalCount, query, userId)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCount, nil
}

func (r *LinkRepo) Save(ctx context.Context, link *model.Link) (int64, error) {
	const op = "storage.sqlite.LinkRepo.Save"

	query := `
		INSERT INTO link(created_by, original, alias, custom_name, expiration_date, clicks_to_expiration) 
		VALUES (:created_by, :original, :alias, :custom_name, :expiration_date, :clicks_to_expiration)
	`

	res, err := r.db.NamedExecContext(ctx, query, entity.ModelToLink(link))

	//if err != nil && errors.Is(err.(sqlite3.Error).ExtendedCode, sqlite3.ErrConstraintUnique) {
	//	return -1, fmt.Errorf("%s: %w", op, core.ErrAliasExists)
	//}
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func (r *LinkRepo) UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error {
	op := "storage.sqlite.LinkRepo.UpdateLastAccess"

	query := `UPDATE link SET last_access_time=?, clicks_count=clicks_count+1 WHERE id=?`

	_, err := r.db.ExecContext(ctx, query, timestamp, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
