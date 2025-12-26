package errors

import (
	"fmt"
	"runtime"
)

// Error types
type ConfigurationError struct {
	Message string
	Cause   error
}

func (e ConfigurationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("configuration error: %s (caused by: %v)", e.Message, e.Cause)
	}
	return fmt.Sprintf("configuration error: %s", e.Message)
}

func (e ConfigurationError) Unwrap() error {
	return e.Cause
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error for field '%s': %s", e.Field, e.Message)
}

type DatabaseError struct {
	Operation string
	Message   string
	Cause     error
}

func (e DatabaseError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("database error during %s: %s (caused by: %v)", e.Operation, e.Message, e.Cause)
	}
	return fmt.Sprintf("database error during %s: %s", e.Operation, e.Message)
}

func (e DatabaseError) Unwrap() error {
	return e.Cause
}

type HTTPError struct {
	StatusCode int
	Message    string
	Cause      error
}

func (e HTTPError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("HTTP error %d: %s (caused by: %v)", e.StatusCode, e.Message, e.Cause)
	}
	return fmt.Sprintf("HTTP error %d: %s", e.StatusCode, e.Message)
}

func (e HTTPError) Unwrap() error {
	return e.Cause
}

// Error wrapping utilities
func WrapConfigurationError(err error, message string) error {
	return ConfigurationError{
		Message: message,
		Cause:   err,
	}
}

func WrapDatabaseError(err error, operation, message string) error {
	return DatabaseError{
		Operation: operation,
		Message:   message,
		Cause:     err,
	}
}

func WrapHTTPError(err error, statusCode int, message string) error {
	return HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Cause:      err,
	}
}

// Stack trace utilities
type StackTraceError struct {
	Message string
	Cause   error
	Stack   []uintptr
}

func (e StackTraceError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s (caused by: %v)", e.Message, e.Cause)
	}
	return e.Message
}

func (e StackTraceError) Unwrap() error {
	return e.Cause
}

func WithStackTrace(err error, message string) error {
	var stack [32]uintptr
	n := runtime.Callers(3, stack[:])

	return StackTraceError{
		Message: message,
		Cause:   err,
		Stack:   stack[:n],
	}
}
