package errors

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Errors map[string]string

func ValidationError(err error, message map[string]string) Errors {
	if err == nil {
		return nil
	}

	validationErr, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}
	out := Errors{}
	for _, fe := range validationErr {
		template := message[fe.Tag()]
		if template == "" {
			template = "%s is invalid"
		}

		if fe.Param() != "" {
			out[fe.Field()] = fmt.Sprintf(template, fe.Field(), fe.Param())
		} else {
			out[fe.Field()] = fmt.Sprintf(template, fe.Field())
		}
	}

	return nil
}

func FormatValidationErrors(errs map[string]string) string {
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
