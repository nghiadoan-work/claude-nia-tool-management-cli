package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantOwner   string
		wantRepo    string
		expectError bool
	}{
		{
			name:      "full HTTPS URL",
			url:       "https://github.com/nghiadoan-work/claude-tools-registry",
			wantOwner: "nghiadoan-work",
			wantRepo:  "claude-tools-registry",
		},
		{
			name:      "full HTTPS URL with .git",
			url:       "https://github.com/nghiadoan-work/claude-tools-registry.git",
			wantOwner: "nghiadoan-work",
			wantRepo:  "claude-tools-registry",
		},
		{
			name:      "HTTP URL",
			url:       "http://github.com/nghiadoan-work/claude-tools-registry",
			wantOwner: "nghiadoan-work",
			wantRepo:  "claude-tools-registry",
		},
		{
			name:      "without protocol",
			url:       "github.com/nghiadoan-work/claude-tools-registry",
			wantOwner: "nghiadoan-work",
			wantRepo:  "claude-tools-registry",
		},
		{
			name:      "simple format",
			url:       "nghiadoan-work/claude-tools-registry",
			wantOwner: "nghiadoan-work",
			wantRepo:  "claude-tools-registry",
		},
		{
			name:      "with trailing slash",
			url:       "https://github.com/nghiadoan-work/claude-tools-registry/",
			wantOwner: "nghiadoan-work",
			wantRepo:  "claude-tools-registry",
		},
		{
			name:        "invalid - missing repo",
			url:         "https://github.com/nghiadoan-work",
			expectError: true,
		},
		{
			name:        "invalid - empty",
			url:         "",
			expectError: true,
		},
		{
			name:        "invalid - just owner",
			url:         "nghiadoan-work",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseGitHubURL(tt.url)

			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantOwner, owner)
				assert.Equal(t, tt.wantRepo, repo)
			}
		})
	}
}
