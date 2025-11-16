package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchCmd(t *testing.T) {
	// These are basic validation tests for argument handling
	// Full integration tests would require mocking the registry service
	t.Run("command definition", func(t *testing.T) {
		assert.Equal(t, "search", searchCmd.Use[:6])
		assert.NotEmpty(t, searchCmd.Short)
		assert.NotEmpty(t, searchCmd.Long)
	})

	t.Run("flags exist", func(t *testing.T) {
		assert.NotNil(t, searchCmd.Flags().Lookup("type"))
		assert.NotNil(t, searchCmd.Flags().Lookup("tag"))
		assert.NotNil(t, searchCmd.Flags().Lookup("author"))
		assert.NotNil(t, searchCmd.Flags().Lookup("json"))
		assert.NotNil(t, searchCmd.Flags().Lookup("regex"))
	})
}

func TestDisplayToolsTable(t *testing.T) {
	tests := []struct {
		name     string
		tools    []*models.ToolInfo
		wantText string
	}{
		{
			name:     "empty list",
			tools:    []*models.ToolInfo{},
			wantText: "No tools found",
		},
		{
			name: "single tool",
			tools: []*models.ToolInfo{
				{
					Name:        "test-agent",
					Type:        models.ToolTypeAgent,
					Version:     "1.0.0",
					Author:      "Test Author",
					Description: "A test agent",
					Downloads:   100,
				},
			},
			wantText: "test-agent",
		},
		{
			name: "multiple tools",
			tools: []*models.ToolInfo{
				{
					Name:        "agent1",
					Type:        models.ToolTypeAgent,
					Version:     "1.0.0",
					Author:      "Author 1",
					Description: "First agent",
					Downloads:   100,
				},
				{
					Name:        "command1",
					Type:        models.ToolTypeCommand,
					Version:     "2.0.0",
					Author:      "Author 2",
					Description: "First command",
					Downloads:   200,
				},
			},
			wantText: "Found 2 tool(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := displayToolsTable(tt.tools)
			require.NoError(t, err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			assert.Contains(t, output, tt.wantText)
		})
	}
}

func TestOutputJSON(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "tool info",
			data: &models.ToolInfo{
				Name:        "test-tool",
				Type:        models.ToolTypeAgent,
				Version:     "1.0.0",
				Author:      "Test Author",
				Description: "Test description",
			},
		},
		{
			name: "tool list",
			data: []*models.ToolInfo{
				{
					Name:    "tool1",
					Type:    models.ToolTypeAgent,
					Version: "1.0.0",
				},
				{
					Name:    "tool2",
					Type:    models.ToolTypeCommand,
					Version: "2.0.0",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := outputJSON(tt.data)
			require.NoError(t, err)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			// Verify it's valid JSON
			var result interface{}
			err = json.Unmarshal(buf.Bytes(), &result)
			assert.NoError(t, err, "output should be valid JSON")
			assert.NotEmpty(t, output)
		})
	}
}
