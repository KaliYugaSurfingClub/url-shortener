package port

import (
	"context"
	"shortener/internal/core/model"
	"time"
)

type LinkStorage interface {
	GetActiveByAlias(ctx context.Context, alias string) (*model.Link, error)
	AliasExists(ctx context.Context, alias string) (bool, error)
	CustomNameExists(ctx context.Context, customName string, userId int64) (bool, error)
	Save(ctx context.Context, link *model.Link) (int64, error)
	UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error
}

type ClickStorage interface {
	Save(ctx context.Context, click *model.Click) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status model.AdStatus) error
}

type UserStorage interface {
	AddToBalance(ctx context.Context, id int64, payment int) error
}

type AliasGenerator interface {
	Generate() string
}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(tx context.Context) error) error
}
