package port

import (
	"context"
	"shortener/internal/core/model"
	"time"
)

type LinkStorage interface {
	GetActiveByAlias(ctx context.Context, alias string) (*model.Link, error)
	GetById(ctx context.Context, id int64) (*model.Link, error)
	Save(ctx context.Context, link model.Link) (*model.Link, error)
	UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error
	GetCountByUserId(ctx context.Context, params model.GetLinksParams) (int64, error)
	GetByUserId(ctx context.Context, params model.GetLinksParams) ([]*model.Link, error)
	DoesLinkBelongUser(ctx context.Context, linkId int64, userId int64) (_ bool, err error)
}

type ClickStorage interface {
	Save(ctx context.Context, click *model.Click) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status model.AdStatus) error
	GetCountByLinkId(ctx context.Context, params model.GetClicksParams) (int64, error)
	GetByLinkId(ctx context.Context, params model.GetClicksParams) ([]*model.Click, error)
}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(tx context.Context) error) error
}

type ClickNotifier interface {
	NotifyOpen(ctx context.Context, userId int64, clickId int64)
	NotifyWatched(ctx context.Context, linkName string, clickId int64)
}

type Payer interface {
	Pay(ctx context.Context, userId int64) error
}

type Generator interface {
	Generate() string
}
