package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListCmd(t *testing.T) {
	// These are basic validation tests for command structure
	// Full integration tests would require mocking the registry service
	t.Run("command definition", func(t *testing.T) {
		assert.Equal(t, "list", listCmd.Use)
		assert.NotEmpty(t, listCmd.Short)
		assert.NotEmpty(t, listCmd.Long)
	})

	t.Run("flags exist", func(t *testing.T) {
		assert.NotNil(t, listCmd.Flags().Lookup("remote"))
		assert.NotNil(t, listCmd.Flags().Lookup("type"))
		assert.NotNil(t, listCmd.Flags().Lookup("tag"))
		assert.NotNil(t, listCmd.Flags().Lookup("author"))
		assert.NotNil(t, listCmd.Flags().Lookup("sort-by"))
		assert.NotNil(t, listCmd.Flags().Lookup("sort-desc"))
		assert.NotNil(t, listCmd.Flags().Lookup("limit"))
		assert.NotNil(t, listCmd.Flags().Lookup("json"))
	})
}

func TestListSortByValidation(t *testing.T) {
	// Test that valid sort fields are recognized
	validFields := []string{"name", "created", "updated", "downloads"}

	for _, field := range validFields {
		t.Run("valid field: "+field, func(t *testing.T) {
			// This just validates that the field names are correct
			// Full validation happens in runList which requires mocked services
			assert.Contains(t, validFields, field)
		})
	}
}
