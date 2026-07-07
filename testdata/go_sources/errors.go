package server

import "fmt"

// AppError is an application error
type AppError struct {
	Code    int
	Message string
}

// Error implements the error interface
func (e *AppError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

// NewAppError creates a new AppError
func NewAppError(code int, msg string) *AppError {
	return &AppError{Code: code, Message: msg}
}

// Wrap wraps an error with context - dead code
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

// Is checks if an error is an AppError - dead code
func Is(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// ErrNotFound is a pre-defined not-found error
var ErrNotFound = NewAppError(404, "not found")

// ErrInternal is a pre-defined internal error - dead code (variable never used)
var ErrInternal = NewAppError(500, "internal error")
