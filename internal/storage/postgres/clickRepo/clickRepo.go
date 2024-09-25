package clickRepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortener/internal/core/model"
	"shortener/internal/storage/transaction"
)

type ClickRepo struct {
	db transaction.Queries
}

func New(db *pgxpool.Pool) *ClickRepo {
	return &ClickRepo{db: transaction.NewQueries(db)}
}

func (r *ClickRepo) Save(ctx context.Context, click *model.Click) (int64, error) {
	const op = "storage.postgres.ClickRepo.Save"

	query := `
		INSERT INTO click(link_id, user_agent, ip, access_time, ad_status) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64

	err := r.db.QueryRow(
		ctx, query, click.LinkId, click.Metadata.UserAgent,
		click.Metadata.IP, click.Metadata.AccessTime, click.Status,
	).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *ClickRepo) UpdateStatus(ctx context.Context, clickId int64, status model.AdStatus) error {
	const op = "storage.postgres.ClickRepo.UpdateStatus"

	query := `UPDATE click SET ad_status = $2 WHERE id = $1`

	_, err := r.db.Exec(ctx, query, clickId, status)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
