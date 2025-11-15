package ui

import (
	"fmt"

	"github.com/fatih/color"
)

// Color output utilities for consistent UX across all commands

var (
	// Success - Green color for success messages
	Success = color.New(color.FgGreen, color.Bold).SprintFunc()

	// Warning - Yellow color for warnings and prompts
	Warning = color.New(color.FgYellow, color.Bold).SprintFunc()

	// Error - Red color for errors and failures
	Error = color.New(color.FgRed, color.Bold).SprintFunc()

	// Info - Blue color for informational messages
	Info = color.New(color.FgBlue, color.Bold).SprintFunc()

	// Highlight - Cyan color for links, paths, and emphasized text
	Highlight = color.New(color.FgCyan).SprintFunc()

	// Bold - Bold text without color
	Bold = color.New(color.Bold).SprintFunc()

	// Faint - Dim text for secondary information
	Faint = color.New(color.Faint).SprintFunc()
)

// PrintSuccess prints a success message with a checkmark
func PrintSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Success("✓"), msg)
}

// PrintError prints an error message with an X mark
func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Error("✗"), msg)
}

// PrintWarning prints a warning message with a warning symbol
func PrintWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Warning("⚠"), msg)
}

// PrintInfo prints an informational message with an info symbol
func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Info("ℹ"), msg)
}

// PrintHint prints a helpful hint for the user
func PrintHint(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", Faint("💡 Hint:"), Faint(msg))
}

// PrintHeader prints a section header
func PrintHeader(text string) {
	fmt.Printf("\n%s\n", Info(text))
	fmt.Println(Faint(repeat("─", len(text))))
}

// repeat returns a string with the character repeated n times
func repeat(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}
	return result
}

// FormatVersion formats a version string with highlighting
func FormatVersion(version string) string {
	return Highlight("v" + version)
}

// FormatToolName formats a tool name with highlighting
func FormatToolName(name string) string {
	return Highlight(name)
}

// FormatPath formats a file path with highlighting
func FormatPath(path string) string {
	return Highlight(path)
}

// FormatURL formats a URL with highlighting
func FormatURL(url string) string {
	return Highlight(url)
}
