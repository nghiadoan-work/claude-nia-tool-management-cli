package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/spf13/cobra"
)

var publishCmd = &cobra.Command{
	Use:   "publish [name]",
	Short: "Publish a tool to the registry",
	Long: `Publish a Claude Code tool to the registry.

This will:
1. Validate the tool directory
2. Generate or update metadata
3. Create a ZIP package
4. Calculate integrity hash
5. Provide instructions for creating a PR to the registry

Examples:
  cntm publish my-agent
  cntm publish my-agent --version 1.0.0
  cntm publish my-agent --version 1.1.0 --changelog "Added new features"
  cntm publish my-agent --force`,
	Args: cobra.ExactArgs(1),
	RunE: runPublish,
}

var (
	publishVersion   string
	publishChangelog string
	publishForce     bool
	publishPath      string
)

func init() {
	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringVar(&publishVersion, "version", "", "Version to publish (required)")
	publishCmd.Flags().StringVar(&publishChangelog, "changelog", "", "Changelog entry for this version")
	publishCmd.Flags().BoolVar(&publishForce, "force", false, "Skip confirmation prompts")
	publishCmd.Flags().StringVar(&publishPath, "path", "", "Custom path to tool directory")
}

func runPublish(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	toolName := args[0]

	// Determine tool path
	var toolPath string
	if publishPath != "" {
		toolPath = publishPath
	} else {
		// Try to find the tool in the default locations
		toolPath = findToolPath(toolName, cfg)
		if toolPath == "" {
			return fmt.Errorf("tool %s not found in local directories\nHint: Use --path to specify a custom location", toolName)
		}
	}

	fmt.Printf("Publishing tool: %s\n", toolName)
	fmt.Printf("Path: %s\n", toolPath)

	// Create services
	basePath := cfg.Local.DefaultPath
	if basePath == "" {
		basePath = ".claude"
	}
	fsManager, err := data.NewFSManager(basePath)
	if err != nil {
		return fmt.Errorf("failed to create fs manager: %w", err)
	}

	owner, repo, err := parseGitHubURL(cfg.Registry.URL)
	if err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}

	githubClient := services.NewGitHubClient(services.GitHubClientConfig{
		Owner:     owner,
		Repo:      repo,
		Branch:    cfg.Registry.Branch,
		AuthToken: cfg.Registry.AuthToken,
	})

	cacheDir := cfg.Local.DefaultPath + "/.cache"
	cacheManager, err := data.NewCacheManager(cacheDir, 3600*time.Second)
	if err != nil {
		return fmt.Errorf("failed to create cache manager: %w", err)
	}

	registryService := services.NewRegistryService(githubClient, cacheManager)

	publisherService, err := services.NewPublisherService(
		fsManager,
		githubClient,
		registryService,
		cfg,
	)
	if err != nil {
		return fmt.Errorf("failed to create publisher service: %w", err)
	}

	// Step 1: Validate tool
	fmt.Println("\nValidating tool...")
	if err := publisherService.ValidateTool(toolPath); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	fmt.Println("Validation passed")

	// Step 2: Read existing metadata
	existingMeta, err := publisherService.ReadExistingMetadata(toolPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	// Step 3: Determine version
	version := publishVersion
	if version == "" {
		if existingMeta != nil && existingMeta.Version != "" {
			// Suggest next version
			suggestedVersion := bumpVersion(existingMeta.Version)
			version, err = promptString(fmt.Sprintf("Version (current: %s)", existingMeta.Version), suggestedVersion)
			if err != nil {
				return err
			}
		} else {
			version, err = promptString("Version", "1.0.0")
			if err != nil {
				return err
			}
		}
	}

	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Step 4: Get changelog entry
	changelog := publishChangelog
	if changelog == "" && !publishForce {
		changelog, err = promptString(fmt.Sprintf("Changelog for %s", version), "")
		if err != nil {
			return err
		}
	}

	// Step 5: Update metadata
	fmt.Println("\nUpdating metadata...")

	// Detect tool type
	toolType, err := detectToolTypeFromPath(toolPath)
	if err != nil {
		return fmt.Errorf("failed to detect tool type: %w", err)
	}

	publishMeta := &services.PublishMetadata{
		Name:    toolName,
		Version: version,
		Type:    models.ToolType(toolType),
	}

	// Copy from existing metadata or prompt
	if existingMeta != nil {
		publishMeta.Author = existingMeta.Author
		publishMeta.Description = existingMeta.Description
		publishMeta.Tags = existingMeta.Tags
		publishMeta.Changelog = existingMeta.Changelog
		publishMeta.Dependencies = existingMeta.Dependencies
	}

	// Ensure required fields
	if publishMeta.Author == "" {
		publishMeta.Author = cfg.Publish.DefaultAuthor
		if publishMeta.Author == "" && !publishForce {
			publishMeta.Author, err = promptString("Author", "")
			if err != nil {
				return err
			}
		}
	}

	if publishMeta.Description == "" && !publishForce {
		publishMeta.Description, err = promptString("Description", "")
		if err != nil {
			return err
		}
	}

	// Add changelog entry
	if publishMeta.Changelog == nil {
		publishMeta.Changelog = make(map[string]string)
	}
	if changelog != "" {
		publishMeta.Changelog[version] = changelog
	} else if _, exists := publishMeta.Changelog[version]; !exists {
		publishMeta.Changelog[version] = "Release " + version
	}

	// Generate metadata.json
	if err := publisherService.GenerateMetadata(toolPath, publishMeta); err != nil {
		return fmt.Errorf("failed to generate metadata: %w", err)
	}
	fmt.Println("Metadata updated")

	// Step 6: Confirm publication
	if !publishForce {
		fmt.Printf("\nReady to publish:\n")
		fmt.Printf("  Tool:    %s\n", toolName)
		fmt.Printf("  Type:    %s\n", toolType)
		fmt.Printf("  Version: %s\n", version)
		fmt.Printf("  Author:  %s\n", publishMeta.Author)
		fmt.Print("\nContinue? (y/n): ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Println("Publication cancelled")
			return nil
		}
	}

	// Step 7: Publish to registry
	fmt.Println("\nPublishing to registry...")
	if err := publisherService.PublishToRegistry(toolPath, version); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	fmt.Println("\nPublication complete!")
	return nil
}

