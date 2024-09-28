package linkCreator

import (
	"context"
	"errors"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
)

type LinkCreator struct {
	storage         port.LinkStorage
	generator       port.Generator
	triesToGenerate int
}

func New(storage port.LinkStorage, generator port.Generator, triesToGenerate int) (*LinkCreator, error) {
	if generator == nil {
		return nil, errors.New("generator can't be nil")
	}

	if triesToGenerate <= 0 {
		return nil, errors.New("triesToGenerate can not be less than 0")
	}

	return &LinkCreator{
		storage,
		generator,
		triesToGenerate,
	}, nil
}

func (c *LinkCreator) Short(ctx context.Context, toSave model.Link) (*model.Link, error) {
	if toSave.Alias == "" {
		return c.generateAndSave(ctx, &toSave)
	}

	return c.save(ctx, &toSave)
}

func (c *LinkCreator) save(ctx context.Context, toSave *model.Link) (saved *model.Link, err error) {
	defer utils.WithinOp("core.manager.LinkCreator.save", &err)

	if toSave.CustomName == "" {
		toSave.CustomName = toSave.Alias
	}

	return c.storage.Save(ctx, *toSave)
}

func (c *LinkCreator) generateAndSave(ctx context.Context, toSave *model.Link) (saved *model.Link, err error) {
	defer utils.WithinOp("core.manager.LinkCreator.generateAndSave", &err)

	noCustomName := toSave.CustomName == ""

	for i := 0; i < c.triesToGenerate; i++ {
		toSave.Alias = c.generator.Generate()

		if noCustomName {
			toSave.CustomName = toSave.Alias
		}

		saved, err = c.storage.Save(ctx, *toSave)
		if err == nil {
			return saved, nil
		}
		if !errors.Is(err, core.ErrAliasExists) && !errors.Is(err, core.ErrCustomNameExists) {
			return nil, err
		}
		if !noCustomName && errors.Is(err, core.ErrCustomNameExists) {
			return nil, err
		}
	}

	return nil, core.ErrCantGenerateInTries
}
