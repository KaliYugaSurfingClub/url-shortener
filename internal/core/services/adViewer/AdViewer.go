package adViewer

import (
	"context"
	"fmt"
	"github.com/thoas/go-funk"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"time"
)

//todo shutdown

type AdViewer struct {
	links       port.LinkStorage
	clicks      port.ClickStorage
	payer       port.Payer
	transactor  port.Transactor
	payErrs     chan error
	cleanerErrs chan error
}

func New(
	linksStorage port.LinkStorage,
	clicksStorage port.ClickStorage,
	payer port.Payer,
	transactor port.Transactor,
	payErrsBuffer int,
) *AdViewer {
	return &AdViewer{
		links:      linksStorage,
		clicks:     clicksStorage,
		transactor: transactor,
		payer:      payer,
		payErrs:    make(chan error, payErrsBuffer),
	}
}

func (v *AdViewer) OnCompleteErrs() <-chan error {
	return v.payErrs
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

	return link.Original, clickId, nil
}

func (v *AdViewer) OnComplete(ctx context.Context, clickId int64) error {
	const op = "core.services.adViewer.OnComplete"

	//todo use cache to optimize this
	//todo or join
	click, err := v.clicks.GetById(ctx, clickId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		//todo change ad status

		link, err := v.links.GetById(ctx, click.LinkId)
		if err != nil {
			v.payErrs <- fmt.Errorf("%s: %w", op, err)
		}

		if err := v.payer.Pay(ctx, link.CreatedBy); err != nil {
			v.payErrs <- fmt.Errorf("%s: %w", op, err)
		}
	}()

	return nil
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

		//todo notify
		//for _, click := range toClose {
		//	if link, ok := v.processedLinks.LoadAndDelete(click.Id); ok {
		//		v.closeAdAndNotify(context.Background(), link.(*model.Link), click.Id)
		//	}
		//}

		time.Sleep(timeout)
	}
}
