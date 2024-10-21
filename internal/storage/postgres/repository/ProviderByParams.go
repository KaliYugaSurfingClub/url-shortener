package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"shortener/errs"
	"shortener/internal/core/model"
	"shortener/internal/storage/postgres/builder"
	"shortener/internal/storage/postgres/transaction"
)

type getEntityByParamsOptions[T any] struct {
	db         transaction.Queries
	query      string
	pagination model.Pagination
	sort       model.Sort
	columns    map[model.SortBy]string
	scanFunc   func(*T, pgx.Row) error
	args       []any
}

func getEntityByParams[T any](ctx context.Context, opt getEntityByParamsOptions[T]) ([]*T, error) {
	const op = "storage.postgres.repository.getEntityByParams"

	entities := make([]*T, 0, opt.pagination.Size)

	query := builder.New(opt.query).
		Sort(opt.columns, opt.sort).
		Paginate(opt.pagination).
		String()

	rows, err := opt.db.Query(ctx, query, opt.args...)
	if err != nil {
		return nil, errs.E(op, err, errs.Database)
	}

	defer rows.Close()

	for rows.Next() {
		entity := new(T)

		if err := opt.scanFunc(entity, rows); err != nil {
			return nil, errs.E(op, err, errs.Database)
		}

		entities = append(entities, entity)
	}

	if err = rows.Err(); err != nil {
		return nil, errs.E(op, err, errs.Database)
	}

	return entities, nil
}
