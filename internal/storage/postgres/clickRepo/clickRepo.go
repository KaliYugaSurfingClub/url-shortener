package clickRepo

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shortener/internal/core/model"
	"shortener/internal/storage/entity"
	"shortener/internal/storage/transaction"
)

type ClickRepo struct {
	db *transaction.Queries
}

func New(db *sqlx.DB) *ClickRepo {
	return &ClickRepo{db: transaction.NewQueries(db)}
}

func (r *ClickRepo) Save(ctx context.Context, click *model.Click) (int64, error) {
	const op = "storage.sqlite.ClickRepo.Save"

	query := `
		INSERT INTO click(link_id, access_time, ip, ad_status) 
		VALUES (:link_id, :access_time, :ip, :ad_status)
	`

	res, err := r.db.NamedExecContext(ctx, query, entity.ClickFromModel(click))
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, _ := res.LastInsertId()

	return id, nil
}

func (r *ClickRepo) UpdateStatus(ctx context.Context, clickId int64, status model.AdStatus) error {
	const op = "storage.sqlite.ClickRepo.UpdateStatus"

	query := `UPDATE click SET ad_status = ? WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, status, clickId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
