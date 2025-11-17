package ui

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Confirm prompts the user for yes/no confirmation
// Returns true if user confirms, false otherwise
// Supports ESC to cancel (returns false)
func Confirm(message string) bool {
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
	}

	result, err := prompt.Run()
	if err != nil {
		// User pressed ESC or Ctrl+C
		return false
	}

	response := strings.TrimSpace(strings.ToLower(result))
	return response == "y" || response == "yes"
}

// ConfirmWithDefault prompts the user for yes/no confirmation with a default value
// Supports ESC to cancel (returns the default value)
func ConfirmWithDefault(message string, defaultYes bool) bool {
	defaultStr := "N"
	if defaultYes {
		defaultStr = "Y"
	}

	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
		Default:   defaultStr,
	}

	result, err := prompt.Run()
	if err != nil {
		// User pressed ESC or Ctrl+C, return default
		return defaultYes
	}

	response := strings.TrimSpace(strings.ToLower(result))

	// If empty, return default
	if response == "" {
		return defaultYes
	}

	return response == "y" || response == "yes"
}

// Prompt prompts the user for input with a message
// Supports ESC to cancel (returns empty string)
func Prompt(message string) string {
	prompt := promptui.Prompt{
		Label: message,
	}

	result, err := prompt.Run()
	if err != nil {
		// User pressed ESC or Ctrl+C
		return ""
	}

	return strings.TrimSpace(result)
}

// PromptWithDefault prompts the user for input with a default value
// Supports ESC to cancel (returns the default value)
func PromptWithDefault(message, defaultValue string) string {
	prompt := promptui.Prompt{
		Label:   message,
		Default: defaultValue,
	}

	result, err := prompt.Run()
	if err != nil {
		// User pressed ESC or Ctrl+C, return default
		return defaultValue
	}

	response := strings.TrimSpace(result)
	if response == "" {
		return defaultValue
	}

	return response
}

// Select prompts the user to select from a list of options using arrow keys
// Supports ESC to cancel (returns -1, "")
func Select(message string, options []string) (int, string) {
	prompt := promptui.Select{
		Label: message,
		Items: options,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | cyan }}",
			Active:   "▸ {{ . | green }}",
			Inactive: "  {{ . }}",
			Selected: "{{ \"✓\" | green }} {{ . }}",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		// User pressed ESC or Ctrl+C
		return -1, ""
	}

	return index, options[index]
}

// ConfirmBulkOperation prompts the user to confirm a bulk operation
// Shows the items that will be affected and asks for confirmation
func ConfirmBulkOperation(operation string, items []string) bool {
	fmt.Printf("\n%s\n", Warning("⚠ Warning: This will "+operation+" the following items:"))

	for _, item := range items {
		fmt.Printf("  - %s\n", Highlight(item))
	}

	fmt.Println()
	return Confirm(fmt.Sprintf("Are you sure you want to %s %d item(s)?", operation, len(items)))
}

// SelectWithArrows prompts the user to select from a list using arrow keys
// Returns the selected index and value, or -1 if cancelled
func SelectWithArrows(label string, items []string) (int, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  10, // Show up to 10 items at a time
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | cyan }}",
			Active:   "▸ {{ . | green }}",
			Inactive: "  {{ . }}",
			Selected: "{{ \"✓\" | green }} {{ . }}",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return index, nil
}
