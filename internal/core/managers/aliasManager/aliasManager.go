package aliasManager

import (
	"context"
	"errors"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type AliasManager struct {
	store           port.LinkStorage
	generator       port.AliasGenerator
	triesToGenerate int
}

func New(store port.LinkStorage, generator port.AliasGenerator, triesToGenerate int) (*AliasManager, error) {
	if generator == nil {
		return nil, errors.New("generator can't be nil")
	}

	if triesToGenerate <= 0 {
		return nil, errors.New("triesToGenerate can not be less than 0")
	}

	return &AliasManager{
		store,
		generator,
		triesToGenerate,
	}, nil
}

func (a *AliasManager) Save(ctx context.Context, link *model.Link) (string, error) {
	exists, err := a.store.CustomNameExists(ctx, link.CustomName, link.CreatedBy)
	if err != nil {
		return "", err
	}
	if exists {
		return "", core.ErrCustomNameExists
	}

	if link.Alias == "" {
		return a.generateAndSave(ctx, *link)
	}

	exists, err = a.store.AliasExists(ctx, link.Alias)
	if err != nil {
		return "", err
	}
	if exists {
		return "", core.ErrAliasExists
	}

	if _, err = a.store.Save(ctx, link); err != nil {
		return "", err
	}

	return link.Alias, nil
}

func (a *AliasManager) generateAndSave(ctx context.Context, link model.Link) (string, error) {
	var i int

	errs := make([]error, a.triesToGenerate)

	for i = 0; i < a.triesToGenerate; i++ {
		link.Alias = a.generator.Generate()

		_, err := a.store.Save(ctx, &link)

		if err == nil {
			return link.Alias, nil
		} else {
			errs = append(errs, err)
		}
	}

	if i == a.triesToGenerate {
		errs = append(errs, core.ErrCantGenerateInTries)
	}

	if len(errs) != 0 {
		return "", errors.Join(errs...)
	}

	return link.Alias, nil
}
