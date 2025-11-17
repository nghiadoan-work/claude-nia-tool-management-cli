package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	// Create flags
	createType string
	createName string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Claude Code tool (agent, command, or skill)",
	Long: `Create a new tool in your local .claude directory.

This command will interactively guide you through creating a new:
  - Agent: Specialized sub-agent for complex tasks
  - Command: Custom slash command for workflows
  - Skill: Knowledge artifact with domain expertise

The command creates the appropriate directory structure and template files
based on best practices for each tool type.`,
	Example: `  cntm create                        # Interactive mode
  cntm create --type agent --name code-reviewer
  cntm create --type command --name test-runner
  cntm create --type skill --name golang-patterns`,
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Create flags
	createCmd.Flags().StringVarP(&createType, "type", "t", "", "type of tool to create (agent, command, skill)")
	createCmd.Flags().StringVarP(&createName, "name", "n", "", "name of the tool")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Determine the .claude directory path
	// If basePath is already .claude or ends with .claude, use it directly
	// Otherwise, join basePath with .claude
	claudeDir := basePath
	if filepath.Base(basePath) != ".claude" {
		claudeDir = filepath.Join(basePath, ".claude")
	}

	// Check if .claude directory exists
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		return fmt.Errorf(".claude directory not found at %s. Run 'cntm init' first to initialize the project", claudeDir)
	}

	// Interactive mode welcome message
	if createType == "" && createName == "" {
		fmt.Println()
		fmt.Println(ui.Highlight("Create a new Claude Code tool"))
		fmt.Println(ui.Faint("Use arrow keys to navigate, Enter to select, Esc to cancel"))
		fmt.Println()
	}

	// Get tool type if not provided
	if createType == "" {
		toolType, err := promptToolType()
		if err != nil {
			fmt.Println()
			fmt.Println(ui.Warning("✗ Cancelled"))
			return nil
		}
		createType = toolType
	} else {
		// Validate provided type
		if !isValidToolType(createType) {
			return fmt.Errorf("invalid tool type: %s (must be: agent, command, or skill)", createType)
		}
	}

	// Get tool name if not provided
	if createName == "" {
		name, err := promptToolName(createType)
		if err != nil {
			fmt.Println()
			fmt.Println(ui.Warning("✗ Cancelled"))
			return nil
		}
		createName = name
	} else {
		// Validate provided name
		if err := validateToolName(createName); err != nil {
			return err
		}
	}

	// Convert name to kebab-case
	createName = toKebabCase(createName)

	// Show what will be created
	fmt.Println()
	fmt.Println(ui.Faint("──────────────────────────────────────"))
	fmt.Printf("  %s %s\n", ui.Faint("Type:"), ui.Highlight(createType))
	fmt.Printf("  %s %s\n", ui.Faint("Name:"), ui.Highlight(createName))
	fmt.Printf("  %s %s\n", ui.Faint("Path:"), ui.Faint(getToolPath(createType, createName)))
	fmt.Println(ui.Faint("──────────────────────────────────────"))
	fmt.Println()

	// Create the tool
	if err := createTool(createType, createName, claudeDir); err != nil {
		return err
	}

	// Success message
	fmt.Println()
	fmt.Println(ui.Success(fmt.Sprintf("✓ Successfully created %s: %s", createType, createName)))
	fmt.Println()
	fmt.Println("Next steps:")

	switch createType {
	case "agent":
		fmt.Printf("  1. Edit .claude/agents/%s/%s.md to define your agent\n", createName, createName)
		fmt.Println("  2. Refer to .claude/AGENT_TEMPLATE_GUIDE.md for guidance")
		fmt.Printf("  3. Use the agent: Claude will invoke it when needed\n")
	case "command":
		fmt.Printf("  1. Edit .claude/commands/%s/*.md to define your command workflow\n", createName)
		fmt.Println("  2. Refer to .claude/COMMAND_TEMPLATE_GUIDE.md for guidance")
		fmt.Printf("  3. Use the command: /%s\n", createName)
	case "skill":
		fmt.Printf("  1. Edit .claude/skills/%s/SKILL.md to define your skill\n", createName)
		fmt.Println("  2. Add examples and reference materials as needed")
		fmt.Println("  3. Refer to .claude/SKILL_TEMPLATE_GUIDE.md for guidance")
		fmt.Printf("  4. Use the skill: Claude will apply it when relevant\n")
	}

	return nil
}

// ToolTypeOption represents a tool type with description
type ToolTypeOption struct {
	Type        string
	Description string
}

