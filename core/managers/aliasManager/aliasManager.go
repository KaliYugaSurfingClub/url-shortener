package aliasManager

import (
	"context"
	"errors"
	"url_shortener/core"
	"url_shortener/core/model"
	"url_shortener/core/port"
)

type AliasManager struct {
	saver           port.LinkSaver
	generator       port.AliasGenerator
	triesToGenerate int
}

func New(saver port.LinkSaver, generator port.AliasGenerator, triesToGenerate int) (*AliasManager, error) {
	if generator == nil {
		return nil, errors.New("generator can't be nil")
	}

	if triesToGenerate <= 0 {
		return nil, errors.New("triesToGenerate can not be less than 0")
	}

	return &AliasManager{
		saver,
		generator,
		triesToGenerate,
	}, nil
}

func (a *AliasManager) Save(ctx context.Context, link model.Link) (string, error) {
	if link.Alias == "" {
		return a.generateAndSave(ctx, link)
	}

	if _, err := a.saver.Save(ctx, link); err != nil {
		return "", err
	}

	return link.Alias, nil
}

func (a *AliasManager) generateAndSave(ctx context.Context, link model.Link) (_ string, err error) {
	var i int

	for i = 0; i < a.triesToGenerate; i++ {
		link.Alias = a.generator.Generate()

		_, err = a.saver.Save(ctx, link)

		if err == nil {
			return link.Alias, nil
		}
	}

	if i == a.triesToGenerate {
		return "", core.ErrCantGenerateInTries
	}

	if err != nil {
		return "", err
	}

	return link.Alias, nil
}
