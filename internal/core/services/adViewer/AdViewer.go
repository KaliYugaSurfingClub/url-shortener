package adViewer

import (
	"context"
	"fmt"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"sync"
)

type AdViewer struct {
	links          port.LinkStorage
	clicks         port.ClickStorage
	notifier       port.ClickNotifier
	payer          port.Payer
	transactor     port.Transactor
	processedLinks sync.Map //todo delete after 10 minutes and mark as closed maybe redis maybe
	//каждые N минут просматривать все ссылки в постгресе те которые открыты давно но не выполнены закрыть и уведомалят
	//сделать батчинг запросов
	errs chan error
}

func New(
	linksStorage port.LinkStorage,
	clicksStorage port.ClickStorage,
	payer port.Payer,
	notifier port.ClickNotifier,
	transactor port.Transactor,
) *AdViewer {
	return &AdViewer{
		links:      linksStorage,
		clicks:     clicksStorage,
		transactor: transactor,
		notifier:   notifier,
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

		link.ClicksCount++

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

	v.processedLinks.Store(clickId, link)
	v.notifier.NotifyOpen(ctx, link, clickId)

	return link.Original, clickId, nil
}

// todo saga transaction
func (v *AdViewer) CompleteAd(ctx context.Context, clickId int64) error { //todo maybe return bool
	const op = "core.services.adViewer.CompleteAd"

	link, ok := v.processedLinks.Load(clickId)
	if !ok {
		return fmt.Errorf("%s: click not found", core.ErrTryToCompleteUnexactingClick)
	}

	v.processedLinks.Delete(clickId)

	go func() {
		link := link.(*model.Link)

		if err := v.payer.Pay(ctx, link.CreatedBy); err != nil {
			v.errs <- fmt.Errorf("%s: %w", op, err)

			if err := v.clicks.UpdateStatus(ctx, clickId, model.AdClosed); err != nil {
				v.errs <- fmt.Errorf("%s: %w", op, err)
				return
			}

			v.notifier.NotifyClosed(ctx, link, clickId)
			return
		}

		if err := v.clicks.UpdateStatus(ctx, clickId, model.AdCompleted); err != nil {
			v.errs <- fmt.Errorf("%s: %w", op, err)
		} else {
			v.notifier.NotifyCompleted(ctx, link, clickId)
		}
	}()

	return nil
}
