package redirectManager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"url_shortener/core/model"
	"url_shortener/core/port"
)

type RedirectToADFunc = func()
type RedirectToOriginalFunc = func(original string)

type RedirectManager struct {
	linksStore         port.LinkStorage
	clicksStore        port.ClickStorage
	transactor         port.Transactor
	redirectToAd       RedirectToADFunc
	redirectToOriginal RedirectToOriginalFunc
}

func New(linksStore port.LinkStorage, clicksStore port.ClickStorage, transactor port.Transactor) *RedirectManager {
	return &RedirectManager{
		linksStore:         linksStore,
		clicksStore:        clicksStore,
		transactor:         transactor,
		redirectToAd:       func() {},
		redirectToOriginal: func(original string) { fmt.Println(original) },
	}
}

func (r *RedirectManager) HandleClick(ctx context.Context, alias string, metadata *model.ClickMetadata) error {
	if alias == "" {
		return errors.New("alias can not be empty")
	}

	return r.transactor.WithinTx(ctx, func(ctx context.Context) error {
		link, err := r.linksStore.GetActiveByAlias(ctx, alias)
		if err != nil {
			return err
		}

		if err = r.linksStore.UpdateLastAccess(ctx, link.Id, metadata.AccessTime); err != nil {
			return err
		}

		clickToSave := &model.Click{
			LinkId:   link.Id,
			Status:   model.AdStarted,
			Metadata: metadata,
		}

		clickId, err := r.clicksStore.Save(ctx, clickToSave)
		if err != nil {
			return err
		}

		r.redirectToAd()

		go r.waitForCompleteAd(link.Original, clickId)

		return nil
	})
}

func (r *RedirectManager) waitForCompleteAd(original string, clickId int64) {
	time.Sleep(1 / 2 * time.Second)

	err := r.clicksStore.UpdateStatus(context.Background(), clickId, model.AdCompleted)
	if err != nil {
		log.Printf("Error updating click status: %v", err)
		return
		//todo update status to closed maybe
	}

	r.redirectToOriginal(original)
}
