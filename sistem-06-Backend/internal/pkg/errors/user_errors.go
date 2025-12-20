package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var UserErrorMessages = map[string]string{
	"required": "%s is required",
	"email":    "%s must be a valid email address",
	"min":      "%s must be at least %s characters",
	"max":      "%s must be at most %s characters",
}

func UserValidationError(err error) map[string]string {
	if err == nil {
		return nil
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)
		for _, fe := range validationErrs {
			template := UserErrorMessages[fe.Tag()]
			if template == "" {
				template = "%s is invalid"
			}
			if fe.Param() != "" {
				errors[fe.Field()] = fmt.Sprintf(template, fe.Field(), fe.Param())
			} else {
				errors[fe.Field()] = fmt.Sprintf(template, fe.Field())
			}
		}
		return errors
	}

	return nil
}

func UserFormatValidationErrors(errs map[string]string) string {
	msg := "{"
	for field, err := range errs {
		msg += fmt.Sprintf("\"%s\": \"%s\",", field, err)
	}
	if len(msg) > 1 {
		msg = msg[:len(msg)-1]
	}
	msg += "}"
	return msg
}
