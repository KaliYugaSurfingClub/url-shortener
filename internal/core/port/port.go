package port

import (
	"context"
	model2 "shortener/internal/core/model"
	"time"
)

type LinkStorage interface {
	GetActiveByAlias(ctx context.Context, alias string) (*model2.Link, error)
	Save(ctx context.Context, link *model2.Link) (int64, error)
	UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error
}

type ClickStorage interface {
	Save(ctx context.Context, click *model2.Click) (int64, error)
	UpdateStatus(ctx context.Context, id int64, status model2.AdStatus) error
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
