package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindToolPath(t *testing.T) {
	tempDir := t.TempDir()

	// Setup test directories
	agentPath := filepath.Join(tempDir, "agents", "test-agent")
	commandPath := filepath.Join(tempDir, "commands", "test-command")
	skillPath := filepath.Join(tempDir, "skills", "test-skill")

	require.NoError(t, os.MkdirAll(agentPath, 0755))
	require.NoError(t, os.MkdirAll(commandPath, 0755))
	require.NoError(t, os.MkdirAll(skillPath, 0755))

	// Create test config
	testCfg := models.NewDefaultConfig()
	testCfg.Local.DefaultPath = tempDir

	tests := []struct {
		name         string
		toolName     string
		expectedPath string
	}{
		{
			name:         "find agent",
			toolName:     "test-agent",
			expectedPath: agentPath,
		},
		{
			name:         "find command",
			toolName:     "test-command",
			expectedPath: commandPath,
		},
		{
			name:         "find skill",
			toolName:     "test-skill",
			expectedPath: skillPath,
		},
		{
			name:         "not found",
			toolName:     "nonexistent",
			expectedPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findToolPath(tt.toolName, testCfg)
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}

func TestDetectToolTypeFromPath(t *testing.T) {
	tests := []struct {
		name         string
		toolPath     string
		expectedType string
		expectError  bool
	}{
		{
			name:         "agent path",
			toolPath:     "/path/to/agents/my-agent",
			expectedType: "agent",
			expectError:  false,
		},
		{
			name:         "command path",
			toolPath:     "/path/to/commands/my-command",
			expectedType: "command",
			expectError:  false,
		},
		{
			name:         "skill path",
			toolPath:     "/path/to/skills/my-skill",
			expectedType: "skill",
			expectError:  false,
		},
		{
			name:         "windows agent path",
			toolPath:     "C:\\path\\to\\agents\\my-agent",
			expectedType: "agent",
			expectError:  false,
		},
		{
			name:        "unknown path",
			toolPath:    "/path/to/unknown/my-tool",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolType, err := detectToolTypeFromPath(tt.toolPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedType, toolType)
			}
		})
	}
}

func TestBumpVersion(t *testing.T) {
	tests := []struct {
		name            string
		currentVersion  string
		expectedVersion string
	}{
		{
			name:            "bump patch version",
			currentVersion:  "1.0.0",
			expectedVersion: "1.0.1",
		},
		{
			name:            "bump higher patch",
			currentVersion:  "1.2.5",
			expectedVersion: "1.2.6",
		},
		{
			name:            "invalid version format",
			currentVersion:  "1.0",
			expectedVersion: "1.0.0",
		},
		{
			name:            "non-numeric patch",
			currentVersion:  "1.0.x",
			expectedVersion: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bumpVersion(tt.currentVersion)
			assert.Equal(t, tt.expectedVersion, result)
		})
	}
}

func TestScanLocalTools(t *testing.T) {
	tempDir := t.TempDir()

	// Setup test directories with various tools
	agentPath1 := filepath.Join(tempDir, "agents", "agent1")
	agentPath2 := filepath.Join(tempDir, "agents", "agent2")
	commandPath := filepath.Join(tempDir, "commands", "cmd1")
	skillPath := filepath.Join(tempDir, "skills", "skill1")

	require.NoError(t, os.MkdirAll(agentPath1, 0755))
	require.NoError(t, os.MkdirAll(agentPath2, 0755))
	require.NoError(t, os.MkdirAll(commandPath, 0755))
	require.NoError(t, os.MkdirAll(skillPath, 0755))

	// Create a file in agents dir (should be ignored)
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "agents", "file.txt"), []byte("test"), 0644))

	// Create test config
	testCfg := models.NewDefaultConfig()
	testCfg.Local.DefaultPath = tempDir

	// Scan tools
	tools, err := scanLocalTools(testCfg)
	require.NoError(t, err)

	// Should find 4 tools (2 agents, 1 command, 1 skill)
	assert.Len(t, tools, 4)

	// Verify tool types
	toolNames := make(map[string]models.ToolType)
	for _, tool := range tools {
		toolNames[tool.Name] = tool.Type
	}

	assert.Equal(t, models.ToolTypeAgent, toolNames["agent1"])
	assert.Equal(t, models.ToolTypeAgent, toolNames["agent2"])
	assert.Equal(t, models.ToolTypeCommand, toolNames["cmd1"])
	assert.Equal(t, models.ToolTypeSkill, toolNames["skill1"])
}

func TestScanLocalTools_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create test config with empty directory
	testCfg := models.NewDefaultConfig()
	testCfg.Local.DefaultPath = tempDir

	// Scan tools
	tools, err := scanLocalTools(testCfg)
	require.NoError(t, err)

	// Should find no tools
	assert.Len(t, tools, 0)
}
