package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func CustomRtRwCodeValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	match, _ := regexp.MatchString(`^[0-9]{3}$`, value)
	return match

}

func CustomPostalCodeValidation(fl validator.FieldLevel) bool {
	postal := fl.Field().String()
	re_pcode := regexp.MustCompile(`^[1-9][0-9]{4}$`)
	return re_pcode.MatchString(postal)
}
