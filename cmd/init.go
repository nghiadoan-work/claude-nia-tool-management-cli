package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/config"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/spf13/cobra"
)

//go:embed templates/AGENT_TEMPLATE_GUIDE.md
var agentTemplateGuide string

//go:embed templates/SKILL_TEMPLATE_GUIDE.md
var skillTemplateGuide string

//go:embed templates/COMMAND_TEMPLATE_GUIDE.md
var commandTemplateGuide string

//go:embed templates/.cntm.env
var cntmEnvTemplate string

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
  - Create template guides for creating tools
  - Create .cntm.env for registry configuration
  - Detect if already initialized and warn (unless --force)

The .claude directory structure:
  .claude/
  ├── agents/                   # Agent tools
  ├── commands/                 # Command tools
  ├── skills/                   # Skill tools
  ├── AGENT_TEMPLATE_GUIDE.md   # Guide for creating agents
  ├── SKILL_TEMPLATE_GUIDE.md   # Guide for creating skills
  ├── COMMAND_TEMPLATE_GUIDE.md # Guide for creating commands
  └── .claude-lock.json         # Installed tools lock file

Project root:
  .cntm.env                      # Environment configuration (registry URL, token, etc.)

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

	// Create template guide files
	if err := createTemplateGuides(claudeDir); err != nil {
		return fmt.Errorf("failed to create template guides: %w", err)
	}
	fmt.Println("  Created template guides")

	// Create .cntm.env file in project root
	projectRoot := filepath.Dir(claudeDir)
	if filepath.Base(absPath) == ".claude" {
		projectRoot = filepath.Dir(absPath)
	}
	if err := createEnvFile(projectRoot); err != nil {
		return fmt.Errorf("failed to create .cntm.env: %w", err)
	}
	fmt.Println("  Created .cntm.env")

	// Success message
	fmt.Println()
	fmt.Println("Successfully initialized Claude tools project!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Configure registry:  Edit .cntm.env to set registry URL and token")
	fmt.Println("  2. Search for tools:    cntm search <query>")
	fmt.Println("  3. Install a tool:      cntm install <tool-name>")
	fmt.Println("  4. Update tools:        cntm update --all")
	fmt.Println("  5. Publish your tool:   cntm publish")
	fmt.Println()
	fmt.Println("Note: Add .cntm.env to .gitignore if it contains sensitive tokens")

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

// createTemplateGuides creates the template guide files in the .claude directory
func createTemplateGuides(claudeDir string) error {
	guides := map[string]string{
		"AGENT_TEMPLATE_GUIDE.md":   agentTemplateGuide,
		"SKILL_TEMPLATE_GUIDE.md":   skillTemplateGuide,
		"COMMAND_TEMPLATE_GUIDE.md": commandTemplateGuide,
	}

	for filename, content := range guides {
		filePath := filepath.Join(claudeDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	return nil
}

// createEnvFile creates the .cntm.env file in the project root
func createEnvFile(projectRoot string) error {
	envPath := filepath.Join(projectRoot, ".cntm.env")

	// Don't overwrite existing .cntm.env file
	if _, err := os.Stat(envPath); err == nil {
		return nil // File already exists, skip
	}

	if err := os.WriteFile(envPath, []byte(cntmEnvTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write .cntm.env: %w", err)
	}

	return nil
}
