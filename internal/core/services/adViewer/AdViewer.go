package adViewer

import (
	"context"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type AdViewer struct {
	links      port.LinkStorage
	clicks     port.ClickStorage
	users      port.UserStorage
	transactor port.Transactor
}

func New(
	linksStorage port.LinkStorage,
	clicksStorage port.ClickStorage,
	userStorage port.UserStorage,
	transactor port.Transactor) *AdViewer {

	return &AdViewer{
		links:      linksStorage,
		clicks:     clicksStorage,
		users:      userStorage,
		transactor: transactor,
	}
}

func (v *AdViewer) RecordClick(ctx context.Context, alias string, metadata model.ClickMetadata) (link *model.Link, clickId int64, err error) {
	err = v.transactor.WithinTx(ctx, func(ctx context.Context) error {
		link, err = v.links.GetActiveByAlias(ctx, alias)
		if err != nil {
			return err
		}

		if err = v.links.UpdateLastAccess(ctx, link.Id, metadata.AccessTime); err != nil {
			return err
		}

		clickToSave := &model.Click{
			LinkId:   link.Id,
			Status:   model.AdStarted,
			Metadata: metadata,
		}

		clickId, err = v.clicks.Save(ctx, clickToSave)
		if err != nil {
			return err
		}

		return nil
	})

	return link, clickId, err
}

func (v *AdViewer) CompleteView(ctx context.Context, clickId int64, userId int64) error {
	return v.transactor.WithinTx(ctx, func(ctx context.Context) error {
		if err := v.clicks.UpdateStatus(ctx, clickId, model.AdCompleted); err != nil {
			return err
		}

		payment := 10

		//todo referal
		if err := v.users.AddToBalance(ctx, userId, payment); err != nil {
			return err
		}

		return nil
	})
}
