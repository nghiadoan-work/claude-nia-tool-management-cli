package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/services"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	// Outdated flags
	outdatedJSON bool
)

// outdatedCmd represents the outdated command
var outdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Check for outdated tools",
	Long: `Check for tools that have available updates in the registry.

This command compares your locally installed tools with the latest versions
available in the remote registry and displays tools that can be updated.`,
	Example: `  cntm outdated              # Show outdated tools in table format
  cntm outdated --json       # Show outdated tools in JSON format`,
	RunE: runOutdated,
}

func init() {
	rootCmd.AddCommand(outdatedCmd)

	// Outdated flags
	outdatedCmd.Flags().BoolVarP(&outdatedJSON, "json", "j", false, "output in JSON format")
}

func runOutdated(cmd *cobra.Command, args []string) error {
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

	// Show progress message
	if !outdatedJSON && verbose {
		fmt.Fprintln(os.Stderr, "Checking for updates...")
	}

	// Check for outdated tools
	outdated, err := updater.CheckOutdated()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w\nHint: Check your internet connection and registry URL", err)
	}

	// Display results
	if outdatedJSON {
		return outputJSON(outdated)
	}

	return displayOutdatedTools(outdated)
}

// displayOutdatedTools displays outdated tools in a table format
func displayOutdatedTools(outdated []services.OutdatedTool) error {
	if len(outdated) == 0 {
		fmt.Println("All tools are up-to-date!")
		return nil
	}

	// Prepare table data
	headers := []string{"Name", "Current Version", "Latest Version", "Type"}
	var rows [][]string

	for _, tool := range outdated {
		rows = append(rows, []string{
			tool.Name,
			tool.CurrentVersion,
			tool.LatestVersion,
			string(tool.Type),
		})
	}

	// Create table with new API
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeader(headers),
	)

	// Add rows
	for _, row := range rows {
		table.Append(row)
	}

	// Render table
	table.Render()

	// Summary
	fmt.Printf("\n%d tool(s) have updates available\n", len(outdated))
	fmt.Println("Run 'cntm update <name>' to update a specific tool")
	fmt.Println("Run 'cntm update --all' to update all tools")

	return nil
}
