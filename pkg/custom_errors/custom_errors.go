package custom_errors

import (
	"fmt"
	"net/http"
)

// Custom error definitions
var (
	ErrFailedToFetchISOStandards = NewCustomError(http.StatusInternalServerError, "Failed to fetch ISO Standards", nil)
	ErrInvalidID                 = NewCustomError(http.StatusBadRequest, "Invalid ISO Standard ID", nil)
	ErrInvalidJSON               = NewCustomError(http.StatusBadRequest, "Invalid JSON format", nil)
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

func InvalidDataType(field string, fieldType string) *CustomError {
	return NewCustomError(http.StatusBadRequest, invalidDataTypeMessage(field, fieldType), nil)
}

func EmptyField(typeName string, typeField string) *CustomError {
	return NewCustomError(http.StatusBadRequest, emptyFieldMessage(typeName, typeField), nil)
}

// func NewCustomError(code int, message string, ctx context.Context) *CustomError {
func NewCustomError(statusCode int, message string, context map[string]interface{}) *CustomError {
	return &CustomError{StatusCode: statusCode, Message: message, Context: context}
}

func invalidDataTypeMessage(field string, fieldType string) string {
	return fmt.Sprintf("Invalid Data - %v must be a %v", field, fieldType)
}

func emptyFieldMessage(typeName string, typeField string) string {
	return fmt.Sprintf("%v %v should not be empty", typeName, typeField)
}
