package services

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
	// DefaultLockFileVersion is the default version for new lock files
	DefaultLockFileVersion = "1.0"

	// LockFilePermission is the file permission for lock files
	LockFilePermission = 0644

	// TempFilePattern is the pattern for temporary files during atomic writes
	TempFilePattern = ".claude-lock-*.tmp"
)

// LockFileService manages the .claude-lock.json file
// It provides thread-safe CRUD operations for installed tools
type LockFileService struct {
	lockFilePath string
	mu           sync.RWMutex // For thread safety
}

// NewLockFileService creates a new LockFileService
func NewLockFileService(lockFilePath string) (*LockFileService, error) {
	if lockFilePath == "" {
		return nil, fmt.Errorf("lock file path cannot be empty")
	}

	return &LockFileService{
		lockFilePath: lockFilePath,
	}, nil
}

// GetLockFilePath returns the lock file path
func (lfs *LockFileService) GetLockFilePath() string {
	return lfs.lockFilePath
}

// Load loads the lock file from disk
// If the file doesn't exist, it creates a default lock file
func (lfs *LockFileService) Load() (*models.LockFile, error) {
	lfs.mu.RLock()
	defer lfs.mu.RUnlock()

	return lfs.loadUnsafe()
}

// loadUnsafe loads without acquiring lock (internal use only)
func (lfs *LockFileService) loadUnsafe() (*models.LockFile, error) {
	// Check if file exists
	data, err := os.ReadFile(lfs.lockFilePath)
	if os.IsNotExist(err) {
		// Create default lock file
		return lfs.createDefaultLockFile(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read lock file: %w", err)
	}

	// Parse JSON
	var lockFile models.LockFile
	if err := json.Unmarshal(data, &lockFile); err != nil {
		return nil, fmt.Errorf("failed to parse lock file: %w", err)
	}

	return &lockFile, nil
}

// Save saves the lock file to disk using atomic operations
// It writes to a temporary file first, then renames it to prevent corruption
func (lfs *LockFileService) Save(lockFile *models.LockFile) error {
	if lockFile == nil {
		return fmt.Errorf("lock file cannot be nil")
	}

	lfs.mu.Lock()
	defer lfs.mu.Unlock()

	return lfs.saveUnsafe(lockFile)
}

// saveUnsafe saves without acquiring lock (internal use only)
func (lfs *LockFileService) saveUnsafe(lockFile *models.LockFile) error {
	// Perform basic validation (skip full validation for empty lock files)
	if lockFile.Version == "" {
		return fmt.Errorf("lock file version cannot be empty")
	}
	if lockFile.Tools == nil {
		return fmt.Errorf("lock file tools cannot be nil")
	}

	// Validate all installed tools
	for name, tool := range lockFile.Tools {
		if name == "" {
			return fmt.Errorf("installed tool name cannot be empty")
		}
		if err := tool.Validate(); err != nil {
			return fmt.Errorf("invalid installed tool %s: %w", name, err)
		}
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(lockFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(lfs.lockFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create lock file directory: %w", err)
	}

	// Write to temporary file in the same directory (for atomic rename)
	tmpFile, err := os.CreateTemp(dir, TempFilePattern)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Ensure cleanup of temp file on error
	defer func() {
		tmpFile.Close()
		// Only remove if file still exists (not renamed)
		os.Remove(tmpPath)
	}()

	// Write data
	if _, err := tmpFile.Write(data); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	// Sync to disk
	if err := tmpFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync lock file: %w", err)
	}

	// Close before rename (required on Windows)
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, lfs.lockFilePath); err != nil {
		return fmt.Errorf("failed to rename lock file: %w", err)
	}

	return nil
}

// AddTool adds a tool to the lock file
func (lfs *LockFileService) AddTool(name string, tool *models.InstalledTool) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}

	lfs.mu.Lock()
	defer lfs.mu.Unlock()

	// Load current lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return fmt.Errorf("failed to load lock file: %w", err)
	}

	// Add tool using model's method
	if err := lockFile.AddTool(name, tool); err != nil {
		return err
	}

	// Save updated lock file
	if err := lfs.saveUnsafe(lockFile); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}

	return nil
}

