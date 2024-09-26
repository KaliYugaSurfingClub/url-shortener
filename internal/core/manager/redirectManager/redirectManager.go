package redirectManager

import (
	"context"
	"errors"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type RedirectManager struct {
	linksStore  port.LinkStorage
	clicksStore port.ClickStorage
	userStore   port.UserStorage
	transactor  port.Transactor
}

func New(
	linksStore port.LinkStorage,
	clicksStore port.ClickStorage,
	userStore port.UserStorage,
	transactor port.Transactor) *RedirectManager {

	return &RedirectManager{
		linksStore:  linksStore,
		clicksStore: clicksStore,
		userStore:   userStore,
		transactor:  transactor,
	}
}

func (r *RedirectManager) Start(ctx context.Context, alias string, metadata *model.ClickMetadata) (link *model.Link, clickId int64, err error) {
	if alias == "" {
		return nil, -1, errors.New("alias can not be empty")
	}

	err = r.transactor.WithinTx(ctx, func(ctx context.Context) error {
		link, err = r.linksStore.GetActiveByAlias(ctx, alias)
		if err != nil {
			return err
		}

		if err = r.linksStore.UpdateLastAccess(ctx, link.Id, metadata.AccessTime); err != nil {
			return err
		}

		clickToSave := &model.Click{
			LinkId:   link.Id,
			Status:   model.AdStarted,
			Metadata: *metadata,
		}

		clickId, err = r.clicksStore.Save(ctx, clickToSave)
		if err != nil {
			return err
		}

		return nil
	})

	return link, clickId, err
}

func (r *RedirectManager) End(ctx context.Context, clickId int64, userId int64) error {
	return r.transactor.WithinTx(ctx, func(ctx context.Context) error {
		if err := r.clicksStore.UpdateStatus(ctx, clickId, model.AdCompleted); err != nil {
			return err
		}

		payment := 10

		if err := r.userStore.AddToBalance(ctx, userId, payment); err != nil {
			return err
		}

		return nil
	})
}

//func (r *RedirectManager) waitForCompleteAd(original string, clickId int64) {
//	time.Sleep(1 / 2 * time.Second)
//
//	err := r.clicksStore.UpdateStatus(context.Background(), clickId, model.AdCompleted)
//	if err != nil {
//		log.Printf("Error updating click status: %v", err)
//		return
//		//todo update status to closed maybe
//	}
//
//	r.redirectToOriginal(original)
//}
