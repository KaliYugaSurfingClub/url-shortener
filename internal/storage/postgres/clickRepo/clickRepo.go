package clickRepo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/storage/transaction"
)

var _ port.ClickStorage = (*ClickRepo)(nil)

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

func (r *ClickRepo) GetCountByLinkId(ctx context.Context, linkId int64, params model.GetClicksParams) (int64, error) {
	const op = "storage.postgres.ClickRepo.GetCountByLinkId"

	query := `SELECT COUNT(*) FROM click WHERE link_id = $1`

	var totalCount int64
	err := r.db.QueryRow(ctx, query, linkId).Scan(&totalCount)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCount, nil
}

func (r *ClickRepo) GetByLinkId(ctx context.Context, linkId int64, params model.GetClicksParams) ([]*model.Click, error) {
	const op = "storage.postgres.ClickRepo.GetByLinkId"

	query := build(`SELECT * FROM click WHERE link_id = $1`).
		Paginate(params.Pagination).
		String()

	rows, err := r.db.Query(ctx, query, linkId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	clicks := make([]*model.Click, 0)

	for rows.Next() {
		click, err := clickFromRow(rows)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		clicks = append(clicks, click)
	}

	return clicks, nil
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
