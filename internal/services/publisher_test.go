package services

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPublisherService(t *testing.T) {
	tempDir := t.TempDir()
	fsManager, _ := data.NewFSManager(tempDir)
	githubClient := NewGitHubClient(GitHubClientConfig{
		Owner:  "test",
		Repo:   "test",
		Branch: "main",
	})
	cacheManager, _ := data.NewCacheManager(tempDir, 3600*time.Second)
	registryService := NewRegistryService(githubClient, cacheManager)

	tests := []struct {
		name        string
		fsManager   *data.FSManager
		github      *GitHubClient
		registry    *RegistryService
		config      *models.Config
		expectError bool
	}{
		{
			name:        "valid inputs",
			fsManager:   fsManager,
			github:      githubClient,
			registry:    registryService,
			config:      models.NewDefaultConfig(),
			expectError: false,
		},
		{
			name:        "nil fs manager",
			fsManager:   nil,
			github:      githubClient,
			registry:    registryService,
			config:      models.NewDefaultConfig(),
			expectError: true,
		},
		{
			name:        "nil github client",
			fsManager:   fsManager,
			github:      nil,
			registry:    registryService,
			config:      models.NewDefaultConfig(),
			expectError: true,
		},
		{
			name:        "nil registry service",
			fsManager:   fsManager,
			github:      githubClient,
			registry:    nil,
			config:      models.NewDefaultConfig(),
			expectError: true,
		},
		{
			name:        "nil config",
			fsManager:   fsManager,
			github:      githubClient,
			registry:    registryService,
			config:      nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewPublisherService(tt.fsManager, tt.github, tt.registry, tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, svc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, svc)
			}
		})
	}
}

