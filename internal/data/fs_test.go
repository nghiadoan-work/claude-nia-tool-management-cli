package data

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFSManager(t *testing.T) {
	tests := []struct {
		name    string
		baseDir string
		wantErr bool
	}{
		{
			name:    "valid base directory",
			baseDir: t.TempDir(),
			wantErr: false,
		},
		{
			name:    "empty base directory",
			baseDir: "",
			wantErr: true,
		},
		{
			name:    "creates non-existent directory",
			baseDir: filepath.Join(t.TempDir(), "new-dir"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm, err := NewFSManager(tt.baseDir)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, fsm)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, fsm)
				if fsm != nil {
					// Verify base directory exists
					_, err := os.Stat(fsm.GetBaseDir())
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestFSManager_ValidatePath(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "path within base directory",
			path:    filepath.Join(baseDir, "subdir", "file.txt"),
			wantErr: false,
		},
		{
			name:    "base directory itself",
			path:    baseDir,
			wantErr: false,
		},
		{
			name:    "path outside base directory",
			path:    filepath.Join(baseDir, "..", "outside"),
			wantErr: true,
		},
		{
			name:    "path with parent references escaping base",
			path:    filepath.Join(baseDir, "subdir", "..", "..", "outside"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsm.ValidatePath(tt.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFSManager_validateZIPPath(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	tests := []struct {
		name    string
		zipPath string
		wantErr bool
	}{
		{
			name:    "normal file path",
			zipPath: "subdir/file.txt",
			wantErr: false,
		},
		{
			name:    "file in root",
			zipPath: "file.txt",
			wantErr: false,
		},
		{
			name:    "nested directories",
			zipPath: "a/b/c/d/file.txt",
			wantErr: false,
		},
		{
			name:    "path traversal with ..",
			zipPath: "../../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "path traversal in middle",
			zipPath: "subdir/../../etc/passwd",
			wantErr: true,
		},
		{
			name:    "absolute path",
			zipPath: "/etc/passwd",
			wantErr: true,
		},
		{
			name:    "path with backslash separator",
			zipPath: "subdir\\file.txt",
			wantErr: false, // Backslashes in paths are treated as normal chars on Unix
		},
		{
			name:    "leading slash",
			zipPath: "/file.txt",
			wantErr: true,
		},
		{
			name:    "leading backslash",
			zipPath: "\\file.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fsm.validateZIPPath(tt.zipPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFSManager_CreateZIP(t *testing.T) {
	// Create source directory with files
	srcDir := t.TempDir()

	// Create test files
	err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644)
	require.NoError(t, err)

	// Create hidden file (should be excluded)
	err = os.WriteFile(filepath.Join(srcDir, ".hidden"), []byte("hidden"), 0644)
	require.NoError(t, err)

	// Create FSManager
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Create ZIP
	zipPath := filepath.Join(baseDir, "test.zip")
	err = fsm.CreateZIP(srcDir, zipPath)
	require.NoError(t, err)

	// Verify ZIP exists
	_, err = os.Stat(zipPath)
	require.NoError(t, err)

	// Verify ZIP contents
	reader, err := zip.OpenReader(zipPath)
	require.NoError(t, err)
	defer reader.Close()

	// Check that we have the expected files
	fileNames := make(map[string]bool)
	for _, file := range reader.File {
		fileNames[file.Name] = true
	}

	// Should include regular files
	assert.True(t, fileNames["file1.txt"], "should contain file1.txt")
	assert.True(t, fileNames["subdir/"], "should contain subdir/")
	assert.True(t, fileNames["subdir/file2.txt"], "should contain subdir/file2.txt")

	// Should exclude hidden files
	assert.False(t, fileNames[".hidden"], "should not contain .hidden file")
}

func TestFSManager_ExtractZIP(t *testing.T) {
	// Create a test ZIP file
	baseDir := t.TempDir()
	zipPath := filepath.Join(baseDir, "test.zip")

	// Create ZIP with test files
	zipFile, err := os.Create(zipPath)
	require.NoError(t, err)

	zipWriter := zip.NewWriter(zipFile)

	// Add a file
	writer, err := zipWriter.Create("file1.txt")
	require.NoError(t, err)
	_, err = writer.Write([]byte("content1"))
	require.NoError(t, err)

	// Add a file in subdirectory
	writer, err = zipWriter.Create("subdir/file2.txt")
	require.NoError(t, err)
	_, err = writer.Write([]byte("content2"))
	require.NoError(t, err)

	err = zipWriter.Close()
	require.NoError(t, err)
	err = zipFile.Close()
	require.NoError(t, err)

	// Create FSManager
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Extract ZIP
	destDir := filepath.Join(baseDir, "extracted")
	err = fsm.ExtractZIP(zipPath, destDir)
	require.NoError(t, err)

	// Verify extracted files
	content, err := os.ReadFile(filepath.Join(destDir, "file1.txt"))
	require.NoError(t, err)
	assert.Equal(t, "content1", string(content))

	content, err = os.ReadFile(filepath.Join(destDir, "subdir", "file2.txt"))
	require.NoError(t, err)
	assert.Equal(t, "content2", string(content))
}

func TestFSManager_ExtractZIP_PathTraversal(t *testing.T) {
	baseDir := t.TempDir()
	zipPath := filepath.Join(baseDir, "malicious.zip")

	// Create ZIP with path traversal attempt
	zipFile, err := os.Create(zipPath)
	require.NoError(t, err)

	zipWriter := zip.NewWriter(zipFile)

	// Try to add a file with path traversal
	writer, err := zipWriter.Create("../../../etc/passwd")
	require.NoError(t, err)
	_, err = writer.Write([]byte("malicious content"))
	require.NoError(t, err)

	err = zipWriter.Close()
	require.NoError(t, err)
	err = zipFile.Close()
	require.NoError(t, err)

	// Create FSManager
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Try to extract - should fail
	destDir := filepath.Join(baseDir, "extracted")
	err = fsm.ExtractZIP(zipPath, destDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "path traversal")
}

func TestFSManager_ExtractZIP_ZipBomb(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Set low limits for testing
	fsm.SetMaxUncompressedSize(1024) // 1KB
	fsm.SetMaxFiles(10)

	zipPath := filepath.Join(baseDir, "bomb.zip")

	// Create ZIP that exceeds file count limit
	zipFile, err := os.Create(zipPath)
	require.NoError(t, err)

	zipWriter := zip.NewWriter(zipFile)

	// Add more files than allowed
	for i := 0; i < 20; i++ {
		writer, err := zipWriter.Create(filepath.Join("file", string(rune(i))+".txt"))
		require.NoError(t, err)
		_, err = writer.Write([]byte("content"))
		require.NoError(t, err)
	}

	err = zipWriter.Close()
	require.NoError(t, err)
	err = zipFile.Close()
	require.NoError(t, err)

	// Try to extract - should fail
	destDir := filepath.Join(baseDir, "extracted")
	err = fsm.ExtractZIP(zipPath, destDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "too many files")
}

func TestFSManager_ExtractZIP_OversizedFile(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Set low limit for testing
	fsm.SetMaxUncompressedSize(100) // 100 bytes

	zipPath := filepath.Join(baseDir, "large.zip")

	// Create ZIP with oversized file
	zipFile, err := os.Create(zipPath)
	require.NoError(t, err)

	zipWriter := zip.NewWriter(zipFile)

	writer, err := zipWriter.Create("large.txt")
	require.NoError(t, err)

	// Write more than the limit
	largeContent := strings.Repeat("a", 200)
	_, err = writer.Write([]byte(largeContent))
	require.NoError(t, err)

	err = zipWriter.Close()
	require.NoError(t, err)
	err = zipFile.Close()
	require.NoError(t, err)

	// Try to extract - should fail
	destDir := filepath.Join(baseDir, "extracted")
	err = fsm.ExtractZIP(zipPath, destDir)
	require.Error(t, err)
}

func TestFSManager_CalculateSHA256(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Create test file
	testFile := filepath.Join(baseDir, "test.txt")
	content := []byte("Hello, World!")
	err = os.WriteFile(testFile, content, 0644)
	require.NoError(t, err)

	// Calculate hash
	hash, err := fsm.CalculateSHA256(testFile)
	require.NoError(t, err)

	// Known SHA256 of "Hello, World!"
	expectedHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	assert.Equal(t, expectedHash, hash)

	// Test with non-existent file
	_, err = fsm.CalculateSHA256(filepath.Join(baseDir, "nonexistent.txt"))
	assert.Error(t, err)
}

func TestFSManager_VerifyIntegrity(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Create test file
	testFile := filepath.Join(baseDir, "test.txt")
	content := []byte("Hello, World!")
	err = os.WriteFile(testFile, content, 0644)
	require.NoError(t, err)

	// Known SHA256 of "Hello, World!"
	correctHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	incorrectHash := "0000000000000000000000000000000000000000000000000000000000000000"

	// Test with correct hash
	err = fsm.VerifyIntegrity(testFile, correctHash)
	assert.NoError(t, err)

	// Test with incorrect hash
	err = fsm.VerifyIntegrity(testFile, incorrectHash)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "integrity check failed")

	// Test case-insensitive comparison
	err = fsm.VerifyIntegrity(testFile, strings.ToUpper(correctHash))
	assert.NoError(t, err)
}

func TestFSManager_EnsureDir(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Test creating nested directories
	testDir := filepath.Join(baseDir, "a", "b", "c")
	err = fsm.EnsureDir(testDir)
	require.NoError(t, err)

	// Verify directory exists
	info, err := os.Stat(testDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Test creating directory outside base - should fail
	outsideDir := filepath.Join(baseDir, "..", "outside")
	err = fsm.EnsureDir(outsideDir)
	assert.Error(t, err)
}

func TestFSManager_RemoveDir(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Create test directory with files
	testDir := filepath.Join(baseDir, "testdir")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(testDir, "file.txt"), []byte("content"), 0644)
	require.NoError(t, err)

	// Remove directory
	err = fsm.RemoveDir(testDir)
	require.NoError(t, err)

	// Verify directory is removed
	_, err = os.Stat(testDir)
	assert.True(t, os.IsNotExist(err))

	// Test removing directory outside base - should fail
	outsideDir := filepath.Join(baseDir, "..", "outside")
	err = fsm.RemoveDir(outsideDir)
	assert.Error(t, err)
}

func TestFSManager_GetDirSize(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Create test directory with files
	testDir := filepath.Join(baseDir, "testdir")
	err = os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create files with known sizes
	err = os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("12345"), 0644) // 5 bytes
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(testDir, "file2.txt"), []byte("1234567890"), 0644) // 10 bytes
	require.NoError(t, err)

	// Calculate size
	size, err := fsm.GetDirSize(testDir)
	require.NoError(t, err)
	assert.Equal(t, int64(15), size) // 5 + 10 = 15 bytes

	// Test with directory outside base - should fail
	outsideDir := filepath.Join(baseDir, "..", "outside")
	_, err = fsm.GetDirSize(outsideDir)
	assert.Error(t, err)
}

func TestFSManager_RoundTrip(t *testing.T) {
	// Create source directory
	srcDir := t.TempDir()

	// Create test files
	err := os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644)
	require.NoError(t, err)

	// Create FSManager
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Create ZIP
	zipPath := filepath.Join(baseDir, "test.zip")
	err = fsm.CreateZIP(srcDir, zipPath)
	require.NoError(t, err)

	// Calculate hash of ZIP
	hash, err := fsm.CalculateSHA256(zipPath)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Verify integrity
	err = fsm.VerifyIntegrity(zipPath, hash)
	require.NoError(t, err)

	// Extract ZIP
	destDir := filepath.Join(baseDir, "extracted")
	err = fsm.ExtractZIP(zipPath, destDir)
	require.NoError(t, err)

	// Verify extracted content matches original
	content, err := os.ReadFile(filepath.Join(destDir, "file1.txt"))
	require.NoError(t, err)
	assert.Equal(t, "content1", string(content))

	content, err = os.ReadFile(filepath.Join(destDir, "subdir", "file2.txt"))
	require.NoError(t, err)
	assert.Equal(t, "content2", string(content))
}

func TestFSManager_EmptyZIP(t *testing.T) {
	baseDir := t.TempDir()
	zipPath := filepath.Join(baseDir, "empty.zip")

	// Create empty ZIP
	zipFile, err := os.Create(zipPath)
	require.NoError(t, err)

	zipWriter := zip.NewWriter(zipFile)
	err = zipWriter.Close()
	require.NoError(t, err)
	err = zipFile.Close()
	require.NoError(t, err)

	// Create FSManager
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Try to extract empty ZIP - should fail
	destDir := filepath.Join(baseDir, "extracted")
	err = fsm.ExtractZIP(zipPath, destDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestFSManager_SettersGetters(t *testing.T) {
	baseDir := t.TempDir()
	fsm, err := NewFSManager(baseDir)
	require.NoError(t, err)

	// Test GetBaseDir
	assert.NotEmpty(t, fsm.GetBaseDir())

	// Test SetMaxUncompressedSize
	fsm.SetMaxUncompressedSize(5000)
	assert.Equal(t, int64(5000), fsm.maxUncompressedSize)

	// Test SetMaxFiles
	fsm.SetMaxFiles(100)
	assert.Equal(t, 100, fsm.maxFiles)

	// Test SetMaxCompressionRatio
	fsm.SetMaxCompressionRatio(50.0)
	assert.Equal(t, 50.0, fsm.maxCompressionRatio)

	// Test that negative values are ignored
	fsm.SetMaxUncompressedSize(-1)
	assert.Equal(t, int64(5000), fsm.maxUncompressedSize) // Should remain unchanged

	fsm.SetMaxFiles(-1)
	assert.Equal(t, 100, fsm.maxFiles) // Should remain unchanged

	fsm.SetMaxCompressionRatio(-1.0)
	assert.Equal(t, 50.0, fsm.maxCompressionRatio) // Should remain unchanged
}
