package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLockFileService(t *testing.T) {
	tests := []struct {
		name         string
		lockFilePath string
		wantErr      bool
	}{
		{
			name:         "valid path",
			lockFilePath: "/tmp/test.lock.json",
			wantErr:      false,
		},
		{
			name:         "empty path",
			lockFilePath: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewLockFileService(tt.lockFilePath)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, svc)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, svc)
				assert.Equal(t, tt.lockFilePath, svc.GetLockFilePath())
			}
		})
	}
}

func TestLockFileService_Load(t *testing.T) {
	t.Run("load non-existent file creates default", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		lockFile, err := svc.Load()
		require.NoError(t, err)
		assert.NotNil(t, lockFile)
		assert.Equal(t, "1.0", lockFile.Version)
		assert.NotNil(t, lockFile.Tools)
		assert.Empty(t, lockFile.Tools)
		assert.NotZero(t, lockFile.UpdatedAt)
	})

	t.Run("load existing lock file", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		// Create a lock file
		expectedLock := &models.LockFile{
			Version:   "1.0",
			UpdatedAt: time.Now(),
			Registry:  "https://github.com/nghiadoan-work/claude-tools-registry",
			Tools: map[string]*models.InstalledTool{
				"code-reviewer": {
					Version:     "1.0.0",
					Type:        models.ToolTypeAgent,
					InstalledAt: time.Now(),
					Source:      "registry",
					Integrity:   "sha256-abc123",
				},
			},
		}

		data, err := json.MarshalIndent(expectedLock, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(lockPath, data, 0644))

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		lockFile, err := svc.Load()
		require.NoError(t, err)
		assert.Equal(t, expectedLock.Version, lockFile.Version)
		assert.Equal(t, expectedLock.Registry, lockFile.Registry)
		assert.Len(t, lockFile.Tools, 1)
		assert.Contains(t, lockFile.Tools, "code-reviewer")
	})

	t.Run("load invalid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		// Write invalid JSON
		require.NoError(t, os.WriteFile(lockPath, []byte("{invalid json}"), 0644))

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		_, err = svc.Load()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse")
	})
}

func TestLockFileService_Save(t *testing.T) {
	t.Run("save lock file creates valid JSON", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		lockFile := &models.LockFile{
			Version:   "1.0",
			UpdatedAt: time.Now(),
			Registry:  "https://github.com/nghiadoan-work/claude-tools-registry",
			Tools: map[string]*models.InstalledTool{
				"code-reviewer": {
					Version:     "1.0.0",
					Type:        models.ToolTypeAgent,
					InstalledAt: time.Now(),
					Source:      "registry",
					Integrity:   "sha256-abc123",
				},
			},
		}

		err = svc.Save(lockFile)
		require.NoError(t, err)

		// Verify file exists
		_, err = os.Stat(lockPath)
		assert.NoError(t, err)

		// Verify JSON is valid and pretty-printed
		data, err := os.ReadFile(lockPath)
		require.NoError(t, err)

		var loaded models.LockFile
		err = json.Unmarshal(data, &loaded)
		require.NoError(t, err)
		assert.Equal(t, lockFile.Version, loaded.Version)
		assert.Equal(t, lockFile.Registry, loaded.Registry)

		// Check indentation (pretty-print)
		assert.Contains(t, string(data), "\n  ")
	})

	t.Run("save with nil lock file", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		err = svc.Save(nil)
		assert.Error(t, err)
	})

	t.Run("atomic write - temp file pattern", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		lockFile := &models.LockFile{
			Version:   "1.0",
			UpdatedAt: time.Now(),
			Registry:  "https://github.com/nghiadoan-work/claude-tools-registry",
			Tools:     make(map[string]*models.InstalledTool),
		}

		err = svc.Save(lockFile)
		require.NoError(t, err)

		// After save, no temp files should exist
		files, err := os.ReadDir(tmpDir)
		require.NoError(t, err)

		for _, file := range files {
			assert.NotContains(t, file.Name(), ".tmp")
			assert.NotContains(t, file.Name(), ".temp")
		}
	})
}

