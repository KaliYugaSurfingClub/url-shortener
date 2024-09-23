package redirectManager

import (
	"context"
	"errors"
	"time"
	"url_shortener/core/model"
	"url_shortener/core/port"
)

//todo вот кто то получил оригиналбную ссылку но нужно ли за это давать награду пользователю
//пересенти вознаграждние в другое место

type RedirectManager struct {
	provider   port.LinkProvider
	updater    port.LinkUpdater
	saver      port.ClickSaver
	transfer   port.RewardTransfer
	transactor port.Transactor
}

func New(provider port.LinkProvider, updater port.LinkUpdater, saver port.ClickSaver, transfer port.RewardTransfer, transactor port.Transactor) *RedirectManager {
	return &RedirectManager{
		provider:   provider,
		updater:    updater,
		saver:      saver,
		transfer:   transfer,
		transactor: transactor,
	}
}

func (r *RedirectManager) Process(ctx context.Context, alias string, click model.Click) (string, error) {
	if alias == "" {
		return "", errors.New("alias can not be empty")
	}

	var original string

	err := r.transactor.WithinTx(ctx, func(ctx context.Context) error {
		link, err := r.provider.GetActualByAlias(ctx, alias)
		if err != nil {
			return err
		}

		if err = r.updater.UpdateLastAccess(ctx, link.Id, time.Now()); err != nil {
			return err
		}

		if err = r.transfer.TransferReward(link.CreatedBy); err != nil {
			return err
		}

		click.LinkId = link.Id
		if _, err = r.saver.Save(ctx, click); err != nil {
			return err
		}

		original = link.Original

		return nil
	})

	if err != nil {
		return "", err
	}

	return original, nil
}
