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

var createCmd = &cobra.Command{
	Use:   "create [type] [name]",
	Short: "Create a new tool locally",
	Long: `Create a new Claude Code tool locally with the specified type and name.

Tool types: agent, command, skill

Examples:
  cntm create agent my-agent
  cntm create command test-runner
  cntm create skill docker-patterns
  cntm create  # Interactive mode`,
	Args: cobra.MaximumNArgs(2),
	RunE: runCreate,
}

var (
	createAuthor      string
	createDescription string
	createTags        string
	createVersion     string
	createInteractive bool
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createAuthor, "author", "", "Tool author name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Tool description")
	createCmd.Flags().StringVar(&createTags, "tags", "", "Comma-separated tags")
	createCmd.Flags().StringVar(&createVersion, "version", "1.0.0", "Initial version")
	createCmd.Flags().BoolVarP(&createInteractive, "interactive", "i", false, "Interactive mode")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var toolType models.ToolType
	var toolName string

	// Determine tool type and name
	if len(args) >= 1 {
		// Try to parse first arg as type
		switch strings.ToLower(args[0]) {
		case "agent", "agents":
			toolType = models.ToolTypeAgent
		case "command", "commands":
			toolType = models.ToolTypeCommand
		case "skill", "skills":
			toolType = models.ToolTypeSkill
		default:
			// If not a type, treat as name and prompt for type
			toolName = args[0]
		}
	}

	if len(args) >= 2 {
		toolName = args[1]
	}

	// Interactive prompts for missing information
	if toolType == "" {
		var err error
		toolType, err = promptToolType()
		if err != nil {
			return err
		}
	}

	if toolName == "" {
		var err error
		toolName, err = promptString("Tool name", "")
		if err != nil {
			return err
		}
		if toolName == "" {
			return fmt.Errorf("tool name cannot be empty")
		}
	}

	// Get or prompt for metadata
	author := createAuthor
	if author == "" {
		// Try to get from config
		author = cfg.Publish.DefaultAuthor
		if author == "" {
			var err error
			author, err = promptString("Author name", "")
			if err != nil {
				return err
			}
		}
	}

	description := createDescription
	if description == "" {
		var err error
		description, err = promptString("Description", "")
		if err != nil {
			return err
		}
	}

	tags := []string{}
	if createTags != "" {
		tags = strings.Split(createTags, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	} else {
		tagsStr, err := promptString("Tags (comma-separated)", "")
		if err != nil {
			return err
		}
		if tagsStr != "" {
			tags = strings.Split(tagsStr, ",")
			for i, tag := range tags {
				tags[i] = strings.TrimSpace(tag)
			}
		}
	}

	version := createVersion
	if version == "" {
		version = "1.0.0"
	}

	// Create tool directory
	toolPath := filepath.Join(cfg.Local.DefaultPath, string(toolType)+"s", toolName)

	// Check if directory already exists
	if _, err := os.Stat(toolPath); err == nil {
		return fmt.Errorf("tool directory already exists: %s", toolPath)
	}

	// Create directory
	if err := os.MkdirAll(toolPath, 0755); err != nil {
		return fmt.Errorf("failed to create tool directory: %w", err)
	}

	fmt.Printf("Creating %s: %s\n", toolType, toolName)
	fmt.Printf("Location: %s\n", toolPath)

	// Create README.md
	if err := createReadme(toolPath, toolType, toolName, description); err != nil {
		return fmt.Errorf("failed to create README: %w", err)
	}

	// Create type-specific files
	if err := createTypeSpecificFiles(toolPath, toolType, toolName, description); err != nil {
		return fmt.Errorf("failed to create type-specific files: %w", err)
	}

	// Create metadata.json
	publishMeta := &services.PublishMetadata{
		Name:        toolName,
		Version:     version,
		Description: description,
		Author:      author,
		Tags:        tags,
		Type:        toolType,
		Changelog: map[string]string{
			version: "Initial release",
		},
	}

	// We need a publisher service to generate metadata
	// Create temporary services
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

	if err := publisherService.GenerateMetadata(toolPath, publishMeta); err != nil {
		return fmt.Errorf("failed to generate metadata: %w", err)
	}

	fmt.Printf("\nSuccessfully created %s: %s\n", toolType, toolName)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("1. Edit %s/README.md with your tool documentation\n", toolPath)
	fmt.Printf("2. Add your tool implementation files\n")
	fmt.Printf("3. Test your tool locally\n")
	fmt.Printf("4. Publish with: cntm publish %s\n", toolName)

	return nil
}

