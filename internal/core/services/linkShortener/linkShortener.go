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
	storage         port.Repository
	generator       port.Generator
	triesToGenerate int
}

func New(storage port.Repository, generator port.Generator, triesToGenerate int) (*LinkShortener, error) {
	if triesToGenerate <= 0 {
		return nil, errors.New("triesToGenerate can not be less than 0")
	}

	return &LinkShortener{
		storage,
		generator,
		triesToGenerate,
	}, nil
}

func (s *LinkShortener) Short(ctx context.Context, toSave model.Link) (_ *model.Link, err error) {
	defer utils.WithinOp("core.manager.LinkShortener.Short", &err)

	if toSave.Alias == "" {
		return s.generateAndSave(ctx, &toSave)
	}

	return s.save(ctx, &toSave)
}

func (s *LinkShortener) save(ctx context.Context, toSave *model.Link) (*model.Link, error) {
	if toSave.CustomName == "" {
		toSave.CustomName = toSave.Alias
	}

	return s.storage.CreateLink(ctx, *toSave)
}

func (s *LinkShortener) generateAndSave(ctx context.Context, toSave *model.Link) (*model.Link, error) {
	noCustomName := toSave.CustomName == ""

	for i := 0; i < s.triesToGenerate; i++ {
		toSave.Alias = s.generator.Generate()

		if noCustomName {
			toSave.CustomName = toSave.Alias
		}

		saved, err := s.storage.CreateLink(ctx, *toSave)
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
