package errors

import (
	"context"
	"net/http"
)

// Custom error definitions
var (
	// ErrFailedToFetchISOStandards = &CustomError{
	// 	StatusCode: http.StatusInternalServerError,
	// 	Message:    "database error",
	// 	Context:    nil,
	// }
	ErrFailedToFetchISOStandards = NewCustomError(http.StatusInternalServerError, "Failed to fetch ISO Standards", nil)
)

// CustomError struct
type CustomError struct {
	StatusCode int             `json:"status"`
	Message    string          `json:"message"`
	Context    context.Context `json:"-"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(code int, message string, ctx context.Context) *CustomError {
	return &CustomError{StatusCode: code, Message: message, Context: ctx}
}