func promptToolType() (models.ToolType, error) {
	fmt.Println("Select tool type:")
	fmt.Println("  1) agent")
	fmt.Println("  2) command")
	fmt.Println("  3) skill")
	fmt.Print("Enter choice (1-3): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	switch input {
	case "1", "agent":
		return models.ToolTypeAgent, nil
	case "2", "command":
		return models.ToolTypeCommand, nil
	case "3", "skill":
		return models.ToolTypeSkill, nil
	default:
		return "", fmt.Errorf("invalid choice: %s", input)
	}
}

func promptString(prompt, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

func createReadme(toolPath string, toolType models.ToolType, toolName, description string) error {
	readmePath := filepath.Join(toolPath, "README.md")

	content := fmt.Sprintf(`# %s

%s

## Type

%s

## Usage

This %s helps with [describe usage here].

## Features

- Feature 1
- Feature 2
- Feature 3

## Configuration

[Describe any configuration options]

## Examples

[Provide usage examples]

## License

[Your license here]
`, toolName, description, toolType, toolType)

	return os.WriteFile(readmePath, []byte(content), 0644)
}

func createTypeSpecificFiles(toolPath string, toolType models.ToolType, toolName, description string) error {
	switch toolType {
	case models.ToolTypeAgent:
		return createAgentFile(toolPath, toolName)
	case models.ToolTypeCommand:
		// Don't create any files for commands - users should create spec.md, apply.md, archive.md manually
		return nil
	case models.ToolTypeSkill:
		return createSkillFile(toolPath, toolName, description)
	default:
		return fmt.Errorf("unknown tool type: %s", toolType)
	}
}

func createAgentFile(toolPath, toolName string) error {
	agentPath := filepath.Join(toolPath, "agent.md")

	content := fmt.Sprintf(`# %s Agent

You are a specialized agent for [describe purpose].

## Capabilities

- Capability 1
- Capability 2
- Capability 3

## Instructions

[Provide detailed instructions for the agent]

## Examples

### Example 1

[Describe example scenario]

### Example 2

[Describe another example scenario]

## Limitations

[Describe any limitations]
`, toolName)

	return os.WriteFile(agentPath, []byte(content), 0644)
}

func createSkillFile(toolPath, toolName, description string) error {
	skillPath := filepath.Join(toolPath, "SKILL.md")

	content := fmt.Sprintf(`---
name: %s
description: %s
---

# %s Skill

A skill for [describe purpose].

## Quick Start

[Provide overview and how to use the skill]

## Implementation Workflow

### Step 1:
### Step 2:

## Knowledge Areas
[Provide Knowledge reference of the skill to file folder ./reference]

- Area 1
- Area 2
- Area 3

## Best Practices

1. Best practice 1
2. Best practice 2
3. Best practice 3

## Patterns

### Pattern 1

[Describe pattern]

`+"```"+`
[Link example here, Code example should be in ./examples folder]
`+"```"+`

### Pattern 2

[Describe pattern]

`+"```"+`
[Link example here, Code example should be in ./examples folder]
`+"```"+`

## Common Pitfalls

- Pitfall 1
- Pitfall 2

## Resources
[Provide Resource reference of the skill to file folder ./reference]

- Resource 1
- Resource 2

`, toolName, description, toolName)

	return os.WriteFile(skillPath, []byte(content), 0644)
}
