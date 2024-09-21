package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"url_shortener/core"
	"url_shortener/core/model"
	"url_shortener/storage/transaction"
)

type ClickRepo struct {
	db *transaction.Queries
}

func NewClickRepo(db *sql.DB) *ClickRepo {
	return &ClickRepo{db: transaction.NewQueries(db)}
}

func (r *ClickRepo) Save(ctx context.Context, click model.Click) (*model.Click, error) {
	const op = "storage.sqlite.ClickRepo.Save"

	stmt := `INSERT INTO click(link_id, access_time, ip, full_ad) VALUES (?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, stmt, click.LinkId, click.AccessTime, click.IP, click.FullAD)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, core.ErrLastInsertId)
	}

	click.Id = id

	return &click, nil
}