func TestLockFileService_AddTool(t *testing.T) {
	t.Run("add tool to empty lock file", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		tool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}

		err = svc.AddTool("code-reviewer", tool)
		require.NoError(t, err)

		// Verify tool was added
		lockFile, err := svc.Load()
		require.NoError(t, err)
		assert.Len(t, lockFile.Tools, 1)
		assert.Contains(t, lockFile.Tools, "code-reviewer")
	})

	t.Run("add multiple tools", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		tool1 := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}

		tool2 := &models.InstalledTool{
			Version:     "2.0.0",
			Type:        models.ToolTypeCommand,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-def456",
		}

		require.NoError(t, svc.AddTool("code-reviewer", tool1))
		require.NoError(t, svc.AddTool("git-helper", tool2))

		lockFile, err := svc.Load()
		require.NoError(t, err)
		assert.Len(t, lockFile.Tools, 2)
		assert.Contains(t, lockFile.Tools, "code-reviewer")
		assert.Contains(t, lockFile.Tools, "git-helper")
	})

	t.Run("add tool with empty name", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		tool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}

		err = svc.AddTool("", tool)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("add invalid tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		// Tool with missing version
		tool := &models.InstalledTool{
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}

		err = svc.AddTool("invalid-tool", tool)
		assert.Error(t, err)
	})
}

