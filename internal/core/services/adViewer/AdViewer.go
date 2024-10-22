package adViewer

import (
	"context"
	"errors"
	"shortener/errs"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

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
		return nil, errs.E(op, err)
	}
	if link.Archived {
		return nil, errs.E(op, errors.New("someone tries to open archived link"), errs.NotExist)
	}

	adSourceId, err := v.adProvider.GetAdByMetadata(ctx, metadata)
	if err != nil {
		return nil, errs.E(op, err)
	}

	clickToSave := model.Click{
		LinkId:     link.Id,
		Metadata:   metadata,
		AdSourceId: adSourceId,
	}

	click, err := v.repo.CreateClick(ctx, clickToSave)
	if err != nil {
		return nil, errs.E(op, err)
	}

	adPage := &model.AdPage{
		AdType:     click.AdType,
		ClickId:    click.Id,
		AdSourceId: adSourceId,
	}

	return adPage, nil
}

func (v *AdViewer) CompleteAd(ctx context.Context, clickId int64) (string, error) {
	const op = "core.services.adViewer.CompleteAd"

	go func() {
		if err := v.payer.Pay(context.Background(), clickId); err != nil {
			v.payErrs <- errs.E(op, err)
		}
	}()

	original, err := v.repo.GetOriginalByClickId(ctx, clickId) //todo if archived
	if err != nil {
		return "", errs.E(op, err)
	}

	return original, nil
}
