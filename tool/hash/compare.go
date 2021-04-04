package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"regexp"
	"strings"
)

var (
	pbkdf2Regex = regexp.MustCompile("^pbkdf2:sha\\d+(:\\d+)?\\$.*\\$.*$")
	keyLength   = 32
	iterations  = []int{150000, 50000}
	ErrMismatchedHashAndPassword = bcrypt.ErrMismatchedHashAndPassword
)

func CompareHashAndPassword(hashedPW, pw string) (err error) {
	switch true {
	case pbkdf2Regex.MatchString(hashedPW):
		if !checkPbkdf2PasswordHash(hashedPW, pw) {
			err = ErrMismatchedHashAndPassword
		}
	default:
		if err = bcrypt.CompareHashAndPassword([]byte(hashedPW), []byte(pw)); err == bcrypt.ErrMismatchedHashAndPassword {
			err = ErrMismatchedHashAndPassword
		}
	}
	return
}

func checkPbkdf2PasswordHash(hash, password string) bool {
	if strings.Count(hash, "$") < 2 {
		return false
	}
	sep := strings.Split(hash, "$")
	for _, iteration := range iterations {
		if sep[2] == hashInternal(sep[1], password, iteration, keyLength) {
			return true
		}
	}
	return false
}

func hashInternal(salt, password string, iterations, keyLen int) string {
	hash := pbkdf2.Key([]byte(password), []byte(salt), iterations, keyLen, sha256.New)
	return hex.EncodeToString(hash)
}
