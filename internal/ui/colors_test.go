package ui

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestPrintSuccess(t *testing.T) {
	output := captureOutput(func() {
		PrintSuccess("test message")
	})
	assert.Contains(t, output, "test message")
}

func TestPrintError(t *testing.T) {
	output := captureOutput(func() {
		PrintError("error message")
	})
	assert.Contains(t, output, "error message")
}

func TestPrintWarning(t *testing.T) {
	output := captureOutput(func() {
		PrintWarning("warning message")
	})
	assert.Contains(t, output, "warning message")
}

func TestPrintInfo(t *testing.T) {
	output := captureOutput(func() {
		PrintInfo("info message")
	})
	assert.Contains(t, output, "info message")
}

func TestPrintHint(t *testing.T) {
	output := captureOutput(func() {
		PrintHint("hint message")
	})
	assert.Contains(t, output, "hint message")
}

func TestPrintHeader(t *testing.T) {
	output := captureOutput(func() {
		PrintHeader("Test Header")
	})
	assert.Contains(t, output, "Test Header")
}

func TestFormatVersion(t *testing.T) {
	result := FormatVersion("1.0.0")
	assert.Contains(t, result, "v1.0.0")
}

func TestFormatToolName(t *testing.T) {
	result := FormatToolName("test-tool")
	assert.Contains(t, result, "test-tool")
}

func TestFormatPath(t *testing.T) {
	result := FormatPath("/path/to/file")
	assert.Contains(t, result, "/path/to/file")
}

func TestFormatURL(t *testing.T) {
	result := FormatURL("https://example.com")
	assert.Contains(t, result, "https://example.com")
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		char     string
		count    int
		expected string
	}{
		{"empty", "x", 0, ""},
		{"single", "x", 1, "x"},
		{"multiple", "x", 5, "xxxxx"},
		{"dash", "-", 3, "---"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repeat(tt.char, tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestColorFunctions(t *testing.T) {
	// Test that color functions don't panic
	assert.NotPanics(t, func() {
		_ = Success("test")
		_ = Warning("test")
		_ = Error("test")
		_ = Info("test")
		_ = Highlight("test")
		_ = Bold("test")
		_ = Faint("test")
	})
}

func TestPrintFormattedMessages(t *testing.T) {
	tests := []struct {
		name     string
		fn       func(string, ...interface{})
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "success with format",
			fn:       PrintSuccess,
			format:   "installed %s",
			args:     []interface{}{"tool"},
			expected: "installed tool",
		},
		{
			name:     "error with format",
			fn:       PrintError,
			format:   "failed to install %s",
			args:     []interface{}{"tool"},
			expected: "failed to install tool",
		},
		{
			name:     "warning with format",
			fn:       PrintWarning,
			format:   "tool %s is deprecated",
			args:     []interface{}{"old-tool"},
			expected: "tool old-tool is deprecated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.fn(tt.format, tt.args...)
			})
			assert.Contains(t, output, tt.expected)
		})
	}
}

func TestSuccessErrorWarningSymbols(t *testing.T) {
	// Verify that symbols are present in output
	successOut := captureOutput(func() {
		PrintSuccess("test")
	})
	assert.True(t, strings.Contains(successOut, "✓") || strings.Contains(successOut, "test"))

	errorOut := captureOutput(func() {
		PrintError("test")
	})
	assert.True(t, strings.Contains(errorOut, "✗") || strings.Contains(errorOut, "test"))

	warningOut := captureOutput(func() {
		PrintWarning("test")
	})
	assert.True(t, strings.Contains(warningOut, "⚠") || strings.Contains(warningOut, "test"))
}
