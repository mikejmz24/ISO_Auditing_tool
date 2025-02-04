package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Custom validation logic
func validateNotBool(fl validator.FieldLevel) bool {
	value := strings.ToLower(strings.TrimSpace(fl.Field().String()))

	// List of possible boolean-like values to reject
	booleanValues := map[string]bool{
		"true":     true,
		"false":    true,
		"yes":      true,
		"off":      true,
		"enabled":  true,
		"disabled": true,
	}

	return !booleanValues[value]
}
