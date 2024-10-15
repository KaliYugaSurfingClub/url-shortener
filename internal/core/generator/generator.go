package generator

import "math/rand"

type AliasGenerator struct {
	alp    []rune
	length int
}

func New(alp []rune, length int) *AliasGenerator {
	return &AliasGenerator{alp: alp, length: length}
}

func (g *AliasGenerator) Generate() string {
	res := make([]rune, g.length)

	for i := 0; i < g.length; i++ {
		res[i] = g.alp[rand.Intn(len(g.alp))]
	}

	return string(res)
}
