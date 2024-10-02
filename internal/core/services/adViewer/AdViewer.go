package adViewer

import (
	"context"
	"errors"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type AdViewer struct {
	linksStore  port.LinkStorage
	clicksStore port.ClickStorage
	userStore   port.UserStorage
	transactor  port.Transactor
}

func New(
	linksStore port.LinkStorage,
	clicksStore port.ClickStorage,
	userStore port.UserStorage,
	transactor port.Transactor) *AdViewer {

	return &AdViewer{
		linksStore:  linksStore,
		clicksStore: clicksStore,
		userStore:   userStore,
		transactor:  transactor,
	}
}

func (v *AdViewer) RecordClick(ctx context.Context, alias string, metadata *model.ClickMetadata) (link *model.Link, clickId int64, err error) {
	if alias == "" {
		return nil, -1, errors.New("alias can not be empty")
	}

	err = v.transactor.WithinTx(ctx, func(ctx context.Context) error {
		link, err = v.linksStore.GetActiveByAlias(ctx, alias)
		if err != nil {
			return err
		}

		if err = v.linksStore.UpdateLastAccess(ctx, link.Id, metadata.AccessTime); err != nil {
			return err
		}

		clickToSave := &model.Click{
			LinkId:   link.Id,
			Status:   model.AdStarted,
			Metadata: *metadata,
		}

		clickId, err = v.clicksStore.Save(ctx, clickToSave)
		if err != nil {
			return err
		}

		return nil
	})

	return link, clickId, err
}

func (v *AdViewer) CompleteView(ctx context.Context, clickId int64, userId int64) error {
	return v.transactor.WithinTx(ctx, func(ctx context.Context) error {
		if err := v.clicksStore.UpdateStatus(ctx, clickId, model.AdCompleted); err != nil {
			return err
		}

		payment := 10

		if err := v.userStore.AddToBalance(ctx, userId, payment); err != nil {
			return err
		}

		return nil
	})
}