func TestLockFileService_RemoveTool(t *testing.T) {
	t.Run("remove existing tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		// Add a tool first
		tool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}
		require.NoError(t, svc.AddTool("code-reviewer", tool))

		// Remove the tool
		err = svc.RemoveTool("code-reviewer")
		require.NoError(t, err)

		// Verify tool was removed
		lockFile, err := svc.Load()
		require.NoError(t, err)
		assert.Empty(t, lockFile.Tools)
	})

	t.Run("remove non-existent tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		err = svc.RemoveTool("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("remove with empty name", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		err = svc.RemoveTool("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})
}

func TestLockFileService_UpdateTool(t *testing.T) {
	t.Run("update existing tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		// Add initial tool
		tool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}
		require.NoError(t, svc.AddTool("code-reviewer", tool))

		// Update tool
		updatedTool := &models.InstalledTool{
			Version:     "2.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-xyz789",
		}
		err = svc.UpdateTool("code-reviewer", updatedTool)
		require.NoError(t, err)

		// Verify update
		retrieved, err := svc.GetTool("code-reviewer")
		require.NoError(t, err)
		assert.Equal(t, "2.0.0", retrieved.Version)
		assert.Equal(t, "sha256-xyz789", retrieved.Integrity)
	})

	t.Run("update non-existent tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		tool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}

		err = svc.UpdateTool("non-existent", tool)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestLockFileService_GetTool(t *testing.T) {
	t.Run("get existing tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		expectedTool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}
		require.NoError(t, svc.AddTool("code-reviewer", expectedTool))

		tool, err := svc.GetTool("code-reviewer")
		require.NoError(t, err)
		assert.Equal(t, expectedTool.Version, tool.Version)
		assert.Equal(t, expectedTool.Type, tool.Type)
		assert.Equal(t, expectedTool.Integrity, tool.Integrity)
	})

	t.Run("get non-existent tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		_, err = svc.GetTool("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("get with empty name", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		_, err = svc.GetTool("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})
}

func TestLockFileService_ListTools(t *testing.T) {
	t.Run("list empty lock file", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		tools, err := svc.ListTools()
		require.NoError(t, err)
		assert.NotNil(t, tools)
		assert.Empty(t, tools)
	})

	t.Run("list multiple tools", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		tool1 := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}

		tool2 := &models.InstalledTool{
			Version:     "2.0.0",
			Type:        models.ToolTypeCommand,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-def456",
		}

		require.NoError(t, svc.AddTool("code-reviewer", tool1))
		require.NoError(t, svc.AddTool("git-helper", tool2))

		tools, err := svc.ListTools()
		require.NoError(t, err)
		assert.Len(t, tools, 2)
		assert.Contains(t, tools, "code-reviewer")
		assert.Contains(t, tools, "git-helper")
	})
}

func TestLockFileService_IsInstalled(t *testing.T) {
	t.Run("check installed tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		tool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-abc123",
		}
		require.NoError(t, svc.AddTool("code-reviewer", tool))

		installed, err := svc.IsInstalled("code-reviewer")
		require.NoError(t, err)
		assert.True(t, installed)
	})

	t.Run("check non-installed tool", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		installed, err := svc.IsInstalled("non-existent")
		require.NoError(t, err)
		assert.False(t, installed)
	})

	t.Run("check with empty name", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		_, err = svc.IsInstalled("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})
}

func TestLockFileService_ConcurrentAccess(t *testing.T) {
	t.Run("concurrent reads and writes", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		// Add initial tool
		initialTool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-initial",
		}
		require.NoError(t, svc.AddTool("initial-tool", initialTool))

		var wg sync.WaitGroup
		numGoroutines := 10

		// Concurrent writes
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				tool := &models.InstalledTool{
					Version:     "1.0.0",
					Type:        models.ToolTypeAgent,
					InstalledAt: time.Now(),
					Source:      "registry",
					Integrity:   "sha256-abc123",
				}

				err := svc.AddTool(fmt.Sprintf("tool-%d", idx), tool)
				assert.NoError(t, err)
			}(i)
		}

		// Concurrent reads
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				_, err := svc.ListTools()
				assert.NoError(t, err)
			}()
		}

		wg.Wait()

		// Verify all tools were added
		tools, err := svc.ListTools()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(tools), numGoroutines)
	})

	t.Run("concurrent updates", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		// Add initial tool
		initialTool := &models.InstalledTool{
			Version:     "1.0.0",
			Type:        models.ToolTypeAgent,
			InstalledAt: time.Now(),
			Source:      "registry",
			Integrity:   "sha256-initial",
		}
		require.NoError(t, svc.AddTool("test-tool", initialTool))

		var wg sync.WaitGroup
		numUpdates := 5

		for i := 0; i < numUpdates; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				tool := &models.InstalledTool{
					Version:     fmt.Sprintf("1.0.%d", idx),
					Type:        models.ToolTypeAgent,
					InstalledAt: time.Now(),
					Source:      "registry",
					Integrity:   fmt.Sprintf("sha256-%d", idx),
				}

				err := svc.UpdateTool("test-tool", tool)
				assert.NoError(t, err)
			}(i)
		}

		wg.Wait()

		// Verify tool still exists and is valid
		tool, err := svc.GetTool("test-tool")
		require.NoError(t, err)
		assert.NotNil(t, tool)
		assert.NotEmpty(t, tool.Version)
	})
}

func TestLockFileService_SetRegistry(t *testing.T) {
	t.Run("set registry URL", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		registryURL := "https://github.com/nghiadoan-work/claude-tools-registry"
		err = svc.SetRegistry(registryURL)
		require.NoError(t, err)

		lockFile, err := svc.Load()
		require.NoError(t, err)
		assert.Equal(t, registryURL, lockFile.Registry)
	})

	t.Run("set empty registry", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		err = svc.SetRegistry("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "registry URL cannot be empty")
	})
}

func TestLockFileService_GetRegistry(t *testing.T) {
	t.Run("get registry from lock file", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		registryURL := "https://github.com/nghiadoan-work/claude-tools-registry"
		require.NoError(t, svc.SetRegistry(registryURL))

		url, err := svc.GetRegistry()
		require.NoError(t, err)
		assert.Equal(t, registryURL, url)
	})

	t.Run("get registry from new lock file", func(t *testing.T) {
		tmpDir := t.TempDir()
		lockPath := filepath.Join(tmpDir, ".claude-lock.json")

		svc, err := NewLockFileService(lockPath)
		require.NoError(t, err)

		url, err := svc.GetRegistry()
		require.NoError(t, err)
		// Default should be empty
		assert.Empty(t, url)
	})
}
