# FSManager API Reference

## Overview

FSManager provides secure file system operations for the cntm tool, specifically focused on safe ZIP archive handling, path validation, and integrity verification.

**Package**: `github.com/nghiadt/claude-nia-tool-management-cli/internal/data`

## Type Definition

```go
type FSManager struct {
    baseDir             string
    maxUncompressedSize int64
    maxFiles            int
    maxCompressionRatio float64
}
```

## Constructor

### NewFSManager

Creates a new FSManager instance with default security settings.

```go
func NewFSManager(baseDir string) (*FSManager, error)
```

**Parameters**:
- `baseDir` - Base directory for all operations (will be created if doesn't exist)

**Returns**:
- `*FSManager` - Configured FSManager instance
- `error` - Error if base directory creation fails or path is invalid

**Example**:
```go
fsm, err := data.NewFSManager("/home/user/.claude")
if err != nil {
    log.Fatal(err)
}
```

**Default Security Settings**:
- Max uncompressed size: 1GB
- Max files: 10,000
- Max compression ratio: 100:1
- Max single file size: 500MB

## ZIP Operations

### ExtractZIP

Extracts a ZIP archive to a destination directory with comprehensive security checks.

```go
func (fs *FSManager) ExtractZIP(zipPath, destPath string) error
```

**Parameters**:
- `zipPath` - Path to ZIP file to extract
- `destPath` - Destination directory (must be within base directory)

**Returns**:
- `error` - Error if extraction fails or security validation fails

**Security Checks**:
- Pre-scans entire ZIP before extraction
- Validates all paths for traversal attempts
- Checks file count and size limits
- Verifies compression ratios
- Rejects symlinks
- Ensures destination is within base directory

**Example**:
```go
err := fsm.ExtractZIP("/tmp/tool.zip", "/home/user/.claude/agents/mytool")
if err != nil {
    log.Printf("Extraction failed: %v", err)
}
```

**Possible Errors**:
- Empty ZIP file
- Too many files (>10,000)
- Total size exceeds limit (>1GB)
- Compression ratio too high (>100:1)
- Path traversal detected
- Invalid destination path

### CreateZIP

Creates a ZIP archive from a source directory.

```go
func (fs *FSManager) CreateZIP(srcPath, zipPath string) error
```

**Parameters**:
- `srcPath` - Source directory to archive
- `zipPath` - Path where ZIP file will be created

**Returns**:
- `error` - Error if creation fails

**Behavior**:
- Excludes hidden files (starting with `.`)
- Uses Deflate compression
- Preserves directory structure
- Creates ZIP with forward-slash separators (cross-platform)

**Example**:
```go
err := fsm.CreateZIP("/home/user/mytool", "/tmp/mytool.zip")
if err != nil {
    log.Printf("ZIP creation failed: %v", err)
}
```

## Path Validation

### ValidatePath

Ensures a path is within the base directory, preventing path traversal.

```go
func (fs *FSManager) ValidatePath(path string) error
```

**Parameters**:
- `path` - Path to validate

**Returns**:
- `error` - Error if path is outside base directory

**Example**:
```go
if err := fsm.ValidatePath("/home/user/.claude/agents/mytool"); err != nil {
    log.Printf("Invalid path: %v", err)
}
```

**Rejects**:
- Paths outside base directory
- Paths with `..` that escape base
- Absolute paths pointing elsewhere

## Integrity Verification

### CalculateSHA256

Calculates the SHA256 hash of a file.

```go
func (fs *FSManager) CalculateSHA256(filePath string) (string, error)
```

**Parameters**:
- `filePath` - Path to file to hash

**Returns**:
- `string` - Hexadecimal SHA256 hash
- `error` - Error if file cannot be read

**Example**:
```go
hash, err := fsm.CalculateSHA256("/tmp/tool.zip")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("SHA256: %s\n", hash)
```

### VerifyIntegrity

Verifies a file's SHA256 hash matches the expected value.

```go
func (fs *FSManager) VerifyIntegrity(filePath, expectedHash string) error
```

**Parameters**:
- `filePath` - Path to file to verify
- `expectedHash` - Expected SHA256 hash (case-insensitive)

**Returns**:
- `error` - Error if hash doesn't match or file cannot be read

**Example**:
```go
err := fsm.VerifyIntegrity("/tmp/tool.zip", "abc123...")
if err != nil {
    log.Printf("Integrity check failed: %v", err)
}
```

## Directory Operations

### EnsureDir

Creates a directory and all parent directories if they don't exist.

```go
func (fs *FSManager) EnsureDir(path string) error
```

**Parameters**:
- `path` - Directory path to create (must be within base directory)

**Returns**:
- `error` - Error if creation fails or path is invalid

**Example**:
```go
err := fsm.EnsureDir("/home/user/.claude/agents/mytool")
```

**Behavior**:
- Creates all parent directories (like `mkdir -p`)
- Sets permissions to 0755
- Validates path is within base directory

### RemoveDir

Removes a directory and all its contents.

```go
func (fs *FSManager) RemoveDir(path string) error
```

**Parameters**:
- `path` - Directory path to remove (must be within base directory)

**Returns**:
- `error` - Error if removal fails or path is invalid

**Example**:
```go
err := fsm.RemoveDir("/home/user/.claude/agents/mytool")
```

**Warning**: Recursively deletes all contents. Use with caution.

### GetDirSize

Calculates the total size of a directory and all its contents.

```go
func (fs *FSManager) GetDirSize(path string) (int64, error)
```

**Parameters**:
- `path` - Directory path (must be within base directory)

**Returns**:
- `int64` - Total size in bytes
- `error` - Error if calculation fails or path is invalid

**Example**:
```go
size, err := fsm.GetDirSize("/home/user/.claude/agents/mytool")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Size: %d bytes\n", size)
```

## Configuration Methods

### GetBaseDir

Returns the base directory path.

```go
func (fs *FSManager) GetBaseDir() string
```

**Returns**:
- `string` - Absolute path to base directory

**Example**:
```go
baseDir := fsm.GetBaseDir()
fmt.Printf("Base: %s\n", baseDir)
```

### SetMaxUncompressedSize

Sets the maximum allowed uncompressed size for ZIP files.

```go
func (fs *FSManager) SetMaxUncompressedSize(size int64)
```

**Parameters**:
- `size` - Maximum size in bytes (must be > 0)

**Example**:
```go
fsm.SetMaxUncompressedSize(500 * 1024 * 1024) // 500MB
```

### SetMaxFiles

Sets the maximum number of files allowed in a ZIP archive.

```go
func (fs *FSManager) SetMaxFiles(count int)
```

**Parameters**:
- `count` - Maximum file count (must be > 0)

**Example**:
```go
fsm.SetMaxFiles(5000)
```

### SetMaxCompressionRatio

Sets the maximum compression ratio allowed.

```go
func (fs *FSManager) SetMaxCompressionRatio(ratio float64)
```

**Parameters**:
- `ratio` - Maximum ratio (e.g., 100.0 for 100:1) (must be > 0)

**Example**:
```go
fsm.SetMaxCompressionRatio(50.0) // 50:1 ratio
```

## Constants

```go
const (
    MaxUncompressedSize int64 = 1024 * 1024 * 1024  // 1GB
    MaxFiles            = 10000                      // 10,000 files
    MaxCompressionRatio = 100.0                      // 100:1 ratio
    MaxSingleFileSize   int64 = 500 * 1024 * 1024   // 500MB
    DefaultDirPerm      = 0755                       // Directory permissions
    DefaultFilePerm     = 0644                       // File permissions
)
```

## Error Types

FSManager returns wrapped errors with context:

```go
// Path validation errors
"path %s is outside base directory %s"
"failed to get absolute path: %w"

// ZIP validation errors
"ZIP file is empty"
"ZIP contains too many files (%d), maximum allowed: %d"
"total uncompressed size (%d bytes) exceeds maximum (%d bytes)"
"compression ratio (%.2f:1) exceeds maximum (%.2f:1), possible ZIP bomb"

// Path traversal errors
"absolute paths are not allowed in ZIP files: %s"
"path traversal detected in ZIP entry: %s"
"path traversal detected: %s escapes destination %s"

// Symlink errors
"symlinks are not allowed in ZIP files: %s"

// Integrity errors
"integrity check failed: expected %s, got %s"
```

## Complete Usage Example

```go
package main

import (
    "log"
    "github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
)

func main() {
    // Initialize FSManager
    fsm, err := data.NewFSManager("/home/user/.claude")
    if err != nil {
        log.Fatal(err)
    }

    // Configure security limits (optional)
    fsm.SetMaxUncompressedSize(500 * 1024 * 1024) // 500MB
    fsm.SetMaxFiles(5000)

    // Create a tool directory
    toolDir := "/home/user/.claude/agents/code-reviewer"
    if err := fsm.EnsureDir(toolDir); err != nil {
        log.Fatal(err)
    }

    // Create ZIP from directory
    zipPath := "/tmp/code-reviewer.zip"
    if err := fsm.CreateZIP(toolDir, zipPath); err != nil {
        log.Fatal(err)
    }

    // Calculate and verify integrity
    hash, err := fsm.CalculateSHA256(zipPath)
    if err != nil {
        log.Fatal(err)
    }

    if err := fsm.VerifyIntegrity(zipPath, hash); err != nil {
        log.Fatal(err)
    }

    // Extract ZIP (with security checks)
    extractDir := "/home/user/.claude/agents/code-reviewer-v2"
    if err := fsm.ExtractZIP(zipPath, extractDir); err != nil {
        log.Fatal(err)
    }

    // Get directory size
    size, err := fsm.GetDirSize(extractDir)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Extracted %d bytes\n", size)
}
```

## Security Best Practices

1. **Always verify integrity** before extracting downloaded ZIPs
2. **Use ExtractZIP** instead of manual extraction (has built-in security)
3. **Don't disable security limits** unless absolutely necessary
4. **Validate all user-provided paths** using ValidatePath()
5. **Keep base directory restrictive** (ideally in user's home)
6. **Log security violations** for monitoring

## Thread Safety

FSManager is **NOT** thread-safe. If you need concurrent access:

```go
type SafeFSManager struct {
    fsm *data.FSManager
    mu  sync.Mutex
}

func (s *SafeFSManager) ExtractZIP(zipPath, destPath string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.fsm.ExtractZIP(zipPath, destPath)
}
```

## Performance Considerations

- **Streaming operations**: Files are not buffered entirely in memory
- **Pre-scan overhead**: ZIP validation adds ~10-50ms for typical tools
- **Hash calculation**: O(n) where n is file size, typically fast
- **Directory walking**: Efficient single-pass implementation

## Integration with cntm

FSManager will be used by:
- **InstallerService**: Extract downloaded tool ZIPs
- **PublisherService**: Create tool ZIPs for upload
- **UpdaterService**: Extract updated tools
- **LockFileService**: Manage lock file operations

Typical flow:
```
Download ZIP → VerifyIntegrity → ExtractZIP → Update Lock File
```
