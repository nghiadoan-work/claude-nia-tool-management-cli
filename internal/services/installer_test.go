package services

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock GitHub Downloader for testing
type mockGitHubDownloader struct {
	downloadFunc  func(url string, size int64, showProgress bool) ([]byte, error)
	downloadError error
	downloadData  []byte
}

func (m *mockGitHubDownloader) DownloadFile(url string, size int64, showProgress bool) ([]byte, error) {
	if m.downloadFunc != nil {
		return m.downloadFunc(url, size, showProgress)
	}
	if m.downloadError != nil {
		return nil, m.downloadError
	}
	return m.downloadData, nil
}

// Mock Registry Service for installer testing
type mockInstallerRegistryService struct {
	tools      map[string]*models.ToolInfo
	fetchError error
}

func (m *mockInstallerRegistryService) GetTool(name string, toolType models.ToolType) (*models.ToolInfo, error) {
	if m.fetchError != nil {
		return nil, m.fetchError
	}
	key := string(toolType) + ":" + name
	if tool, ok := m.tools[key]; ok {
		return tool, nil
	}
	return nil, fmt.Errorf("tool %s not found", name)
}

func (m *mockInstallerRegistryService) GetRegistry() (*models.Registry, error) {
	return &models.Registry{}, nil
}

// setupTestInstaller creates a test installer with all dependencies
func setupTestInstaller(t *testing.T) (*InstallerService, string, func()) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "installer-test-*")
	require.NoError(t, err)

	baseDir := filepath.Join(tempDir, ".claude")
	lockFilePath := filepath.Join(baseDir, ".claude-lock.json")

	// Create FSManager
	fsManager, err := data.NewFSManager(baseDir)
	require.NoError(t, err)

	// Create LockFileService
	lockFileService, err := NewLockFileService(lockFilePath)
	require.NoError(t, err)

	// Create config
	config := &models.Config{
		Registry: models.RegistryConfig{
			URL:    "https://github.com/test/registry",
			Branch: "main",
		},
		Local: models.LocalConfig{
			DefaultPath: baseDir,
		},
	}

	// Create mock GitHub client with test ZIP data
	githubClient := &mockGitHubDownloader{
		downloadData: createTestZIP(t),
	}

	// Create mock registry service
	registryService := &mockInstallerRegistryService{
		tools: make(map[string]*models.ToolInfo),
	}

	// Add test tool to registry
	registryService.tools["agent:test-agent"] = &models.ToolInfo{
		Name:        "test-agent",
		Version:     "1.0.0",
		Description: "Test agent",
		Type:        models.ToolTypeAgent,
		Author:      "test",
		File:        "tools/agents/test-agent.zip",
		Size:        1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create installer service
	installer, err := NewInstallerService(
		githubClient,
		registryService,
		fsManager,
		lockFileService,
		config,
	)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return installer, baseDir, cleanup
}

