package data

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// MaxUncompressedSize is the maximum allowed uncompressed size (1GB)
	MaxUncompressedSize int64 = 1024 * 1024 * 1024 // 1GB

	// MaxFiles is the maximum number of files allowed in a ZIP
	MaxFiles = 10000

	// MaxCompressionRatio is the maximum compression ratio allowed (100:1)
	// A higher ratio could indicate a ZIP bomb
	MaxCompressionRatio = 100.0

	// MaxSingleFileSize is the maximum size of a single uncompressed file (500MB)
	MaxSingleFileSize int64 = 500 * 1024 * 1024 // 500MB

	// DefaultDirPerm is the default permission for created directories
	DefaultDirPerm = 0755

	// DefaultFilePerm is the default permission for extracted files
	DefaultFilePerm = 0644
)

// FSManager handles file system operations for tool installation
type FSManager struct {
	baseDir             string
	maxUncompressedSize int64
	maxFiles            int
	maxCompressionRatio float64
}

// NewFSManager creates a new FSManager with default security settings
func NewFSManager(baseDir string) (*FSManager, error) {
	if baseDir == "" {
		return nil, fmt.Errorf("base directory cannot be empty")
	}

	// Ensure base directory exists
	if err := os.MkdirAll(baseDir, DefaultDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	// Clean and convert to absolute path
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	return &FSManager{
		baseDir:             absBaseDir,
		maxUncompressedSize: MaxUncompressedSize,
		maxFiles:            MaxFiles,
		maxCompressionRatio: MaxCompressionRatio,
	}, nil
}

// ExtractZIP extracts a ZIP file to the destination path with security checks
func (fs *FSManager) ExtractZIP(zipPath, destPath string) error {
	// Validate inputs
	if zipPath == "" {
		return fmt.Errorf("zip path cannot be empty")
	}
	if destPath == "" {
		return fmt.Errorf("destination path cannot be empty")
	}

	// Ensure destination is within base directory
	if err := fs.ValidatePath(destPath); err != nil {
		return fmt.Errorf("invalid destination path: %w", err)
	}

	// Open ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer reader.Close()

	// Pre-scan ZIP for security threats
	if err := fs.validateZIPContents(&reader.Reader); err != nil {
		return fmt.Errorf("ZIP validation failed: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(destPath, DefaultDirPerm); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract files
	for _, file := range reader.File {
		if err := fs.extractFile(file, destPath); err != nil {
			return fmt.Errorf("failed to extract file %s: %w", file.Name, err)
		}
	}

	return nil
}

// validateZIPContents performs security checks on ZIP contents before extraction
func (fs *FSManager) validateZIPContents(reader *zip.Reader) error {
	if len(reader.File) == 0 {
		return fmt.Errorf("ZIP file is empty")
	}

	if len(reader.File) > fs.maxFiles {
		return fmt.Errorf("ZIP contains too many files (%d), maximum allowed: %d", len(reader.File), fs.maxFiles)
	}

	var totalUncompressedSize int64
	var totalCompressedSize int64

	for _, file := range reader.File {
		// Check for path traversal
		if err := fs.validateZIPPath(file.Name); err != nil {
			return err
		}

		// Check for symlinks (they have mode with ModeSymlink bit set)
		if file.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlinks are not allowed in ZIP files: %s", file.Name)
		}

		// Accumulate sizes
		totalUncompressedSize += int64(file.UncompressedSize64)
		totalCompressedSize += int64(file.CompressedSize64)

		// Check single file size
		if int64(file.UncompressedSize64) > MaxSingleFileSize {
			return fmt.Errorf("file %s is too large (%d bytes), maximum allowed: %d bytes",
				file.Name, file.UncompressedSize64, MaxSingleFileSize)
		}
	}

	// Check total uncompressed size
	if totalUncompressedSize > fs.maxUncompressedSize {
		return fmt.Errorf("total uncompressed size (%d bytes) exceeds maximum (%d bytes)",
			totalUncompressedSize, fs.maxUncompressedSize)
	}

	// Check compression ratio to detect ZIP bombs
	if totalCompressedSize > 0 {
		ratio := float64(totalUncompressedSize) / float64(totalCompressedSize)
		if ratio > fs.maxCompressionRatio {
			return fmt.Errorf("compression ratio (%.2f:1) exceeds maximum (%.2f:1), possible ZIP bomb",
				ratio, fs.maxCompressionRatio)
		}
	}

	return nil
}

// validateZIPPath checks for path traversal attempts in ZIP entry names
func (fs *FSManager) validateZIPPath(zipPath string) error {
	// Clean the path
	cleanPath := filepath.Clean(zipPath)

	// Check for absolute paths
	if filepath.IsAbs(cleanPath) {
		return fmt.Errorf("absolute paths are not allowed in ZIP files: %s", zipPath)
	}

	// Check for parent directory references
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("path traversal detected in ZIP entry: %s", zipPath)
	}

	// Check for leading slashes or backslashes
	if strings.HasPrefix(zipPath, "/") || strings.HasPrefix(zipPath, "\\") {
		return fmt.Errorf("paths cannot start with / or \\: %s", zipPath)
	}

	return nil
}

// extractFile extracts a single file from a ZIP archive
func (fs *FSManager) extractFile(file *zip.File, destPath string) error {
	// Validate and clean the file path
	if err := fs.validateZIPPath(file.Name); err != nil {
		return err
	}

	// Build the full destination path
	destFilePath := filepath.Join(destPath, file.Name)

	// Double-check the path is still within destPath (defense in depth)
	if !strings.HasPrefix(destFilePath, filepath.Clean(destPath)+string(os.PathSeparator)) {
		return fmt.Errorf("path traversal detected: %s escapes destination %s", destFilePath, destPath)
	}

	// Handle directories
	if file.FileInfo().IsDir() {
		return os.MkdirAll(destFilePath, DefaultDirPerm)
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(destFilePath), DefaultDirPerm); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Open the file from ZIP
	srcFile, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in ZIP: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.OpenFile(destFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, DefaultFilePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy with size limit check
	written, err := io.CopyN(destFile, srcFile, MaxSingleFileSize+1)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to copy file data: %w", err)
	}

	// Verify size matches expected
	if written > MaxSingleFileSize {
		// Clean up the oversized file
		destFile.Close()
		os.Remove(destFilePath)
		return fmt.Errorf("file exceeded maximum size during extraction: %s", file.Name)
	}

	return nil
}

// CreateZIP creates a ZIP archive from a directory
func (fs *FSManager) CreateZIP(srcPath, zipPath string) error {
	// Validate inputs
	if srcPath == "" {
		return fmt.Errorf("source path cannot be empty")
	}
	if zipPath == "" {
		return fmt.Errorf("zip path cannot be empty")
	}

	// Ensure source directory exists
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("failed to stat source path: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source path is not a directory: %s", srcPath)
	}

	// Create ZIP file
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create ZIP file: %w", err)
	}
	defer zipFile.Close()

	// Create ZIP writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Get absolute source path for relative path calculation
	absSrcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute source path: %w", err)
	}

	// Walk the directory and add files
	err = filepath.Walk(absSrcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories (except the root)
		if info.Name() != filepath.Base(absSrcPath) && strings.HasPrefix(info.Name(), ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(absSrcPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Normalize path separators for ZIP (use forward slashes)
		zipPath := filepath.ToSlash(relPath)

		// Create ZIP entry header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create ZIP header: %w", err)
		}
		header.Name = zipPath

		// Set compression method
		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store
		} else {
			header.Method = zip.Deflate
		}

		// Create writer for this entry
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create ZIP entry: %w", err)
		}

		// If it's a file, copy its contents
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()

			if _, err := io.Copy(writer, file); err != nil {
				return fmt.Errorf("failed to write file to ZIP: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	return nil
}

// ValidatePath ensures a path is within the base directory (prevents path traversal)
func (fs *FSManager) ValidatePath(path string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Clean the path
	cleanPath := filepath.Clean(absPath)

	// Check if path is within base directory
	relPath, err := filepath.Rel(fs.baseDir, cleanPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	// If the relative path starts with "..", it's outside the base directory
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("path %s is outside base directory %s", path, fs.baseDir)
	}

	return nil
}

// CalculateSHA256 calculates the SHA256 hash of a file
func (fs *FSManager) CalculateSHA256(filePath string) (string, error) {
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create hash
	hash := sha256.New()

	// Copy file data to hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	// Get the hash sum
	hashSum := hash.Sum(nil)

	// Convert to hex string
	return hex.EncodeToString(hashSum), nil
}

// VerifyIntegrity verifies a file's SHA256 hash matches the expected value
func (fs *FSManager) VerifyIntegrity(filePath, expectedHash string) error {
	actualHash, err := fs.CalculateSHA256(filePath)
	if err != nil {
		return fmt.Errorf("failed to calculate hash: %w", err)
	}

	// Compare hashes (case-insensitive)
	if !strings.EqualFold(actualHash, expectedHash) {
		return fmt.Errorf("integrity check failed: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

// EnsureDir creates a directory and all parent directories if they don't exist
func (fs *FSManager) EnsureDir(path string) error {
	// Validate path is within base directory
	if err := fs.ValidatePath(path); err != nil {
		return err
	}

	// Create directory
	if err := os.MkdirAll(path, DefaultDirPerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return nil
}

// RemoveDir removes a directory and all its contents
func (fs *FSManager) RemoveDir(path string) error {
	// Validate path is within base directory
	if err := fs.ValidatePath(path); err != nil {
		return err
	}

	// Remove directory
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}

	return nil
}

// GetBaseDir returns the base directory
func (fs *FSManager) GetBaseDir() string {
	return fs.baseDir
}

// SetMaxUncompressedSize sets the maximum uncompressed size for ZIP files
func (fs *FSManager) SetMaxUncompressedSize(size int64) {
	if size > 0 {
		fs.maxUncompressedSize = size
	}
}

// SetMaxFiles sets the maximum number of files allowed in a ZIP
func (fs *FSManager) SetMaxFiles(count int) {
	if count > 0 {
		fs.maxFiles = count
	}
}

// SetMaxCompressionRatio sets the maximum compression ratio allowed
func (fs *FSManager) SetMaxCompressionRatio(ratio float64) {
	if ratio > 0 {
		fs.maxCompressionRatio = ratio
	}
}

// GetDirSize calculates the total size of a directory
func (fs *FSManager) GetDirSize(path string) (int64, error) {
	// Validate path is within base directory
	if err := fs.ValidatePath(path); err != nil {
		return 0, err
	}

	var totalSize int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to calculate directory size: %w", err)
	}

	return totalSize, nil
}
