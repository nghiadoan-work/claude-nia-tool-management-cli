package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
)

// This example demonstrates how to use the LockFileService to manage installed tools
func main() {
	// Create a temporary directory for demonstration
	tmpDir, err := os.MkdirTemp("", "cntm-example-")
	if err != nil {
		log.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	lockFilePath := filepath.Join(tmpDir, ".claude-lock.json")
	fmt.Printf("Lock file path: %s\n\n", lockFilePath)

	// Create a new LockFileService
	lockService, err := services.NewLockFileService(lockFilePath)
	if err != nil {
		log.Fatalf("Failed to create lock service: %v", err)
	}

	// Step 1: Load lock file (creates default if doesn't exist)
	fmt.Println("=== Step 1: Load Lock File ===")
	lockFile, err := lockService.Load()
	if err != nil {
		log.Fatalf("Failed to load lock file: %v", err)
	}
	fmt.Printf("Loaded lock file version: %s\n", lockFile.Version)
	fmt.Printf("Tools count: %d\n\n", len(lockFile.Tools))

	// Step 2: Set registry URL
	fmt.Println("=== Step 2: Set Registry URL ===")
	registryURL := "https://github.com/nghiadoan-work/claude-tools-registry"
	err = lockService.SetRegistry(registryURL)
	if err != nil {
		log.Fatalf("Failed to set registry: %v", err)
	}
	fmt.Printf("Registry set to: %s\n\n", registryURL)

	// Step 3: Add tools
	fmt.Println("=== Step 3: Add Tools ===")

	// Add code-reviewer agent
	codeReviewer := &models.InstalledTool{
		Version:     "1.2.0",
		Type:        models.ToolTypeAgent,
		InstalledAt: time.Now(),
		Source:      "registry",
		Integrity:   "sha256-abcdef123456",
	}
	err = lockService.AddTool("code-reviewer", codeReviewer)
	if err != nil {
		log.Fatalf("Failed to add code-reviewer: %v", err)
	}
	fmt.Println("Added: code-reviewer@1.2.0")

	// Add git-helper command
	gitHelper := &models.InstalledTool{
		Version:     "2.0.1",
		Type:        models.ToolTypeCommand,
		InstalledAt: time.Now(),
		Source:      "registry",
		Integrity:   "sha256-xyz789abc",
	}
	err = lockService.AddTool("git-helper", gitHelper)
	if err != nil {
		log.Fatalf("Failed to add git-helper: %v", err)
	}
	fmt.Println("Added: git-helper@2.0.1")

	// Add test-runner skill
	testRunner := &models.InstalledTool{
		Version:     "3.1.0",
		Type:        models.ToolTypeSkill,
		InstalledAt: time.Now(),
		Source:      "registry",
		Integrity:   "sha256-skill123",
	}
	err = lockService.AddTool("test-runner", testRunner)
	if err != nil {
		log.Fatalf("Failed to add test-runner: %v", err)
	}
	fmt.Println("Added: test-runner@3.1.0\n")

	// Step 4: List all tools
	fmt.Println("=== Step 4: List All Tools ===")
	tools, err := lockService.ListTools()
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for name, tool := range tools {
		fmt.Printf("  - %s@%s [%s]\n", name, tool.Version, tool.Type)
	}
	fmt.Println()

	// Step 5: Check if specific tool is installed
	fmt.Println("=== Step 5: Check Installation Status ===")
	installed, err := lockService.IsInstalled("code-reviewer")
	if err != nil {
		log.Fatalf("Failed to check installation: %v", err)
	}
	fmt.Printf("code-reviewer installed: %v\n", installed)

	installed, err = lockService.IsInstalled("non-existent-tool")
	if err != nil {
		log.Fatalf("Failed to check installation: %v", err)
	}
	fmt.Printf("non-existent-tool installed: %v\n\n", installed)

	// Step 6: Get specific tool
	fmt.Println("=== Step 6: Get Tool Details ===")
	tool, err := lockService.GetTool("git-helper")
	if err != nil {
		log.Fatalf("Failed to get tool: %v", err)
	}
	fmt.Printf("Tool: git-helper\n")
	fmt.Printf("  Version: %s\n", tool.Version)
	fmt.Printf("  Type: %s\n", tool.Type)
	fmt.Printf("  Installed: %s\n", tool.InstalledAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Integrity: %s\n\n", tool.Integrity)

	// Step 7: Update a tool
	fmt.Println("=== Step 7: Update Tool ===")
	updatedCodeReviewer := &models.InstalledTool{
		Version:     "1.3.0",
		Type:        models.ToolTypeAgent,
		InstalledAt: time.Now(),
		Source:      "registry",
		Integrity:   "sha256-updated456",
	}
	err = lockService.UpdateTool("code-reviewer", updatedCodeReviewer)
	if err != nil {
		log.Fatalf("Failed to update tool: %v", err)
	}
	fmt.Println("Updated code-reviewer from 1.2.0 to 1.3.0\n")

	// Verify update
	tool, err = lockService.GetTool("code-reviewer")
	if err != nil {
		log.Fatalf("Failed to get updated tool: %v", err)
	}
	fmt.Printf("Verified: code-reviewer now at version %s\n\n", tool.Version)

	// Step 8: Remove a tool
	fmt.Println("=== Step 8: Remove Tool ===")
	err = lockService.RemoveTool("test-runner")
	if err != nil {
		log.Fatalf("Failed to remove tool: %v", err)
	}
	fmt.Println("Removed: test-runner\n")

	// List remaining tools
	tools, err = lockService.ListTools()
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	fmt.Printf("Remaining tools (%d):\n", len(tools))
	for name, tool := range tools {
		fmt.Printf("  - %s@%s [%s]\n", name, tool.Version, tool.Type)
	}
	fmt.Println()

	// Step 9: Show lock file content
	fmt.Println("=== Step 9: Lock File Content ===")
	data, err := os.ReadFile(lockFilePath)
	if err != nil {
		log.Fatalf("Failed to read lock file: %v", err)
	}
	fmt.Println(string(data))

	fmt.Println("\n=== Example Complete ===")
	fmt.Printf("Lock file created at: %s\n", lockFilePath)
	fmt.Println("Note: This file will be deleted when the example exits.")
}
