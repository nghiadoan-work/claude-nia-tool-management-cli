package version

import (
	"strings"
	"testing"
)

func TestGetInfo(t *testing.T) {
	info := GetInfo()

	if info.Version == "" {
		t.Error("Version should not be empty")
	}

	if info.GitCommit == "" {
		t.Error("GitCommit should not be empty")
	}

	if info.BuildDate == "" {
		t.Error("BuildDate should not be empty")
	}

	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}
}

func TestInfo_String(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2025-11-15",
		GoVersion: "1.21.0",
	}

	result := info.String()
	if result != "1.0.0" {
		t.Errorf("Expected '1.0.0', got '%s'", result)
	}
}

func TestInfo_LongString(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2025-11-15",
		GoVersion: "1.21.0",
	}

	result := info.LongString()

	// Check that all fields are present
	expectedParts := []string{
		"cntm version 1.0.0",
		"Git commit: abc123",
		"Build date: 2025-11-15",
		"Go version: 1.21.0",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected LongString to contain '%s', but got: %s", part, result)
		}
	}
}
