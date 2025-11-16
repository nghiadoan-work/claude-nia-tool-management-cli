package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// parseGitHubURL extracts owner and repo from a GitHub URL
// Supports formats:
//   - https://github.com/owner/repo
//   - https://github.com/owner/repo.git
//   - github.com/owner/repo
func parseGitHubURL(url string) (owner, repo string, err error) {
	// Remove common prefixes
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "github.com/")
	url = strings.TrimSuffix(url, ".git")
	url = strings.Trim(url, "/")

	// Split by /
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format, expected: github.com/owner/repo")
	}

	return parts[0], parts[1], nil
}

// promptString prompts the user for a string input with an optional default value
func promptString(prompt, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}
