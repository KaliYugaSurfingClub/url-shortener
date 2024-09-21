package storage

import (
	"context"
	"database/sql"
	"fmt"
	"url_shortener/storage/transaction"
)

type ClickRepo struct {
	db *transaction.Queries
}

func NewClickRepo(db *sql.DB) *ClickRepo {
	return &ClickRepo{db: transaction.NewQueries(db)}
}

func (r *ClickRepo) Save(ctx context.Context, linkID int64) error {
	const op = "storage.sqlite.ClickRepo.Save"

	if _, err := r.db.ExecContext(ctx, `INSERT INTO clicks(link_id) VALUES (?)`, linkID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