// findToolPath searches for a tool in the default local directories
func findToolPath(toolName string, cfg *models.Config) string {
	baseDir := cfg.Local.DefaultPath

	// Check in agents
	agentPath := filepath.Join(baseDir, "agents", toolName)
	if _, err := os.Stat(agentPath); err == nil {
		return agentPath
	}

	// Check in commands
	commandPath := filepath.Join(baseDir, "commands", toolName)
	if _, err := os.Stat(commandPath); err == nil {
		return commandPath
	}

	// Check in skills
	skillPath := filepath.Join(baseDir, "skills", toolName)
	if _, err := os.Stat(skillPath); err == nil {
		return skillPath
	}

	return ""
}

// detectToolTypeFromPath detects the tool type from its path
func detectToolTypeFromPath(toolPath string) (string, error) {
	absPath, err := filepath.Abs(toolPath)
	if err != nil {
		return "", err
	}

	pathLower := strings.ToLower(absPath)

	if strings.Contains(pathLower, "/agents/") || strings.Contains(pathLower, "\\agents\\") {
		return "agent", nil
	}
	if strings.Contains(pathLower, "/commands/") || strings.Contains(pathLower, "\\commands\\") {
		return "command", nil
	}
	if strings.Contains(pathLower, "/skills/") || strings.Contains(pathLower, "\\skills\\") {
		return "skill", nil
	}

	return "", fmt.Errorf("could not detect tool type from path")
}

// bumpVersion suggests the next version based on auto version bump config
func bumpVersion(currentVersion string) string {
	// Simple version bumping (patch version)
	// In a full implementation, this would use semver logic
	parts := strings.Split(currentVersion, ".")
	if len(parts) != 3 {
		return "1.0.0"
	}

	// For now, just suggest incrementing patch version
	major := parts[0]
	minor := parts[1]
	patch := "0"

	// Try to parse patch version and increment
	var patchNum int
	if _, err := fmt.Sscanf(parts[2], "%d", &patchNum); err == nil {
		patchNum++
		patch = fmt.Sprintf("%d", patchNum)
	}

	return fmt.Sprintf("%s.%s.%s", major, minor, patch)
}
