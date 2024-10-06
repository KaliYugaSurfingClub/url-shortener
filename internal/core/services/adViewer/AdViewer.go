package adViewer

import (
	"context"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type AdViewer struct {
	links      port.LinkStorage
	clicks     port.ClickStorage
	transactor port.Transactor
}

func New(
	linksStorage port.LinkStorage,
	clicksStorage port.ClickStorage,
	transactor port.Transactor) *AdViewer {

	return &AdViewer{
		links:      linksStorage,
		clicks:     clicksStorage,
		transactor: transactor,
	}
}

func (v *AdViewer) OnClick(ctx context.Context, alias string, metadata model.ClickMetadata) (original string, clickId int64, err error) {
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
		return "", 0, err
	}

	return link.Original, clickId, nil
}
