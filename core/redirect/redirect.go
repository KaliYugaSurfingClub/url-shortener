package redirect

import (
	"context"
	"errors"
	"time"
)

type linkProvider interface {
	GetOriginalByAlias(ctx context.Context, alias string) (int64, string, error)
}

type linkUpdater interface {
	UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error
}

type clickSaver interface {
	Save(ctx context.Context, aliasId int64) error
}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(tx context.Context) error) error
}

type redirectFunc = func(url string)

type Redirect struct {
	provider   linkProvider
	updater    linkUpdater
	saver      clickSaver
	redirect   redirectFunc
	transactor Transactor
}

func New(provider linkProvider, updater linkUpdater, saver clickSaver, redirect redirectFunc, transactor Transactor) *Redirect {
	return &Redirect{
		provider:   provider,
		updater:    updater,
		saver:      saver,
		redirect:   redirect,
		transactor: transactor,
	}
}

func (r *Redirect) To(ctx context.Context, alias string) error {
	if alias == "" {
		return errors.New("alias can not be empty")
	}

	return r.transactor.WithinTx(ctx, func(ctx context.Context) error {
		id, original, err := r.provider.GetOriginalByAlias(ctx, alias)
		if err != nil {
			return err
		}

		if err = r.updater.UpdateLastAccess(ctx, id, time.Now()); err != nil {
			return err
		}

		if err = r.saver.Save(ctx, id); err != nil {
			return err
		}

		r.redirect(original)

		return nil
	})
}
