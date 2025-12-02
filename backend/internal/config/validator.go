package config

import (
	"backend-sistem06.com/internal/validation"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("RT_RW", validation.CustomRtRwCodeValidation)
	validate.RegisterValidation("postal_code", validation.CustomPostalCodeValidation)
	return validate
}
