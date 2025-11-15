package ui

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		cliErr   *CLIError
		expected string
	}{
		{
			name: "error with wrapped error",
			cliErr: &CLIError{
				Type:    ErrorTypeNotFound,
				Message: "tool not found",
				Err:     errors.New("underlying error"),
			},
			expected: "tool not found: underlying error",
		},
		{
			name: "error without wrapped error",
			cliErr: &CLIError{
				Type:    ErrorTypeValidation,
				Message: "validation failed",
			},
			expected: "validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cliErr.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCLIError_Unwrap(t *testing.T) {
	underlyingErr := errors.New("underlying error")
	cliErr := &CLIError{
		Type:    ErrorTypeNetwork,
		Message: "network error",
		Err:     underlyingErr,
	}

	unwrapped := cliErr.Unwrap()
	assert.Equal(t, underlyingErr, unwrapped)
}

func TestCLIError_Print(t *testing.T) {
	cliErr := &CLIError{
		Type:    ErrorTypeNotFound,
		Message: "tool not found",
		Err:     errors.New("underlying"),
		Hint:    "try searching for the tool",
	}

	// Should not panic
	assert.NotPanics(t, func() {
		cliErr.Print()
	})
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("tool", "run cntm search")
	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeNotFound, err.Type)
	assert.Contains(t, err.Message, "tool not found")
	assert.Equal(t, "run cntm search", err.Hint)
}

func TestNewNetworkError(t *testing.T) {
	underlyingErr := errors.New("connection failed")
	err := NewNetworkError("download", underlyingErr)

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeNetwork, err.Type)
	assert.Contains(t, err.Message, "download")
	assert.Equal(t, underlyingErr, err.Err)
	assert.NotEmpty(t, err.Hint)
}

func TestNewAuthError(t *testing.T) {
	underlyingErr := errors.New("invalid token")
	err := NewAuthError(underlyingErr)

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeAuth, err.Type)
	assert.Contains(t, err.Message, "Authentication")
	assert.Equal(t, underlyingErr, err.Err)
	assert.NotEmpty(t, err.Hint)
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("invalid input", "check the format")

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "invalid input", err.Message)
	assert.Equal(t, "check the format", err.Hint)
}

func TestNewIntegrityError(t *testing.T) {
	err := NewIntegrityError("file.zip")

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeIntegrity, err.Type)
	assert.Contains(t, err.Message, "file.zip")
	assert.NotEmpty(t, err.Hint)
}

func TestNewAlreadyExistsError(t *testing.T) {
	err := NewAlreadyExistsError("tool", "use --force to overwrite")

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypeAlreadyExists, err.Type)
	assert.Contains(t, err.Message, "already exists")
	assert.Equal(t, "use --force to overwrite", err.Hint)
}

func TestNewPermissionError(t *testing.T) {
	err := NewPermissionError("write to", "/protected/path")

	assert.NotNil(t, err)
	assert.Equal(t, ErrorTypePermission, err.Type)
	assert.Contains(t, err.Message, "write to")
	assert.Contains(t, err.Message, "/protected/path")
	assert.NotEmpty(t, err.Hint)
}

func TestHandleError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "no error",
			err:          nil,
			expectedCode: 0,
		},
		{
			name: "cli error",
			err: &CLIError{
				Type:    ErrorTypeNotFound,
				Message: "not found",
			},
			expectedCode: 1,
		},
		{
			name:         "generic error",
			err:          errors.New("generic error"),
			expectedCode: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := HandleError(tt.err)
			assert.Equal(t, tt.expectedCode, code)
		})
	}
}

func TestErrorTypes(t *testing.T) {
	// Verify all error type constants are defined
	types := []ErrorType{
		ErrorTypeNotFound,
		ErrorTypeNetwork,
		ErrorTypeAuth,
		ErrorTypeValidation,
		ErrorTypeIntegrity,
		ErrorTypeAlreadyExists,
		ErrorTypePermission,
	}

	for i, et := range types {
		assert.Equal(t, ErrorType(i), et)
	}
}
