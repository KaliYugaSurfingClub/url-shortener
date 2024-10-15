package transaction

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	db *pgxpool.Pool
}

func NewQueries(db *pgxpool.Pool) Queries {
	return Queries{db: db}
}

func (q *Queries) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRow(ctx, query, args...)
	}

	return q.db.QueryRow(ctx, query, args...)
}

func (q *Queries) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.Query(ctx, query, args...)
	}

	return q.db.Query(ctx, query, args...)
}

func (q *Queries) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.Exec(ctx, query, args...)
	}

	return q.db.Exec(ctx, query, args...)
}

func (q *Queries) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	if tx := extractTx(ctx); tx != nil {
		return tx.SendBatch(ctx, b)
	}

	return q.db.SendBatch(ctx, b)
}
