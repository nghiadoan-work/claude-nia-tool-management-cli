package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToolArg(t *testing.T) {
	tests := []struct {
		name            string
		arg             string
		expectedName    string
		expectedVersion string
	}{
		{
			name:            "name only",
			arg:             "code-reviewer",
			expectedName:    "code-reviewer",
			expectedVersion: "",
		},
		{
			name:            "name with version",
			arg:             "code-reviewer@1.0.0",
			expectedName:    "code-reviewer",
			expectedVersion: "1.0.0",
		},
		{
			name:            "name with semantic version",
			arg:             "git-helper@2.3.1",
			expectedName:    "git-helper",
			expectedVersion: "2.3.1",
		},
		{
			name:            "name with version containing dots",
			arg:             "test-tool@1.0.0-beta.1",
			expectedName:    "test-tool",
			expectedVersion: "1.0.0-beta.1",
		},
		{
			name:            "multiple @ signs (only first is delimiter)",
			arg:             "tool@1.0@test",
			expectedName:    "tool",
			expectedVersion: "1.0@test",
		},
		{
			name:            "empty string",
			arg:             "",
			expectedName:    "",
			expectedVersion: "",
		},
		{
			name:            "@ only",
			arg:             "@",
			expectedName:    "",
			expectedVersion: "",
		},
		{
			name:            "name with trailing @",
			arg:             "tool@",
			expectedName:    "tool",
			expectedVersion: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, version := parseToolArg(tt.arg)
			assert.Equal(t, tt.expectedName, name, "name should match")
			assert.Equal(t, tt.expectedVersion, version, "version should match")
		})
	}
}

func TestInstallCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true, // cobra.MinimumNArgs(1) will fail
		},
		{
			name:    "single tool",
			args:    []string{"code-reviewer"},
			wantErr: false,
		},
		{
			name:    "tool with version",
			args:    []string{"code-reviewer@1.0.0"},
			wantErr: false,
		},
		{
			name:    "multiple tools",
			args:    []string{"tool1", "tool2", "tool3"},
			wantErr: false,
		},
		{
			name:    "mixed tools with and without versions",
			args:    []string{"tool1@1.0", "tool2", "tool3@2.1.0"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test argument validation
			err := installCmd.Args(installCmd, tt.args)
			if tt.wantErr {
				assert.Error(t, err, "should return error for invalid arguments")
			} else {
				assert.NoError(t, err, "should not return error for valid arguments")
			}
		})
	}
}

func TestInstallCmdFlags(t *testing.T) {
	// Test that flags are defined
	assert.NotNil(t, installCmd.Flags().Lookup("force"), "should have --force flag")
	assert.NotNil(t, installCmd.Flags().Lookup("path"), "should have --path flag")

	// Test flag shortcuts
	forceFlag := installCmd.Flags().Lookup("force")
	assert.Equal(t, "f", forceFlag.Shorthand, "force flag should have -f shorthand")
}

func TestInstallCmdMetadata(t *testing.T) {
	// Test command metadata
	assert.Equal(t, "install", installCmd.Use[:7], "command name should be install")
	assert.NotEmpty(t, installCmd.Short, "should have short description")
	assert.NotEmpty(t, installCmd.Long, "should have long description")
	assert.NotEmpty(t, installCmd.Example, "should have examples")
}

func TestToolSpec(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		version  string
	}{
		{
			name:     "basic spec",
			toolName: "code-reviewer",
			version:  "1.0.0",
		},
		{
			name:     "no version",
			toolName: "git-helper",
			version:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := toolSpec{
				name:    tt.toolName,
				version: tt.version,
			}

			assert.Equal(t, tt.toolName, spec.name, "name should match")
			assert.Equal(t, tt.version, spec.version, "version should match")
		})
	}
}
