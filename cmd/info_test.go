package cmd

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfoCmd(t *testing.T) {
	// These are basic validation tests for command structure
	// Full integration tests would require mocking the registry service
	t.Run("command definition", func(t *testing.T) {
		assert.Equal(t, "info", infoCmd.Use[:4])
		assert.NotEmpty(t, infoCmd.Short)
		assert.NotEmpty(t, infoCmd.Long)
	})

	t.Run("flags exist", func(t *testing.T) {
		assert.NotNil(t, infoCmd.Flags().Lookup("type"))
		assert.NotNil(t, infoCmd.Flags().Lookup("json"))
	})
}

func TestDisplayToolInfo(t *testing.T) {
	tool := &models.ToolInfo{
		Name:        "test-agent",
		Type:        models.ToolTypeAgent,
		Version:     "1.0.0",
		Author:      "Test Author",
		Description: "A comprehensive test agent for testing purposes",
		Tags:        []string{"test", "agent", "demo"},
		Downloads:   1234,
		Size:        1024 * 1024, // 1 MB
		File:        "agents/test-agent.zip",
		CreatedAt:   time.Now().Add(-24 * time.Hour),
		UpdatedAt:   time.Now(),
	}

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := displayToolInfo(tool)
	require.NoError(t, err)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify all fields are displayed
	assert.Contains(t, output, "test-agent")
	assert.Contains(t, output, "1.0.0")
	assert.Contains(t, output, "agent")
	assert.Contains(t, output, "Test Author")
	assert.Contains(t, output, "A comprehensive test agent")
	assert.Contains(t, output, "test, agent, demo")
	assert.Contains(t, output, "1234")
	assert.Contains(t, output, "1.0 MB")
	assert.Contains(t, output, "agents/test-agent.zip")
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "bytes",
			bytes:    500,
			expected: "500 B",
		},
		{
			name:     "kilobytes",
			bytes:    1024,
			expected: "1.0 KB",
		},
		{
			name:     "megabytes",
			bytes:    1024 * 1024,
			expected: "1.0 MB",
		},
		{
			name:     "gigabytes",
			bytes:    1024 * 1024 * 1024,
			expected: "1.0 GB",
		},
		{
			name:     "fractional KB",
			bytes:    1536,
			expected: "1.5 KB",
		},
		{
			name:     "fractional MB",
			bytes:    1024*1024 + 512*1024,
			expected: "1.5 MB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}
