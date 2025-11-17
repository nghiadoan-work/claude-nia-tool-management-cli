package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/ui"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/spf13/cobra"
)

var (
	// Install flags
	installForce bool
	installPath  string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install [tool-name][@<version>] [tool-name2] [...]",
	Short: "Install Claude Code tools from the registry",
	Long: `Install one or more tools from the remote registry.

By default, this command installs the latest version of a tool.
You can specify a version using the @version syntax.

If no arguments are provided, the command will run in interactive mode
and guide you through selecting a tool to install.

Installation locations:
  - Agents:   .claude/agents/<name>/
  - Commands: .claude/commands/<name>/
  - Skills:   .claude/skills/<name>/`,
	Example: `  cntm install                            # Interactive mode
  cntm install code-reviewer              # Install latest version
  cntm install code-reviewer@1.0.0        # Install specific version
  cntm install agent1 agent2 agent3       # Install multiple tools
  cntm install --force code-reviewer      # Force reinstall
  cntm install --path /custom code-reviewer # Custom install path`,
	RunE: runInstall,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Install flags
	installCmd.Flags().BoolVarP(&installForce, "force", "f", false, "force reinstall even if already installed")
	installCmd.Flags().StringVar(&installPath, "path", "", "custom installation path (overrides default .claude directory)")
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return ui.NewValidationError(
			"Failed to load configuration",
			"Run 'cntm init' to initialize the project or check your config file",
		)
	}

	// Determine base path for installation
	installBasePath := basePath
	if installPath != "" {
		installBasePath = installPath
	}

	// Override config with custom path if specified
	if installPath != "" {
		cfg.Local.DefaultPath = installPath
	}

	// Parse GitHub URL to get owner and repo
	owner, repo, err := parseGitHubURL(cfg.Registry.URL)
	if err != nil {
		return ui.NewValidationError(
			"Invalid registry URL in configuration",
			fmt.Sprintf("Check the registry URL in your config: %s", ui.FormatURL(cfg.Registry.URL)),
		)
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
	fsManager, err := data.NewFSManager(installBasePath)
	if err != nil {
		return fmt.Errorf("failed to create file system manager: %w", err)
	}

	lockFilePath := filepath.Join(installBasePath, ".claude-lock.json")
	lockFileService, err := services.NewLockFileService(lockFilePath)
	if err != nil {
		return fmt.Errorf("failed to create lock file service: %w", err)
	}
	lockFileService.SetRegistry(cfg.Registry.URL)

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

	// Parse tool arguments or run interactive mode
	var toolsToInstall []toolSpec
	isInteractive := len(args) == 0

	if isInteractive {
		// Interactive mode
		toolSpec, err := selectToolInteractivelyForInstall(registryService)
		if err != nil {
			// Check if it's a cancellation (Ctrl+C or Ctrl+D)
			if errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrEOF) {
				fmt.Println()
				fmt.Println(ui.Warning("✗ Installation cancelled"))
				return nil
			}
			return fmt.Errorf("interactive selection failed: %w", err)
		}
		toolsToInstall = append(toolsToInstall, *toolSpec)
	} else {
		// Parse from arguments
		for _, arg := range args {
			name, version := parseToolArg(arg)
			toolsToInstall = append(toolsToInstall, toolSpec{
				name:    name,
				version: version,
			})
		}
	}

	// Install tools
	successCount := 0
	skipCount := 0
	failCount := 0

	for _, spec := range toolsToInstall {
		// Check if already installed (unless force is set or in interactive mode)
		// In interactive mode, automatically reinstall if already installed
		if !installForce && !isInteractive {
			installed, err := installer.IsInstalled(spec.name)
			if err == nil && installed {
				// Check version
				installedVersion, err := installer.GetInstalledVersion(spec.name)
				if err == nil {
					// If no version specified or version matches, skip
					if spec.version == "" || spec.version == installedVersion {
						ui.PrintWarning("Tool %s is already installed (version %s)",
							ui.FormatToolName(spec.name),
							ui.FormatVersion(installedVersion))
						ui.PrintHint("Use --force to reinstall")
						fmt.Println()
						skipCount++
						continue
					}
				}
			}
		}

		// Install the tool
		var err error
		displayName := spec.name
		if spec.version != "" {
			displayName = spec.name + "@" + spec.version
			err = installer.InstallWithVersion(spec.name, spec.version)
		} else {
			err = installer.Install(spec.name)
		}

		if err != nil {
			ui.PrintError("Failed to install %s", ui.FormatToolName(displayName))
			// Print the actual error message (includes available versions if version not found)
			fmt.Fprintf(os.Stderr, "  Error: %s\n", err.Error())
			if strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "Available versions") {
				ui.PrintHint("Run 'cntm search %s' to find similar tools", spec.name)
			} else if strings.Contains(err.Error(), "network") || strings.Contains(err.Error(), "connection") {
				ui.PrintHint("Check your internet connection and try again")
			}
			fmt.Fprintln(os.Stderr)
			failCount++
			continue
		}

		successCount++
		fmt.Println() // Add spacing between tools
	}

	// Display summary for multiple tools
	if len(toolsToInstall) > 1 {
		ui.PrintHeader("Installation Summary")
		if successCount > 0 {
			ui.PrintSuccess("%d tool(s) installed", successCount)
		}
		if skipCount > 0 {
			ui.PrintWarning("%d tool(s) skipped (already installed)", skipCount)
		}
		if failCount > 0 {
			ui.PrintError("%d tool(s) failed to install", failCount)
		}
		fmt.Println()
	}

	// Return error if any installations failed
	if failCount > 0 {
		return ui.NewValidationError(
			fmt.Sprintf("%d tool(s) failed to install", failCount),
			"Check the errors above for details",
		)
	}

	return nil
}

