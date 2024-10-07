package clickRepo

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/storage/postgres/builder"
	"shortener/internal/storage/transaction"
	"time"
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

func (r *ClickRepo) GetCountByLinkId(ctx context.Context, params model.GetClicksParams) (int64, error) {
	const op = "storage.postgres.ClickRepo.GetCountByLinkId"

	query := `SELECT COUNT(*) FROM click WHERE link_id = $1`

	var totalCount int64
	err := r.db.QueryRow(ctx, query, params.LinkId).Scan(&totalCount)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return totalCount, nil
}

// todo what if pagination = {0, 0}
func (r *ClickRepo) GetByLinkId(ctx context.Context, params model.GetClicksParams) ([]*model.Click, error) {
	const op = "storage.postgres.ClickRepo.GetByLinkId"

	query := builder.New(`SELECT * FROM click WHERE link_id = $1`).Paginate(params.Pagination).String()

	rows, err := r.db.Query(ctx, query, params.LinkId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	clicks := make([]*model.Click, 0, params.Pagination.Size)

	for rows.Next() {
		click, err := clickFromRow(rows)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		clicks = append(clicks, click)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
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

func (r *ClickRepo) BatchUpdateStatus(ctx context.Context, clicksIds []int64, status model.AdStatus) error {
	const op = "storage.postgres.ClickRepo.BatchUpdateStatus"

	batch := new(pgx.Batch)

	query := `UPDATE click SET ad_status = $2 WHERE id = $1`

	for _, id := range clicksIds {
		batch.Queue(query, id, status)
	}

	errs := new(multierror.Error)
	br := r.db.SendBatch(ctx, batch)

	for i := 0; i < len(clicksIds); i++ {
		_, err := br.Exec()
		errs = multierror.Append(errs, err)
	}

	if err := errs.ErrorOrNil(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// todo duration
func (r *ClickRepo) GetExpiredClickSessions(ctx context.Context, sessionLifetime time.Duration, count int64) ([]*model.Click, error) {
	const op = "storage.postgres.ClickRepo.GetExpiredClickSessions"

	query := `SELECT * FROM click WHERE CURRENT_TIMESTAMP - access_time >= $1 LIMIT $2`

	rows, err := r.db.Query(ctx, query, sessionLifetime, count)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	clicks := make([]*model.Click, 0, count)

	for rows.Next() {
		click, err := clickFromRow(rows)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		clicks = append(clicks, click)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return clicks, nil
}
