package validate

import (
	"github.com/go-playground/validator/v10"
	"log"
)

var DBValidator *validator.Validate

func init() {
	DBValidator = validator.New()

	if err := DBValidator.RegisterValidation("uuid", isValidateUUID); err != nil { log.Fatal(err) }
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