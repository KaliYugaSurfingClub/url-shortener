package sqlite

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"url_shortener/core/model"
	"url_shortener/storage/transaction"
)

type ClickRepo struct {
	db *transaction.Queries
}

func NewClickRepo(db *sqlx.DB) *ClickRepo {
	return &ClickRepo{db: transaction.NewQueries(db)}
}

func (r *ClickRepo) Save(ctx context.Context, click model.Click) (int64, error) {
	const op = "storage.sqlite.ClickRepo.Save"

	stmt := `INSERT INTO click(link_id, access_time, ip, full_ad) VALUES (?, ?, ?, ?)`

	var id int64

	err := r.db.GetContext(ctx, &id, stmt, click.LinkId, click.AccessTime, click.IP, click.FullAD)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
