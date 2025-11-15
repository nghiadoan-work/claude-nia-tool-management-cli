package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
)

// PublisherService handles tool publishing operations
type PublisherService struct {
	fsManager       *data.FSManager
	githubClient    *GitHubClient
	registryService *RegistryService
	config          *models.Config
}

// PublishMetadata represents metadata for publishing a tool
type PublishMetadata struct {
	Name        string
	Version     string
	Description string
	Author      string
	Tags        []string
	Type        models.ToolType
	Changelog   map[string]string
	Dependencies []string
}

// NewPublisherService creates a new PublisherService
func NewPublisherService(
	fsManager *data.FSManager,
	githubClient *GitHubClient,
	registryService *RegistryService,
	config *models.Config,
) (*PublisherService, error) {
	if fsManager == nil {
		return nil, fmt.Errorf("fs manager cannot be nil")
	}
	if githubClient == nil {
		return nil, fmt.Errorf("github client cannot be nil")
	}
	if registryService == nil {
		return nil, fmt.Errorf("registry service cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	return &PublisherService{
		fsManager:       fsManager,
		githubClient:    githubClient,
		registryService: registryService,
		config:          config,
	}, nil
}

// ValidateTool validates a tool directory before publishing
func (ps *PublisherService) ValidateTool(toolPath string) error {
	if toolPath == "" {
		return fmt.Errorf("tool path cannot be empty")
	}

	// Check directory exists
	info, err := os.Stat(toolPath)
	if err != nil {
		return fmt.Errorf("tool directory does not exist: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("tool path is not a directory: %s", toolPath)
	}

	// Check for README.md
	readmePath := filepath.Join(toolPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		return fmt.Errorf("README.md not found in tool directory")
	}

	// Determine tool type from path
	toolType, err := ps.detectToolType(toolPath)
	if err != nil {
		return fmt.Errorf("failed to detect tool type: %w", err)
	}

	// Validate tool type-specific files
	if err := ps.validateToolTypeFiles(toolPath, toolType); err != nil {
		return fmt.Errorf("tool type validation failed: %w", err)
	}

	// Check for sensitive files that should not be published
	sensitiveFiles := []string{".git", ".env", ".DS_Store", "node_modules", "credentials.json"}
	for _, sensitiveFile := range sensitiveFiles {
		sensitivePath := filepath.Join(toolPath, sensitiveFile)
		if _, err := os.Stat(sensitivePath); err == nil {
			return fmt.Errorf("sensitive file/directory found: %s (should be excluded)", sensitiveFile)
		}
	}

	return nil
}

// detectToolType detects the tool type from the directory path
func (ps *PublisherService) detectToolType(toolPath string) (models.ToolType, error) {
	absPath, err := filepath.Abs(toolPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if path contains type indicators
	pathLower := strings.ToLower(absPath)

	if strings.Contains(pathLower, "/agents/") || strings.Contains(pathLower, "\\agents\\") {
		return models.ToolTypeAgent, nil
	}
	if strings.Contains(pathLower, "/commands/") || strings.Contains(pathLower, "\\commands\\") {
		return models.ToolTypeCommand, nil
	}
	if strings.Contains(pathLower, "/skills/") || strings.Contains(pathLower, "\\skills\\") {
		return models.ToolTypeSkill, nil
	}

	// Check for metadata.json to determine type
	metadataPath := filepath.Join(toolPath, "metadata.json")
	if _, err := os.Stat(metadataPath); err == nil {
		data, err := os.ReadFile(metadataPath)
		if err == nil {
			var metadata models.ToolMetadata
			if err := json.Unmarshal(data, &metadata); err == nil {
				// Check if there's a type field in custom metadata
				if typeStr, ok := metadata.Custom["type"]; ok {
					return models.ToolType(typeStr), nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not detect tool type from path or metadata")
}

// validateToolTypeFiles validates type-specific files
func (ps *PublisherService) validateToolTypeFiles(toolPath string, toolType models.ToolType) error {
	switch toolType {
	case models.ToolTypeAgent:
		// Agents should have agent.md or similar
		agentFile := filepath.Join(toolPath, "agent.md")
		if _, err := os.Stat(agentFile); os.IsNotExist(err) {
			// Agent file is optional, just warn
			fmt.Printf("Warning: agent.md not found (optional)\n")
		}
	case models.ToolTypeCommand:
		// Commands should have command.md or similar
		commandFile := filepath.Join(toolPath, "command.md")
		if _, err := os.Stat(commandFile); os.IsNotExist(err) {
			fmt.Printf("Warning: command.md not found (optional)\n")
		}
	case models.ToolTypeSkill:
		// Skills should have skill.md or similar
		skillFile := filepath.Join(toolPath, "skill.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			fmt.Printf("Warning: skill.md not found (optional)\n")
		}
	}

	return nil
}

// GenerateMetadata generates or updates metadata.json for a tool
func (ps *PublisherService) GenerateMetadata(toolPath string, meta *PublishMetadata) error {
	if toolPath == "" {
		return fmt.Errorf("tool path cannot be empty")
	}
	if meta == nil {
		return fmt.Errorf("publish metadata cannot be nil")
	}

	// Validate metadata
	if meta.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if meta.Version == "" {
		return fmt.Errorf("tool version cannot be empty")
	}
	if meta.Author == "" {
		return fmt.Errorf("tool author cannot be empty")
	}
	if meta.Description == "" {
		return fmt.Errorf("tool description cannot be empty")
	}

	// Create ToolMetadata
	toolMetadata := &models.ToolMetadata{
		Author:       meta.Author,
		Tags:         meta.Tags,
		Description:  meta.Description,
		Version:      meta.Version,
		Dependencies: meta.Dependencies,
		Changelog:    meta.Changelog,
		Custom: map[string]string{
			"type": string(meta.Type),
		},
	}

	// Convert to JSON
	data, err := json.MarshalIndent(toolMetadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Write to metadata.json
	metadataPath := filepath.Join(toolPath, "metadata.json")
	if err := os.WriteFile(metadataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write metadata.json: %w", err)
	}

	return nil
}

// CreatePackage creates a ZIP package from a tool directory
func (ps *PublisherService) CreatePackage(toolPath, outputPath string) (string, error) {
	if toolPath == "" {
		return "", fmt.Errorf("tool path cannot be empty")
	}
	if outputPath == "" {
		return "", fmt.Errorf("output path cannot be empty")
	}

	// Validate tool before packaging
	if err := ps.ValidateTool(toolPath); err != nil {
		return "", fmt.Errorf("tool validation failed: %w", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create ZIP file
	if err := ps.fsManager.CreateZIP(toolPath, outputPath); err != nil {
		return "", fmt.Errorf("failed to create ZIP: %w", err)
	}

	// Calculate SHA256 hash
	hash, err := ps.fsManager.CalculateSHA256(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return hash, nil
}

// PublishToRegistry publishes a tool to the registry
// This creates a PR to the registry repository
func (ps *PublisherService) PublishToRegistry(toolPath, version string) error {
	if toolPath == "" {
		return fmt.Errorf("tool path cannot be empty")
	}
	if version == "" {
		return fmt.Errorf("version cannot be empty")
	}

	// Step 1: Validate tool
	if err := ps.ValidateTool(toolPath); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Step 2: Detect tool type and name
	toolType, err := ps.detectToolType(toolPath)
	if err != nil {
		return fmt.Errorf("failed to detect tool type: %w", err)
	}

	toolName := filepath.Base(toolPath)

	// Step 3: Create package
	tempDir, err := os.MkdirTemp("", "cntm-publish-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	zipPath := filepath.Join(tempDir, fmt.Sprintf("%s.zip", toolName))
	hash, err := ps.CreatePackage(toolPath, zipPath)
	if err != nil {
		return fmt.Errorf("failed to create package: %w", err)
	}

	// Step 4: Create ToolInfo for registry
	toolInfo := &models.ToolInfo{
		Name:        toolName,
		Version:     version,
		Type:        toolType,
		File:        fmt.Sprintf("tools/%ss/%s.zip", toolType, toolName),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Load metadata if exists
	metadataPath := filepath.Join(toolPath, "metadata.json")
	if data, err := os.ReadFile(metadataPath); err == nil {
		var metadata models.ToolMetadata
		if err := json.Unmarshal(data, &metadata); err == nil {
			toolInfo.Author = metadata.Author
			toolInfo.Description = metadata.Description
			toolInfo.Tags = metadata.Tags
		}
	}

	// Get ZIP file size
	zipInfo, err := os.Stat(zipPath)
	if err != nil {
		return fmt.Errorf("failed to stat ZIP file: %w", err)
	}
	toolInfo.Size = zipInfo.Size()

	// Step 5: Create pull request (simplified for now)
	fmt.Printf("\nTool packaged successfully!\n")
	fmt.Printf("  Tool:    %s\n", toolName)
	fmt.Printf("  Type:    %s\n", toolType)
	fmt.Printf("  Version: %s\n", version)
	fmt.Printf("  Size:    %d bytes\n", toolInfo.Size)
	fmt.Printf("  Hash:    %s\n", hash)
	fmt.Printf("  Package: %s\n", zipPath)
	fmt.Printf("\nTo complete publishing:\n")
	fmt.Printf("1. Upload %s to registry repository\n", zipPath)
	fmt.Printf("2. Update registry.json with the tool information\n")
	fmt.Printf("3. Create a pull request to the registry\n")

	return nil
}

// CreatePullRequest creates a PR to the registry repository
// Simplified implementation - in production, this would use GitHub API
func (ps *PublisherService) CreatePullRequest(tool *models.ToolInfo, zipPath string) error {
	// This is a placeholder for the full GitHub PR workflow
	// In a complete implementation, this would:
	// 1. Fork the registry repo (if not already forked)
	// 2. Create a new branch
	// 3. Upload the ZIP file
	// 4. Update registry.json
	// 5. Commit changes
	// 6. Create pull request

	fmt.Printf("Creating pull request for %s@%s...\n", tool.Name, tool.Version)
	fmt.Printf("Note: Automated PR creation is not yet implemented.\n")
	fmt.Printf("Please manually create a PR to the registry repository.\n")

	return nil
}

// ReadExistingMetadata reads metadata.json from a tool directory
func (ps *PublisherService) ReadExistingMetadata(toolPath string) (*models.ToolMetadata, error) {
	metadataPath := filepath.Join(toolPath, "metadata.json")

	data, err := os.ReadFile(metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No existing metadata
		}
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var metadata models.ToolMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	return &metadata, nil
}
