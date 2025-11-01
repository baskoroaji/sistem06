package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) map[string]string {
	if err != nil {
		return nil
	}
	errors := make(map[string]string)
	if ValidationErr, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ValidationErr {
			template := Messages[fe.Tag()]
			if template == "" {
				template = "%s is invalid"
			}
			if fe.Param() != "" {
				errors[fe.Field()] = fmt.Sprintf(template, fe.Field(), fe.Param())
			} else {
				errors[fe.Field()] = fmt.Sprintf(template, fe.Field())
			}
		}
	}
	return errors
}