// createTestZIP creates a minimal valid ZIP file for testing
func createTestZIP(t *testing.T) []byte {
	// Create a temp directory with a test file
	tempDir, err := os.MkdirTemp("", "zip-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	// Create FSManager to create ZIP
	fsManager, err := data.NewFSManager(tempDir)
	require.NoError(t, err)

	// Create ZIP
	zipPath := filepath.Join(tempDir, "test.zip")
	err = fsManager.CreateZIP(tempDir, zipPath)
	require.NoError(t, err)

	// Read ZIP data
	zipData, err := os.ReadFile(zipPath)
	require.NoError(t, err)

	return zipData
}

func TestNewInstallerService(t *testing.T) {
	tests := []struct {
		name        string
		setupMocks  func() (GitHubDownloader, RegistryServiceInterface, FSManagerInterface, LockFileServiceInterface, *models.Config)
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid initialization",
			setupMocks: func() (GitHubDownloader, RegistryServiceInterface, FSManagerInterface, LockFileServiceInterface, *models.Config) {
				tempDir, _ := os.MkdirTemp("", "test-*")
				baseDir := filepath.Join(tempDir, ".claude")
				lockFilePath := filepath.Join(baseDir, ".claude-lock.json")

				fsManager, _ := data.NewFSManager(baseDir)
				lockFileService, _ := NewLockFileService(lockFilePath)
				config := &models.Config{
					Registry: models.RegistryConfig{URL: "https://test.com", Branch: "main"},
					Local:    models.LocalConfig{DefaultPath: baseDir},
				}
				githubClient := &mockGitHubDownloader{downloadData: []byte("test")}
				registryService := &mockInstallerRegistryService{tools: make(map[string]*models.ToolInfo)}

				return githubClient, registryService, fsManager, lockFileService, config
			},
			expectError: false,
		},
		{
			name: "nil github client",
			setupMocks: func() (GitHubDownloader, RegistryServiceInterface, FSManagerInterface, LockFileServiceInterface, *models.Config) {
				tempDir, _ := os.MkdirTemp("", "test-*")
				baseDir := filepath.Join(tempDir, ".claude")
				lockFilePath := filepath.Join(baseDir, ".claude-lock.json")

				fsManager, _ := data.NewFSManager(baseDir)
				lockFileService, _ := NewLockFileService(lockFilePath)
				config := &models.Config{
					Registry: models.RegistryConfig{URL: "https://test.com", Branch: "main"},
					Local:    models.LocalConfig{DefaultPath: baseDir},
				}
				registryService := &mockInstallerRegistryService{tools: make(map[string]*models.ToolInfo)}

				return nil, registryService, fsManager, lockFileService, config
			},
			expectError: true,
			errorMsg:    "github client cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gc, rs, fm, lfs, cfg := tt.setupMocks()
			defer func() {
				if cfg != nil && cfg.Local.DefaultPath != "" {
					os.RemoveAll(filepath.Dir(cfg.Local.DefaultPath))
				}
			}()

			installer, err := NewInstallerService(gc, rs, fm, lfs, cfg)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, installer)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, installer)
			}
		})
	}
}

func TestInstaller_Install(t *testing.T) {
	t.Run("install new tool successfully", func(t *testing.T) {
		installer, baseDir, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.Install("test-agent")
		assert.NoError(t, err)

		// Verify tool is in lock file
		installed, err := installer.IsInstalled("test-agent")
		assert.NoError(t, err)
		assert.True(t, installed)

		// Verify installation directory exists
		destDir := filepath.Join(baseDir, "agents", "test-agent")
		_, err = os.Stat(destDir)
		assert.NoError(t, err)

		// Verify tool info in lock file
		version, err := installer.GetInstalledVersion("test-agent")
		assert.NoError(t, err)
		assert.Equal(t, "1.0.0", version)
	})

	t.Run("skip already installed tool with same version", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		// Install first time
		err := installer.Install("test-agent")
		assert.NoError(t, err)

		// Install again - should skip
		err = installer.Install("test-agent")
		assert.NoError(t, err)
	})

	t.Run("tool not found in registry", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.Install("nonexistent-tool")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find tool")
	})

	t.Run("empty tool name", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.Install("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool name cannot be empty")
	})
}

func TestInstaller_InstallWithVersion(t *testing.T) {
	t.Run("install specific version", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.InstallWithVersion("test-agent", "1.0.0")
		assert.NoError(t, err)

		version, err := installer.GetInstalledVersion("test-agent")
		assert.NoError(t, err)
		assert.Equal(t, "1.0.0", version)
	})

	t.Run("version mismatch", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.InstallWithVersion("test-agent", "2.0.0")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "requested version 2.0.0 not found")
	})
}

