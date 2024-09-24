package port

import (
	"context"
	"time"
	"url_shortener/core/model"
)

type LinkProvider interface {
	GetActiveByAlias(ctx context.Context, alias string) (*model.Link, error)
}

type LinkSaver interface {
	Save(ctx context.Context, link model.Link) (int64, error)
}

type LinkUpdater interface {
	UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error
}

type ClickSaver interface {
	Save(ctx context.Context, click model.Click) (int64, error)
}

type RewardTransfer interface {
	TransferReward(userId int64) error
}

type AliasGenerator interface {
	Generate() string
}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(tx context.Context) error) error
}
