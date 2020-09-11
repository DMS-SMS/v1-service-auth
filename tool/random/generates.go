package random

import (
	"math/rand"
	"strconv"
	"time"
)

var (
	intLetters = []rune("0123456789")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func StringConsistOfIntWithLength(length int) string {
	randomRuneArr := make([]rune, length)
	for i := range randomRuneArr {
		randomRuneArr[i] = intLetters[rand.Intn(len(intLetters))]
	}
	return string(randomRuneArr)
}

func Int64WithLength(length int) int64 {
	randomString := StringConsistOfIntWithLength(length)
	stringToInt, _ := strconv.Atoi(randomString)
	return int64(stringToInt)
}