package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadt/claude-nia-tool-management-cli/internal/services"
	"github.com/spf13/cobra"
)

var (
	// Remove flags
	removeYes bool
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove <tool-name> [tool-name2] [...]",
	Aliases: []string{"uninstall", "rm"},
	Short:   "Remove installed tools",
	Long: `Remove one or more installed tools from the local .claude directory.

This command will:
  - Remove the tool directory from .claude/<type>/<name>/
  - Update the .claude-lock.json file
  - Prompt for confirmation before removal (unless --yes is used)

Examples:
  cntm remove code-reviewer           # Remove with confirmation
  cntm remove tool1 tool2 tool3       # Remove multiple tools
  cntm remove --yes old-agent         # Remove without confirmation
  cntm uninstall code-reviewer        # Using alias
  cntm rm code-reviewer               # Using short alias`,
	Args: cobra.MinimumNArgs(1),
	RunE: runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Remove flags
	removeCmd.Flags().BoolVarP(&removeYes, "yes", "y", false, "skip confirmation prompts")
}

func runRemove(cmd *cobra.Command, args []string) error {
	// Initialize services
	lockFilePath := filepath.Join(basePath, ".claude-lock.json")
	lockFileService, err := services.NewLockFileService(lockFilePath)
	if err != nil {
		return fmt.Errorf("failed to create lock file service: %w", err)
	}

	fsManager, err := data.NewFSManager(basePath)
	if err != nil {
		return fmt.Errorf("failed to create file system manager: %w", err)
	}

	// Get list of installed tools
	installedTools, err := lockFileService.ListTools()
	if err != nil {
		return fmt.Errorf("failed to list installed tools: %w\nHint: No tools installed? Use 'cntm list' to see installed tools", err)
	}

	// Validate that all tools exist
	var toolsToRemove []string
	for _, toolName := range args {
		if _, exists := installedTools[toolName]; !exists {
			fmt.Fprintf(os.Stderr, "Warning: Tool '%s' is not installed, skipping\n", toolName)
			continue
		}
		toolsToRemove = append(toolsToRemove, toolName)
	}

	if len(toolsToRemove) == 0 {
		return fmt.Errorf("no valid tools to remove\nHint: Use 'cntm list' to see installed tools")
	}

	// Show what will be removed
	fmt.Println("The following tools will be removed:")
	for _, toolName := range toolsToRemove {
		tool := installedTools[toolName]
		fmt.Printf("  - %s@%s (%s)\n", toolName, tool.Version, tool.Type)
	}
	fmt.Println()

	// Confirmation prompt (unless --yes)
	if !removeYes {
		if !promptConfirmation("Are you sure you want to remove these tools?") {
			fmt.Println("Removal cancelled")
			return nil
		}
		fmt.Println()
	}

	// Remove each tool
	successCount := 0
	failCount := 0

	for _, toolName := range toolsToRemove {
		tool := installedTools[toolName]

		// Construct tool directory path
		toolDir := filepath.Join(basePath, string(tool.Type)+"s", toolName)

		// Remove tool directory from file system
		if err := fsManager.RemoveDir(toolDir); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to remove directory for %s: %v\n", toolName, err)
			failCount++
			continue
		}

		// Remove tool from lock file
		if err := lockFileService.RemoveTool(toolName); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update lock file for %s: %v\n", toolName, err)
			fmt.Fprintf(os.Stderr, "Warning: Directory was removed but lock file not updated\n")
			failCount++
			continue
		}

		fmt.Printf("Successfully removed %s@%s\n", toolName, tool.Version)
		successCount++
	}

	// Display summary for multiple tools
	if len(toolsToRemove) > 1 {
		fmt.Println()
		fmt.Println("---")
		fmt.Printf("Summary: %d removed, %d failed\n", successCount, failCount)
	}

	// Return error if any removals failed
	if failCount > 0 {
		return fmt.Errorf("%d tool(s) failed to remove", failCount)
	}

	return nil
}
