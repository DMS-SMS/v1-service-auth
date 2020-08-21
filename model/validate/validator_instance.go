package validate

import (
	"github.com/go-playground/validator/v10"
	"log"
)

var DBValidator *validator.Validate

func init() {
	DBValidator = validator.New()

	if err := DBValidator.RegisterValidation("uuid", uuidValidateFunc); err != nil { log.Fatal(err) }
}

func uuidValidateFunc(fl validator.FieldLevel) bool {
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