package errors

import (
	"fmt"
	"net/http"
)

// Custom error types
var (
	ErrNoStartingPoint   = NewBusinessError("no valid starting point found")
	ErrCircularRoute     = NewBusinessError("circular route detected")
	ErrDisconnectedRoute = NewBusinessError("disconnected route found")
)

// AppError represents application-specific errors
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewBusinessError creates a new business logic error
func NewBusinessError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Type:    "business_error",
	}
}

// NewValidationError creates a new validation error
func NewValidationError(format string, args ...interface{}) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf(format, args...),
		Type:    "validation_error",
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Type:    "internal_error",
	}
}
