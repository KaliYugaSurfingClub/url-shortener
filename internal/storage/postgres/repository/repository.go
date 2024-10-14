package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/storage/postgres"
	"shortener/internal/storage/postgres/transaction"
	"shortener/internal/utils"
)

type Repository struct {
	db transaction.Queries
	transaction.Transactor
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: transaction.NewQueries(db)}
}

const createLinkQuery = `
	INSERT INTO link(person_id, original, alias, custom_name)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
`

func (r *Repository) CreateLink(ctx context.Context, link model.Link) (*model.Link, error) {
	const op = "storage.postgres.repository.CreateLink"

	row := r.db.QueryRow(
		ctx, createLinkQuery, link.CreatedBy,
		link.Original, link.Alias, link.CustomName,
	)

	err := row.Scan(&link.Id, &link.CreatedAt)

	if name, ok := postgres.ParseConstraintError(err); ok {
		switch name {
		case "link_alias_key":
			return nil, fmt.Errorf("%s: %w", op, core.ErrAliasExists)
		case "link_custom_name_person_id_key":
			return nil, fmt.Errorf("%s: %w", op, core.ErrCustomNameExists)
		default:
			return nil, fmt.Errorf("%s: unexpected constraint error %w", op, err)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &link, nil
}

const GetLinkByIdQuery = `
	SELECT id, person_id, original, alias, custom_name, archived, created_at
	FROM link 
	WHERE alias = $1 
	LIMIT 1
`

func (r *Repository) GetLinkByAlias(ctx context.Context, alias string) (*model.Link, error) {
	const op = "storage.postgres.repository.GetLinkByAlias"

	link := new(model.Link)

	err := r.db.QueryRow(ctx, GetLinkByIdQuery, alias).Scan(
		&link.Id, &link.CreatedBy, &link.Original, &link.Alias,
		&link.CustomName, &link.Archived, &link.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, core.ErrLinkNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return link, nil
}

const GetLinksQuery = `
	SELECT  
	    l.id, l.person_id, l.original, l.alias, l.custom_name, l.archived, l.created_at,
		COUNT(c.id)        AS clicks_count,
		MAX(c.access_time) AS last_access_time
	FROM (
		SELECT id, person_id, original, alias, custom_name, archived, created_at
		FROM link
		WHERE person_id = $1 AND archived = $2
	)
    AS l LEFT JOIN click AS c ON l.id = c.link_id
	GROUP BY l.id, l.person_id, l.original, l.alias, l.custom_name, l.archived, l.created_at
`

func (r *Repository) GetLinksByParams(ctx context.Context, params model.GetLinksParams) (_ []*model.Link, err error) {
	defer utils.WithinOp("storage.postgres.repository.GetLinksByParams", &err)

	scanFunc := func(link *model.Link, row pgx.Row) error {
		return row.Scan(
			&link.Id, &link.CreatedBy, &link.Original, &link.Alias, &link.CustomName,
			&link.Archived, &link.CreatedAt, &link.ClicksCount, &link.LastAccessTime,
		)
	}

	opt := getEntityByParamsOptions[model.Link]{
		db:         r.db,
		query:      GetLinksQuery,
		args:       []any{params.UserId, params.Archived},
		pagination: params.Pagination,
		sort:       params.Sort,
		columns:    SortLinksBy,
		scanFunc:   scanFunc,
	}

	return getEntityByParams(ctx, opt)
}

const GetLinksCountQuery = `
	SELECT COUNT(*) FROM link WHERE person_id = $1 AND archived = $2 
`

func (r *Repository) GetLinksCountByParams(ctx context.Context, params model.GetLinksParams) (count int64, err error) {
	const op = "storage.postgres.repository.GetLinksCountByParams"

	err = r.db.QueryRow(ctx, GetLinksCountQuery, params.UserId, params.Archived).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}

const doesLinkBelongsUserQuery = `
	SELECT EXISTS (
		SELECT 1 FROM link WHERE id = $1 AND person_id = $2
	);
`

func (r *Repository) DoesLinkBelongsToUser(ctx context.Context, linkId, userId int64) (belongs bool, err error) {
	const op = "storage.postgres.repository.DoesLinkBelongsToUser"

	err = r.db.QueryRow(ctx, doesLinkBelongsUserQuery, linkId, userId).Scan(&belongs)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("%s: %w", op, core.ErrLinkNotFound)
	}
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return belongs, nil
}

const deleteLinkQuery = `
	DELETE FROM link WHERE id = $1
`

func (r *Repository) DeleteLink(ctx context.Context, linkId int64) error {
	const op = "storage.postgres.repository.DeleteLink"

	_, err := r.db.Exec(ctx, deleteLinkQuery, linkId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

const createClickQuery = `
	WITH inserted_click AS (
		INSERT INTO click (link_id, ad_source_id, user_agent, ip, access_time)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, ad_source_id
	)
	SELECT inserted_click.id, ad_source.type
	FROM inserted_click
	JOIN ad_source ON ad_source.id = inserted_click.ad_source_id
`

func (r *Repository) CreateClick(ctx context.Context, click model.Click) (*model.Click, error) {
	const op = "storage.postgres.repository.CreateClick"

	row := r.db.QueryRow(
		ctx, createClickQuery, click.LinkId, click.AdSourceId,
		click.Metadata.UserAgent, click.Metadata.IP, click.Metadata.AccessTime,
	)

	err := row.Scan(&click.Id, &click.AdType)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &click, nil
}

const getClicksByParamsQuery = `
	SELECT
		c.id, c.link_id, c.ad_source_id, c.user_agent, c.ip, c.access_time, s.type AS ad_type,
		CASE
			WHEN r.click_id IS NOT NULL THEN 'completed'::click_status
			WHEN r.click_id IS NULL THEN 'opened'::click_status
		END AS status
	FROM
		(
			SELECT id, link_id, ad_source_id, user_agent, ip, access_time
			FROM click
			WHERE link_id = $1
		) AS c
		LEFT JOIN click_reward AS r ON c.id = r.click_id
		JOIN ad_source AS s ON c.ad_source_id = s.id
`

func (r *Repository) GetClicksByParams(ctx context.Context, params model.GetClicksParams) (_ []*model.Click, err error) {
	defer utils.WithinOp("storage.postgres.repository.GetClicksByParams", &err)

	scanFunc := func(click *model.Click, row pgx.Row) error {
		return row.Scan(
			&click.Id, &click.LinkId, &click.AdSourceId, &click.Metadata.UserAgent,
			&click.Metadata.IP, &click.Metadata.AccessTime, &click.AdType, &click.Status,
		)
	}

	opt := getEntityByParamsOptions[model.Click]{
		db:         r.db,
		query:      getClicksByParamsQuery,
		args:       []any{params.LinkId},
		pagination: params.Pagination,
		sort:       params.Sort,
		columns:    SortClicksBy,
		scanFunc:   scanFunc,
	}

	return getEntityByParams(ctx, opt)
}

const getClicksCountQuery = `
	SELECT COUNT(*) FROM click WHERE link_id = $1
`

func (r *Repository) GetClicksCountByParams(ctx context.Context, params model.GetClicksParams) (count int64, err error) {
	const op = "storage.postgres.repository.GetClicksCount"

	err = r.db.QueryRow(ctx, getClicksCountQuery, params.LinkId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return count, nil
}
