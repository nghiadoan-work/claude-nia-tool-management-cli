package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Note: These tests are limited as they require interactive input
// In a real scenario, we'd use dependency injection to mock the reader

func TestConfirmBulkOperation(t *testing.T) {
	// This function requires user input, so we can only test it doesn't panic
	// In a production environment, we would refactor to allow dependency injection
	items := []string{"item1", "item2", "item3"}

	assert.NotPanics(t, func() {
		// Note: This will hang waiting for input if run interactively
		// In CI/CD, stdin would be closed and it would return false
		_ = ConfirmBulkOperation("delete", items)
	}, "ConfirmBulkOperation should not panic")
}

func TestPromptFunctions(t *testing.T) {
	// Test that prompt functions exist and have correct signatures
	// Actual testing would require mocking stdin

	t.Run("function signatures", func(t *testing.T) {
		// These tests just verify the functions compile
		var _ func(string) bool = Confirm
		var _ func(string, bool) bool = ConfirmWithDefault
		var _ func(string) string = Prompt
		var _ func(string, string) string = PromptWithDefault
		var _ func(string, []string) (int, string) = Select
	})
}
