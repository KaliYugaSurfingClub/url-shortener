package aliasManager

import (
	"context"
	"errors"
	"fmt"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
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

func (a *AliasManager) Save(ctx context.Context, link *model.Link) (_ string, err error) {
	defer utils.WithinOp("core.manager.AliasManager.Save", &err)

	err = handleExits(a.store.CustomNameExists(ctx, link.CustomName, link.CreatedBy))
	if err != nil {
		return "", err
	}

	if link.Alias == "" {
		return a.generateAndSave(ctx, *link)
	}

	if err = handleExits(a.store.AliasExists(ctx, link.Alias)); err != nil {
		return "", err
	}

	if _, err = a.store.Save(ctx, link); err != nil {
		return "", err
	}

	return link.Alias, nil
}

func (a *AliasManager) generateAndSave(ctx context.Context, link model.Link) (_ string, err error) {
	defer utils.WithinOp("core.manager.AliasManager.generateAndSave", &err)

	errs := make([]error, a.triesToGenerate+1)

	var i int

	for i = 0; i < a.triesToGenerate; i++ {
		link.Alias = a.generator.Generate()

		_, err = a.store.Save(ctx, &link)

		if err == nil {
			return link.Alias, nil
		}

		errs = append(errs, fmt.Errorf("%w (generate - %s)", err, link.Alias))
	}

	if i == a.triesToGenerate {
		errs = append(errs, core.ErrCantGenerateInTries)
	}

	if len(errs) > 0 {
		return "", errors.Join(errs...)
	}

	return link.Alias, nil
}

func handleExits(exists bool, err error) error {
	if err != nil {
		return err
	}
	if exists {
		return core.ErrAliasExists
	}

	return nil
}
