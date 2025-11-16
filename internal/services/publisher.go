package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/internal/data"
	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
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
	Name         string
	Version      string
	Description  string
	Author       string
	Tags         []string
	Type         models.ToolType
	Changelog    map[string]string
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

	// Check for README.md (optional, but recommended)
	readmePath := filepath.Join(toolPath, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		fmt.Printf("Warning: README.md not found (recommended for documentation)\n")
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
		// Skills should have SKILL.md or similar
		skillFile := filepath.Join(toolPath, "SKILL.md")
		if _, err := os.Stat(skillFile); os.IsNotExist(err) {
			fmt.Printf("Warning: SKILL.md not found (optional)\n")
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

	// Generate default author if empty
	if meta.Author == "" {
		meta.Author = "Anonymous"
		fmt.Printf("Info: Generated default author: %s\n", meta.Author)
	}

	// Generate default description if empty
	if meta.Description == "" {
		meta.Description = fmt.Sprintf("A %s tool for Claude Code", meta.Type)
		fmt.Printf("Info: Generated default description: %s\n", meta.Description)
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
	// Convert version to filename format (1.0.0 -> v1-0-0)
	versionFileName := versionToFileName(version)

	// Get ZIP file size
	zipInfo, err := os.Stat(zipPath)
	if err != nil {
		return fmt.Errorf("failed to stat ZIP file: %w", err)
	}

	// Create VersionInfo for this specific version
	versionInfo := &models.VersionInfo{
		File:      fmt.Sprintf("tools/%ss/%s/%s.zip", toolType, toolName, versionFileName),
		Size:      zipInfo.Size(),
		CreatedAt: time.Now(),
	}

	// Load metadata if exists
	metadataPath := filepath.Join(toolPath, "metadata.json")
	var toolAuthor, toolDescription string
	var toolTags []string
	if data, err := os.ReadFile(metadataPath); err == nil {
		var metadata models.ToolMetadata
		if err := json.Unmarshal(data, &metadata); err == nil {
			toolAuthor = metadata.Author
			toolDescription = metadata.Description
			toolTags = metadata.Tags
			// Add changelog for this version if available
			if changelog, ok := metadata.Changelog[version]; ok {
				versionInfo.Changelog = changelog
			}
		}
	}

	// Create ToolInfo structure (will be used for registry update)
	toolInfo := &models.ToolInfo{
		Name:          toolName,
		LatestVersion: version,
		Type:          toolType,
		Author:        toolAuthor,
		Description:   toolDescription,
		Tags:          toolTags,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Versions: map[string]*models.VersionInfo{
			version: versionInfo,
		},
	}

	// Print package info
	fmt.Printf("\nTool packaged successfully!\n")
	fmt.Printf("  Tool:    %s\n", toolName)
	fmt.Printf("  Type:    %s\n", toolType)
	fmt.Printf("  Version: %s\n", version)
	fmt.Printf("  Size:    %d bytes\n", versionInfo.Size)
	fmt.Printf("  Hash:    %s\n", hash)
	fmt.Printf("  Package: %s\n", zipPath)

	// Step 5: Create pull request if configured
	if ps.config.Publish.CreatePR {
		fmt.Printf("\nCreating pull request to registry...\n")

		// Read ZIP file for upload
		zipData, err := os.ReadFile(zipPath)
		if err != nil {
			return fmt.Errorf("failed to read ZIP file: %w", err)
		}

		if err := ps.CreatePullRequest(toolInfo, zipData, hash); err != nil {
			return fmt.Errorf("failed to create pull request: %w", err)
		}

		fmt.Printf("\nPublication complete!\n")
	} else {
		fmt.Printf("\nTo complete publishing:\n")
		fmt.Printf("1. Upload %s to registry repository\n", zipPath)
		fmt.Printf("2. Update registry.json with the tool information\n")
		fmt.Printf("3. Create a pull request to the registry\n")
		fmt.Printf("\nTip: Set 'create_pr: true' in config to automate this process\n")
	}

	return nil
}

// CreatePullRequest creates a PR to the registry repository
func (ps *PublisherService) CreatePullRequest(tool *models.ToolInfo, zipData []byte, hash string) error {
	// Check if we have a GitHub token (should be auto-detected by GitHubClient)
	if ps.githubClient.authToken == "" {
		return fmt.Errorf(`GitHub authentication required for automated PR creation

Please authenticate using one of these methods:
1. Install and login to GitHub CLI:
   brew install gh
   gh auth login

2. Set environment variable:
   export GITHUB_TOKEN=your_token_here

3. Add to config file (~/.claude-tools-config.yaml):
   registry:
     auth_token: your_token_here

Get a token from: https://github.com/settings/tokens (needs 'repo' scope)`)
	}

	// Parse registry URL to get owner and repo
	owner, repo, err := ParseRepoURL(ps.config.Registry.URL)
	if err != nil {
		return fmt.Errorf("failed to parse registry URL: %w", err)
	}

	fmt.Printf("  Registry: %s/%s\n", owner, repo)

	// Step 1: Get authenticated user
	username, err := ps.githubClient.GetAuthenticatedUser()
	if err != nil {
		return fmt.Errorf("failed to get authenticated user: %w", err)
	}
	fmt.Printf("  User: %s\n", username)

	// Step 2: Fork repository if needed
	fmt.Printf("  Checking fork...\n")
	defaultBranch, err := ps.githubClient.GetDefaultBranch(username, repo)
	if err != nil {
		// Fork doesn't exist, create it
		fmt.Printf("  Creating fork...\n")
		fork, err := ps.githubClient.ForkRepository(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to fork repository: %w", err)
		}
		defaultBranch = fork.GetDefaultBranch()
		fmt.Printf("  Fork created\n")
	} else {
		fmt.Printf("  Fork exists\n")
	}

	// Step 3: Create a new branch
	branchName := fmt.Sprintf("publish-%s-%s", tool.Name, tool.LatestVersion)
	fmt.Printf("  Creating branch: %s\n", branchName)

	err = ps.githubClient.CreateBranch(username, repo, branchName, defaultBranch)
	if err != nil {
		// Branch might already exist, that's okay
		fmt.Printf("  Branch already exists or created\n")
	}

	// Step 4: Upload ZIP file with versioned path
	versionFileName := versionToFileName(tool.LatestVersion)
	zipFilePath := fmt.Sprintf("tools/%ss/%s/%s.zip", tool.Type, tool.Name, versionFileName)
	fmt.Printf("  Uploading: %s\n", zipFilePath)

	err = ps.githubClient.UploadFile(
		username,
		repo,
		zipFilePath,
		branchName,
		zipData,
		fmt.Sprintf("Add %s v%s", tool.Name, tool.LatestVersion),
	)
	if err != nil {
		return fmt.Errorf("failed to upload ZIP file: %w", err)
	}

	// Step 5: Update or create registry.json
	fmt.Printf("  Updating registry.json\n")

	// Try to fetch current registry.json
	var registry models.Registry
	registryData, err := ps.githubClient.FetchFile("registry.json")
	if err != nil {
		// registry.json doesn't exist, create a new one
		fmt.Printf("  Creating new registry.json\n")
		registry = models.Registry{
			Version:   "1.0.0",
			UpdatedAt: time.Now(),
			Tools:     make(map[models.ToolType][]*models.ToolInfo),
		}
	} else {
		// registry.json exists, parse it
		if err := json.Unmarshal(registryData, &registry); err != nil {
			return fmt.Errorf("failed to parse registry.json: %w", err)
		}
	}

	// Initialize tools map if needed
	if registry.Tools == nil {
		registry.Tools = make(map[models.ToolType][]*models.ToolInfo)
	}

	// Get tools of this type
	toolsOfType := registry.Tools[tool.Type]

	// Check if tool already exists
	found := false
	for i, existingTool := range toolsOfType {
		if existingTool.Name == tool.Name {
			// Tool exists - append the new version to its versions map
			if existingTool.Versions == nil {
				existingTool.Versions = make(map[string]*models.VersionInfo)
			}

			// Add the new version
			existingTool.Versions[tool.LatestVersion] = tool.Versions[tool.LatestVersion]

			// Update latest_version to the new version
			existingTool.LatestVersion = tool.LatestVersion

			// Update metadata (in case description, author, tags changed)
			existingTool.Description = tool.Description
			existingTool.Author = tool.Author
			existingTool.Tags = tool.Tags
			existingTool.UpdatedAt = time.Now()

			toolsOfType[i] = existingTool
			found = true
			break
		}
	}

	// If not found, add it as a new tool
	if !found {
		tool.CreatedAt = time.Now()
		toolsOfType = append(toolsOfType, tool)
	}

	// Update the map
	registry.Tools[tool.Type] = toolsOfType
	registry.UpdatedAt = time.Now()

	// Marshal updated registry
	updatedRegistryData, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry.json: %w", err)
	}

	// Upload updated registry.json
	err = ps.githubClient.UploadFile(
		username,
		repo,
		"registry.json",
		branchName,
		updatedRegistryData,
		fmt.Sprintf("Update registry for %s v%s", tool.Name, tool.LatestVersion),
	)
	if err != nil {
		return fmt.Errorf("failed to update registry.json: %w", err)
	}

	// Step 6: Create pull request
	fmt.Printf("  Creating pull request\n")

	prTitle := fmt.Sprintf("Publish %s v%s", tool.Name, tool.LatestVersion)
	prBody := fmt.Sprintf(`## Tool Publication

**Name:** %s
**Version:** %s
**Type:** %s
**Author:** %s

**Description:** %s

**File:** %s
**Size:** %d bytes
**Hash:** %s

---
*This PR was automatically generated by cntm*
`, tool.Name, tool.LatestVersion, tool.Type, tool.Author, tool.Description, zipFilePath, tool.Versions[tool.LatestVersion].Size, hash)

	headBranch := fmt.Sprintf("%s:%s", username, branchName)
	pr, err := ps.githubClient.CreatePullRequest(owner, repo, prTitle, prBody, headBranch, defaultBranch)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	fmt.Printf("\n✓ Pull request created: %s\n", pr.GetHTMLURL())

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

// versionToFileName converts a semantic version to a filename-safe format
// Examples:
//   1.0.0 -> v1-0-0
//   2.1.3 -> v2-1-3
//   1.0.0-beta -> v1-0-0-beta
func versionToFileName(version string) string {
	// Replace dots with dashes
	fileName := strings.ReplaceAll(version, ".", "-")
	// Add 'v' prefix if not present
	if !strings.HasPrefix(fileName, "v") {
		fileName = "v" + fileName
	}
	return fileName
}
