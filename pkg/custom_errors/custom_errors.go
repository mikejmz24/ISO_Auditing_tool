package custom_errors

import (
	"fmt"
	"net/http"
)

// Custom error definitions
var (
	ErrInvalidJSON     = NewCustomError(http.StatusBadRequest, "Invalid JSON format", nil)
	ErrInvalidFormData = NewCustomError(http.StatusBadRequest, "Invalid Form Data", nil)
)

// CustomError struct
type CustomError struct {
	StatusCode int               `json:"status"`
	Errors     map[string]string `json:"errors"`
}

func (e *CustomError) Error() string {
	// Return the first error message we find
	for _, v := range e.Errors {
		return v
	}

	// Return an empty string or some default message if no message is found
	return ""
}

func IsABool(field string) *CustomError {
	return NewCustomError(http.StatusBadRequest, isABoolMessage(field), nil)
}

func FailedToFetch(objectType string) *CustomError {
	return NewCustomError(http.StatusInternalServerError, failedToFetchMessage(objectType), nil)
}

func EmptyData(objectType string) *CustomError {
	return NewCustomError(http.StatusBadRequest, emptyDataMessage(objectType), nil)
}

func InvalidDataType(field string, fieldType string) *CustomError {
	return NewCustomError(http.StatusBadRequest, invalidDataTypeMessage(field, fieldType), nil)
}

func InvalidID(objectType string) *CustomError {
	return NewCustomError(http.StatusBadRequest, invalidIDMessage(objectType), nil)
}

func NotFound(objectType string) *CustomError {
	return NewCustomError(http.StatusNotFound, notFoundMessage(objectType), nil)
}

func EmptyField(typeName string, typeField string) *CustomError {
	return NewCustomError(http.StatusBadRequest, emptyFieldMessage(typeName, typeField), nil)
}

func MissingField(fieldName string) *CustomError {
	return NewCustomError(http.StatusBadRequest, missingFieldMessage(fieldName), nil)
}

func MinFieldCharacters(fieldName string, chars int) *CustomError {
	return NewCustomError(http.StatusBadRequest, minFieldCharactersMessage(fieldName, chars), nil)
}

func MaxFieldCharacters(fieldName string, chars int) *CustomError {
	return NewCustomError(http.StatusBadRequest, maxFieldCharactersMessage(fieldName, chars), nil)
}

// func NewCustomError(code int, message string, ctx context.Context) *CustomError {
func NewCustomError(statusCode int, message string, context map[string]interface{}) *CustomError {
	return &CustomError{StatusCode: statusCode, Errors: map[string]string{
		"message": message, // Use ca fixed key "validation" instead of the message itself
	}}
}

func isABoolMessage(field string) string {
	return fmt.Sprintf("%v should not be a bool", field)
}

func minFieldCharactersMessage(field string, chars int) string {
	return fmt.Sprintf("%v has to be at least %v characters", field, chars)
}

func maxFieldCharactersMessage(field string, chars int) string {
	return fmt.Sprintf("%v cannot be longer than %v characters", field, chars)
}

func failedToFetchMessage(objectType string) string {
	return fmt.Sprintf("Failed to fetch %v", objectType)
}

func emptyDataMessage(objectType string) string {
	return fmt.Sprintf("Invalid data - %v cannot be empty", objectType)
}

func invalidDataTypeMessage(field string, fieldType string) string {
	return fmt.Sprintf("Invalid Data - %v must be a %v", field, fieldType)
}

func invalidIDMessage(objectType string) string {
	return fmt.Sprintf("Invalid %v ID", objectType)
}

func emptyFieldMessage(typeName string, typeField string) string {
	return fmt.Sprintf("%v %v should not be empty", typeName, typeField)
}

func notFoundMessage(objecType string) string {
	return fmt.Sprintf("%v not found", objecType)
}

func missingFieldMessage(fieldName string) string {
	return fmt.Sprintf("Missing required field %v", fieldName)
}