// promptToolType prompts the user to select a tool type
func promptToolType() (string, error) {
	options := []ToolTypeOption{
		{
			Type:        "agent",
			Description: "Specialized sub-agent for complex, multi-step tasks",
		},
		{
			Type:        "command",
			Description: "Custom slash command for workflows and automation",
		},
		{
			Type:        "skill",
			Description: "Knowledge artifact with domain expertise and patterns",
		},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Type | green | bold }}: {{ .Description | faint }}",
		Inactive: "  {{ .Type | cyan }}: {{ .Description | faint }}",
		Selected: "{{ \"✓\" | green | bold }} {{ .Type | green }}",
	}

	prompt := promptui.Select{
		Label:     ui.Info("What type of tool do you want to create?"),
		Items:     options,
		Size:      3,
		Templates: templates,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return options[index].Type, nil
}

// promptToolName prompts the user to enter a tool name
func promptToolName(toolType string) (string, error) {
	// Print helper message
	fmt.Println()
	fmt.Printf("  %s You can use spaces, hyphens, or underscores - they'll be converted to kebab-case\n", ui.Faint("ℹ"))
	fmt.Printf("  %s Example: \"Code Reviewer\" → \"code-reviewer\"\n", ui.Faint("ℹ"))
	fmt.Println()

	validate := func(input string) error {
		// Only validate, don't print anything
		return validateToolName(input)
	}

	prompt := promptui.Prompt{
		Label:    ui.Info(fmt.Sprintf("Enter %s name", toolType)),
		Validate: validate,
		Templates: &promptui.PromptTemplates{
			Prompt:  "{{ . }} ",
			Valid:   "{{ . }} ",
			Invalid: "{{ . }} ",
			Success: "{{ \"✓\" | green }} {{ . | faint }} ",
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// Show what the kebab-case version will be after input is complete
	kebab := toKebabCase(result)
	if result != kebab {
		fmt.Printf("  %s Will be created as: %s\n\n", ui.Faint("→"), ui.Highlight(kebab))
	}

	return result, nil
}

// validateToolName validates the tool name
func validateToolName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	// Check for valid characters (letters, numbers, hyphens, spaces, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z0-9-_ ]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("name must contain only letters, numbers, hyphens, spaces, or underscores")
	}

	// After conversion to kebab-case, check it doesn't start or end with hyphen
	kebab := toKebabCase(name)
	if strings.HasPrefix(kebab, "-") || strings.HasSuffix(kebab, "-") {
		return fmt.Errorf("name cannot start or end with a hyphen")
	}

	return nil
}

// isValidToolType checks if the tool type is valid
func isValidToolType(toolType string) bool {
	return toolType == "agent" || toolType == "command" || toolType == "skill"
}

// toKebabCase converts a string to kebab-case
func toKebabCase(s string) string {
	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Convert to lowercase
	s = strings.ToLower(s)

	// Remove any consecutive hyphens
	re := regexp.MustCompile(`-+`)
	s = re.ReplaceAllString(s, "-")

	return s
}

// createTool creates the tool directory and template files
func createTool(toolType, name, claudeDir string) error {
	switch toolType {
	case "agent":
		return createAgent(name, claudeDir)
	case "command":
		return createCommand(name, claudeDir)
	case "skill":
		return createSkill(name, claudeDir)
	default:
		return fmt.Errorf("unknown tool type: %s", toolType)
	}
}

// createAgent creates a new agent
func createAgent(name, claudeDir string) error {
	agentDir := filepath.Join(claudeDir, "agents", name)

	// Check if agent already exists
	if _, err := os.Stat(agentDir); err == nil {
		return fmt.Errorf("agent '%s' already exists", name)
	}

	// Create agent directory
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return fmt.Errorf("failed to create agent directory: %w", err)
	}

	// Create agent file
	agentFile := filepath.Join(agentDir, name+".md")
	agentTemplate := fmt.Sprintf(`---
name: %s
description: Brief description of what this agent does and when to use it
tools: Read, Write, Edit, Bash, Grep, Glob
model: inherit
---

# %s

## Purpose
Describe what this agent does in 1-2 sentences.

## Instructions
When invoked, you should:

1. **First Action**
   - Detail about the action
   - Additional context or requirements

2. **Second Action**
   - Detail about the action
   - Additional context or requirements

3. **Final Action**
   - Detail about the action
   - What to return or output

## Guidelines
- Behavioral rule or priority
- Constraint or limitation
- Best practice to follow

## Output Format
Describe how the agent should structure its output or response.

## Scope
This agent WILL:
- Capability 1
- Capability 2

This agent WILL NOT:
- Limitation 1
- Limitation 2

## Error Handling
- **Error Type**: How to handle this error
- **Edge Case**: How to handle this case
`, name, toTitleCase(name))

	if err := os.WriteFile(agentFile, []byte(agentTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write agent file: %w", err)
	}

	fmt.Printf("  Created .claude/agents/%s/%s.md\n", name, name)
	return nil
}

