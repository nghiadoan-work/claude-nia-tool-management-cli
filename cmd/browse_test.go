package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatRelativeTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "just now",
			time:     now.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-1 * time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "5 minutes ago",
			time:     now.Add(-5 * time.Minute),
			expected: "5 minutes ago",
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-1 * time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "3 hours ago",
			time:     now.Add(-3 * time.Hour),
			expected: "3 hours ago",
		},
		{
			name:     "1 day ago",
			time:     now.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "3 days ago",
			time:     now.Add(-3 * 24 * time.Hour),
			expected: "3 days ago",
		},
		{
			name:     "1 week ago",
			time:     now.Add(-7 * 24 * time.Hour),
			expected: "1 week ago",
		},
		{
			name:     "2 weeks ago",
			time:     now.Add(-14 * 24 * time.Hour),
			expected: "2 weeks ago",
		},
		{
			name:     "1 month ago",
			time:     now.Add(-30 * 24 * time.Hour),
			expected: "1 month ago",
		},
		{
			name:     "3 months ago",
			time:     now.Add(-90 * 24 * time.Hour),
			expected: "3 months ago",
		},
		{
			name:     "1 year ago",
			time:     now.Add(-365 * 24 * time.Hour),
			expected: "1 year ago",
		},
		{
			name:     "2 years ago",
			time:     now.Add(-730 * 24 * time.Hour),
			expected: "2 years ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRelativeTime(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatRelativeTime_EdgeCases(t *testing.T) {
	now := time.Now()

	// Test boundary conditions
	tests := []struct {
		name     string
		time     time.Time
		contains string // Check if result contains this string
	}{
		{
			name:     "59 seconds",
			time:     now.Add(-59 * time.Second),
			contains: "just now",
		},
		{
			name:     "59 minutes",
			time:     now.Add(-59 * time.Minute),
			contains: "minutes ago",
		},
		{
			name:     "23 hours",
			time:     now.Add(-23 * time.Hour),
			contains: "hours ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatRelativeTime(tt.time)
			assert.Contains(t, result, tt.contains)
		})
	}
}
