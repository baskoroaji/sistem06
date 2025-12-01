package config

import (
	"backend-sistem06.com/utils"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("RT_RW")
	validate.RegisterValidation("postal_code", utils.CustomPostalCodeValidation)
	return validate
}
