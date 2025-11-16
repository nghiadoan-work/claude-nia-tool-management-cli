package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
)

func main() {
	fmt.Println("=== Installer Service Example ===\n")

	// Setup paths
	baseDir := filepath.Join(os.TempDir(), "cntm-installer-example")
	defer os.RemoveAll(baseDir)

	lockFilePath := filepath.Join(baseDir, ".claude-lock.json")

	// Create config
	config := &models.Config{
		Registry: models.RegistryConfig{
			URL:    "https://github.com/nghiadoan-work/claude-tools-registry",
			Branch: "main",
		},
		Local: models.LocalConfig{
			DefaultPath: baseDir,
		},
	}

	// Initialize services
	githubClient := services.NewGitHubClient(services.GitHubClientConfig{
		Owner:  "nghiadoan-work",
		Repo:   "claude-tools-registry",
		Branch: "main",
	})

	registryService := services.NewRegistryServiceWithoutCache(githubClient)

	cacheManager, err := data.NewCacheManager(filepath.Join(baseDir, ".cache"), 3600*time.Second)
	if err != nil {
		fmt.Printf("Error creating cache manager: %v\n", err)
		return
	}

	registryServiceWithCache := services.NewRegistryService(githubClient, cacheManager)

	fsManager, err := data.NewFSManager(baseDir)
	if err != nil {
		fmt.Printf("Error creating FSManager: %v\n", err)
		return
	}

	lockFileService, err := services.NewLockFileService(lockFilePath)
	if err != nil {
		fmt.Printf("Error creating LockFileService: %v\n", err)
		return
	}

	// Create installer service
	installer, err := services.NewInstallerService(
		githubClient,
		registryServiceWithCache,
		fsManager,
		lockFileService,
		config,
	)
	if err != nil {
		fmt.Printf("Error creating installer: %v\n", err)
		return
	}

	// Example 1: List available tools
	fmt.Println("1. Fetching available tools from registry...")
	registry, err := registryService.GetRegistry()
	if err != nil {
		fmt.Printf("Error fetching registry: %v\n", err)
		return
	}

	fmt.Printf("Registry version: %s\n", registry.Version)
	fmt.Printf("Total tool types: %d\n\n", len(registry.Tools))

	// Example 2: Search for a specific tool
	fmt.Println("2. Searching for 'code-reviewer' tool...")
	tool, err := registryServiceWithCache.GetTool("code-reviewer", models.ToolTypeAgent)
	if err != nil {
		fmt.Printf("Error finding tool: %v\n", err)
		// Continue with other examples
	} else {
		fmt.Printf("Found: %s v%s\n", tool.Name, tool.Version)
		fmt.Printf("Description: %s\n", tool.Description)
		fmt.Printf("Author: %s\n", tool.Author)
		fmt.Printf("Size: %d bytes\n\n", tool.Size)

		// Example 3: Install the tool
		fmt.Println("3. Installing tool...")
		err = installer.Install("code-reviewer")
		if err != nil {
			fmt.Printf("Error installing tool: %v\n", err)
		} else {
			fmt.Println("Installation successful!\n")

			// Example 4: Verify installation
			fmt.Println("4. Verifying installation...")
			err = installer.VerifyInstallation("code-reviewer")
			if err != nil {
				fmt.Printf("Verification failed: %v\n", err)
			} else {
				fmt.Println("Verification successful!\n")
			}

			// Example 5: Check installed version
			fmt.Println("5. Checking installed version...")
			version, err := installer.GetInstalledVersion("code-reviewer")
			if err != nil {
				fmt.Printf("Error getting version: %v\n", err)
			} else {
				fmt.Printf("Installed version: %s\n\n", version)
			}

			// Example 6: List all installed tools
			fmt.Println("6. Listing installed tools...")
			installed, err := installer.GetInstalledTools()
			if err != nil {
				fmt.Printf("Error listing tools: %v\n", err)
			} else {
				fmt.Printf("Total installed tools: %d\n", len(installed))
				for name, tool := range installed {
					fmt.Printf("  - %s v%s (%s)\n", name, tool.Version, tool.Type)
				}
				fmt.Println()
			}
		}
	}

	// Example 7: Install multiple tools
	fmt.Println("7. Installing multiple tools...")
	toolNames := []string{"git-helper", "test-writer"}
	results, errors := installer.InstallMultiple(toolNames)

	fmt.Printf("Installation results:\n")
	for _, result := range results {
		status := "SUCCESS"
		if !result.Success {
			status = "FAILED"
		}
		fmt.Printf("  - %s: %s - %s\n", result.ToolName, status, result.Message)
	}

	if len(errors) > 0 {
		fmt.Printf("\nEncountered %d error(s) during installation\n", len(errors))
	}
	fmt.Println()

	// Example 8: Uninstall a tool
	fmt.Println("8. Uninstalling 'code-reviewer'...")
	err = installer.Uninstall("code-reviewer")
	if err != nil {
		fmt.Printf("Error uninstalling: %v\n", err)
	} else {
		fmt.Println("Uninstallation successful!")
	}

	fmt.Println("\n=== Example Complete ===")
	fmt.Printf("Installation directory: %s\n", baseDir)
	fmt.Printf("Lock file: %s\n", lockFilePath)
}
