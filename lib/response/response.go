package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func ValidationError(errs validator.ValidationErrors) []string {
	var errMessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("%s is a required field", err.Field()))
		case "gt":
			errMessages = append(errMessages, fmt.Sprintf("field %s must be greater than %v", err.Field(), err.Value()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}
	return errMessages
}
