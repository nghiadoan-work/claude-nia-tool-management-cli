package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/services"
	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/spf13/cobra"
)

var (
	// Info flags
	infoType string
	infoJSON bool
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info <name>",
	Short: "Display detailed information about a tool",
	Long: `Display detailed information about a specific tool from the registry.

This command shows all available metadata for a tool including:
  - Name, version, and type
  - Author and description
  - Tags and download count
  - File size and location
  - Creation and update timestamps

Examples:
  cntm info code-reviewer               # Show info for code-reviewer
  cntm info git-helper --type agent     # Show info with specific type
  cntm info code-reviewer --json        # Output in JSON format`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	// Info flags
	infoCmd.Flags().StringVarP(&infoType, "type", "t", "", "tool type (agent, command, skill) - auto-detected if not specified")
	infoCmd.Flags().BoolVarP(&infoJSON, "json", "j", false, "output in JSON format")
}

func runInfo(cmd *cobra.Command, args []string) error {
	toolName := args[0]

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

	// Show progress message
	if !infoJSON && verbose {
		fmt.Fprintln(os.Stderr, "Fetching tool information...")
	}

	var tool *models.ToolInfo

	// If type is specified, get tool directly
	if infoType != "" {
		toolType := models.ToolType(infoType)
		if err := toolType.Validate(); err != nil {
			return fmt.Errorf("invalid tool type: %w", err)
		}

		tool, err = registryService.GetTool(toolName, toolType)
		if err != nil {
			return fmt.Errorf("tool not found: %w\nHint: Try searching with 'cntm search %s'", err, toolName)
		}
	} else {
		// Try to find the tool by searching all types
		tool, err = findToolByName(registryService, toolName)
		if err != nil {
			return fmt.Errorf("tool not found: %w\nHint: Try searching with 'cntm search %s' or specify --type", err, toolName)
		}
	}

	// Display results
	if infoJSON {
		return outputJSON(tool)
	}

	return displayToolInfo(tool)
}

// findToolByName searches for a tool by name across all types
func findToolByName(registryService *services.RegistryService, name string) (*models.ToolInfo, error) {
	// Try each tool type
	types := []models.ToolType{
		models.ToolTypeAgent,
		models.ToolTypeCommand,
		models.ToolTypeSkill,
	}

	for _, toolType := range types {
		tool, err := registryService.GetTool(name, toolType)
		if err == nil {
			return tool, nil
		}
	}

	return nil, fmt.Errorf("tool '%s' not found in registry", name)
}

// displayToolInfo displays detailed tool information in a readable format
func displayToolInfo(tool *models.ToolInfo) error {
	fmt.Println()
	fmt.Printf("  Name:        %s\n", tool.Name)
	fmt.Printf("  Version:     %s\n", tool.Version)
	fmt.Printf("  Type:        %s\n", tool.Type)
	fmt.Printf("  Author:      %s\n", tool.Author)
	fmt.Println()
	fmt.Printf("  Description: %s\n", tool.Description)
	fmt.Println()

	if len(tool.Tags) > 0 {
		fmt.Printf("  Tags:        %s\n", strings.Join(tool.Tags, ", "))
	}

	fmt.Printf("  Downloads:   %d\n", tool.Downloads)
	fmt.Printf("  Size:        %s\n", formatBytes(tool.Size))
	fmt.Println()
	fmt.Printf("  Created:     %s\n", tool.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Updated:     %s\n", tool.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Printf("  File:        %s\n", tool.File)
	fmt.Println()

	return nil
}

// formatBytes converts bytes to a human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
