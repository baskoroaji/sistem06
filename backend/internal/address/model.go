package address

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type AddressEntity struct {
	ID         int    `json:"id"`
	Jalan      string `json:"jalan"`
	RT         string `json:"RT"`
	RW         string `json:"RW"`
	Kota       string `json:"Kota"`
	PostalCode string `json:"PostalCode"`
}

type AddressRequest struct {
	Jalan      string `json:"jalan" validate:"required"`
	RT         string `json:"RT" validate:"required,RT_RW"`
	RW         string `json:"RW" validate:"required,RT_RW"`
	Kota       string `json:"Kota" validate:"required"`
	PostalCode string `json:"PostalCode" validate:"required,postal_code"`
}

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
