package validate

import (
	"github.com/go-playground/validator/v10"
	"log"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var DBValidator *validator.Validate

func init() {
	DBValidator = validator.New()

	if err := DBValidator.RegisterValidation("uuid", isValidateUUID);        err != nil { log.Fatal(err) } // 문자열 전용
	if err := DBValidator.RegisterValidation("korean", isKoreanString);      err != nil { log.Fatal(err) } // 문자열 전용
	if err := DBValidator.RegisterValidation("phone_number", isPhoneNumber); err != nil { log.Fatal(err) } // 문자열 전용
	if err := DBValidator.RegisterValidation("range", isWithinRange);        err != nil { log.Fatal(err) } // 정수 전용
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

func isPhoneNumber(fl validator.FieldLevel) bool {
	return phoneNumberRegex.MatchString(fl.Field().String())
}

func isWithinRange(fl validator.FieldLevel) bool {
	_range := strings.Split(fl.Param(), "~")
	if len(_range) != 2 {
		log.Fatal("please set param of range like (int)~(int)")
	}

	start, err := strconv.Atoi(_range[0])
	if err != nil {
		log.Fatalf("please set param of range like (int)~(int), err: %v", err)
	}
	end, err := strconv.Atoi(_range[1])
	if err != nil {
		log.Fatalf("please set param of range like (int)~(int), err: %v", err)
	}

	field := int(fl.Field().Int())
	return field >= start && field <= end
}