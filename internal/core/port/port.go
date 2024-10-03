package port

import (
	"context"
	"shortener/internal/core/model"
	"time"
)

type LinkStorage interface {
	GetActiveByAlias(ctx context.Context, alias string) (*model.Link, error)
	GetById(ctx context.Context, id int64) (*model.Link, error)
	AliasExists(ctx context.Context, alias string) (bool, error)
	CustomNameExists(ctx context.Context, customName string, userId int64) (bool, error)
	Save(ctx context.Context, link model.Link) (*model.Link, error)
	UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error
	GetCountByUserId(ctx context.Context, userId int64, params model.LinkFilter) (int64, error)
	GetByUserId(ctx context.Context, userId int64, params model.GetLinksParams) ([]*model.Link, error)
}

type ClickStorage interface {
	Save(ctx context.Context, click *model.Click) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status model.AdStatus) error
	GetCountByLinkId(ctx context.Context, linkId int64, params model.GetClicksParams) (int64, error)
	GetByLinkId(ctx context.Context, linkId int64, params model.GetClicksParams) ([]*model.Click, error)
}

type UserStorage interface {
	AddToBalance(ctx context.Context, id int64, payment int) error
}

type Generator interface {
	Generate() string
}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(tx context.Context) error) error
}
