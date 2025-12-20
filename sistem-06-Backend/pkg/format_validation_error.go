package pkg

import "fmt"

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
