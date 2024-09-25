package aliasManager

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"shortener/internal/core"
	"shortener/internal/core/model"
	"testing"
	"time"
)

var someErr = errors.New("some error")

type saver struct {
	data         map[string]struct{}
	saveNumber   int
	errorNumbers map[int]struct{}
}

func (s *saver) GetActiveByAlias(ctx context.Context, alias string) (*model.Link, error) {
	panic("implement me")
}

func (s *saver) AliasExists(ctx context.Context, alias string) (bool, error) {
	return false, nil
}

func (s *saver) CustomNameExists(ctx context.Context, customName string, userId int64) (bool, error) {
	return false, nil
}

func (s *saver) UpdateLastAccess(ctx context.Context, id int64, timestamp time.Time) error {
	panic("implement me")
}

func newSaver(exists ...string) *saver {
	mp := make(map[string]struct{})

	for _, alias := range exists {
		mp[alias] = struct{}{}
	}

	return &saver{data: mp}
}

func (s *saver) withErrors(n ...int) *saver {
	s.errorNumbers = make(map[int]struct{})

	for _, i := range n {
		s.errorNumbers[i] = struct{}{}
	}

	return s
}

func (s *saver) Save(_ context.Context, link *model.Link) (int64, error) {
	if _, ok := s.errorNumbers[s.saveNumber]; ok {
		s.saveNumber++
		return 0, someErr
	}

	s.saveNumber++

	if _, ok := s.data[link.Alias]; ok {
		return 0, core.ErrAliasExists
	}

	s.data[link.Alias] = struct{}{}

	return 0, nil
}

type serialGenerator struct {
	aliases []string
}

func newSerialGenerator(aliases ...string) *serialGenerator {
	return &serialGenerator{aliases}
}

func (g *serialGenerator) Generate() string {
	alias := g.aliases[0]
	g.aliases = g.aliases[1:]

	return alias
}

func TestWithAlias(t *testing.T) {
	manager, _ := New(newSaver(), newSerialGenerator(), 1)

	toSave := "abc"

	alias, err := manager.Save(context.Background(), &model.Link{Alias: toSave})
	require.NoError(t, err)
	assert.Equal(t, toSave, alias)
}

func TestFirstNoAlias(t *testing.T) {
	willGenerated := "abc"

	manager, _ := New(newSaver(), newSerialGenerator(willGenerated), 3)

	alias, err := manager.Save(context.Background(), &model.Link{})
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

	alias, err := manager.Save(context.Background(), &model.Link{})

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

	alias, err := manager.Save(context.Background(), &model.Link{})
	assert.Equal(t, true, errors.Is(err, core.ErrCantGenerateInTries) && errors.Is(err, core.ErrAliasExists), err)
	assert.Equal(t, "", alias)
}

func TestErrorWhileGeneratingSuccess(t *testing.T) {
	alreadyExistAlias := "abc"
	willGenerated := "123"

	manager, _ := New(
		newSaver(alreadyExistAlias).withErrors(1),
		newSerialGenerator(alreadyExistAlias, willGenerated, willGenerated),
		3,
	)

	alias, err := manager.Save(context.Background(), &model.Link{})
	require.NoError(t, err)
	assert.Equal(t, willGenerated, alias)
}

func TestErrorWhileGeneratingFail(t *testing.T) {
	alreadyExistAlias := "abc"
	willGenerated := "123"

	manager, _ := New(
		newSaver(alreadyExistAlias).withErrors(1, 2),
		newSerialGenerator(alreadyExistAlias, willGenerated, willGenerated),
		3,
	)

	alias, err := manager.Save(context.Background(), &model.Link{})
	assert.Equal(t, true,
		errors.Is(err, core.ErrCantGenerateInTries) &&
			errors.Is(err, someErr) &&
			errors.Is(err, core.ErrAliasExists),
		err,
	)
	assert.Equal(t, "", alias)
}