// createCommand creates a new command
func createCommand(name, claudeDir string) error {
	commandDir := filepath.Join(claudeDir, "commands", name)

	// Check if command already exists
	if _, err := os.Stat(commandDir); err == nil {
		return fmt.Errorf("command '%s' already exists", name)
	}

	// Create command directory
	if err := os.MkdirAll(commandDir, 0755); err != nil {
		return fmt.Errorf("failed to create command directory: %w", err)
	}

	// Create command file
	commandFile := filepath.Join(commandDir, name+".md")
	commandTemplate := fmt.Sprintf(`---
name: %s
description: Brief description of what this command does
---

# %s

## Usage
Describe when and how to use this command.

## Command Behavior
When invoked, this command will:

1. **Action 1**
   - Detail about what happens
   - Expected input or context

2. **Action 2**
   - Detail about what happens
   - How it processes information

3. **Action 3**
   - Detail about what happens
   - What output is produced

## Examples
Provide examples of using this command:

**Example 1: Basic usage**
` + "```" + `
/%s
` + "```" + `

**Example 2: Advanced usage**
` + "```" + `
/%s --option value
` + "```" + `

## Notes
- Important considerations
- Edge cases to be aware of
- Dependencies or requirements
`, name, toTitleCase(name), name, name)

	if err := os.WriteFile(commandFile, []byte(commandTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write command file: %w", err)
	}

	fmt.Printf("  Created .claude/commands/%s/%s.md\n", name, name)
	return nil
}

// createSkill creates a new skill
func createSkill(name, claudeDir string) error {
	skillDir := filepath.Join(claudeDir, "skills", name)

	// Check if skill already exists
	if _, err := os.Stat(skillDir); err == nil {
		return fmt.Errorf("skill '%s' already exists", name)
	}

	// Create skill directory
	if err := os.MkdirAll(skillDir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	// Create examples subdirectory
	examplesDir := filepath.Join(skillDir, "examples")
	if err := os.MkdirAll(examplesDir, 0755); err != nil {
		return fmt.Errorf("failed to create examples directory: %w", err)
	}

	// Create skill file
	skillFile := filepath.Join(skillDir, "SKILL.md")
	skillTemplate := fmt.Sprintf(`---
name: %s
description: Brief description of what knowledge or expertise this skill provides
---

# %s

## Quick Start
Provide a brief overview and quick usage guide.

## Overview
Detailed description of the skill's domain and what it covers:
- Key concept 1
- Key concept 2
- Key concept 3

## Core Concepts

### Concept 1
Explanation of the first key concept.

### Concept 2
Explanation of the second key concept.

## Implementation Patterns

### Pattern 1: Pattern Name
**When to use**: Describe the use case

**Example**:
` + "```" + `
// Code example here
` + "```" + `

**Explanation**: Why this pattern works and when to use it.

### Pattern 2: Pattern Name
**When to use**: Describe the use case

**Example**:
` + "```" + `
// Code example here
` + "```" + `

**Explanation**: Why this pattern works and when to use it.

## Best Practices
- Best practice 1
- Best practice 2
- Best practice 3

## Common Pitfalls
- **Pitfall 1**: What to avoid and why
- **Pitfall 2**: What to avoid and why

## Troubleshooting
**Problem**: Common issue description
**Solution**: How to resolve it

**Problem**: Another common issue
**Solution**: How to resolve it

## Additional Resources
- Resource 1
- Resource 2
`, name, toTitleCase(name))

	if err := os.WriteFile(skillFile, []byte(skillTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write skill file: %w", err)
	}

	// Create examples README
	examplesReadme := filepath.Join(examplesDir, "README.md")
	examplesTemplate := fmt.Sprintf(`# %s Examples

This directory contains code examples and usage patterns for the %s skill.

## Examples

### Example 1: [Description]
File: ` + "`example-1.ext`" + `

Description of what this example demonstrates.

### Example 2: [Description]
File: ` + "`example-2.ext`" + `

Description of what this example demonstrates.

## How to Use These Examples
Instructions on how to apply these examples in real projects.
`, toTitleCase(name), name)

	if err := os.WriteFile(examplesReadme, []byte(examplesTemplate), 0644); err != nil {
		return fmt.Errorf("failed to write examples README: %w", err)
	}

	fmt.Printf("  Created .claude/skills/%s/SKILL.md\n", name)
	fmt.Printf("  Created .claude/skills/%s/examples/\n", name)
	return nil
}

// toTitleCase converts kebab-case to Title Case
func toTitleCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// getToolPath returns the relative path where the tool will be created
func getToolPath(toolType, name string) string {
	switch toolType {
	case "agent":
		return fmt.Sprintf(".claude/agents/%s/%s.md", name, name)
	case "command":
		return fmt.Sprintf(".claude/commands/%s/%s.md", name, name)
	case "skill":
		return fmt.Sprintf(".claude/skills/%s/SKILL.md", name)
	default:
		return ""
	}
}
