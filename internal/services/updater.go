package services

import (
	"fmt"
	"strings"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"golang.org/x/mod/semver"
)

// OutdatedTool represents a tool that has an available update
type OutdatedTool struct {
	Name           string
	CurrentVersion string
	LatestVersion  string
	Type           models.ToolType
}

// UpdateResult represents the result of updating a single tool
type UpdateResult struct {
	ToolName       string
	Success        bool
	Error          error
	OldVersion     string
	NewVersion     string
	Skipped        bool // If already up-to-date
	Message        string
}

// UpdaterService handles tool update operations
type UpdaterService struct {
	registryService RegistryServiceInterface
	lockFileService LockFileServiceInterface
	installerService *InstallerService
}

// NewUpdaterService creates a new UpdaterService
func NewUpdaterService(
	registryService RegistryServiceInterface,
	lockFileService LockFileServiceInterface,
	installerService *InstallerService,
) (*UpdaterService, error) {
	if registryService == nil {
		return nil, fmt.Errorf("registry service cannot be nil")
	}
	if lockFileService == nil {
		return nil, fmt.Errorf("lock file service cannot be nil")
	}
	if installerService == nil {
		return nil, fmt.Errorf("installer service cannot be nil")
	}

	return &UpdaterService{
		registryService:  registryService,
		lockFileService:  lockFileService,
		installerService: installerService,
	}, nil
}

// CheckOutdated checks for tools that have available updates
func (us *UpdaterService) CheckOutdated() ([]OutdatedTool, error) {
	// Get all installed tools
	installedTools, err := us.lockFileService.ListTools()
	if err != nil {
		return nil, fmt.Errorf("failed to list installed tools: %w", err)
	}

	if len(installedTools) == 0 {
		return []OutdatedTool{}, nil
	}

	// Get latest registry
	registry, err := us.registryService.GetRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry: %w", err)
	}

	var outdated []OutdatedTool

	// Check each installed tool
	for name, installedTool := range installedTools {
		// Find the tool in the registry
		latestTool, err := registry.GetTool(name, installedTool.Type)
		if err != nil {
			// Tool no longer exists in registry, skip it
			continue
		}

		// Compare versions
		cmp := us.CompareVersions(installedTool.Version, latestTool.Version)
		if cmp < 0 {
			// Current version is older than latest
			outdated = append(outdated, OutdatedTool{
				Name:           name,
				CurrentVersion: installedTool.Version,
				LatestVersion:  latestTool.Version,
				Type:           installedTool.Type,
			})
		}
	}

	return outdated, nil
}

// Update updates a specific tool to the latest version
func (us *UpdaterService) Update(toolName string) (*UpdateResult, error) {
	if toolName == "" {
		return nil, fmt.Errorf("tool name cannot be empty")
	}

	result := &UpdateResult{
		ToolName: toolName,
	}

	// Step 1: Check if tool is installed
	installedTool, err := us.lockFileService.GetTool(toolName)
	if err != nil {
		result.Error = fmt.Errorf("tool not installed: %w", err)
		result.Success = false
		return result, result.Error
	}
	result.OldVersion = installedTool.Version

	// Step 2: Get latest version from registry
	latestTool, err := us.registryService.GetTool(toolName, installedTool.Type)
	if err != nil {
		result.Error = fmt.Errorf("tool not found in registry: %w", err)
		result.Success = false
		return result, result.Error
	}
	result.NewVersion = latestTool.Version

	// Step 3: Compare versions
	cmp := us.CompareVersions(installedTool.Version, latestTool.Version)
	if cmp >= 0 {
		// Already up-to-date or newer
		result.Skipped = true
		result.Success = true
		result.Message = fmt.Sprintf("already up-to-date (version %s)", installedTool.Version)
		return result, nil
	}

	// Step 4: Use InstallerService to install the new version
	// The installer will handle backing up, extracting, and updating the lock file
	if err := us.installerService.InstallWithVersion(toolName, latestTool.Version); err != nil {
		result.Error = fmt.Errorf("update failed: %w", err)
		result.Success = false
		return result, result.Error
	}

	result.Success = true
	result.Message = fmt.Sprintf("updated from %s to %s", result.OldVersion, result.NewVersion)
	return result, nil
}

// UpdateAll updates all outdated tools
func (us *UpdaterService) UpdateAll() ([]UpdateResult, []error) {
	// Get all outdated tools
	outdated, err := us.CheckOutdated()
	if err != nil {
		return nil, []error{fmt.Errorf("failed to check for updates: %w", err)}
	}

	if len(outdated) == 0 {
		return []UpdateResult{}, nil
	}

	var results []UpdateResult
	var errors []error

	// Update each outdated tool
	for _, tool := range outdated {
		result, err := us.Update(tool.Name)
		if result != nil {
			results = append(results, *result)
		}
		if err != nil && !result.Skipped {
			errors = append(errors, err)
		}
	}

	return results, errors
}

// CompareVersions compares two semantic version strings
// Returns: -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func (us *UpdaterService) CompareVersions(v1, v2 string) int {
	// Add "v" prefix if not present for semver compatibility
	if v1 != "" && !strings.HasPrefix(v1, "v") {
		v1 = "v" + v1
	}
	if v2 != "" && !strings.HasPrefix(v2, "v") {
		v2 = "v" + v2
	}

	// Use semver.Compare which returns -1, 0, or 1
	return semver.Compare(v1, v2)
}

// GetOutdatedCount returns the number of tools with available updates
func (us *UpdaterService) GetOutdatedCount() (int, error) {
	outdated, err := us.CheckOutdated()
	if err != nil {
		return 0, err
	}
	return len(outdated), nil
}

// IsOutdated checks if a specific tool is outdated
func (us *UpdaterService) IsOutdated(toolName string) (bool, error) {
	if toolName == "" {
		return false, fmt.Errorf("tool name cannot be empty")
	}

	// Get installed tool
	installedTool, err := us.lockFileService.GetTool(toolName)
	if err != nil {
		return false, fmt.Errorf("tool not installed: %w", err)
	}

	// Get latest version from registry
	latestTool, err := us.registryService.GetTool(toolName, installedTool.Type)
	if err != nil {
		return false, fmt.Errorf("tool not found in registry: %w", err)
	}

	// Compare versions
	cmp := us.CompareVersions(installedTool.Version, latestTool.Version)
	return cmp < 0, nil
}

// GetInstalledVersion returns the currently installed version of a tool
func (us *UpdaterService) GetInstalledVersion(toolName string) (string, error) {
	if toolName == "" {
		return "", fmt.Errorf("tool name cannot be empty")
	}

	installedTool, err := us.lockFileService.GetTool(toolName)
	if err != nil {
		return "", fmt.Errorf("tool not installed: %w", err)
	}

	return installedTool.Version, nil
}

// GetLatestVersion returns the latest available version of a tool from the registry
func (us *UpdaterService) GetLatestVersion(toolName string) (string, error) {
	if toolName == "" {
		return "", fmt.Errorf("tool name cannot be empty")
	}

	// First, get the installed tool to determine its type
	installedTool, err := us.lockFileService.GetTool(toolName)
	if err != nil {
		return "", fmt.Errorf("tool not installed: %w", err)
	}

	// Get the tool from registry
	latestTool, err := us.registryService.GetTool(toolName, installedTool.Type)
	if err != nil {
		return "", fmt.Errorf("tool not found in registry: %w", err)
	}

	return latestTool.Version, nil
}
