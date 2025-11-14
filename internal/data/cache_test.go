package data

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestRegistry creates a test registry for testing
func createTestRegistry() *models.Registry {
	return &models.Registry{
		Version:   "1.0.0",
		UpdatedAt: time.Now(),
		Tools: map[models.ToolType][]*models.ToolInfo{
			models.ToolTypeAgent: {
				{
					Name:        "test-agent",
					Version:     "1.0.0",
					Description: "A test agent",
					Type:        models.ToolTypeAgent,
					Author:      "test-author",
					Tags:        []string{"test", "agent"},
					File:        "agents/test-agent.zip",
					Size:        1024,
					Downloads:   100,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
		},
	}
}

func TestNewCacheManager(t *testing.T) {
	tests := []struct {
		name      string
		cacheDir  string
		ttl       time.Duration
		wantError bool
	}{
		{
			name:      "valid cache directory and TTL",
			cacheDir:  filepath.Join(t.TempDir(), "cache"),
			ttl:       1 * time.Hour,
			wantError: false,
		},
		{
			name:      "empty cache directory uses default",
			cacheDir:  "",
			ttl:       1 * time.Hour,
			wantError: false,
		},
		{
			name:      "zero TTL uses default",
			cacheDir:  filepath.Join(t.TempDir(), "cache"),
			ttl:       0,
			wantError: false,
		},
		{
			name:      "negative TTL uses default",
			cacheDir:  filepath.Join(t.TempDir(), "cache"),
			ttl:       -1 * time.Hour,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm, err := NewCacheManager(tt.cacheDir, tt.ttl)

			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, cm)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cm)

				if tt.cacheDir != "" {
					assert.Equal(t, tt.cacheDir, cm.GetCacheDir())
				}

				if tt.ttl > 0 {
					assert.Equal(t, tt.ttl, cm.GetTTL())
				} else {
					assert.Equal(t, DefaultCacheTTL, cm.GetTTL())
				}

				// Cleanup
				if cm != nil {
					_ = cm.Clear()
				}
			}
		})
	}
}

func TestCacheManager_SetRegistry(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("set valid registry", func(t *testing.T) {
		registry := createTestRegistry()
		err := cm.SetRegistry(registry)
		assert.NoError(t, err)

		// Verify registry file exists
		registryPath := filepath.Join(cacheDir, RegistryCacheFileName)
		assert.FileExists(t, registryPath)

		// Verify metadata file exists
		metadataPath := filepath.Join(cacheDir, MetadataFileName)
		assert.FileExists(t, metadataPath)
	})

	t.Run("set nil registry", func(t *testing.T) {
		err := cm.SetRegistry(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "registry cannot be nil")
	})

	t.Run("set invalid registry", func(t *testing.T) {
		invalidRegistry := &models.Registry{
			Version: "", // Invalid: empty version
			Tools:   map[models.ToolType][]*models.ToolInfo{},
		}
		err := cm.SetRegistry(invalidRegistry)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid registry")
	})
}

func TestCacheManager_GetRegistry(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("get cached registry", func(t *testing.T) {
		// Set registry
		registry := createTestRegistry()
		err := cm.SetRegistry(registry)
		require.NoError(t, err)

		// Get registry
		cachedRegistry, err := cm.GetRegistry()
		assert.NoError(t, err)
		assert.NotNil(t, cachedRegistry)
		assert.Equal(t, registry.Version, cachedRegistry.Version)
		assert.Equal(t, len(registry.Tools), len(cachedRegistry.Tools))
	})

	t.Run("get non-existent cache", func(t *testing.T) {
		// Clear cache first
		err := cm.Clear()
		require.NoError(t, err)

		// Try to get registry
		cachedRegistry, err := cm.GetRegistry()
		assert.Error(t, err)
		assert.Nil(t, cachedRegistry)
	})

	t.Run("get expired cache", func(t *testing.T) {
		// Create cache manager with very short TTL
		shortTTLCM, err := NewCacheManager(filepath.Join(t.TempDir(), "short-ttl-cache"), 100*time.Millisecond)
		require.NoError(t, err)
		defer shortTTLCM.Clear()

		// Set registry
		registry := createTestRegistry()
		err = shortTTLCM.SetRegistry(registry)
		require.NoError(t, err)

		// Wait for cache to expire
		time.Sleep(200 * time.Millisecond)

		// Try to get expired registry
		cachedRegistry, err := shortTTLCM.GetRegistry()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cache expired")
		assert.Nil(t, cachedRegistry)
	})
}

