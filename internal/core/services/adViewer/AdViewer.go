package adViewer

import (
	"context"
	"fmt"
	"github.com/thoas/go-funk"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"sync"
	"time"
)

//todo shutdown

type AdViewer struct {
	links          port.LinkStorage
	clicks         port.ClickStorage
	notifier       port.ClickNotifier
	payer          port.Payer
	transactor     port.Transactor
	processedLinks sync.Map //todo delete after 10 minutes and mark as closed maybe redis maybe
	//каждые N минут просматривать все ссылки в постгресе те которые открыты давно но не выполнены закрыть и уведомалят
	//сделать батчинг запросов
	onCompleteErrs chan error
	cleanerErrs    chan error
}

func New(
	linksStorage port.LinkStorage,
	clicksStorage port.ClickStorage,
	payer port.Payer,
	notifier port.ClickNotifier,
	transactor port.Transactor,
	onCompleteBuffer int,
	cleanerBuffer int,
) *AdViewer {
	return &AdViewer{
		links:          linksStorage,
		clicks:         clicksStorage,
		transactor:     transactor,
		notifier:       notifier,
		payer:          payer,
		onCompleteErrs: make(chan error, onCompleteBuffer),
		cleanerErrs:    make(chan error, cleanerBuffer),
	}
}

func (v *AdViewer) OnCompleteErrs() <-chan error {
	return v.onCompleteErrs
}

func (v *AdViewer) CleanerErrs() <-chan error {
	return v.cleanerErrs
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
		link.LastAccessTime = &metadata.AccessTime

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

func (v *AdViewer) OnComplete(ctx context.Context, clickId int64) bool {
	link, ok := v.processedLinks.LoadAndDelete(clickId)
	if !ok {
		return false
	}

	//todo put to the queue
	//or in startClosing check any of expired sessions has payment (one query)
	go func() {
		link := link.(*model.Link)

		if err := v.payer.Pay(ctx, link.CreatedBy); err != nil {
			v.closeAd(ctx, link, clickId) //not paid, not closed, not notified - OK, not paid closed and notified - OK
			return
		}

		v.completeAd(ctx, link, clickId) //paid, not completed, not notified - BAD, paid, completed. notified - OK //todo
	}()

	return true
}

func (v *AdViewer) StartCleaningExpiredSessions(sessionLifetime time.Duration, timeout time.Duration, batchSize int64) {
	//todo check if payment was done and mark as completed
	for {
		toClose, err := v.clicks.GetExpiredClickSessions(context.Background(), sessionLifetime, batchSize)
		if err != nil {
			v.cleanerErrs <- err
		}

		toCloseIds := funk.Map(toClose, func(click *model.Click) int64 { return click.Id }).([]int64)

		if err = v.clicks.BatchUpdateStatus(context.Background(), toCloseIds, model.AdClosed); err != nil {
			v.cleanerErrs <- err
		}

		//for _, click := range toClose {
		//	if link, ok := v.processedLinks.LoadAndDelete(click.Id); ok {
		//		v.closeAd(context.Background(), link.(*model.Link), click.Id)
		//	}
		//}

		time.Sleep(timeout)
	}
}

func (v *AdViewer) closeAd(ctx context.Context, link *model.Link, clickId int64) {
	const op = "core.services.adViewer.CloseAd"

	err := v.clicks.UpdateStatus(ctx, clickId, model.AdClosed)
	if err != nil {
		v.onCompleteErrs <- fmt.Errorf("%s: %w", op, err)
	} else {
		v.notifier.NotifyClosed(ctx, link, clickId)
	}
}

func (v *AdViewer) completeAd(ctx context.Context, link *model.Link, clickId int64) {
	const op = "core.services.adViewer.OnComplete"

	err := v.clicks.UpdateStatus(ctx, clickId, model.AdCompleted)
	if err != nil {
		v.onCompleteErrs <- fmt.Errorf("%s: %w", op, err)
	} else {
		v.notifier.NotifyClosed(ctx, link, clickId)
	}
}
