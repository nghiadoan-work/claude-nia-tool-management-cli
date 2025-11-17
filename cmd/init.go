package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/spf13/cobra"
)

//go:embed templates/AGENT_TEMPLATE_GUIDE.md
var agentTemplateGuide string

//go:embed templates/SKILL_TEMPLATE_GUIDE.md
var skillTemplateGuide string

//go:embed templates/COMMAND_TEMPLATE_GUIDE.md
var commandTemplateGuide string

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

Configuration:
  - Global config: ~/.claude-tools-config.yaml (optional)
  - Project config: .claude-tools-config.yaml (optional, overrides global)

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

	claudeDirExists := false
	if _, err := os.Stat(claudeDir); err == nil {
		claudeDirExists = true
		if initForce {
			fmt.Printf("Warning: Reinitializing existing .claude directory at %s\n\n", claudeDir)
		} else {
			fmt.Printf("Checking existing .claude directory at %s\n", claudeDir)
		}
	}

	// Create .claude directory
	if !claudeDirExists {
		fmt.Printf("Initializing Claude tools project at %s\n", claudeDir)
	}

	// Create main directory if it doesn't exist
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}
	if !claudeDirExists {
		fmt.Println("  Created .claude/")
	}

	// Create subdirectories (only if they don't exist)
	subdirs := []string{"agents", "commands", "skills"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(claudeDir, subdir)
		if _, err := os.Stat(subdirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(subdirPath, 0755); err != nil {
				return fmt.Errorf("failed to create %s directory: %w", subdir, err)
			}
			fmt.Printf("  Created .claude/%s/\n", subdir)
		}
	}

	// Initialize lock file (only if it doesn't exist or force flag is set)
	lockFilePath := filepath.Join(claudeDir, ".claude-lock.json")
	if _, err := os.Stat(lockFilePath); os.IsNotExist(err) || initForce {
		if err := initializeLockFile(lockFilePath); err != nil {
			return fmt.Errorf("failed to initialize lock file: %w", err)
		}
		fmt.Println("  Created .claude-lock.json")
	}

	// Create template guide files (only if they don't exist or force flag is set)
	if err := createTemplateGuides(claudeDir); err != nil {
		return fmt.Errorf("failed to create template guides: %w", err)
	}

	// Create config template (in project root, not .claude directory)
	projectRoot := filepath.Dir(claudeDir)
	if filepath.Base(absPath) == ".claude" {
		projectRoot = filepath.Dir(absPath)
	}
	configPath := filepath.Join(projectRoot, ".claude-tools-config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) || initForce {
		if err := createConfigTemplate(configPath); err != nil {
			return fmt.Errorf("failed to create config template: %w", err)
		}
		fmt.Println("  Created .claude-tools-config.yaml template")
	}

	// Success message
	fmt.Println()
	fmt.Println("Successfully initialized Claude tools project!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Configure registry:  Edit .claude-tools-config.yaml and add your registry URL")
	fmt.Println("  2. Search for tools:    cntm search <query>")
	fmt.Println("  3. Install a tool:      cntm install <tool-name>")
	fmt.Println("  4. Update tools:        cntm update --all")
	fmt.Println("  5. Publish your tool:   cntm publish")

	return nil
}

// initializeLockFile creates an empty lock file with proper structure
func initializeLockFile(path string) error {
	// Create empty lock file with no registry URL
	// User will configure registry URL in .claude-tools-config.yaml
	lockFile := &models.LockFile{
		Version:   "1.0",
		UpdatedAt: time.Now(),
		Registry:  "", // Will be populated from config when tools are installed
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

	createdAny := false
	for filename, content := range guides {
		filePath := filepath.Join(claudeDir, filename)
		// Only create if doesn't exist or force flag is set
		if _, err := os.Stat(filePath); os.IsNotExist(err) || initForce {
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", filename, err)
			}
			createdAny = true
		}
	}

	if createdAny {
		fmt.Println("  Created template guides")
	}

	return nil
}

// createConfigTemplate creates a template .claude-tools-config.yaml file
func createConfigTemplate(path string) error {
	template := `# Claude Tools Configuration
# This file configures the Claude tools package manager (cntm)

# Registry configuration - specify where to fetch tools from
registry:
  url: ""  # REQUIRED: Add your registry URL here (e.g., https://github.com/your-org/your-registry)
  branch: main
  auth_token: ""  # Optional: GitHub Personal Access Token for private repositories

# Local configuration
local:
  default_path: .claude
  auto_update_check: true
  update_check_interval: 86400  # Check for updates every 24 hours (in seconds)

# Publishing configuration
publish:
  default_author: ""  # Optional: Your name or organization
  auto_version_bump: patch  # Options: patch, minor, major
  create_pr: true  # Create pull request when publishing
`

	if err := os.WriteFile(path, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to write config template: %w", err)
	}

	return nil
}
