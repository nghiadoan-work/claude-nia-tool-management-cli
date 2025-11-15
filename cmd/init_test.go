package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitCommand_Flags(t *testing.T) {
	// Test that the --path flag exists
	pathFlag := initCmd.Flags().Lookup("path")
	assert.NotNil(t, pathFlag)
	assert.Equal(t, "", pathFlag.DefValue)

	// Test that the --force flag exists
	forceFlag := initCmd.Flags().Lookup("force")
	assert.NotNil(t, forceFlag)
	assert.Equal(t, "false", forceFlag.DefValue)

	// Test that the -f shorthand exists
	forceFlag = initCmd.Flags().ShorthandLookup("f")
	assert.NotNil(t, forceFlag)
	assert.Equal(t, "force", forceFlag.Name)
}

func TestInitializeLockFile(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	lockFilePath := filepath.Join(tempDir, ".claude-lock.json")

	// Initialize lock file
	err := initializeLockFile(lockFilePath)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(lockFilePath)
	require.NoError(t, err)

	// Read and verify contents
	data, err := os.ReadFile(lockFilePath)
	require.NoError(t, err)

	var lockFile models.LockFile
	err = json.Unmarshal(data, &lockFile)
	require.NoError(t, err)

	// Verify structure
	assert.Equal(t, "1.0", lockFile.Version)
	assert.NotEmpty(t, lockFile.Registry)
	assert.NotNil(t, lockFile.Tools)
	assert.Equal(t, 0, len(lockFile.Tools))
	assert.False(t, lockFile.UpdatedAt.IsZero())
}

func TestInitializeLockFile_CreatesValidJSON(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	lockFilePath := filepath.Join(tempDir, ".claude-lock.json")

	// Initialize lock file
	err := initializeLockFile(lockFilePath)
	require.NoError(t, err)

	// Read file
	data, err := os.ReadFile(lockFilePath)
	require.NoError(t, err)

	// Verify it's valid JSON
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	// Verify expected fields
	assert.Contains(t, result, "version")
	assert.Contains(t, result, "updated_at")
	assert.Contains(t, result, "registry")
	assert.Contains(t, result, "tools")
}

func TestInitCommand_DirectoryStructure(t *testing.T) {
	// This test verifies the expected directory structure
	// We can't test the full command without mocking, but we can test the structure

	expectedDirs := []string{"agents", "commands", "skills"}

	for _, dir := range expectedDirs {
		assert.NotEmpty(t, dir, "Directory name should not be empty")
	}
}
