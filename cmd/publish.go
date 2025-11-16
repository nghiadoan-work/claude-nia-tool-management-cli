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
	Use:   "publish [type] [name]",
	Short: "Publish a tool to the registry",
	Long: `Publish a Claude Code tool to the registry.

Tool types: agent, command, skill

This will:
1. Validate the tool directory
2. Generate or update metadata
3. Create a ZIP package
4. Calculate integrity hash
5. Provide instructions for creating a PR to the registry

Examples:
  cntm publish                      # Interactive mode - choose from available tools
  cntm publish agent my-agent
  cntm publish skill docker-patterns --version 1.0.0
  cntm publish command test-runner --version 1.1.0 --changelog "Added new features"
  cntm publish agent code-reviewer --force`,
	Args: cobra.RangeArgs(0, 2),
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

	var toolType models.ToolType
	var toolName string
	var toolPath string

	// Interactive mode: no arguments provided
	if len(args) == 0 {
		// Scan for available tools
		tools, err := scanLocalTools(cfg)
		if err != nil {
			return fmt.Errorf("failed to scan local tools: %w", err)
		}

		if len(tools) == 0 {
			return fmt.Errorf("no tools found in %s\nCreate a tool first with: cntm create", cfg.Local.DefaultPath)
		}

		// Let user select a tool
		selectedTool, err := selectToolInteractively(tools)
		if err != nil {
			return err
		}

		toolType = selectedTool.Type
		toolName = selectedTool.Name
		toolPath = selectedTool.Path

		fmt.Printf("\nSelected: %s (%s)\n", toolName, toolType)
	} else if len(args) == 2 {
		// Explicit mode: type and name provided
		toolTypeStr := strings.ToLower(args[0])
		toolName = args[1]

		// Validate tool type
		switch toolTypeStr {
		case "agent", "agents":
			toolType = models.ToolTypeAgent
		case "command", "commands":
			toolType = models.ToolTypeCommand
		case "skill", "skills":
			toolType = models.ToolTypeSkill
		default:
			return fmt.Errorf("invalid tool type: %s\nValid types: agent, command, skill", toolTypeStr)
		}

		// Determine tool path
		if publishPath != "" {
			toolPath = publishPath
		} else {
			// Build path based on type and name
			toolPath = filepath.Join(cfg.Local.DefaultPath, string(toolType)+"s", toolName)

			// Check if tool exists
			if _, err := os.Stat(toolPath); os.IsNotExist(err) {
				return fmt.Errorf("tool %s not found at %s\nHint: Use --path to specify a custom location", toolName, toolPath)
			}
		}
	} else {
		return fmt.Errorf("invalid arguments\nUsage: cntm publish [type] [name] OR cntm publish (interactive)")
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

	publishMeta := &services.PublishMetadata{
		Name:    toolName,
		Version: version,
		Type:    toolType,
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
		fmt.Printf("  Type:    %s\n", string(toolType))
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

// toolInfo represents information about a local tool
type toolInfo struct {
	Name string
	Type models.ToolType
	Path string
}

// scanLocalTools scans the local directories for available tools
func scanLocalTools(cfg *models.Config) ([]toolInfo, error) {
	baseDir := cfg.Local.DefaultPath
	var tools []toolInfo

	// Tool types to scan
	toolTypes := []struct {
		dir      string
		toolType models.ToolType
	}{
		{"agents", models.ToolTypeAgent},
		{"commands", models.ToolTypeCommand},
		{"skills", models.ToolTypeSkill},
	}

	for _, tt := range toolTypes {
		dirPath := filepath.Join(baseDir, tt.dir)

		// Check if directory exists
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			continue
		}

		// Read directory entries
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			continue // Skip on error
		}

		// Add each subdirectory as a tool
		for _, entry := range entries {
			if entry.IsDir() {
				tools = append(tools, toolInfo{
					Name: entry.Name(),
					Type: tt.toolType,
					Path: filepath.Join(dirPath, entry.Name()),
				})
			}
		}
	}

	return tools, nil
}

// selectToolInteractively presents a menu for selecting a tool to publish
func selectToolInteractively(tools []toolInfo) (*toolInfo, error) {
	if len(tools) == 0 {
		return nil, fmt.Errorf("no tools found in local directories")
	}

	fmt.Println("\nAvailable tools to publish:")
	fmt.Println()

	// Group by type for better display
	agentTools := []toolInfo{}
	commandTools := []toolInfo{}
	skillTools := []toolInfo{}

	for _, tool := range tools {
		switch tool.Type {
		case models.ToolTypeAgent:
			agentTools = append(agentTools, tool)
		case models.ToolTypeCommand:
			commandTools = append(commandTools, tool)
		case models.ToolTypeSkill:
			skillTools = append(skillTools, tool)
		}
	}

	// Display tools with numbers
	index := 1
	toolMap := make(map[int]toolInfo)

	if len(agentTools) > 0 {
		fmt.Println("Agents:")
		for _, tool := range agentTools {
			fmt.Printf("  %d) %s\n", index, tool.Name)
			toolMap[index] = tool
			index++
		}
		fmt.Println()
	}

	if len(commandTools) > 0 {
		fmt.Println("Commands:")
		for _, tool := range commandTools {
			fmt.Printf("  %d) %s\n", index, tool.Name)
			toolMap[index] = tool
			index++
		}
		fmt.Println()
	}

	if len(skillTools) > 0 {
		fmt.Println("Skills:")
		for _, tool := range skillTools {
			fmt.Printf("  %d) %s\n", index, tool.Name)
			toolMap[index] = tool
			index++
		}
		fmt.Println()
	}

	// Prompt for selection
	fmt.Printf("Select tool to publish (1-%d): ", len(tools))

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	// Parse selection
	var selection int
	if _, err := fmt.Sscanf(input, "%d", &selection); err != nil {
		return nil, fmt.Errorf("invalid selection: %s", input)
	}

	// Validate selection
	selectedTool, exists := toolMap[selection]
	if !exists {
		return nil, fmt.Errorf("invalid selection: %d (must be between 1 and %d)", selection, len(tools))
	}

	return &selectedTool, nil
}
