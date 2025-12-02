package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func CustomPostalCodeValidation(fl validator.FieldLevel) bool {
	postal := fl.Field().String()
	re_pcode := regexp.MustCompile(`^[1-9][0-9]{4}$`)
	return re_pcode.MatchString(postal)
}