func TestInstaller_InstallMultiple(t *testing.T) {
	t.Run("install multiple tools successfully", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		// Add another tool to registry
		regService := installer.registryService.(*mockInstallerRegistryService)
		regService.tools["command:test-command"] = &models.ToolInfo{
			Name:        "test-command",
			Version:     "1.0.0",
			Description: "Test command",
			Type:        models.ToolTypeCommand,
			Author:      "test",
			File:        "tools/commands/test-command.zip",
			Size:        1024,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		results, errors := installer.InstallMultiple([]string{"test-agent", "test-command"})
		assert.Len(t, results, 2)
		assert.Len(t, errors, 0)

		for _, result := range results {
			assert.True(t, result.Success)
			assert.NoError(t, result.Error)
		}
	})

	t.Run("install multiple with some failures", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		results, errors := installer.InstallMultiple([]string{"test-agent", "nonexistent"})
		assert.Len(t, results, 2)
		assert.Len(t, errors, 1)

		// First should succeed
		assert.True(t, results[0].Success)
		// Second should fail
		assert.False(t, results[1].Success)
	})

	t.Run("empty tool list", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		results, errors := installer.InstallMultiple([]string{})
		assert.Nil(t, results)
		assert.Len(t, errors, 1)
	})
}

func TestInstaller_VerifyInstallation(t *testing.T) {
	t.Run("verify valid installation", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		// Install tool
		err := installer.Install("test-agent")
		require.NoError(t, err)

		// Verify
		err = installer.VerifyInstallation("test-agent")
		assert.NoError(t, err)
	})

	t.Run("verify non-installed tool", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.VerifyInstallation("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool not found in lock file")
	})

	t.Run("verify with missing directory", func(t *testing.T) {
		installer, baseDir, cleanup := setupTestInstaller(t)
		defer cleanup()

		// Install tool
		err := installer.Install("test-agent")
		require.NoError(t, err)

		// Remove directory but keep lock file entry
		destDir := filepath.Join(baseDir, "agents", "test-agent")
		os.RemoveAll(destDir)

		// Verify should fail
		err = installer.VerifyInstallation("test-agent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "installation directory does not exist")
	})

	t.Run("empty tool name", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.VerifyInstallation("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool name cannot be empty")
	})
}

func TestInstaller_Uninstall(t *testing.T) {
	t.Run("uninstall successfully", func(t *testing.T) {
		installer, baseDir, cleanup := setupTestInstaller(t)
		defer cleanup()

		// Install tool first
		err := installer.Install("test-agent")
		require.NoError(t, err)

		// Verify installed
		installed, err := installer.IsInstalled("test-agent")
		require.NoError(t, err)
		require.True(t, installed)

		// Uninstall
		err = installer.Uninstall("test-agent")
		assert.NoError(t, err)

		// Verify not installed
		installed, err = installer.IsInstalled("test-agent")
		assert.NoError(t, err)
		assert.False(t, installed)

		// Verify directory removed
		destDir := filepath.Join(baseDir, "agents", "test-agent")
		_, err = os.Stat(destDir)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("uninstall non-installed tool", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.Uninstall("nonexistent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool not installed")
	})

	t.Run("empty tool name", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		err := installer.Uninstall("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool name cannot be empty")
	})
}

func TestInstaller_DownloadTool(t *testing.T) {
	t.Run("successful download", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		tool := &models.ToolInfo{
			Name:    "test-tool",
			Version: "1.0.0",
			Type:    models.ToolTypeAgent,
			File:    "tools/agents/test-tool.zip",
			Size:    100,
		}

		tempDir, err := os.MkdirTemp("", "download-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		destPath := filepath.Join(tempDir, "test.zip")
		err = installer.downloadTool(tool, destPath)
		assert.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(destPath)
		assert.NoError(t, err)
	})

	t.Run("download failure", func(t *testing.T) {
		installer, _, cleanup := setupTestInstaller(t)
		defer cleanup()

		// Replace GitHub client with one that errors
		githubClient := &mockGitHubDownloader{
			downloadError: fmt.Errorf("network error"),
		}
		installer.githubClient = githubClient

		tool := &models.ToolInfo{
			Name:    "test-tool",
			Version: "1.0.0",
			Type:    models.ToolTypeAgent,
			File:    "tools/agents/test-tool.zip",
			Size:    100,
		}

		tempDir, err := os.MkdirTemp("", "download-test-*")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)

		destPath := filepath.Join(tempDir, "test.zip")
		err = installer.downloadTool(tool, destPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "download failed")
	})
}

