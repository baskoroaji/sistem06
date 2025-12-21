package errors

var Messages = map[string]string{
	"required":    "%s is required",
	"postal_code": "%s must be exact %s characters",
	"RT_RW":       "%s must be exact %s characters",
}

// func ValidationError(err error) map[string]string {
// 	if err == nil {
// 		return nil
// 	}

// 	if validationErrs, ok := err.(validator.ValidationErrors); ok {
// 		errors := make(map[string]string)
// 		for _, fe := range validationErrs {
// 			template := Messages[fe.Tag()]
// 			if template == "" {
// 				template = "%s is invalid"
// 			}
// 			if fe.Param() != "" {
// 				errors[fe.Field()] = fmt.Sprintf(template, fe.Field(), fe.Param())
// 			} else {
// 				errors[fe.Field()] = fmt.Sprintf(template, fe.Field())
// 			}
// 		}
// 		return errors
// 	}

// 	return nil
// }
