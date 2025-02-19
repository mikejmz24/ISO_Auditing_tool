package custom_errors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode represents a unique identifier for each type of error
type ErrorCode string

// Define error codes as constants for better type safety and performance
const (
	ErrCodeInvalidJSON     ErrorCode = "INVALID_JSON"
	ErrCodeInvalidFormData ErrorCode = "INVALID_FORM_DATA"
	ErrCodeInvalidDataType ErrorCode = "INVALID_DATA_TYPE"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeEmptyField      ErrorCode = "EMPTY_FIELD"
	ErrCodeIsABool         ErrorCode = "IS_A_BOOL"
	ErrCodeFailedToFetch   ErrorCode = "FAILED_TO_FETCH"
	ErrCodeEmptyData       ErrorCode = "EMPTY_DATA"
	ErrCodeInvalidID       ErrorCode = "INVALID_ID"
	ErrCodeMissingField    ErrorCode = "MISSING_FIELD"
	ErrCodeMinChars        ErrorCode = "MIN_CHARACTERS"
	ErrCodeMaxChars        ErrorCode = "MAX_CHARACTERS"
	ErrCodeInvalidData     ErrorCode = "INVALID_DATA"
)

// Predefined errors for common cases
var (
	ErrInvalidJSON     = NewError(context.Background(), ErrCodeInvalidJSON, "Invalid JSON format", http.StatusBadRequest, nil)
	ErrInvalidFormData = NewError(context.Background(), ErrCodeInvalidFormData, "Invalid form data format", http.StatusBadRequest, nil)
	ErrInvalidData     = NewError(context.Background(), ErrCodeInvalidData, "Invalid data", http.StatusBadRequest, nil)
)

// CustomError represents a structured error with context and metadata
type CustomError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	StatusCode int                    `json:"status"`
	Context    map[string]interface{} `json:"-"`
	Err        error                  `json:"-"` // Wrapped error
}

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements the unwrap interface for error chains
func (e *CustomError) Unwrap() error {
	return e.Err
}

// WithContext adds context to an existing error
func (e *CustomError) WithContext(ctx map[string]interface{}) *CustomError {
	e.Context = ctx
	return e
}

// NewError creates a new CustomError with context
func NewError(ctx context.Context, code ErrorCode, message string, statusCode int, err error) *CustomError {
	ctxValues := make(map[string]interface{})
	if requestID := ctx.Value("request_id"); requestID != nil {
		ctxValues["request_id"] = requestID
	}

	return &CustomError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Context:    ctxValues,
		Err:        err,
	}
}

// Original methods updated with context and error codes
func IsABool(ctx context.Context, field string) *CustomError {
	return NewError(ctx, ErrCodeIsABool, fmt.Sprintf("%v should not be a bool", field), http.StatusBadRequest, nil)
}

func FailedToFetch(ctx context.Context, objectType string) *CustomError {
	return NewError(ctx, ErrCodeFailedToFetch, fmt.Sprintf("Failed to fetch %v", objectType), http.StatusInternalServerError, nil)
}

func EmptyData(ctx context.Context, objectType string) *CustomError {
	return NewError(ctx, ErrCodeEmptyData, fmt.Sprintf("Invalid data - %v cannot be empty", objectType), http.StatusBadRequest, nil)
}

func InvalidDataType(ctx context.Context, field, fieldType string) *CustomError {
	return NewError(ctx, ErrCodeInvalidDataType, fmt.Sprintf("Invalid Data - %v must be a %v", field, fieldType), http.StatusBadRequest, nil)
}

func InvalidID(ctx context.Context, objectType string) *CustomError {
	return NewError(ctx, ErrCodeInvalidID, fmt.Sprintf("Invalid %v ID", objectType), http.StatusBadRequest, nil)
}

func NotFound(ctx context.Context, objectType string) *CustomError {
	return NewError(ctx, ErrCodeNotFound, fmt.Sprintf("%v not found", objectType), http.StatusNotFound, nil)
}

func EmptyField(ctx context.Context, typeName, typeField string) *CustomError {
	return NewError(ctx, ErrCodeEmptyField, fmt.Sprintf("%v %v should not be empty", typeName, typeField), http.StatusBadRequest, nil)
}

func MissingField(ctx context.Context, fieldName string) *CustomError {
	return NewError(ctx, ErrCodeMissingField, fmt.Sprintf("Missing required field %v", fieldName), http.StatusBadRequest, nil)
}

func MinFieldCharacters(ctx context.Context, fieldName string, chars int) *CustomError {
	return NewError(ctx, ErrCodeMinChars, fmt.Sprintf("%v has to be at least %v characters", fieldName, chars), http.StatusBadRequest, nil)
}

func MaxFieldCharacters(ctx context.Context, fieldName string, chars int) *CustomError {
	return NewError(ctx, ErrCodeMaxChars, fmt.Sprintf("%v cannot be longer than %v characters", fieldName, chars), http.StatusBadRequest, nil)
}

// For backward compatibility with existing code
type ErrorResponse struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ToResponse converts a CustomError to an API response
func (e *CustomError) ToResponse() ErrorResponse {
	return ErrorResponse{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Context,
	}
}

// Helper function to check if an error is of a specific type
func IsErrorCode(err error, code ErrorCode) bool {
	var customErr *CustomError
	if ok := As(err, &customErr); ok {
		return customErr.Code == code
	}
	return false
}

// As is a helper function that wraps errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