// RemoveTool removes a tool from the lock file
func (lfs *LockFileService) RemoveTool(name string) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	lfs.mu.Lock()
	defer lfs.mu.Unlock()

	// Load current lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return fmt.Errorf("failed to load lock file: %w", err)
	}

	// Remove tool using model's method
	if err := lockFile.RemoveTool(name); err != nil {
		return err
	}

	// Save updated lock file
	if err := lfs.saveUnsafe(lockFile); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}

	return nil
}

// UpdateTool updates a tool in the lock file
func (lfs *LockFileService) UpdateTool(name string, tool *models.InstalledTool) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if tool == nil {
		return fmt.Errorf("tool cannot be nil")
	}

	lfs.mu.Lock()
	defer lfs.mu.Unlock()

	// Load current lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return fmt.Errorf("failed to load lock file: %w", err)
	}

	// Check if tool exists
	if _, exists := lockFile.Tools[name]; !exists {
		return fmt.Errorf("tool %s not found in lock file", name)
	}

	// Update tool using model's AddTool method (which updates if exists)
	if err := lockFile.AddTool(name, tool); err != nil {
		return err
	}

	// Save updated lock file
	if err := lfs.saveUnsafe(lockFile); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}

	return nil
}

// GetTool retrieves a tool from the lock file
func (lfs *LockFileService) GetTool(name string) (*models.InstalledTool, error) {
	if name == "" {
		return nil, fmt.Errorf("tool name cannot be empty")
	}

	lfs.mu.RLock()
	defer lfs.mu.RUnlock()

	// Load lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return nil, fmt.Errorf("failed to load lock file: %w", err)
	}

	// Get tool using model's method
	return lockFile.GetTool(name)
}

// ListTools returns all installed tools
func (lfs *LockFileService) ListTools() (map[string]*models.InstalledTool, error) {
	lfs.mu.RLock()
	defer lfs.mu.RUnlock()

	// Load lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return nil, fmt.Errorf("failed to load lock file: %w", err)
	}

	// Return copy of tools map to prevent external modification
	tools := make(map[string]*models.InstalledTool, len(lockFile.Tools))
	for name, tool := range lockFile.Tools {
		tools[name] = tool
	}

	return tools, nil
}

// IsInstalled checks if a tool is installed
func (lfs *LockFileService) IsInstalled(name string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf("tool name cannot be empty")
	}

	lfs.mu.RLock()
	defer lfs.mu.RUnlock()

	// Load lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return false, fmt.Errorf("failed to load lock file: %w", err)
	}

	_, exists := lockFile.Tools[name]
	return exists, nil
}

// SetRegistry sets the registry URL in the lock file
func (lfs *LockFileService) SetRegistry(registryURL string) error {
	if registryURL == "" {
		return fmt.Errorf("registry URL cannot be empty")
	}

	lfs.mu.Lock()
	defer lfs.mu.Unlock()

	// Load current lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return fmt.Errorf("failed to load lock file: %w", err)
	}

	// Update registry
	lockFile.Registry = registryURL
	lockFile.UpdatedAt = time.Now()

	// Save updated lock file
	if err := lfs.saveUnsafe(lockFile); err != nil {
		return fmt.Errorf("failed to save lock file: %w", err)
	}

	return nil
}

// GetRegistry returns the registry URL from the lock file
func (lfs *LockFileService) GetRegistry() (string, error) {
	lfs.mu.RLock()
	defer lfs.mu.RUnlock()

	// Load lock file
	lockFile, err := lfs.loadUnsafe()
	if err != nil {
		return "", fmt.Errorf("failed to load lock file: %w", err)
	}

	return lockFile.Registry, nil
}

// createDefaultLockFile creates a default lock file
func (lfs *LockFileService) createDefaultLockFile() *models.LockFile {
	return &models.LockFile{
		Version:   DefaultLockFileVersion,
		UpdatedAt: time.Now(),
		Registry:  "", // Will be set when first tool is installed
		Tools:     make(map[string]*models.InstalledTool),
	}
}
