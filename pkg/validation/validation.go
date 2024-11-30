package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidataionError(errs validator.ValidationErrors) string {
	var errsMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errsMsgs = append(errsMsgs, fmt.Sprintf("field %s is a required field", err.Field()))

		default:
			errsMsgs = append(errsMsgs, fmt.Sprintf("field %s is not valid"))
		}
	}

	return strings.Join(errsMsgs, ", ")
}
