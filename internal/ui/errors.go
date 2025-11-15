package ui

import (
	"fmt"
	"os"
)

// ErrorType represents the type of error
type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota
	ErrorTypeNetwork
	ErrorTypeAuth
	ErrorTypeValidation
	ErrorTypeIntegrity
	ErrorTypeAlreadyExists
	ErrorTypePermission
)

// CLIError represents a user-friendly CLI error with hints
type CLIError struct {
	Type    ErrorType
	Message string
	Err     error
	Hint    string
}

// Error implements the error interface
func (e *CLIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements the error unwrapping interface
func (e *CLIError) Unwrap() error {
	return e.Err
}

// Print prints the error with formatting and hints
func (e *CLIError) Print() {
	PrintError("%s", e.Message)

	if e.Err != nil {
		fmt.Fprintf(os.Stderr, "%s %v\n", Faint("Error:"), e.Err)
	}

	if e.Hint != "" {
		PrintHint("%s", e.Hint)
	}
}

// NewNotFoundError creates a new "not found" error
func NewNotFoundError(item string, hint string) *CLIError {
	return &CLIError{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf("%s not found", item),
		Hint:    hint,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(operation string, err error) *CLIError {
	return &CLIError{
		Type:    ErrorTypeNetwork,
		Message: fmt.Sprintf("Network error during %s", operation),
		Err:     err,
		Hint:    "Check your internet connection and try again",
	}
}

// NewAuthError creates a new authentication error
func NewAuthError(err error) *CLIError {
	return &CLIError{
		Type:    ErrorTypeAuth,
		Message: "Authentication failed",
		Err:     err,
		Hint:    "Check your GitHub token in the config file or CNTM_GITHUB_TOKEN environment variable",
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string, hint string) *CLIError {
	return &CLIError{
		Type:    ErrorTypeValidation,
		Message: message,
		Hint:    hint,
	}
}

// NewIntegrityError creates a new integrity check error
func NewIntegrityError(file string) *CLIError {
	return &CLIError{
		Type:    ErrorTypeIntegrity,
		Message: fmt.Sprintf("Integrity check failed for %s", file),
		Hint:    "The downloaded file may be corrupted. Try again or contact the tool author.",
	}
}

// NewAlreadyExistsError creates a new "already exists" error
func NewAlreadyExistsError(item string, hint string) *CLIError {
	return &CLIError{
		Type:    ErrorTypeAlreadyExists,
		Message: fmt.Sprintf("%s already exists", item),
		Hint:    hint,
	}
}

// NewPermissionError creates a new permission error
func NewPermissionError(operation string, path string) *CLIError {
	return &CLIError{
		Type:    ErrorTypePermission,
		Message: fmt.Sprintf("Permission denied: cannot %s %s", operation, path),
		Hint:    "Check file permissions and try running with appropriate privileges",
	}
}

// HandleError handles an error by printing it with appropriate formatting
// Returns an exit code
func HandleError(err error) int {
	if err == nil {
		return 0
	}

	// Check if it's a CLIError
	if cliErr, ok := err.(*CLIError); ok {
		cliErr.Print()
		return 1
	}

	// Generic error
	PrintError("An error occurred")
	fmt.Fprintf(os.Stderr, "%s %v\n", Faint("Error:"), err)
	return 1
}
