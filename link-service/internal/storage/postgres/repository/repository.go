package repository

import (
	"context"
	"errors"
	"github.com/KaliYugaSurfingClub/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/storage/postgres"
	"shortener/internal/storage/postgres/transaction"
)

type Repository struct {
	queries transaction.Queries
	transaction.Transactor
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		queries:    transaction.NewQueries(db),
		Transactor: transaction.NewTransactor(db),
	}
}

const createLinkQuery = `
	INSERT INTO link(person_id, original, alias, custom_name)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
`

func (r *Repository) CreateLink(ctx context.Context, link model.Link) (*model.Link, error) {
	const op errs.Op = "storage.postgres.repository.CreateLink"

	row := r.queries.QueryRow(
		ctx, createLinkQuery, link.CreatedBy,
		link.Original, link.Alias, link.CustomName,
	)

	err := row.Scan(&link.Id, &link.CreatedAt)

	if name, ok := postgres.ParseConstraintError(err); ok {
		switch name {
		case "link_alias_key":
			return nil, errs.E(op, err, errs.Exist, core.AliasExistsCode)
		case "link_custom_name_person_id_key":
			return nil, errs.E(op, err, errs.Exist, core.CustomNameExistsCode)
		default:
			return nil, errs.E(op, err, errs.Unanticipated)
		}
	}

	if err != nil {
		return nil, errs.E(op, err, errs.Database)
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
	const op errs.Op = "storage.postgres.repository.GetLinkByAlias"

	link := new(model.Link)

	err := r.queries.QueryRow(ctx, GetLinkByIdQuery, alias).Scan(
		&link.Id, &link.CreatedBy, &link.Original, &link.Alias,
		&link.CustomName, &link.Archived, &link.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.E(op, err, errs.NotExist)
	}
	if err != nil {
		return nil, errs.E(op, err, errs.Database)
	}

	return link, nil
}

const GetLinkByClickIdQuery = `
	SELECT l.id, l.person_id, l.original, l.alias, l.custom_name, l.archived, l.created_at
	FROM 
	( SELECT link_id FROM click WHERE id = $1 ) AS c
	JOIN link AS l ON l.id = link_id
`

func (r *Repository) GetOriginalByClickId(ctx context.Context, clickId int64) (*model.Link, error) {
	const op errs.Op = "storage.postgres.repository.GetOriginalByClickId"

	link := new(model.Link)

	err := r.queries.QueryRow(ctx, GetLinkByClickIdQuery, clickId).Scan(
		&link.Id, &link.CreatedBy, &link.Original, &link.Alias,
		&link.CustomName, &link.Archived, &link.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errs.E(op, err, errs.NotExist)
	}
	if err != nil {
		return nil, errs.E(op, err, errs.Database)
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

func (r *Repository) GetLinksByParams(ctx context.Context, params model.GetLinksParams) ([]*model.Link, error) {
	const op errs.Op = "storage.postgres.repository.GetLinksByParams"

	scanFunc := func(link *model.Link, row pgx.Row) error {
		return row.Scan(
			&link.Id, &link.CreatedBy, &link.Original, &link.Alias, &link.CustomName,
			&link.Archived, &link.CreatedAt, &link.ClicksCount, &link.LastAccessTime,
		)
	}

	opt := getEntityByParamsOptions[model.Link]{
		db:         r.queries,
		query:      GetLinksQuery,
		args:       []any{params.UserId, params.Archived},
		pagination: params.Pagination,
		sort:       params.Sort,
		columns:    SortLinksBy,
		scanFunc:   scanFunc,
	}

	links, err := getEntityByParams(ctx, opt)
	if err != nil {
		return nil, errs.E(op, err)
	}

	return links, nil
}

const GetLinksCountQuery = `
	SELECT COUNT(*) FROM link WHERE person_id = $1 AND archived = $2 
`

func (r *Repository) GetLinksCountByParams(ctx context.Context, params model.GetLinksParams) (count int64, err error) {
	const op errs.Op = "storage.postgres.repository.GetLinksCountByParams"

	err = r.queries.QueryRow(ctx, GetLinksCountQuery, params.UserId, params.Archived).Scan(&count)
	if err != nil {
		return 0, errs.E(op, err, errs.Database)
	}

	return count, nil
}

const doesLinkBelongsUserQuery = `
	SELECT EXISTS (
		SELECT 1 FROM link WHERE id = $1 AND person_id = $2
	);
`

func (r *Repository) DoesLinkBelongsToUser(ctx context.Context, linkId, userId int64) (belongs bool, err error) {
	const op errs.Op = "storage.postgres.repository.DoesLinkBelongsToUser"

	err = r.queries.QueryRow(ctx, doesLinkBelongsUserQuery, linkId, userId).Scan(&belongs)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, errs.E(op, err, errs.Database)
	}

	return belongs, nil
}

const deleteLinkQuery = `
	DELETE FROM link WHERE id = $1
`

func (r *Repository) DeleteLink(ctx context.Context, linkId int64) error {
	const op errs.Op = "storage.postgres.repository.DeleteLink"

	_, err := r.queries.Exec(ctx, deleteLinkQuery, linkId)
	if err != nil {
		return errs.E(op, err, errs.Database)
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
	const op errs.Op = "storage.postgres.repository.CreateClick"

	row := r.queries.QueryRow(
		ctx, createClickQuery, click.LinkId, click.AdSourceId,
		click.Metadata.UserAgent, click.Metadata.IP, click.Metadata.AccessTime,
	)

	err := row.Scan(&click.Id, &click.AdType)
	if err != nil {
		return nil, errs.E(op, err, errs.Database)
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

func (r *Repository) GetClicksByParams(ctx context.Context, params model.GetClicksParams) ([]*model.Click, error) {
	const op errs.Op = "storage.postgres.repository.GetClicksByParams"

	scanFunc := func(click *model.Click, row pgx.Row) error {
		return row.Scan(
			&click.Id, &click.LinkId, &click.AdSourceId, &click.Metadata.UserAgent,
			&click.Metadata.IP, &click.Metadata.AccessTime, &click.AdType, &click.Status,
		)
	}

	opt := getEntityByParamsOptions[model.Click]{
		db:         r.queries,
		query:      getClicksByParamsQuery,
		args:       []any{params.LinkId},
		pagination: params.Pagination,
		sort:       params.Sort,
		columns:    SortClicksBy,
		scanFunc:   scanFunc,
	}

	clicks, err := getEntityByParams(ctx, opt)
	if err != nil {
		return nil, errs.E(op, err)
	}

	return clicks, nil
}

const getClicksCountQuery = `
	SELECT COUNT(*) FROM click WHERE link_id = $1
`

func (r *Repository) GetClicksCountByParams(ctx context.Context, params model.GetClicksParams) (count int64, err error) {
	const op errs.Op = "storage.postgres.repository.GetClicksCount"

	err = r.queries.QueryRow(ctx, getClicksCountQuery, params.LinkId).Scan(&count)
	if err != nil {
		return 0, errs.E(op, err, errs.Database)
	}

	return count, nil
}
