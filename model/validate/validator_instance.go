package validate

import (
	"github.com/go-playground/validator/v10"
	"log"
	"unicode"
	"unicode/utf8"
)

var DBValidator *validator.Validate

func init() {
	DBValidator = validator.New()

	if err := DBValidator.RegisterValidation("uuid", isValidateUUID); err != nil { log.Fatal(err) }
	if err := DBValidator.RegisterValidation("korean", isKoreanString); err != nil { log.Fatal(err) }
}

func isValidateUUID(fl validator.FieldLevel) bool {
	switch fl.Param() {
	case "student":
		return studentUUIDRegex.MatchString(fl.Field().String())
	case "teacher":
		return teacherUUIDRegex.MatchString(fl.Field().String())
	case "parent":
		return parentUUIDRegex.MatchString(fl.Field().String())
	}
	return false
}

func isKoreanString(fl validator.FieldLevel) bool {
	b := []byte(fl.Field().String())
	var idx int

	for {
		r, size := utf8.DecodeRune(b[idx:])
		if size == 0 { break }
		if !unicode.Is(unicode.Hangul, r) { return false }
		idx += size
	}
	return true
}