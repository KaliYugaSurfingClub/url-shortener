package AliasManager

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"url_shortener/core"
	"url_shortener/core/model"
)

var someErr = errors.New("some error")

type saver struct {
	data      map[string]struct{}
	beforeErr int
}

func newSaver(exists ...string) *saver {
	mp := make(map[string]struct{}, len(exists))

	for _, alias := range exists {
		mp[alias] = struct{}{}
	}

	return &saver{data: mp, beforeErr: -1}
}

func (s *saver) WithErrorAfter(n int) *saver {
	s.beforeErr = n
	return s
}

func (s *saver) Save(_ context.Context, link model.Link) (int64, error) {
	if s.beforeErr == 0 {
		return 0, someErr
	}

	if _, ok := s.data[link.Alias]; ok {
		return 0, core.ErrAliasExists
	}

	s.data[link.Alias] = struct{}{}

	s.beforeErr--

	return 0, nil
}

type serialGenerator struct {
	aliases []string
}

func newSerialGenerator(aliases ...string) *serialGenerator {
	return &serialGenerator{aliases}
}

func (g *serialGenerator) Generate() string {
	alias := g.aliases[len(g.aliases)-1]
	g.aliases = g.aliases[:len(g.aliases)-1]

	return alias
}

func TestWithAlias(t *testing.T) {
	manager, _ := New(newSaver(), newSerialGenerator(), 1)

	toSave := "abc"

	alias, err := manager.Save(context.Background(), model.Link{Alias: toSave})
	require.NoError(t, err)
	assert.Equal(t, toSave, alias)
}

func TestFirstNoAlias(t *testing.T) {
	willGenerated := "abc"

	manager, _ := New(newSaver(), newSerialGenerator(willGenerated), 3)

	alias, err := manager.Save(context.Background(), model.Link{})
	require.NoError(t, err)
	assert.Equal(t, willGenerated, alias)
}

func TestGenerateAlreadyExistsAliasSuccess(t *testing.T) {
	alreadyExistAlias := "abc"
	willGenerated := "123"

	manager, _ := New(
		newSaver(alreadyExistAlias),
		newSerialGenerator(alreadyExistAlias, alreadyExistAlias, willGenerated),
		3,
	)

	alias, err := manager.Save(context.Background(), model.Link{})

	require.NoError(t, err)
	assert.Equal(t, willGenerated, alias)
}

func TestCannotGenerateAliasInTries(t *testing.T) {
	alreadyExistAlias := "abc"

	manager, _ := New(
		newSaver(alreadyExistAlias),
		newSerialGenerator(alreadyExistAlias, alreadyExistAlias, alreadyExistAlias, alreadyExistAlias),
		3,
	)

	alias, err := manager.Save(context.Background(), model.Link{})
	assert.Equal(t, true, errors.Is(err, core.ErrCantGenerateInTries), err)
	assert.Equal(t, "", alias)
}

func TestErrorWhileGeneratingSuccess(t *testing.T) {
	alreadyExistAlias := "abc"
	willGenerated := "123"

	manager, _ := New(
		newSaver(alreadyExistAlias).WithErrorAfter(1),
		newSerialGenerator(alreadyExistAlias, willGenerated, willGenerated),
		3,
	)

	alias, err := manager.Save(context.Background(), model.Link{})
	require.NoError(t, err)
	assert.Equal(t, willGenerated, alias)
}

func TestErrorWhileGeneratingFail(t *testing.T) {
	alreadyExistAlias := "abc"
	willGenerated := "123"

	manager, _ := New(
		newSaver(alreadyExistAlias).WithErrorAfter(1),
		newSerialGenerator(alreadyExistAlias, willGenerated, willGenerated),
		3,
	)

	alias, err := manager.Save(context.Background(), model.Link{})
	assert.Equal(t, true, errors.Is(err, core.ErrCantGenerateInTries), err)
	assert.Equal(t, willGenerated, alias)
}
