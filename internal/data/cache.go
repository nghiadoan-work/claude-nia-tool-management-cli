package data

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
)

const (
	// DefaultCacheTTL is the default time-to-live for cache entries (1 hour)
	DefaultCacheTTL = 1 * time.Hour

	// CacheDirName is the name of the cache directory
	CacheDirName = ".claude-tools-cache"

	// RegistryCacheFileName is the name of the cached registry file
	RegistryCacheFileName = "registry.json"

	// MetadataFileName is the name of the cache metadata file
	MetadataFileName = "metadata.json"
)

// CacheMetadata stores metadata about cached data
type CacheMetadata struct {
	CachedAt  time.Time     `json:"cached_at"`
	ExpiresAt time.Time     `json:"expires_at"`
	TTL       time.Duration `json:"ttl"`
	ETag      string        `json:"etag,omitempty"` // For HTTP cache validation
}

// CacheManager manages local caching of registry data
type CacheManager struct {
	cacheDir string
	ttl      time.Duration
	mu       sync.RWMutex
}

// NewCacheManager creates a new CacheManager
func NewCacheManager(cacheDir string, ttl time.Duration) (*CacheManager, error) {
	if cacheDir == "" {
		// Use default cache directory in user's home
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		cacheDir = filepath.Join(homeDir, CacheDirName)
	}

	if ttl <= 0 {
		ttl = DefaultCacheTTL
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &CacheManager{
		cacheDir: cacheDir,
		ttl:      ttl,
	}, nil
}

// GetRegistry retrieves the cached registry if it exists and is not expired
func (cm *CacheManager) GetRegistry() (*models.Registry, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Check if cache is valid
	metadata, err := cm.getMetadata()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache metadata: %w", err)
	}

	// Check if cache is expired
	if time.Now().After(metadata.ExpiresAt) {
		return nil, fmt.Errorf("cache expired")
	}

	// Read cached registry
	registryPath := filepath.Join(cm.cacheDir, RegistryCacheFileName)
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cached registry: %w", err)
	}

	var registry models.Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached registry: %w", err)
	}

	return &registry, nil
}

// SetRegistry caches the registry data with TTL-based expiration
func (cm *CacheManager) SetRegistry(registry *models.Registry) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if registry == nil {
		return fmt.Errorf("registry cannot be nil")
	}

	// Validate registry before caching
	if err := registry.Validate(); err != nil {
		return fmt.Errorf("invalid registry: %w", err)
	}

	// Marshal registry to JSON
	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	// Write registry to cache
	registryPath := filepath.Join(cm.cacheDir, RegistryCacheFileName)
	if err := os.WriteFile(registryPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry cache: %w", err)
	}

	// Create and save metadata
	now := time.Now()
	metadata := &CacheMetadata{
		CachedAt:  now,
		ExpiresAt: now.Add(cm.ttl),
		TTL:       cm.ttl,
	}

	if err := cm.saveMetadata(metadata); err != nil {
		return fmt.Errorf("failed to save cache metadata: %w", err)
	}

	return nil
}

// IsValid checks if the cache exists and is not expired
func (cm *CacheManager) IsValid() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	metadata, err := cm.getMetadata()
	if err != nil {
		return false
	}

	// Check if cache is expired
	return time.Now().Before(metadata.ExpiresAt)
}

// Invalidate removes the cached registry and metadata
func (cm *CacheManager) Invalidate() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Remove registry cache file
	registryPath := filepath.Join(cm.cacheDir, RegistryCacheFileName)
	if err := os.Remove(registryPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove registry cache: %w", err)
	}

	// Remove metadata file
	metadataPath := filepath.Join(cm.cacheDir, MetadataFileName)
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove metadata: %w", err)
	}

	return nil
}

// Clear removes all cache files and the cache directory
func (cm *CacheManager) Clear() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Remove entire cache directory
	if err := os.RemoveAll(cm.cacheDir); err != nil {
		return fmt.Errorf("failed to clear cache directory: %w", err)
	}

	// Recreate cache directory
	if err := os.MkdirAll(cm.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate cache directory: %w", err)
	}

	return nil
}

// GetMetadata returns the cache metadata
func (cm *CacheManager) GetMetadata() (*CacheMetadata, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.getMetadata()
}

// getMetadata is the internal (non-locking) version of GetMetadata
func (cm *CacheManager) getMetadata() (*CacheMetadata, error) {
	metadataPath := filepath.Join(cm.cacheDir, MetadataFileName)
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata CacheMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return &metadata, nil
}

// saveMetadata is the internal (non-locking) version for saving metadata
func (cm *CacheManager) saveMetadata(metadata *CacheMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataPath := filepath.Join(cm.cacheDir, MetadataFileName)
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// SetTTL updates the TTL for future cache entries
func (cm *CacheManager) SetTTL(ttl time.Duration) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if ttl > 0 {
		cm.ttl = ttl
	}
}

// GetTTL returns the current TTL setting
func (cm *CacheManager) GetTTL() time.Duration {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.ttl
}

// GetCacheDir returns the cache directory path
func (cm *CacheManager) GetCacheDir() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return cm.cacheDir
}

// GetCacheSize returns the total size of cached data in bytes
func (cm *CacheManager) GetCacheSize() (int64, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var totalSize int64

	err := filepath.Walk(cm.cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate cache size: %w", err)
	}

	return totalSize, nil
}
