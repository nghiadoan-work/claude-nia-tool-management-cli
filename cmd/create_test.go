package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateReadme(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		toolType    models.ToolType
		toolName    string
		description string
	}{
		{
			name:        "agent readme",
			toolType:    models.ToolTypeAgent,
			toolName:    "test-agent",
			description: "A test agent",
		},
		{
			name:        "command readme",
			toolType:    models.ToolTypeCommand,
			toolName:    "test-command",
			description: "A test command",
		},
		{
			name:        "skill readme",
			toolType:    models.ToolTypeSkill,
			toolName:    "test-skill",
			description: "A test skill",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolPath := filepath.Join(tempDir, tt.toolName)
			require.NoError(t, os.MkdirAll(toolPath, 0755))

			err := createReadme(toolPath, tt.toolType, tt.toolName, tt.description)
			assert.NoError(t, err)

			// Verify README.md was created
			readmePath := filepath.Join(toolPath, "README.md")
			assert.FileExists(t, readmePath)

			// Verify content
			content, err := os.ReadFile(readmePath)
			require.NoError(t, err)
			assert.Contains(t, string(content), tt.toolName)
			assert.Contains(t, string(content), tt.description)
			assert.Contains(t, string(content), string(tt.toolType))
		})
	}
}

func TestCreateAgentFile(t *testing.T) {
	tempDir := t.TempDir()
	toolPath := filepath.Join(tempDir, "test-agent")
	require.NoError(t, os.MkdirAll(toolPath, 0755))

	err := createAgentFile(toolPath, "test-agent")
	assert.NoError(t, err)

	agentPath := filepath.Join(toolPath, "agent.md")
	assert.FileExists(t, agentPath)

	content, err := os.ReadFile(agentPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test-agent")
	assert.Contains(t, string(content), "Capabilities")
}

func TestCreateCommandFile(t *testing.T) {
	tempDir := t.TempDir()
	toolPath := filepath.Join(tempDir, "test-command")
	require.NoError(t, os.MkdirAll(toolPath, 0755))

	err := createCommandFile(toolPath, "test-command")
	assert.NoError(t, err)

	commandPath := filepath.Join(toolPath, "command.md")
	assert.FileExists(t, commandPath)

	content, err := os.ReadFile(commandPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test-command")
	assert.Contains(t, string(content), "Syntax")
	assert.Contains(t, string(content), "Examples")
}

func TestCreateSkillFile(t *testing.T) {
	tempDir := t.TempDir()
	toolPath := filepath.Join(tempDir, "test-skill")
	require.NoError(t, os.MkdirAll(toolPath, 0755))

	err := createSkillFile(toolPath, "test-skill")
	assert.NoError(t, err)

	skillPath := filepath.Join(toolPath, "skill.md")
	assert.FileExists(t, skillPath)

	content, err := os.ReadFile(skillPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test-skill")
	assert.Contains(t, string(content), "Knowledge Areas")
	assert.Contains(t, string(content), "Best Practices")
}

func TestCreateTypeSpecificFiles(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		toolType     models.ToolType
		expectedFile string
	}{
		{
			name:         "agent",
			toolType:     models.ToolTypeAgent,
			expectedFile: "agent.md",
		},
		{
			name:         "command",
			toolType:     models.ToolTypeCommand,
			expectedFile: "command.md",
		},
		{
			name:         "skill",
			toolType:     models.ToolTypeSkill,
			expectedFile: "skill.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toolPath := filepath.Join(tempDir, "test-"+string(tt.toolType))
			require.NoError(t, os.MkdirAll(toolPath, 0755))

			err := createTypeSpecificFiles(toolPath, tt.toolType, "test-tool")
			assert.NoError(t, err)

			expectedPath := filepath.Join(toolPath, tt.expectedFile)
			assert.FileExists(t, expectedPath)
		})
	}
}
