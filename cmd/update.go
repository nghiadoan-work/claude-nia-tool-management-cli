package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/ui"
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

If no arguments are provided, the command will run in interactive mode
and guide you through selecting tools to update.

Examples:
  cntm update                        # Interactive mode
  cntm update code-reviewer          # Update specific tool
  cntm update --all                  # Update all outdated tools
  cntm update --all --yes            # Update all without confirmation`,
	Example: `  cntm update                        # Interactive mode
  cntm update code-reviewer          # Update specific tool
  cntm update --all                  # Update all outdated tools
  cntm update --all --yes            # Update all without confirmation
  cntm update code-reviewer --yes    # Update without confirmation`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Either provide a tool name, use --all, or run interactive
		if updateAll && len(args) > 0 {
			return fmt.Errorf("cannot specify tool name with --all flag")
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

	registryService := services.NewRegistryServiceWithoutCache(githubClient)

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

	// Interactive mode if no arguments
	if len(args) == 0 {
		return runUpdateInteractive(updater)
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
		ui.PrintInfo("Tool %s is already up-to-date", ui.FormatToolName(toolName))
		return nil
	}

	// Get current and latest version for confirmation
	installedTool, err := updater.GetInstalledVersion(toolName)
	if err != nil {
		return ui.NewValidationError(
			fmt.Sprintf("Failed to get installed version for %s", ui.FormatToolName(toolName)),
			"Ensure the tool is properly installed",
		)
	}

	latestVersion, err := updater.GetLatestVersion(toolName)
	if err != nil {
		return ui.NewNetworkError("fetching latest version", err)
	}

	// Confirmation prompt (unless --yes)
	if !updateYes {
		ui.PrintInfo("Updating %s from %s to %s",
			ui.FormatToolName(toolName),
			ui.FormatVersion(installedTool),
			ui.FormatVersion(latestVersion))

		if !ui.Confirm("Are you sure you want to continue?") {
			ui.PrintWarning("Update cancelled")
			return nil
		}
	}

	// Perform update
	result, err := updater.Update(toolName)
	if err != nil {
		ui.PrintError("Update failed for %s", ui.FormatToolName(toolName))
		ui.PrintHint("Try running 'cntm install --force %s' to force reinstall", toolName)
		return err
	}

	if result.Skipped {
		ui.PrintInfo("Tool %s is %s", ui.FormatToolName(toolName), result.Message)
	} else {
		ui.PrintSuccess("%s", result.Message)
	}

	return nil
}

// runUpdateAll updates all outdated tools
func runUpdateAll(updater *services.UpdaterService) error {
	// Check for outdated tools
	sp := ui.NewSpinner("Checking for outdated tools...")
	sp.Start()

	outdated, err := updater.CheckOutdated()
	sp.Stop()

	if err != nil {
		return ui.NewNetworkError("checking for updates", err)
	}

	if len(outdated) == 0 {
		ui.PrintSuccess("All tools are up-to-date!")
		return nil
	}

	// Display outdated tools
	ui.PrintInfo("Found %d outdated tool(s):", len(outdated))
	for _, tool := range outdated {
		fmt.Printf("  - %s: %s → %s\n",
			ui.FormatToolName(tool.Name),
			ui.FormatVersion(tool.CurrentVersion),
			ui.FormatVersion(tool.LatestVersion))
	}
	fmt.Println()

	// Confirmation prompt (unless --yes)
	if !updateYes {
		if !ui.Confirm("Update all tools?") {
			ui.PrintWarning("Update cancelled")
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
				ui.PrintInfo("%s: %s", ui.FormatToolName(result.ToolName), result.Message)
				skipCount++
			} else {
				ui.PrintSuccess("%s", result.Message)
				successCount++
			}
		} else {
			ui.PrintError("Failed to update %s", ui.FormatToolName(result.ToolName))
			failCount++
		}
		fmt.Println()
	}

	// Summary
	ui.PrintHeader("Update Summary")
	if successCount > 0 {
		ui.PrintSuccess("%d tool(s) updated", successCount)
	}
	if skipCount > 0 {
		ui.PrintInfo("%d tool(s) skipped (already up-to-date)", skipCount)
	}
	if failCount > 0 {
		ui.PrintError("%d tool(s) failed to update", failCount)
	}

	// Return error if any updates failed
	if len(errors) > 0 {
		return ui.NewValidationError(
			fmt.Sprintf("%d tool(s) failed to update", len(errors)),
			"Check the errors above for details",
		)
	}

	return nil
}

// runUpdateInteractive presents an interactive menu for selecting tools to update
func runUpdateInteractive(updater *services.UpdaterService) error {
	fmt.Println()
	ui.PrintHeader("Interactive Tool Update")
	fmt.Println()

	// Check for outdated tools
	ui.PrintInfo("Checking for outdated tools...")
	outdated, err := updater.CheckOutdated()
	if err != nil {
		return ui.NewNetworkError("checking for updates", err)
	}

	if len(outdated) == 0 {
		ui.PrintSuccess("All tools are up-to-date!")
		return nil
	}

	fmt.Println()
	ui.PrintInfo("Found %d outdated tool(s)", len(outdated))
	fmt.Println()

	// Create selection options
	options := make([]string, len(outdated)+1)
	options[0] = fmt.Sprintf("Update all %d tools", len(outdated))

	for i, tool := range outdated {
		options[i+1] = fmt.Sprintf("%-20s  %s → %s",
			tool.Name,
			tool.CurrentVersion,
			tool.LatestVersion)
	}

	// Let user select
	selectedIdx, err := ui.SelectWithArrows("Select what to update", options)
	if err != nil {
		// Check if it's a cancellation (Ctrl+C or Ctrl+D)
		if errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrEOF) {
			fmt.Println()
			fmt.Println(ui.Warning("✗ Update cancelled"))
			return nil
		}
		return fmt.Errorf("selection failed: %w", err)
	}

	fmt.Println()

	// Handle selection
	if selectedIdx == 0 {
		// Update all
		return runUpdateAll(updater)
	}

	// Update specific tool
	selectedTool := outdated[selectedIdx-1]
	return runUpdateSingle(updater, selectedTool.Name)
}
