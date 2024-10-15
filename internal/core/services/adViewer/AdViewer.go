package adViewer

import (
	"context"
	"fmt"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

//todo shutdown

type AdViewer struct {
	repo       port.Repository
	payer      port.ClickPayer
	adProvider port.AdProvider
	payErrs    chan error
}

func New(repo port.Repository, payer port.ClickPayer, adProvider port.AdProvider) *AdViewer {
	return &AdViewer{
		repo:       repo,
		payer:      payer,
		adProvider: adProvider,
		payErrs:    make(chan error),
	}
}

func (v *AdViewer) OnCompleteErrs() <-chan error {
	return v.payErrs
}

func (v *AdViewer) GetAdPage(ctx context.Context, alias string, metadata model.ClickMetadata) (*model.AdPage, error) {
	const op = "core.services.adViewer.GetAdPage"

	link, err := v.repo.GetLinkByAlias(ctx, alias)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if link.Archived {
		return nil, fmt.Errorf("%s: %w", op, core.ErrOpenArchivedLink)
	}

	adSourceId, err := v.adProvider.GetAdByMetadata(ctx, metadata)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	clickToSave := model.Click{
		LinkId:     link.Id,
		Metadata:   metadata,
		AdSourceId: adSourceId,
	}

	click, err := v.repo.CreateClick(ctx, clickToSave)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	adPage := &model.AdPage{
		AdType:     click.AdType,
		Original:   link.Original,
		ClickId:    click.Id,
		AdSourceId: adSourceId,
	}

	return adPage, nil
}

func (v *AdViewer) CompleteAd(_ context.Context, clickId int64) {
	const op = "core.services.adViewer.CompleteAd"

	if err := v.payer.Pay(context.Background(), clickId); err != nil {
		v.payErrs <- fmt.Errorf("%s: %w", op, err)
	}
}
