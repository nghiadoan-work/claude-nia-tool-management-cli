package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
)

// Example demonstrating FSManager usage for safe ZIP operations
func main() {
	// Create a temporary base directory for our operations
	baseDir := filepath.Join(os.TempDir(), "cntm-fs-example")
	defer os.RemoveAll(baseDir)

	// Create FSManager
	fsm, err := data.NewFSManager(baseDir)
	if err != nil {
		log.Fatalf("Failed to create FSManager: %v", err)
	}

	fmt.Println("=== FSManager Example ===")
	fmt.Printf("Base directory: %s\n\n", fsm.GetBaseDir())

	// Example 1: Create a sample tool directory
	fmt.Println("1. Creating sample tool directory...")
	toolDir := filepath.Join(baseDir, "sample-tool")
	if err := fsm.EnsureDir(toolDir); err != nil {
		log.Fatalf("Failed to create tool directory: %v", err)
	}

	// Add sample files
	toolFiles := map[string]string{
		"metadata.json": `{"name": "sample-tool", "version": "1.0.0", "type": "agent"}`,
		"README.md":     "# Sample Tool\n\nThis is a sample tool for demonstration.",
		"main.go":       "package main\n\nfunc main() {\n\tprintln(\"Hello from sample-tool\")\n}",
	}

	for filename, content := range toolFiles {
		filePath := filepath.Join(toolDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to write file %s: %v", filename, err)
		}
	}
	fmt.Printf("   Created tool directory with %d files\n\n", len(toolFiles))

	// Example 2: Create a ZIP archive from the tool directory
	fmt.Println("2. Creating ZIP archive...")
	zipPath := filepath.Join(baseDir, "sample-tool.zip")
	if err := fsm.CreateZIP(toolDir, zipPath); err != nil {
		log.Fatalf("Failed to create ZIP: %v", err)
	}

	// Get ZIP file info
	zipInfo, _ := os.Stat(zipPath)
	fmt.Printf("   Created: %s (Size: %d bytes)\n\n", zipPath, zipInfo.Size())

	// Example 3: Calculate SHA256 integrity hash
	fmt.Println("3. Calculating SHA256 hash...")
	hash, err := fsm.CalculateSHA256(zipPath)
	if err != nil {
		log.Fatalf("Failed to calculate hash: %v", err)
	}
	fmt.Printf("   SHA256: %s\n\n", hash)

	// Example 4: Verify integrity
	fmt.Println("4. Verifying ZIP integrity...")
	if err := fsm.VerifyIntegrity(zipPath, hash); err != nil {
		log.Fatalf("Integrity check failed: %v", err)
	}
	fmt.Println("   Integrity verified successfully!")
	fmt.Println()

	// Example 5: Extract the ZIP to a new location
	fmt.Println("5. Extracting ZIP archive...")
	extractDir := filepath.Join(baseDir, "extracted")
	if err := fsm.ExtractZIP(zipPath, extractDir); err != nil {
		log.Fatalf("Failed to extract ZIP: %v", err)
	}

	// List extracted files
	extractedFiles := []string{}
	filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			relPath, _ := filepath.Rel(extractDir, path)
			extractedFiles = append(extractedFiles, relPath)
		}
		return nil
	})
	fmt.Printf("   Extracted %d files:\n", len(extractedFiles))
	for _, file := range extractedFiles {
		fmt.Printf("   - %s\n", file)
	}
	fmt.Println()

	// Example 6: Calculate directory size
	fmt.Println("6. Calculating directory sizes...")
	toolSize, _ := fsm.GetDirSize(toolDir)
	extractedSize, _ := fsm.GetDirSize(extractDir)
	fmt.Printf("   Original tool directory: %d bytes\n", toolSize)
	fmt.Printf("   Extracted directory: %d bytes\n", extractedSize)
	fmt.Println()

	// Example 7: Demonstrate security - path traversal prevention
	fmt.Println("7. Security demonstration - Path traversal prevention...")
	maliciousPath := filepath.Join(baseDir, "..", "escape-attempt")
	if err := fsm.ValidatePath(maliciousPath); err != nil {
		fmt.Printf("   Path traversal attempt blocked: %v\n", err)
	} else {
		fmt.Println("   WARNING: Path traversal not blocked!")
	}
	fmt.Println()

	// Example 8: Clean up - remove extracted directory
	fmt.Println("8. Cleaning up extracted directory...")
	if err := fsm.RemoveDir(extractDir); err != nil {
		log.Fatalf("Failed to remove directory: %v", err)
	}
	fmt.Println("   Removed successfully")
	fmt.Println()

	fmt.Println("=== Example Complete ===")
	fmt.Println("\nKey Security Features Demonstrated:")
	fmt.Println("- Path traversal prevention")
	fmt.Println("- ZIP bomb protection (size limits)")
	fmt.Println("- File integrity verification (SHA256)")
	fmt.Println("- Safe extraction with validation")
	fmt.Println("- Base directory sandboxing")
}
