package transaction

import (
	"context"
	"database/sql"
)

type Queries struct {
	db *sql.DB
}

func NewQueries(db *sql.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return q.db.QueryRowContext(ctx, query, args...)
}

func (q *Queries) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}

	return q.db.QueryContext(ctx, query, args...)
}

func (q *Queries) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.PrepareContext(ctx, query)
	}

	return q.db.PrepareContext(ctx, query)
}

func (q *Queries) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}

	return q.db.ExecContext(ctx, query, args...)
}
