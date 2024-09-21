package transaction

import (
	"context"
	"database/sql"
	"fmt"
)

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}

	return nil
}

type Transactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) *Transactor {
	return &Transactor{db: db}
}

func (t *Transactor) WithinTx(ctx context.Context, tFunc func(ctx context.Context) error) (err error) {
	tx, err := t.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("panic occurred: %v", r)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	if err = tFunc(injectTx(ctx, tx)); err != nil {
		return fmt.Errorf("execute transaction: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
