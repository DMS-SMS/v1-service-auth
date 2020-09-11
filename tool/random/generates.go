package random

import "math/rand"

var (
	intLetters = []rune("0123456789")
)

func StringConsistOfIntWithLength(length int) string {
	s := make([]rune, length)
	for i := range s {
		s[i] = intLetters[rand.Intn(len(intLetters))]
	}
	return string(s)
}
