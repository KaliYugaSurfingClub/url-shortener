package adViewer

import (
	"context"
	"fmt"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"sync"
)

type AdViewer struct {
	links              port.LinkStorage
	clicks             port.ClickStorage
	notifier           port.ClickNotifier
	payer              port.Payer
	transactor         port.Transactor
	clickToLinkCreator sync.Map
}

func New(
	linksStorage port.LinkStorage,
	clicksStorage port.ClickStorage,
	transactor port.Transactor,
	payer port.Payer,
) *AdViewer {
	return &AdViewer{
		links:      linksStorage,
		clicks:     clicksStorage,
		transactor: transactor,
		payer:      payer,
	}
}

func (v *AdViewer) OnClick(ctx context.Context, alias string, metadata model.ClickMetadata) (original string, clickId int64, err error) {
	const op = "core.services.adViewer.OnClick"

	var link *model.Link

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

	if err != nil {
		return "", 0, fmt.Errorf("%s: %w", op, err)
	}

	v.clickToLinkCreator.Store(clickId, link.CreatedBy)
	v.notifier.NotifyOpen(ctx, link.CreatedBy, clickId)

	return link.Original, clickId, nil
}

func (v *AdViewer) CompleteAd(ctx context.Context, clickId int64) error {
	const op = "core.services.adViewer.CompleteAd"

	if err := v.clicks.UpdateStatus(ctx, clickId, model.AdCompleted); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	userId, ok := v.clickToLinkCreator.Load(clickId)
	if !ok {
		return fmt.Errorf("%s: click not found", op) //todo
	}

	v.notifier.NotifyOpen(ctx, userId.(int64), clickId)

	v.clickToLinkCreator.Delete(clickId)

	//todo pay

	return nil
}