func TestCacheManager_IsValid(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("valid cache", func(t *testing.T) {
		// Set registry
		registry := createTestRegistry()
		err := cm.SetRegistry(registry)
		require.NoError(t, err)

		// Check if valid
		assert.True(t, cm.IsValid())
	})

	t.Run("no cache", func(t *testing.T) {
		// Clear cache
		err := cm.Clear()
		require.NoError(t, err)

		// Check if valid
		assert.False(t, cm.IsValid())
	})

	t.Run("expired cache", func(t *testing.T) {
		// Create cache manager with very short TTL
		shortTTLCM, err := NewCacheManager(filepath.Join(t.TempDir(), "expired-cache"), 100*time.Millisecond)
		require.NoError(t, err)
		defer shortTTLCM.Clear()

		// Set registry
		registry := createTestRegistry()
		err = shortTTLCM.SetRegistry(registry)
		require.NoError(t, err)

		// Initially valid
		assert.True(t, shortTTLCM.IsValid())

		// Wait for expiration
		time.Sleep(200 * time.Millisecond)

		// Should be invalid now
		assert.False(t, shortTTLCM.IsValid())
	})
}

func TestCacheManager_Invalidate(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("invalidate existing cache", func(t *testing.T) {
		// Set registry
		registry := createTestRegistry()
		err := cm.SetRegistry(registry)
		require.NoError(t, err)

		// Verify cache is valid
		assert.True(t, cm.IsValid())

		// Invalidate cache
		err = cm.Invalidate()
		assert.NoError(t, err)

		// Verify cache is no longer valid
		assert.False(t, cm.IsValid())

		// Verify files are removed
		registryPath := filepath.Join(cacheDir, RegistryCacheFileName)
		assert.NoFileExists(t, registryPath)

		metadataPath := filepath.Join(cacheDir, MetadataFileName)
		assert.NoFileExists(t, metadataPath)
	})

	t.Run("invalidate non-existent cache", func(t *testing.T) {
		// Invalidate when nothing is cached
		err := cm.Invalidate()
		assert.NoError(t, err) // Should not error
	})
}

func TestCacheManager_Clear(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)

	t.Run("clear cache", func(t *testing.T) {
		// Set registry
		registry := createTestRegistry()
		err := cm.SetRegistry(registry)
		require.NoError(t, err)

		// Verify cache directory exists
		_, err = os.Stat(cacheDir)
		assert.NoError(t, err)

		// Clear cache
		err = cm.Clear()
		assert.NoError(t, err)

		// Verify cache directory still exists but is empty
		_, err = os.Stat(cacheDir)
		assert.NoError(t, err)

		// Verify cache is not valid
		assert.False(t, cm.IsValid())
	})
}

func TestCacheManager_SetTTL(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("set valid TTL", func(t *testing.T) {
		newTTL := 2 * time.Hour
		cm.SetTTL(newTTL)
		assert.Equal(t, newTTL, cm.GetTTL())
	})

	t.Run("set zero TTL", func(t *testing.T) {
		originalTTL := cm.GetTTL()
		cm.SetTTL(0)
		assert.Equal(t, originalTTL, cm.GetTTL()) // Should not change
	})

	t.Run("set negative TTL", func(t *testing.T) {
		originalTTL := cm.GetTTL()
		cm.SetTTL(-1 * time.Hour)
		assert.Equal(t, originalTTL, cm.GetTTL()) // Should not change
	})
}

