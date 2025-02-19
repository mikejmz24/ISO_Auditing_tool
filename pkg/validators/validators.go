package validators

import (
	"ISO_Auditing_Tool/pkg/custom_errors"
	"context"
	"fmt"
	"sync"

	"github.com/go-playground/validator/v10"
)

// Global Validator Instance
var (
	validate *validator.Validate
	once     sync.Once
)

func InitValidator() *validator.Validate {
	once.Do(func() {
		validate = validator.New()
		// Register custom validation functions
		if err := validate.RegisterValidation("not_boolean", validateNotBool); err != nil {
			panic(fmt.Sprintf("Failed to register validator: %v", err))
		}
	})
	return validate
}

// GetValidator returns the initialized validator instance
func GetValidator() *validator.Validate {
	if validate == nil {
		InitValidator()
	}
	return validate
}

// ValidateStruct validates a struct and returns a structured custom error
func ValidateStruct(data interface{}) *custom_errors.CustomError {
	validate := GetValidator()
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	// // Empty string
	// if data.Name == "" {
	// 	return custom_errors.EmptyField("string", "name")
	// }
	//
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		// Process validation errors
		// if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			fieldName := e.Field()

			switch e.Tag() {
			case "required":
				return custom_errors.EmptyField(context.TODO(), "string", fieldName)
			case "min":
				return custom_errors.MinFieldCharacters(context.TODO(), fieldName, 2)
			case "max":
				return custom_errors.MaxFieldCharacters(context.TODO(), fieldName, 100)
			case "not_boolean":
				return custom_errors.IsABool(context.TODO(), fieldName)
			}
		}
	}
	// return custom_errors.NewError(context.TODO(), custom_errors.ErrInvalidData, "Unexpected validation error", http.StatusInternalServerError, err)
	return custom_errors.ErrInvalidFormData
}
