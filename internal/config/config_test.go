package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfigService(t *testing.T) {
	config := models.NewDefaultConfig()
	service := NewConfigService(config)

	assert.NotNil(t, service)
	assert.Equal(t, config, service.GetConfig())
}

func TestLoadConfig_DefaultOnly(t *testing.T) {
	// Load config with no files (should return defaults)
	config, err := LoadConfig("")

	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "https://github.com/nghiadt/claude-tools-registry", config.Registry.URL)
	assert.Equal(t, "main", config.Registry.Branch)
	assert.Equal(t, ".claude", config.Local.DefaultPath)
	assert.True(t, config.Local.AutoUpdateCheck)
	assert.Equal(t, 86400, config.Local.UpdateCheckInterval)
}

func TestLoadConfig_FromFile(t *testing.T) {
	// Create temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `registry:
  url: https://github.com/test/custom-registry
  branch: develop
  auth_token: test-token
local:
  default_path: .custom-claude
  auto_update_check: false
  update_check_interval: 3600
publish:
  default_author: Test Author
  auto_version_bump: minor
  create_pr: false
`

	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	config, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, "https://github.com/test/custom-registry", config.Registry.URL)
	assert.Equal(t, "develop", config.Registry.Branch)
	assert.Equal(t, "test-token", config.Registry.AuthToken)
	assert.Equal(t, ".custom-claude", config.Local.DefaultPath)
	assert.False(t, config.Local.AutoUpdateCheck)
	assert.Equal(t, 3600, config.Local.UpdateCheckInterval)
	assert.Equal(t, "Test Author", config.Publish.DefaultAuthor)
	assert.Equal(t, "minor", config.Publish.AutoVersionBump)
	assert.False(t, config.Publish.CreatePR)
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid-config.yaml")

	// Write invalid YAML
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
	require.NoError(t, err)

	_, err = LoadConfig(configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load config")
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	_, err := LoadConfig("/non/existent/path/config.yaml")
	assert.Error(t, err)
}

func TestMergeConfig(t *testing.T) {
	target := models.NewDefaultConfig()
	source := &models.Config{
		Registry: models.RegistryConfig{
			URL:       "https://github.com/custom/registry",
			Branch:    "custom-branch",
			AuthToken: "custom-token",
		},
		Local: models.LocalConfig{
			DefaultPath:         ".custom",
			AutoUpdateCheck:     false,
			UpdateCheckInterval: 7200,
		},
		Publish: models.PublishConfig{
			DefaultAuthor:   "Custom Author",
			AutoVersionBump: "major",
			CreatePR:        false,
		},
	}

	mergeConfig(target, source)

	assert.Equal(t, "https://github.com/custom/registry", target.Registry.URL)
	assert.Equal(t, "custom-branch", target.Registry.Branch)
	assert.Equal(t, "custom-token", target.Registry.AuthToken)
	assert.Equal(t, ".custom", target.Local.DefaultPath)
	assert.False(t, target.Local.AutoUpdateCheck)
	assert.Equal(t, 7200, target.Local.UpdateCheckInterval)
	assert.Equal(t, "Custom Author", target.Publish.DefaultAuthor)
	assert.Equal(t, "major", target.Publish.AutoVersionBump)
	assert.False(t, target.Publish.CreatePR)
}

func TestMergeConfig_PartialOverride(t *testing.T) {
	target := models.NewDefaultConfig()
	originalURL := target.Registry.URL

	source := &models.Config{
		Registry: models.RegistryConfig{
			Branch: "new-branch",
		},
		Local:   models.LocalConfig{},
		Publish: models.PublishConfig{},
	}

	mergeConfig(target, source)

	// URL should remain unchanged
	assert.Equal(t, originalURL, target.Registry.URL)
	// Branch should be updated
	assert.Equal(t, "new-branch", target.Registry.Branch)
}

func TestApplyEnvOverrides(t *testing.T) {
	config := models.NewDefaultConfig()

	// Set environment variables
	os.Setenv("CNTM_REGISTRY_URL", "https://github.com/env/registry")
	os.Setenv("CNTM_REGISTRY_BRANCH", "env-branch")
	os.Setenv("CNTM_REGISTRY_TOKEN", "env-token")
	os.Setenv("CNTM_DEFAULT_PATH", ".env-claude")
	os.Setenv("CNTM_AUTO_UPDATE", "false")
	os.Setenv("CNTM_DEFAULT_AUTHOR", "Env Author")
	os.Setenv("CNTM_AUTO_VERSION_BUMP", "major")

	defer func() {
		os.Unsetenv("CNTM_REGISTRY_URL")
		os.Unsetenv("CNTM_REGISTRY_BRANCH")
		os.Unsetenv("CNTM_REGISTRY_TOKEN")
		os.Unsetenv("CNTM_DEFAULT_PATH")
		os.Unsetenv("CNTM_AUTO_UPDATE")
		os.Unsetenv("CNTM_DEFAULT_AUTHOR")
		os.Unsetenv("CNTM_AUTO_VERSION_BUMP")
	}()

	applyEnvOverrides(config)

	assert.Equal(t, "https://github.com/env/registry", config.Registry.URL)
	assert.Equal(t, "env-branch", config.Registry.Branch)
	assert.Equal(t, "env-token", config.Registry.AuthToken)
	assert.Equal(t, ".env-claude", config.Local.DefaultPath)
	assert.False(t, config.Local.AutoUpdateCheck)
	assert.Equal(t, "Env Author", config.Publish.DefaultAuthor)
	assert.Equal(t, "major", config.Publish.AutoVersionBump)
}

func TestApplyEnvOverrides_AutoUpdateVariations(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"true string", "true", true},
		{"1 string", "1", true},
		{"false string", "false", false},
		{"0 string", "0", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := models.NewDefaultConfig()
			config.Local.AutoUpdateCheck = false // Start with false

			if tt.envValue != "" {
				os.Setenv("CNTM_AUTO_UPDATE", tt.envValue)
				defer os.Unsetenv("CNTM_AUTO_UPDATE")
			}

			applyEnvOverrides(config)

			if tt.envValue != "" {
				assert.Equal(t, tt.expected, config.Local.AutoUpdateCheck)
			} else {
				// Empty should not change the value
				assert.False(t, config.Local.AutoUpdateCheck)
			}
		})
	}
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "subdir", "config.yaml")

	config := &models.Config{
		Registry: models.RegistryConfig{
			URL:       "https://github.com/test/registry",
			Branch:    "main",
			AuthToken: "test-token",
		},
		Local: models.LocalConfig{
			DefaultPath:         ".claude",
			AutoUpdateCheck:     true,
			UpdateCheckInterval: 86400,
		},
		Publish: models.PublishConfig{
			DefaultAuthor:   "Test",
			AutoVersionBump: "patch",
			CreatePR:        true,
		},
	}

	err := SaveConfig(config, configPath)
	require.NoError(t, err)

	// Verify file was created
	assert.FileExists(t, configPath)

	// Load and verify content
	loadedConfig, err := LoadConfig(configPath)
	require.NoError(t, err)

	assert.Equal(t, config.Registry.URL, loadedConfig.Registry.URL)
	assert.Equal(t, config.Registry.Branch, loadedConfig.Registry.Branch)
	assert.Equal(t, config.Local.DefaultPath, loadedConfig.Local.DefaultPath)
}

func TestSaveConfig_InvalidConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	// Invalid config (missing required fields)
	config := &models.Config{
		Registry: models.RegistryConfig{},
		Local:    models.LocalConfig{},
		Publish:  models.PublishConfig{},
	}

	err := SaveConfig(config, configPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid config")
}

func TestGetGlobalConfigPath(t *testing.T) {
	path, err := GetGlobalConfigPath()
	require.NoError(t, err)

	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".claude-tools-config.yaml")

	assert.Equal(t, expectedPath, path)
}

func TestGetProjectConfigPath(t *testing.T) {
	path, err := GetProjectConfigPath()
	require.NoError(t, err)

	currentDir, _ := os.Getwd()
	expectedPath := filepath.Join(currentDir, ".claude-tools-config.yaml")

	assert.Equal(t, expectedPath, path)
}

func TestLoadConfig_ConfigPrecedence(t *testing.T) {
	// This test verifies the precedence: ENV > Specific File > Project > Global > Default
	tmpDir := t.TempDir()

	// Create a specific config file
	specificPath := filepath.Join(tmpDir, "specific.yaml")
	specificContent := `registry:
  url: https://github.com/specific/registry
  branch: specific-branch
`
	err := os.WriteFile(specificPath, []byte(specificContent), 0644)
	require.NoError(t, err)

	// Set environment variable (should have highest priority)
	os.Setenv("CNTM_REGISTRY_URL", "https://github.com/env/registry")
	defer os.Unsetenv("CNTM_REGISTRY_URL")

	config, err := LoadConfig(specificPath)
	require.NoError(t, err)

	// ENV variable should override file
	assert.Equal(t, "https://github.com/env/registry", config.Registry.URL)
	// But file should set branch
	assert.Equal(t, "specific-branch", config.Registry.Branch)
}
