package validate

import (
	"github.com/go-playground/validator/v10"
)

var DBValidator *validator.Validate

func init() {
	DBValidator = validator.New()
}
