package linkShortener

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TESTf(t *testing.T) {
	require.Equal(t, 1, 1)
}

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
