package clickRepo

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

func New(db *sqlx.DB) *ClickRepo {
	return &ClickRepo{db: transaction.NewQueries(db)}
}

func (r *ClickRepo) Save(ctx context.Context, click model.Click) (int64, error) {
	const op = "storage.sqlite.ClickRepo.Save"

	query := `INSERT INTO click(link_id, access_time, ip, full_ad) VALUES (?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, query, click.LinkId, click.AccessTime, click.IP, click.FullAD)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}
