package custom_errors

import (
	"net/http"
)

// Custom error definitions
var (
	ErrFailedToFetchISOStandards = NewCustomError(http.StatusInternalServerError, "Failed to fetch ISO Standards", nil)
	ErrInvalidID                 = NewCustomError(http.StatusBadRequest, "Invalid ISO Standard ID", nil)
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

// func NewCustomError(code int, message string, ctx context.Context) *CustomError {
func NewCustomError(statusCode int, message string, context map[string]interface{}) *CustomError {
	return &CustomError{StatusCode: statusCode, Message: message, Context: context}
}
