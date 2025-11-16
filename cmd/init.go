package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/spf13/cobra"
)

var (
	// Init flags
	initPath  string
	initForce bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Claude tools project",
	Long: `Initialize a new project with .claude directory structure.

This command will:
  - Create .claude/ directory
  - Create subdirectories: agents/, commands/, skills/
  - Initialize .claude-lock.json with empty tool list
  - Detect if already initialized and warn (unless --force)

The .claude directory structure:
  .claude/
  ├── agents/          # Agent tools
  ├── commands/        # Command tools
  ├── skills/          # Skill tools
  └── .claude-lock.json # Installed tools lock file

Examples:
  cntm init                         # Initialize in current directory
  cntm init --path /custom/path     # Initialize at custom location
  cntm init --force                 # Reinitialize even if exists`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Init flags
	initCmd.Flags().StringVar(&initPath, "path", "", "custom path for .claude directory (default: current directory)")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "force initialization even if .claude exists")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Determine the base path
	initBasePath := basePath
	if initPath != "" {
		initBasePath = initPath
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(initBasePath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check if .claude directory already exists
	claudeDir := absPath
	if filepath.Base(absPath) != ".claude" {
		claudeDir = filepath.Join(absPath, ".claude")
	}

	if _, err := os.Stat(claudeDir); err == nil {
		// Directory exists
		if !initForce {
			return fmt.Errorf(".claude directory already exists at %s\nUse --force to reinitialize", claudeDir)
		}
		fmt.Printf("Warning: Reinitializing existing .claude directory at %s\n\n", claudeDir)
	}

	// Create .claude directory
	fmt.Printf("Initializing Claude tools project at %s\n", claudeDir)

	// Create main directory
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}
	fmt.Println("  Created .claude/")

	// Create subdirectories
	subdirs := []string{"agents", "commands", "skills"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(claudeDir, subdir)
		if err := os.MkdirAll(subdirPath, 0755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", subdir, err)
		}
		fmt.Printf("  Created .claude/%s/\n", subdir)
	}

	// Initialize lock file
	lockFilePath := filepath.Join(claudeDir, ".claude-lock.json")
	if err := initializeLockFile(lockFilePath); err != nil {
		return fmt.Errorf("failed to initialize lock file: %w", err)
	}
	fmt.Println("  Created .claude-lock.json")

	// Success message
	fmt.Println()
	fmt.Println("Successfully initialized Claude tools project!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Search for tools:    cntm search <query>")
	fmt.Println("  2. Browse tools:        cntm browse")
	fmt.Println("  3. Install a tool:      cntm install <tool-name>")
	fmt.Println("  4. List installed:      cntm list")

	return nil
}

// initializeLockFile creates an empty lock file with proper structure
func initializeLockFile(path string) error {
	// Load config to get registry URL
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		// Use default if config not available
		cfg = &models.Config{
			Registry: models.RegistryConfig{
				URL: "https://github.com/nghiadoan-work/claude-tools-registry",
			},
		}
	}

	// Create empty lock file
	lockFile := &models.LockFile{
		Version:   "1.0",
		UpdatedAt: time.Now(),
		Registry:  cfg.Registry.URL,
		Tools:     make(map[string]*models.InstalledTool),
	}

	// Write to file
	data, err := json.MarshalIndent(lockFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock file: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %w", err)
	}

	return nil
}
