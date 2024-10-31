package linkShortener

import (
	"context"
	"github.com/KaliYugaSurfingClub/pkg/errs"
	"link-service/internal/core"
	"link-service/internal/core/model"
	"link-service/internal/core/port"
)

type LinkShortener struct {
	storage         port.Repository
	generator       port.Generator
	triesToGenerate int
}

func New(storage port.Repository, generator port.Generator, triesToGenerate int) *LinkShortener {
	return &LinkShortener{
		storage,
		generator,
		triesToGenerate,
	}
}

func (s *LinkShortener) Short(ctx context.Context, toSave model.Link) (*model.Link, error) {
	const op errs.Op = "core.services.LinkShortener.Short"

	saveFunc := s.saveWithAlias
	if toSave.Alias == "" {
		saveFunc = s.saveWithoutAlias
	}

	link, err := saveFunc(ctx, &toSave)
	if err != nil {
		return nil, errs.E(op, err)
	}

	return link, nil
}

func (s *LinkShortener) saveWithAlias(ctx context.Context, toSave *model.Link) (link *model.Link, err error) {
	const op errs.Op = "core.services.LinkShortener.saveWithAlias"

	if toSave.CustomName == "" {
		toSave.CustomName = toSave.Alias
	}

	if link, err = s.storage.CreateLink(ctx, *toSave); err != nil {
		return nil, errs.E(op, err)
	}

	return link, nil
}

func (s *LinkShortener) saveWithoutAlias(ctx context.Context, toSave *model.Link) (*model.Link, error) {
	const op errs.Op = "core.services.LinkShortener.saveWithoutAlias"

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
