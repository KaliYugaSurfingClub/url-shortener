package port

import (
	"context"
	"link-service/internal/core/model"
)

type Repository interface {
	CreateLink(ctx context.Context, link model.Link) (*model.Link, error)
	GetLinkByAlias(ctx context.Context, alias string) (*model.Link, error)
	GetOriginalByClickId(ctx context.Context, clickId int64) (*model.Link, error)
	GetLinksByParams(ctx context.Context, params model.GetLinksParams) ([]*model.Link, error)
	GetLinksCountByParams(ctx context.Context, params model.GetLinksParams) (int64, error)
	DoesLinkBelongsToUser(ctx context.Context, linkId int64, userId int64) (bool, error)
	DeleteLink(ctx context.Context, linkId int64) error

	CreateClick(ctx context.Context, click model.Click) (*model.Click, error)
	GetClicksByParams(ctx context.Context, params model.GetClicksParams) ([]*model.Click, error)
	GetClicksCountByParams(ctx context.Context, params model.GetClicksParams) (int64, error)

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
