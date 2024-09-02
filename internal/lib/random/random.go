package random

import (
	"math/rand"
	"slices"
	"time"
)

func SliceAlp(start rune, end rune) []rune {
	res := make([]rune, start-end+1)

	for it := start; it <= end; it++ {
		res = append(res, it)
	}

	return res
}

func AlphaNumAlp() []rune {
	return slices.Concat(
		SliceAlp('0', '9'),
		SliceAlp('a', 'z'),
		SliceAlp('A', 'Z'),
	)
}

func NewRandomString(length int, chars []rune) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	res := make([]rune, length)

	for i := range res {
		res[i] = chars[rnd.Intn(len(chars))]
	}

	return string(res)
}