func TestCacheManager_GetMetadata(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("get metadata after caching", func(t *testing.T) {
		// Set registry
		registry := createTestRegistry()
		before := time.Now()
		err := cm.SetRegistry(registry)
		require.NoError(t, err)
		after := time.Now()

		// Get metadata
		metadata, err := cm.GetMetadata()
		assert.NoError(t, err)
		assert.NotNil(t, metadata)

		// Verify metadata fields
		assert.True(t, metadata.CachedAt.After(before) || metadata.CachedAt.Equal(before))
		assert.True(t, metadata.CachedAt.Before(after) || metadata.CachedAt.Equal(after))
		assert.Equal(t, cm.GetTTL(), metadata.TTL)
		assert.True(t, metadata.ExpiresAt.After(metadata.CachedAt))
	})

	t.Run("get metadata when no cache exists", func(t *testing.T) {
		// Clear cache
		err := cm.Clear()
		require.NoError(t, err)

		// Try to get metadata
		metadata, err := cm.GetMetadata()
		assert.Error(t, err)
		assert.Nil(t, metadata)
	})
}

func TestCacheManager_GetCacheSize(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("get cache size with cached data", func(t *testing.T) {
		// Set registry
		registry := createTestRegistry()
		err := cm.SetRegistry(registry)
		require.NoError(t, err)

		// Get cache size
		size, err := cm.GetCacheSize()
		assert.NoError(t, err)
		assert.Greater(t, size, int64(0))
	})

	t.Run("get cache size when empty", func(t *testing.T) {
		// Clear cache
		err := cm.Clear()
		require.NoError(t, err)

		// Get cache size
		size, err := cm.GetCacheSize()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), size)
	})
}

func TestCacheManager_Concurrency(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 1*time.Hour)
	require.NoError(t, err)
	defer cm.Clear()

	t.Run("concurrent reads and writes", func(t *testing.T) {
		registry := createTestRegistry()

		// Set initial registry
		err := cm.SetRegistry(registry)
		require.NoError(t, err)

		// Run concurrent operations
		done := make(chan bool)
		operations := 10

		// Concurrent reads
		for i := 0; i < operations; i++ {
			go func() {
				_, _ = cm.GetRegistry()
				done <- true
			}()
		}

		// Concurrent validity checks
		for i := 0; i < operations; i++ {
			go func() {
				_ = cm.IsValid()
				done <- true
			}()
		}

		// Wait for all operations
		for i := 0; i < operations*2; i++ {
			<-done
		}

		// Verify cache is still valid
		assert.True(t, cm.IsValid())
	})
}

func TestCacheManager_RealWorldScenario(t *testing.T) {
	cacheDir := filepath.Join(t.TempDir(), "cache")
	cm, err := NewCacheManager(cacheDir, 500*time.Millisecond)
	require.NoError(t, err)
	defer cm.Clear()

	registry := createTestRegistry()

	// First access: Cache miss
	_, err = cm.GetRegistry()
	assert.Error(t, err, "Should fail on first access")

	// Cache the registry
	err = cm.SetRegistry(registry)
	require.NoError(t, err)

	// Second access: Cache hit
	cachedReg, err := cm.GetRegistry()
	assert.NoError(t, err)
	assert.NotNil(t, cachedReg)
	assert.Equal(t, registry.Version, cachedReg.Version)

	// Verify it's valid
	assert.True(t, cm.IsValid())

	// Wait for expiration
	time.Sleep(600 * time.Millisecond)

	// Third access: Cache expired
	_, err = cm.GetRegistry()
	assert.Error(t, err, "Should fail after expiration")
	assert.False(t, cm.IsValid())

	// Re-cache
	err = cm.SetRegistry(registry)
	require.NoError(t, err)

	// Should work again
	cachedReg, err = cm.GetRegistry()
	assert.NoError(t, err)
	assert.NotNil(t, cachedReg)
	assert.True(t, cm.IsValid())

	// Manual invalidation
	err = cm.Invalidate()
	assert.NoError(t, err)

	// Should fail after invalidation
	_, err = cm.GetRegistry()
	assert.Error(t, err)
	assert.False(t, cm.IsValid())
}
