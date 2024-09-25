package transaction

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}

	return nil
}

type Transactor struct {
	db *pgxpool.Pool
}

func NewTransactor(db *pgxpool.Pool) Transactor {
	return Transactor{db: db}
}

func (t Transactor) WithinTx(ctx context.Context, tFunc func(ctx context.Context) error) (err error) {
	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			tx.Rollback(ctx)
		}
	}()

	//todo Rollback with Timeout????
	if err = tFunc(injectTx(ctx, tx)); err != nil {
		return fmt.Errorf("execute transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
