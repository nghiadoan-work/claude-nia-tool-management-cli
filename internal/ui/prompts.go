package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm prompts the user for yes/no confirmation
// Returns true if user confirms, false otherwise
func Confirm(message string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s %s %s ", Warning("?"), message, Faint("[y/N]"))

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// ConfirmWithDefault prompts the user for yes/no confirmation with a default value
func ConfirmWithDefault(message string, defaultYes bool) bool {
	reader := bufio.NewReader(os.Stdin)

	prompt := "[y/N]"
	if defaultYes {
		prompt = "[Y/n]"
	}

	fmt.Printf("%s %s %s ", Warning("?"), message, Faint(prompt))

	response, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}

	response = strings.TrimSpace(strings.ToLower(response))

	// If empty, return default
	if response == "" {
		return defaultYes
	}

	return response == "y" || response == "yes"
}

// Prompt prompts the user for input with a message
func Prompt(message string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s %s ", Info("?"), message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return ""
	}

	return strings.TrimSpace(response)
}

// PromptWithDefault prompts the user for input with a default value
func PromptWithDefault(message, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s %s %s ", Info("?"), message, Faint(fmt.Sprintf("[%s]", defaultValue)))

	response, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue
	}

	response = strings.TrimSpace(response)
	if response == "" {
		return defaultValue
	}

	return response
}

// Select prompts the user to select from a list of options
func Select(message string, options []string) (int, string) {
	fmt.Printf("%s %s\n", Info("?"), message)

	for i, option := range options {
		fmt.Printf("  %s %s\n", Faint(fmt.Sprintf("%d)", i+1)), option)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(Faint("Select: "))

	response, err := reader.ReadString('\n')
	if err != nil {
		return -1, ""
	}

	response = strings.TrimSpace(response)

	// Try to parse as number
	var selected int
	_, err = fmt.Sscanf(response, "%d", &selected)
	if err != nil || selected < 1 || selected > len(options) {
		return -1, ""
	}

	return selected - 1, options[selected-1]
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
