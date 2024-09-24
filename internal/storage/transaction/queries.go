package transaction

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	db *sqlx.DB
}

func NewQueries(db *sqlx.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) QueryRowContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRowxContext(ctx, query, args...)
	}

	return q.db.QueryRowxContext(ctx, query, args...)
}

func (q *Queries) QueryContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryxContext(ctx, query, args...)
	}

	return q.db.QueryxContext(ctx, query, args...)
}

func (q *Queries) PrepareContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.PreparexContext(ctx, query)
	}

	return q.db.PreparexContext(ctx, query)
}

func (q *Queries) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}

	return q.db.ExecContext(ctx, query, args...)
}

func (q *Queries) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := extractTx(ctx); tx != nil {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return q.db.GetContext(ctx, dest, query, args...)
}

func (q *Queries) NamedExecContext(ctx context.Context, query string, dest any) (sql.Result, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.NamedExecContext(ctx, query, dest)
	}

	return q.db.NamedExecContext(ctx, query, dest)
}

func (q *Queries) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	if tx := extractTx(ctx); tx != nil {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return q.db.SelectContext(ctx, dest, query, args...)
}