func TestInstaller_GetInstallPath(t *testing.T) {
	installer, baseDir, cleanup := setupTestInstaller(t)
	defer cleanup()

	tests := []struct {
		name     string
		toolName string
		toolType models.ToolType
		expected string
	}{
		{
			name:     "agent path",
			toolName: "test-agent",
			toolType: models.ToolTypeAgent,
			expected: filepath.Join(baseDir, "agents", "test-agent"),
		},
		{
			name:     "command path",
			toolName: "test-command",
			toolType: models.ToolTypeCommand,
			expected: filepath.Join(baseDir, "commands", "test-command"),
		},
		{
			name:     "skill path",
			toolName: "test-skill",
			toolType: models.ToolTypeSkill,
			expected: filepath.Join(baseDir, "skills", "test-skill"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := installer.getInstallPath(tt.toolName, tt.toolType)
			assert.Equal(t, tt.expected, path)
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"bytes", 500, "500 bytes"},
		{"kilobytes", 1536, "1.50 KB"},
		{"megabytes", 5242880, "5.00 MB"},
		{"gigabytes", 2147483648, "2.00 GB"},
		{"zero", 0, "0 bytes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBytes(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInstaller_GetInstalledTools(t *testing.T) {
	installer, _, cleanup := setupTestInstaller(t)
	defer cleanup()

	// Install some tools
	err := installer.Install("test-agent")
	require.NoError(t, err)

	tools, err := installer.GetInstalledTools()
	assert.NoError(t, err)
	assert.Len(t, tools, 1)
	assert.Contains(t, tools, "test-agent")
}

func TestInstaller_IsInstalled(t *testing.T) {
	installer, _, cleanup := setupTestInstaller(t)
	defer cleanup()

	// Not installed initially
	installed, err := installer.IsInstalled("test-agent")
	assert.NoError(t, err)
	assert.False(t, installed)

	// Install tool
	err = installer.Install("test-agent")
	require.NoError(t, err)

	// Should be installed now
	installed, err = installer.IsInstalled("test-agent")
	assert.NoError(t, err)
	assert.True(t, installed)
}

func TestInstaller_BuildDownloadURL(t *testing.T) {
	installer, _, cleanup := setupTestInstaller(t)
	defer cleanup()

	url := installer.buildDownloadURL("tools/agents/test.zip")
	assert.Contains(t, url, "raw.githubusercontent.com")
	assert.Contains(t, url, "tools/agents/test.zip")
	assert.Contains(t, url, "main") // branch
}

func TestInstaller_UpdateExistingTool(t *testing.T) {
	installer, baseDir, cleanup := setupTestInstaller(t)
	defer cleanup()

	// Install version 1.0.0
	err := installer.Install("test-agent")
	require.NoError(t, err)

	// Update registry to have version 2.0.0
	regService := installer.registryService.(*mockInstallerRegistryService)
	regService.tools["agent:test-agent"] = &models.ToolInfo{
		Name:        "test-agent",
		Version:     "2.0.0",
		Description: "Test agent",
		Type:        models.ToolTypeAgent,
		Author:      "test",
		File:        "tools/agents/test-agent.zip",
		Size:        1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Install again (should update)
	err = installer.Install("test-agent")
	assert.NoError(t, err)

	// Verify version updated
	version, err := installer.GetInstalledVersion("test-agent")
	assert.NoError(t, err)
	assert.Equal(t, "2.0.0", version)

	// Verify installation directory still exists
	destDir := filepath.Join(baseDir, "agents", "test-agent")
	_, err = os.Stat(destDir)
	assert.NoError(t, err)
}