func TestValidateTool(t *testing.T) {
	// Create test directories
	tempDir := t.TempDir()

	// Create valid agent directory
	validAgentDir := filepath.Join(tempDir, "agents", "test-agent")
	require.NoError(t, os.MkdirAll(validAgentDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(validAgentDir, "README.md"), []byte("# Test Agent"), 0644))

	// Create directory without README
	noReadmeDir := filepath.Join(tempDir, "agents", "no-readme")
	require.NoError(t, os.MkdirAll(noReadmeDir, 0755))

	// Create directory with sensitive files
	sensitiveDir := filepath.Join(tempDir, "agents", "sensitive")
	require.NoError(t, os.MkdirAll(sensitiveDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(sensitiveDir, "README.md"), []byte("# Test"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(sensitiveDir, ".env"), []byte("SECRET=123"), 0644))

	fsManager, _ := data.NewFSManager(tempDir)
	githubClient := NewGitHubClient(GitHubClientConfig{Owner: "test", Repo: "test", Branch: "main"})
	cacheManager, _ := data.NewCacheManager(tempDir, 3600*time.Second)
	registryService := NewRegistryService(githubClient, cacheManager)

	ps, err := NewPublisherService(
		fsManager,
		githubClient,
		registryService,
		models.NewDefaultConfig(),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		toolPath    string
		expectError bool
	}{
		{
			name:        "valid tool directory",
			toolPath:    validAgentDir,
			expectError: false,
		},
		{
			name:        "empty path",
			toolPath:    "",
			expectError: true,
		},
		{
			name:        "non-existent directory",
			toolPath:    filepath.Join(tempDir, "nonexistent"),
			expectError: true,
		},
		{
			name:        "no README",
			toolPath:    noReadmeDir,
			expectError: true,
		},
		{
			name:        "sensitive files",
			toolPath:    sensitiveDir,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ps.ValidateTool(tt.toolPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDetectToolType(t *testing.T) {
	tempDir := t.TempDir()

	fsManager, _ := data.NewFSManager(tempDir)
	githubClient := NewGitHubClient(GitHubClientConfig{Owner: "test", Repo: "test", Branch: "main"})
	cacheManager, _ := data.NewCacheManager(tempDir, 3600*time.Second)
	registryService := NewRegistryService(githubClient, cacheManager)

	ps, err := NewPublisherService(
		fsManager,
		githubClient,
		registryService,
		models.NewDefaultConfig(),
	)
	require.NoError(t, err)

	tests := []struct {
		name         string
		setupPath    string
		expectedType models.ToolType
		expectError  bool
	}{
		{
			name:         "agent directory",
			setupPath:    filepath.Join(tempDir, "agents", "test-agent"),
			expectedType: models.ToolTypeAgent,
			expectError:  false,
		},
		{
			name:         "command directory",
			setupPath:    filepath.Join(tempDir, "commands", "test-command"),
			expectedType: models.ToolTypeCommand,
			expectError:  false,
		},
		{
			name:         "skill directory",
			setupPath:    filepath.Join(tempDir, "skills", "test-skill"),
			expectedType: models.ToolTypeSkill,
			expectError:  false,
		},
		{
			name:        "unknown type",
			setupPath:   filepath.Join(tempDir, "unknown", "test-unknown"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(tt.setupPath, 0755))

			toolType, err := ps.detectToolType(tt.setupPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedType, toolType)
			}
		})
	}
}

func TestGenerateMetadata(t *testing.T) {
	tempDir := t.TempDir()
	toolPath := filepath.Join(tempDir, "test-tool")
	require.NoError(t, os.MkdirAll(toolPath, 0755))

	fsManager, _ := data.NewFSManager(tempDir)
	githubClient := NewGitHubClient(GitHubClientConfig{Owner: "test", Repo: "test", Branch: "main"})
	cacheManager, _ := data.NewCacheManager(tempDir, 3600*time.Second)
	registryService := NewRegistryService(githubClient, cacheManager)

	ps, err := NewPublisherService(
		fsManager,
		githubClient,
		registryService,
		models.NewDefaultConfig(),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		toolPath    string
		metadata    *PublishMetadata
		expectError bool
	}{
		{
			name:     "valid metadata",
			toolPath: toolPath,
			metadata: &PublishMetadata{
				Name:        "test-tool",
				Version:     "1.0.0",
				Description: "Test tool",
				Author:      "Test Author",
				Tags:        []string{"test", "tool"},
				Type:        models.ToolTypeAgent,
				Changelog: map[string]string{
					"1.0.0": "Initial release",
				},
			},
			expectError: false,
		},
		{
			name:        "empty path",
			toolPath:    "",
			metadata:    &PublishMetadata{},
			expectError: true,
		},
		{
			name:        "nil metadata",
			toolPath:    toolPath,
			metadata:    nil,
			expectError: true,
		},
		{
			name:     "missing name",
			toolPath: toolPath,
			metadata: &PublishMetadata{
				Version:     "1.0.0",
				Description: "Test",
				Author:      "Test",
			},
			expectError: true,
		},
		{
			name:     "missing version",
			toolPath: toolPath,
			metadata: &PublishMetadata{
				Name:        "test",
				Description: "Test",
				Author:      "Test",
			},
			expectError: true,
		},
		{
			name:     "missing author",
			toolPath: toolPath,
			metadata: &PublishMetadata{
				Name:        "test",
				Version:     "1.0.0",
				Description: "Test",
			},
			expectError: true,
		},
		{
			name:     "missing description",
			toolPath: toolPath,
			metadata: &PublishMetadata{
				Name:    "test",
				Version: "1.0.0",
				Author:  "Test",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ps.GenerateMetadata(tt.toolPath, tt.metadata)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify metadata.json was created
				metadataPath := filepath.Join(tt.toolPath, "metadata.json")
				assert.FileExists(t, metadataPath)

				// Verify content
				data, err := os.ReadFile(metadataPath)
				require.NoError(t, err)
				// Name is not in metadata.json - it's the directory name
				assert.Contains(t, string(data), tt.metadata.Author)
				assert.Contains(t, string(data), tt.metadata.Version)
			}
		})
	}
}

func TestCreatePackage(t *testing.T) {
	tempDir := t.TempDir()

	// Create a valid tool directory
	toolPath := filepath.Join(tempDir, "agents", "test-agent")
	require.NoError(t, os.MkdirAll(toolPath, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(toolPath, "README.md"), []byte("# Test"), 0644))

	outputPath := filepath.Join(tempDir, "output", "test-agent.zip")

	fsManager, _ := data.NewFSManager(tempDir)
	githubClient := NewGitHubClient(GitHubClientConfig{Owner: "test", Repo: "test", Branch: "main"})
	cacheManager, _ := data.NewCacheManager(tempDir, 3600*time.Second)
	registryService := NewRegistryService(githubClient, cacheManager)

	ps, err := NewPublisherService(
		fsManager,
		githubClient,
		registryService,
		models.NewDefaultConfig(),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		toolPath    string
		outputPath  string
		expectError bool
	}{
		{
			name:        "valid package creation",
			toolPath:    toolPath,
			outputPath:  outputPath,
			expectError: false,
		},
		{
			name:        "empty tool path",
			toolPath:    "",
			outputPath:  outputPath,
			expectError: true,
		},
		{
			name:        "empty output path",
			toolPath:    toolPath,
			outputPath:  "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := ps.CreatePackage(tt.toolPath, tt.outputPath)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			}
		})
	}
}

func TestReadExistingMetadata(t *testing.T) {
	tempDir := t.TempDir()

	// Create directory with metadata
	toolWithMetadata := filepath.Join(tempDir, "with-metadata")
	require.NoError(t, os.MkdirAll(toolWithMetadata, 0755))
	metadataJSON := `{
		"author": "Test Author",
		"tags": ["test"],
		"description": "Test description",
		"version": "1.0.0"
	}`
	require.NoError(t, os.WriteFile(
		filepath.Join(toolWithMetadata, "metadata.json"),
		[]byte(metadataJSON),
		0644,
	))

	// Create directory without metadata
	toolWithoutMetadata := filepath.Join(tempDir, "without-metadata")
	require.NoError(t, os.MkdirAll(toolWithoutMetadata, 0755))

	fsManager, _ := data.NewFSManager(tempDir)
	githubClient := NewGitHubClient(GitHubClientConfig{Owner: "test", Repo: "test", Branch: "main"})
	cacheManager, _ := data.NewCacheManager(tempDir, 3600*time.Second)
	registryService := NewRegistryService(githubClient, cacheManager)

	ps, err := NewPublisherService(
		fsManager,
		githubClient,
		registryService,
		models.NewDefaultConfig(),
	)
	require.NoError(t, err)

	tests := []struct {
		name        string
		toolPath    string
		expectNil   bool
		expectError bool
	}{
		{
			name:        "existing metadata",
			toolPath:    toolWithMetadata,
			expectNil:   false,
			expectError: false,
		},
		{
			name:        "no metadata",
			toolPath:    toolWithoutMetadata,
			expectNil:   true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata, err := ps.ReadExistingMetadata(tt.toolPath)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectNil {
					assert.Nil(t, metadata)
				} else {
					assert.NotNil(t, metadata)
					assert.Equal(t, "Test Author", metadata.Author)
					assert.Equal(t, "1.0.0", metadata.Version)
				}
			}
		})
	}
}
