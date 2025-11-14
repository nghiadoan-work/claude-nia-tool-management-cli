package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/services"
	"github.com/spf13/cobra"
)

var (
	// Update flags
	updateAll bool
	updateYes bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [tool-name]",
	Short: "Update tools to the latest version",
	Long: `Update one or all tools to the latest version available in the registry.

By default, this command will prompt for confirmation before updating.
Use --yes to skip the confirmation prompt.

Examples:
  cntm update code-reviewer          # Update specific tool
  cntm update --all                  # Update all outdated tools
  cntm update --all --yes            # Update all without confirmation`,
	Example: `  cntm update code-reviewer          # Update specific tool
  cntm update --all                  # Update all outdated tools
  cntm update --all --yes            # Update all without confirmation
  cntm update code-reviewer --yes    # Update without confirmation`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Either provide a tool name or use --all
		if updateAll && len(args) > 0 {
			return fmt.Errorf("cannot specify tool name with --all flag")
		}
		if !updateAll && len(args) == 0 {
			return fmt.Errorf("requires a tool name or --all flag")
		}
		return nil
	},
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Update flags
	updateCmd.Flags().BoolVar(&updateAll, "all", false, "update all outdated tools")
	updateCmd.Flags().BoolVarP(&updateYes, "yes", "y", false, "skip confirmation prompts")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Parse GitHub URL to get owner and repo
	owner, repo, err := parseGitHubURL(cfg.Registry.URL)
	if err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}

	// Initialize services
	githubClient := services.NewGitHubClient(services.GitHubClientConfig{
		Owner:     owner,
		Repo:      repo,
		Branch:    cfg.Registry.Branch,
		AuthToken: cfg.Registry.AuthToken,
	})

	cacheManager, err := data.NewCacheManager(basePath, 3600*time.Second) // 1 hour TTL
	if err != nil {
		return fmt.Errorf("failed to create cache manager: %w", err)
	}
	registryService := services.NewRegistryService(githubClient, cacheManager)

	// Initialize FSManager and LockFileService
	fsManager, err := data.NewFSManager(basePath)
	if err != nil {
		return fmt.Errorf("failed to create file system manager: %w", err)
	}

	lockFilePath := filepath.Join(basePath, ".claude-lock.json")
	lockFileService, err := services.NewLockFileService(lockFilePath)
	if err != nil {
		return fmt.Errorf("failed to create lock file service: %w", err)
	}

	// Initialize InstallerService
	installer, err := services.NewInstallerService(
		githubClient,
		registryService,
		fsManager,
		lockFileService,
		cfg,
	)
	if err != nil {
		return fmt.Errorf("failed to create installer service: %w", err)
	}

	// Initialize UpdaterService
	updater, err := services.NewUpdaterService(
		registryService,
		lockFileService,
		installer,
	)
	if err != nil {
		return fmt.Errorf("failed to create updater service: %w", err)
	}

	// Execute update
	if updateAll {
		return runUpdateAll(updater)
	}

	// Update specific tool
	toolName := args[0]
	return runUpdateSingle(updater, toolName)
}

// runUpdateSingle updates a single tool
func runUpdateSingle(updater *services.UpdaterService, toolName string) error {
	// Check if tool is outdated
	outdated, err := updater.IsOutdated(toolName)
	if err != nil {
		return fmt.Errorf("failed to check tool status: %w", err)
	}

	if !outdated {
		fmt.Printf("Tool %s is already up-to-date\n", toolName)
		return nil
	}

	// Get current and latest version for confirmation
	installedTool, err := updater.GetInstalledVersion(toolName)
	if err != nil {
		return fmt.Errorf("failed to get installed version: %w", err)
	}

	latestVersion, err := updater.GetLatestVersion(toolName)
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	// Confirmation prompt (unless --yes)
	if !updateYes {
		fmt.Printf("Updating %s from %s to %s\n", toolName, installedTool, latestVersion)
		if !promptConfirmation("Are you sure?") {
			fmt.Println("Update cancelled")
			return nil
		}
	}

	// Perform update
	result, err := updater.Update(toolName)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if result.Skipped {
		fmt.Printf("Tool %s is %s\n", toolName, result.Message)
	} else {
		fmt.Printf("Successfully %s\n", result.Message)
	}

	return nil
}

// runUpdateAll updates all outdated tools
func runUpdateAll(updater *services.UpdaterService) error {
	// Check for outdated tools
	fmt.Println("Checking for outdated tools...")
	outdated, err := updater.CheckOutdated()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if len(outdated) == 0 {
		fmt.Println("All tools are up-to-date!")
		return nil
	}

	// Display outdated tools
	fmt.Printf("Found %d outdated tool(s):\n", len(outdated))
	for _, tool := range outdated {
		fmt.Printf("  - %s: %s → %s\n", tool.Name, tool.CurrentVersion, tool.LatestVersion)
	}
	fmt.Println()

	// Confirmation prompt (unless --yes)
	if !updateYes {
		if !promptConfirmation("Update all tools?") {
			fmt.Println("Update cancelled")
			return nil
		}
		fmt.Println()
	}

	// Update all tools
	results, errors := updater.UpdateAll()

	// Display results
	successCount := 0
	skipCount := 0
	failCount := 0

	for _, result := range results {
		if result.Success {
			if result.Skipped {
				fmt.Printf("%s: %s\n", result.ToolName, result.Message)
				skipCount++
			} else {
				fmt.Printf("Successfully %s\n", result.Message)
				successCount++
			}
		} else {
			fmt.Fprintf(os.Stderr, "Failed to update %s: %v\n", result.ToolName, result.Error)
			failCount++
		}
		fmt.Println()
	}

	// Summary
	fmt.Println("---")
	fmt.Printf("Summary: %d updated, %d skipped, %d failed\n", successCount, skipCount, failCount)

	// Return error if any updates failed
	if len(errors) > 0 {
		return fmt.Errorf("%d tool(s) failed to update", len(errors))
	}

	return nil
}

// promptConfirmation prompts the user for a yes/no confirmation
func promptConfirmation(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
