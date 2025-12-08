package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func CustomRtRwCodeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	match, _ := regexp.MatchString(`^[0-9]{3}$`, value)
	return match

}
