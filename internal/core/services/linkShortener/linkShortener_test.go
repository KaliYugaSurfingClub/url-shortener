package linkShortener

//
//import (
//	"context"
//	"errors"
//	"github.com/stretchr/testify/require"
//	"shortener/internal/core"
//	"shortener/internal/core/model"
//	"slices"
//	"testing"
//)
//
//var someErr = errors.New("some error")
//
//type linkStorage struct {
//	data []*model.Link
//}
//
//func newStorage() *linkStorage {
//	return &linkStorage{data: make([]*model.Link, 0)}
//}
//
//func (s *linkStorage) Save(_ context.Context, toSave model.Link) (saved *model.Link, err error) {
//	aliasExists := slices.ContainsFunc(s.data, func(link *model.Link) bool {
//		return link.Alias == toSave.Alias
//	})
//
//	if aliasExists {
//		return nil, core.ErrAliasExists
//	}
//
//	nameExists := slices.ContainsFunc(s.data, func(link *model.Link) bool {
//		return link.CustomName == toSave.CustomName && link.CreatedBy == toSave.CreatedBy
//	})
//
//	if nameExists {
//		return nil, core.ErrCustomNameExists
//	}
//
//	s.data = append(s.data, &toSave)
//	return &toSave, nil
//}
//
//type serialGenerator struct {
//	aliases []string
//}
//
//func newSerialGenerator(aliases ...string) *serialGenerator {
//	return &serialGenerator{aliases}
//}
//
//func (g *serialGenerator) Generate() string {
//	alias := g.aliases[0]
//	g.aliases = g.aliases[1:]
//
//	return alias
//}
//
//func TestWithAliasAndName(t *testing.T) {
//	generator := newSerialGenerator()
//	storage := newStorage()
//	shortener, _ := New(storage, generator, 1)
//
//	link := model.Link{
//		Alias:      "abc",
//		CustomName: "123",
//	}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.NoError(t, err)
//	require.Equal(t, link.Alias, saved.Alias)
//	require.Equal(t, link.CustomName, saved.CustomName)
//}
//
//func TestWithAliasWithoutName(t *testing.T) {
//	generator := newSerialGenerator()
//	storage := newStorage()
//	shortener, _ := New(storage, generator, 1)
//
//	link := model.Link{
//		Alias: "abc",
//	}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.NoError(t, err)
//	require.Equal(t, link.Alias, saved.Alias)
//	require.Equal(t, link.Alias, saved.CustomName)
//}
//
//func TestWithoutAliasWithName(t *testing.T) {
//	generatedAlias := "abc"
//
//	generator := newSerialGenerator(generatedAlias)
//	storage := newStorage()
//	shortener, _ := New(storage, generator, 1)
//
//	link := model.Link{
//		CustomName: "123",
//	}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.NoError(t, err)
//	require.Equal(t, generatedAlias, saved.Alias)
//	require.Equal(t, link.CustomName, saved.CustomName)
//}
//
//func TestWithoutAliasWithoutName(t *testing.T) {
//	generatedAlias := "abc"
//
//	generator := newSerialGenerator(generatedAlias)
//	storage := newStorage()
//	shortener, _ := New(storage, generator, 1)
//
//	link := model.Link{}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.NoError(t, err)
//	require.Equal(t, generatedAlias, saved.Alias)
//	require.Equal(t, generatedAlias, saved.CustomName)
//}
//
//func TestWithoutAliasWithoutNameManyTries(t *testing.T) {
//	generatedAlias := "abc"
//	existingAlias := "alias"
//	existsLink := model.Link{Alias: existingAlias}
//
//	generator := newSerialGenerator(existingAlias, existingAlias, existingAlias, generatedAlias)
//	storage := newStorage()
//	storage.Save(context.Background(), existsLink)
//
//	shortener, _ := New(storage, generator, 4)
//
//	link := model.Link{}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.NoError(t, err)
//	require.Equal(t, generatedAlias, saved.Alias)
//	require.Equal(t, generatedAlias, saved.CustomName)
//}
//
//func TestGenerateExistingNameButNotExistingAlias(t *testing.T) {
//	existingName := "123"
//	generatedAlias := "abc"
//
//	generator := newSerialGenerator(existingName, generatedAlias)
//	storage := newStorage()
//	storage.Save(context.Background(), model.Link{CustomName: existingName})
//
//	shortener, _ := New(storage, generator, 2)
//
//	link := model.Link{}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.NoError(t, err)
//	require.Equal(t, generatedAlias, saved.Alias)
//	require.Equal(t, generatedAlias, saved.CustomName)
//}
//
//func TestWithoutAliasWithExistingNameFail(t *testing.T) {
//	existingName := "123"
//
//	generator := newSerialGenerator("abc", "abc")
//	storage := newStorage()
//	storage.Save(context.Background(), model.Link{CustomName: existingName})
//
//	shortener, _ := New(storage, generator, 2)
//
//	link := model.Link{CustomName: existingName}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.Equal(t, true, errors.Is(err, core.ErrCustomNameExists), err)
//	require.Equal(t, true, saved == nil)
//}
//
//func TestGenerateNameAsAnotherUser(t *testing.T) {
//	existingNameAsAnotherUser := "123"
//
//	generator := newSerialGenerator(existingNameAsAnotherUser)
//	storage := newStorage()
//	storage.Save(context.Background(), model.Link{CreatedBy: 1, CustomName: existingNameAsAnotherUser})
//
//	shortener, _ := New(storage, generator, 1)
//
//	link := model.Link{}
//
//	saved, err := shortener.Short(context.Background(), link)
//	require.NoError(t, err)
//	require.Equal(t, existingNameAsAnotherUser, saved.Alias)
//	require.Equal(t, existingNameAsAnotherUser, saved.CustomName)
//}