// toolSpec represents a parsed tool specification
type toolSpec struct {
	name    string
	version string
}

// parseToolArg parses a tool argument in the format "name[@version]"
func parseToolArg(arg string) (name, version string) {
	parts := strings.SplitN(arg, "@", 2)
	name = parts[0]
	if len(parts) > 1 {
		version = parts[1]
	}
	return
}

// selectToolInteractivelyForInstall guides the user through selecting a tool to install
func selectToolInteractivelyForInstall(registryService *services.RegistryService) (*toolSpec, error) {
	fmt.Println()
	ui.PrintHeader("Interactive Tool Installation")
	fmt.Println()

	// Step 1: Select tool type
	ui.PrintInfo("Step 1: Select tool type")
	typeOptions := []string{
		"agent   - AI agents for specialized tasks",
		"command - Command-line tools and workflows",
		"skill   - Reusable skills and capabilities",
	}

	typeIdx, err := ui.SelectWithArrows("Select tool type", typeOptions)
	if err != nil {
		return nil, err  // Return original error to preserve error type
	}

	var toolType string
	switch typeIdx {
	case 0:
		toolType = "agent"
	case 1:
		toolType = "command"
	case 2:
		toolType = "skill"
	}

	fmt.Println()

	// Step 2: Fetch and display available tools
	ui.PrintInfo("Step 2: Select tool")
	fmt.Printf("Fetching available %ss from registry...\n", toolType)

	registry, err := registryService.GetRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry: %w", err)
	}

	tools := registry.Tools[models.ToolType(toolType)]
	if len(tools) == 0 {
		return nil, fmt.Errorf("no %ss found in registry", toolType)
	}

	// Create tool selection options
	toolOptions := make([]string, len(tools))
	for i, tool := range tools {
		toolOptions[i] = fmt.Sprintf("%-20s - %s (latest: %s)",
			tool.Name,
			tool.Description,
			tool.LatestVersion)
	}

	toolIdx, err := ui.SelectWithArrows(fmt.Sprintf("Select %s to install", toolType), toolOptions)
	if err != nil {
		return nil, err  // Return original error to preserve error type
	}

	selectedTool := tools[toolIdx]
	fmt.Println()

	// Step 3: Select version
	ui.PrintInfo("Step 3: Select version")

	versions := selectedTool.ListVersions()
	versionOptions := make([]string, len(versions)+1)
	versionOptions[0] = fmt.Sprintf("%s (latest)", selectedTool.LatestVersion)
	for i, v := range versions {
		if v == selectedTool.LatestVersion {
			continue // Skip latest since it's already at index 0
		}
		versionOptions[i+1] = v
	}

	// Remove empty slots
	cleanVersionOptions := []string{}
	for _, v := range versionOptions {
		if v != "" {
			cleanVersionOptions = append(cleanVersionOptions, v)
		}
	}

	versionIdx, err := ui.SelectWithArrows("Select version to install", cleanVersionOptions)
	if err != nil {
		return nil, err  // Return original error to preserve error type
	}

	var selectedVersion string
	if versionIdx == 0 {
		selectedVersion = "" // Empty means latest
	} else {
		selectedVersion = versions[versionIdx-1]
	}

	fmt.Println()
	ui.PrintSuccess("Selected: %s@%s", selectedTool.Name,
		func() string {
			if selectedVersion == "" {
				return selectedTool.LatestVersion
			}
			return selectedVersion
		}())
	fmt.Println()

	return &toolSpec{
		name:    selectedTool.Name,
		version: selectedVersion,
	}, nil
}
