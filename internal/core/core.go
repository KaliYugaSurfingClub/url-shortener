package core

import (
	"errors"
	"time"
)

type linkProvider interface {
	GetOriginalByAlias(alias string) (int, string, error)
}

type linkUpdater interface {
	UpdateLastAccess(ID int, timestamp time.Time) error
}

type clickSaver interface {
	Save(aliasID int) error
}

type redirectFunc = func(url string)

func Redirect(alias string, provider linkProvider, updater linkUpdater, saver clickSaver, redirect redirectFunc) error {
	if alias == "" {
		return errors.New("alias can not be empty")
	}

	id, original, err := provider.GetOriginalByAlias(alias)
	if err != nil {
		return err
	}

	if err = updater.UpdateLastAccess(id, time.Now()); err != nil {
		return err
	}

	if err = saver.Save(id); err != nil {
		return err
	}

	redirect(original)

	return nil
}
