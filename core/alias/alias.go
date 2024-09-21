package alias

import (
	"context"
	"errors"
	"url_shortener/core"
)

type linkProvider interface {
	GetOriginalByAlias(ctx context.Context, alias string) (int64, string, error)
}

type linkSaver interface {
	Save(ctx context.Context, original string, alias string) (int64, error)
}

type generator interface {
	Generate() string
}

type Alias struct {
	saver           linkSaver
	generator       generator
	triesToGenerate int
}

func New(saver linkSaver, generator generator, triesToGenerate int) *Alias {
	return &Alias{
		saver:           saver,
		generator:       generator,
		triesToGenerate: triesToGenerate,
	}
}

func (a *Alias) Save(ctx context.Context, original string, alias string) (string, error) {
	if alias == "" {
		var err error

		for i := 0; i < a.triesToGenerate; i++ {
			alias = a.generator.Generate()

			_, err = a.saver.Save(ctx, original, alias)

			if err == nil || errors.Is(err, core.ErrLastInsertId) {
				return alias, nil
			}
		}

		//todo log if generating in allowed number of tries was failed

		return "", err
	}

	if _, err := a.saver.Save(ctx, original, alias); err != nil {
		return "", err
	}

	return alias, nil
}
