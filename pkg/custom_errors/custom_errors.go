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
	StatusCode int                    `json:"status"`
	Message    string                 `json:"message"`
	Context    map[string]interface{} `json:"-"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func FailedToFetch(objectType string) *CustomError {
	return NewCustomError(http.StatusInternalServerError, failedToFetchMessage(objectType), nil)
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

// func NewCustomError(code int, message string, ctx context.Context) *CustomError {
func NewCustomError(statusCode int, message string, context map[string]interface{}) *CustomError {
	return &CustomError{StatusCode: statusCode, Message: message, Context: context}
}

func failedToFetchMessage(objectType string) string {
	return fmt.Sprintf("Failed to fetch %v", objectType)
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
