package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {
	//todo pass alp for test
	//option 1 in default case pass const that impl here
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	chars := []rune(
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
			"abcdefghijklmnopqrstuvwxyz" +
			"0123456789",
	)

	// todo it is test value
	chars = []rune("ab")

	res := make([]rune, length)

	for i := range res {
		res[i] = chars[rnd.Intn(len(chars))]
	}

	return string(res)
}
