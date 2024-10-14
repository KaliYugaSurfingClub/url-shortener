package port

import (
	"context"
	"shortener/internal/core/model"
)

type Repository interface {
	CreateLink(ctx context.Context, link model.Link) (*model.Link, error)
	GetLinkByAlias(ctx context.Context, alias string) (*model.Link, error)
	GetLinksByParams(ctx context.Context, params model.GetLinksParams) (_ []*model.Link, err error)
	GetLinksCountByParams(ctx context.Context, params model.GetLinksParams) (count int64, err error)
	DoesLinkBelongsToUser(ctx context.Context, linkId int64, userId int64) (belongs bool, err error)
	DeleteLink(ctx context.Context, linkId int64) error

	CreateClick(ctx context.Context, click model.Click) (*model.Click, error)
	GetClicksByParams(ctx context.Context, params model.GetClicksParams) (_ []*model.Click, err error)
	GetClicksCountByParams(ctx context.Context, params model.GetClicksParams) (count int64, err error)

	WithinTx(ctx context.Context, fn func(tx context.Context) error) error
}

type ClickPayer interface {
	Pay(ctx context.Context, clickId int64) error
}

type AdProvider interface {
	GetAdByMetadata(ctx context.Context, metadata model.ClickMetadata) (int64, error)
}

type Generator interface {
	Generate() string
}
