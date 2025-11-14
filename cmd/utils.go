package cmd

import (
	"fmt"
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
