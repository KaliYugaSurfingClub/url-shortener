package linkShortener

import (
	"context"
	"shortener/errs"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"shortener/internal/core/port"
)

type LinkShortener struct {
	storage         port.Repository
	generator       port.Generator
	triesToGenerate int
}

func New(storage port.Repository, generator port.Generator, triesToGenerate int) (*LinkShortener, error) {
	return &LinkShortener{
		storage,
		generator,
		triesToGenerate,
	}, nil
}

func (s *LinkShortener) Short(ctx context.Context, toSave model.Link) (link *model.Link, err error) {
	const op errs.Op = "core.services.LinkShortener.Short"

	if toSave.Alias == "" {
		return s.generateAndSave(ctx, &toSave)
	}

	if link, err = s.save(ctx, &toSave); err != nil {
		return nil, errs.E(op, err)
	}

	return link, nil
}

func (s *LinkShortener) save(ctx context.Context, toSave *model.Link) (link *model.Link, err error) {
	const op errs.Op = "core.services.LinkShortener.save"

	if toSave.CustomName == "" {
		toSave.CustomName = toSave.Alias
	}

	if link, err = s.storage.CreateLink(ctx, *toSave); err != nil {
		return nil, errs.E(op, err)
	}

	return link, nil
}

func (s *LinkShortener) generateAndSave(ctx context.Context, toSave *model.Link) (*model.Link, error) {
	const op errs.Op = "core.services.LinkShortener.generateAndSave"

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
		//user already has link with this customName
		if !noCustomName && err.(*errs.Error).Code == core.CustomNameExistsCode {
			return nil, errs.E(op, err)
		}
		//some internal error
		if errs.KindIs(err, errs.Database) {
			return nil, errs.E(op, err)
		}
	}

	return nil, errs.E(op, "failed generating alias in tries", errs.Internal)
}
