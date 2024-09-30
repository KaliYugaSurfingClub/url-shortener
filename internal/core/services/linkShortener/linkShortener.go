package linkShortener

import (
	"context"
	"errors"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
	"shortener/internal/utils"
)

type LinkShortener struct {
	storage         port.LinkStorage
	generator       port.Generator
	triesToGenerate int
}

func New(storage port.LinkStorage, generator port.Generator, triesToGenerate int) (*LinkShortener, error) {
	if generator == nil {
		return nil, errors.New("generator can't be nil")
	}

	if triesToGenerate <= 0 {
		return nil, errors.New("triesToGenerate can not be less than 0")
	}

	return &LinkShortener{
		storage,
		generator,
		triesToGenerate,
	}, nil
}

func (s *LinkShortener) Short(ctx context.Context, toSave model.Link) (*model.Link, error) {
	if toSave.Alias == "" {
		return s.generateAndSave(ctx, &toSave)
	}

	return s.save(ctx, &toSave)
}

func (s *LinkShortener) save(ctx context.Context, toSave *model.Link) (saved *model.Link, err error) {
	defer utils.WithinOp("core.manager.LinkShortener.save", &err)

	if toSave.CustomName == "" {
		toSave.CustomName = toSave.Alias
	}

	return s.storage.Save(ctx, *toSave)
}

func (s *LinkShortener) generateAndSave(ctx context.Context, toSave *model.Link) (saved *model.Link, err error) {
	defer utils.WithinOp("core.manager.LinkShortener.generateAndSave", &err)

	noCustomName := toSave.CustomName == ""

	for i := 0; i < s.triesToGenerate; i++ {
		toSave.Alias = s.generator.Generate()

		if noCustomName {
			toSave.CustomName = toSave.Alias
		}

		saved, err = s.storage.Save(ctx, *toSave)
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