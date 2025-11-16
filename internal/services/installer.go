package services

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/schollz/progressbar/v3"
)

// RegistryServiceInterface defines the methods needed from RegistryService
type RegistryServiceInterface interface {
	GetTool(name string, toolType models.ToolType) (*models.ToolInfo, error)
	GetRegistry() (*models.Registry, error)
}

// GitHubDownloader defines the methods needed for downloading files
type GitHubDownloader interface {
	DownloadFile(url string, size int64, showProgress bool) ([]byte, error)
}

// FSManagerInterface defines the methods needed from FSManager
type FSManagerInterface interface {
	ExtractZIP(zipPath, destPath string) error
	CalculateSHA256(filePath string) (string, error)
	RemoveDir(path string) error
}

// LockFileServiceInterface defines the methods needed from LockFileService
type LockFileServiceInterface interface {
	GetTool(name string) (*models.InstalledTool, error)
	AddTool(name string, tool *models.InstalledTool) error
	RemoveTool(name string) error
	ListTools() (map[string]*models.InstalledTool, error)
	IsInstalled(name string) (bool, error)
	GetRegistry() (string, error)
	SetRegistry(registryURL string) error
}

// InstallerService handles tool installation operations
type InstallerService struct {
	githubClient    GitHubDownloader
	registryService RegistryServiceInterface
	fsManager       FSManagerInterface
	lockFileService LockFileServiceInterface
	config          *models.Config
	baseDir         string // Base directory for installations (.claude)
}

// InstallResult represents the result of a single tool installation
type InstallResult struct {
	ToolName string
	Success  bool
	Error    error
	Skipped  bool // If already installed with same version
	Message  string
}

// NewInstallerService creates a new InstallerService
func NewInstallerService(
	githubClient GitHubDownloader,
	registryService RegistryServiceInterface,
	fsManager FSManagerInterface,
	lockFileService LockFileServiceInterface,
	config *models.Config,
) (*InstallerService, error) {
	if githubClient == nil {
		return nil, fmt.Errorf("github client cannot be nil")
	}
	if registryService == nil {
		return nil, fmt.Errorf("registry service cannot be nil")
	}
	if fsManager == nil {
		return nil, fmt.Errorf("fs manager cannot be nil")
	}
	if lockFileService == nil {
		return nil, fmt.Errorf("lock file service cannot be nil")
	}
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	baseDir := filepath.Join(config.Local.DefaultPath)
	if baseDir == "" {
		baseDir = ".claude"
	}

	// Ensure base directory exists
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute base directory: %w", err)
	}

	if err := os.MkdirAll(absBaseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %w", err)
	}

	return &InstallerService{
		githubClient:    githubClient,
		registryService: registryService,
		fsManager:       fsManager,
		lockFileService: lockFileService,
		config:          config,
		baseDir:         absBaseDir,
	}, nil
}

// Install installs a tool by name, using the latest version from the registry
func (ins *InstallerService) Install(toolName string) error {
	return ins.InstallWithVersion(toolName, "")
}

// InstallWithVersion installs a specific version of a tool
// If version is empty, installs the latest version
func (ins *InstallerService) InstallWithVersion(toolName, version string) error {
	if toolName == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	// Step 1: Search for the tool in the registry (try all types)
	tool, err := ins.findTool(toolName)
	if err != nil {
		return fmt.Errorf("failed to find tool: %w\nHint: Run 'cntm search %s' to verify the tool exists", err, toolName)
	}

	// Step 2: Determine which version to install
	versionToInstall := version
	if versionToInstall == "" {
		versionToInstall = tool.LatestVersion
	}

	// Validate that the requested version exists
	versionInfo, err := tool.GetVersion(versionToInstall)
	if err != nil {
		return fmt.Errorf("version %s not found for tool %s\nAvailable versions: %v",
			versionToInstall, toolName, tool.ListVersions())
	}

	// Step 3: Check if already installed with same version
	installedTool, err := ins.lockFileService.GetTool(toolName)
	if err == nil && installedTool != nil {
		if installedTool.Version == versionToInstall {
			fmt.Printf("Tool %s@%s is already installed, skipping\n", toolName, versionToInstall)
			return nil
		}
		fmt.Printf("Updating %s from %s to %s\n", toolName, installedTool.Version, versionToInstall)
	} else {
		fmt.Printf("Installing %s@%s\n", toolName, versionToInstall)
	}

	// Step 4: Install the tool
	if err := ins.installToolWithVersion(tool, versionToInstall, versionInfo); err != nil {
		return fmt.Errorf("failed to install tool: %w", err)
	}

	fmt.Printf("Successfully installed %s@%s\n", toolName, versionToInstall)
	return nil
}

// InstallMultiple installs multiple tools sequentially
// Returns a slice of results for each tool and a slice of errors
func (ins *InstallerService) InstallMultiple(toolNames []string) ([]InstallResult, []error) {
	if len(toolNames) == 0 {
		return nil, []error{fmt.Errorf("no tools specified")}
	}

	results := make([]InstallResult, 0, len(toolNames))
	var errors []error

	for _, toolName := range toolNames {
		result := InstallResult{
			ToolName: toolName,
		}

		err := ins.Install(toolName)
		if err != nil {
			result.Success = false
			result.Error = err
			result.Message = err.Error()
			errors = append(errors, err)
		} else {
			result.Success = true
			result.Message = "installed successfully"
		}

		results = append(results, result)
	}

	return results, errors
}

// VerifyInstallation verifies that a tool is correctly installed
func (ins *InstallerService) VerifyInstallation(toolName string) error {
	if toolName == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	// Step 1: Check if tool is in lock file
	installedTool, err := ins.lockFileService.GetTool(toolName)
	if err != nil {
		return fmt.Errorf("tool not found in lock file: %w", err)
	}

	// Step 2: Check if directory exists
	destDir := ins.getInstallPath(toolName, installedTool.Type)
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		return fmt.Errorf("installation directory does not exist: %s", destDir)
	}

	// Step 3: Verify directory is not empty
	entries, err := os.ReadDir(destDir)
	if err != nil {
		return fmt.Errorf("failed to read installation directory: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("installation directory is empty")
	}

	return nil
}

// Uninstall removes a tool from the system
func (ins *InstallerService) Uninstall(toolName string) error {
	if toolName == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	// Step 1: Check if tool is installed
	installedTool, err := ins.lockFileService.GetTool(toolName)
	if err != nil {
		return fmt.Errorf("tool not installed: %w", err)
	}

	// Step 2: Remove installation directory
	destDir := ins.getInstallPath(toolName, installedTool.Type)
	if err := ins.fsManager.RemoveDir(destDir); err != nil {
		return fmt.Errorf("failed to remove installation directory: %w", err)
	}

	// Step 3: Remove from lock file
	if err := ins.lockFileService.RemoveTool(toolName); err != nil {
		return fmt.Errorf("failed to update lock file: %w", err)
	}

	fmt.Printf("Successfully uninstalled %s\n", toolName)
	return nil
}

// findTool searches for a tool in the registry by trying all tool types
func (ins *InstallerService) findTool(toolName string) (*models.ToolInfo, error) {
	// Try each tool type
	types := []models.ToolType{
		models.ToolTypeAgent,
		models.ToolTypeCommand,
		models.ToolTypeSkill,
	}

	for _, toolType := range types {
		tool, err := ins.registryService.GetTool(toolName, toolType)
		if err == nil {
			return tool, nil
		}
	}

	return nil, fmt.Errorf("tool %s not found in registry", toolName)
}

// installToolWithVersion performs the actual installation of a tool with a specific version
func (ins *InstallerService) installToolWithVersion(tool *models.ToolInfo, version string, versionInfo *models.VersionInfo) error {
	// Create a temporary directory for download
	tempDir, err := os.MkdirTemp("", "cntm-install-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Cleanup temp dir

	// Step 1: Download the ZIP file
	zipPath := filepath.Join(tempDir, tool.Name+".zip")
	if err := ins.downloadToolVersion(tool.Name, versionInfo, zipPath); err != nil {
		return fmt.Errorf("failed to download tool: %w", err)
	}

	// Step 2: Verify integrity if hash is available
	if versionInfo.File != "" {
		// Note: The registry doesn't currently include SHA256 hashes in VersionInfo
		// We calculate it after download for storage in lock file
	}

	// Calculate hash for lock file
	hash, err := ins.fsManager.CalculateSHA256(zipPath)
	if err != nil {
		return fmt.Errorf("failed to calculate integrity hash: %w", err)
	}

	// Step 3: Determine installation directory
	destDir := ins.getInstallPath(tool.Name, tool.Type)

	// Step 4: If updating, backup the old installation
	var backupDir string
	if _, err := os.Stat(destDir); err == nil {
		backupDir = destDir + ".backup"
		if err := os.Rename(destDir, backupDir); err != nil {
			return fmt.Errorf("failed to backup existing installation: %w", err)
		}
		// Cleanup backup on success
		defer func() {
			if backupDir != "" {
				os.RemoveAll(backupDir)
			}
		}()
	}

	// Step 5: Extract ZIP to destination
	if err := ins.fsManager.ExtractZIP(zipPath, destDir); err != nil {
		// Rollback: restore backup if it exists
		if backupDir != "" {
			os.RemoveAll(destDir)
			os.Rename(backupDir, destDir)
		}
		return fmt.Errorf("failed to extract ZIP: %w", err)
	}

	// Step 6: Update lock file
	installedTool := &models.InstalledTool{
		Version:     version,
		Type:        tool.Type,
		InstalledAt: time.Now(),
		Source:      "registry",
		Integrity:   hash,
	}

	if err := ins.lockFileService.AddTool(tool.Name, installedTool); err != nil {
		// Rollback: remove installed directory and restore backup
		ins.fsManager.RemoveDir(destDir)
		if backupDir != "" {
			os.Rename(backupDir, destDir)
		}
		return fmt.Errorf("failed to update lock file: %w", err)
	}

	// Step 7: Update registry URL in lock file if not set
	currentRegistry, _ := ins.lockFileService.GetRegistry()
	if currentRegistry == "" {
		ins.lockFileService.SetRegistry(ins.config.Registry.URL)
	}

	return nil
}

// downloadToolVersion downloads a specific version of a tool's ZIP file from GitHub
func (ins *InstallerService) downloadToolVersion(toolName string, versionInfo *models.VersionInfo, destPath string) error {
	// Construct the raw GitHub URL for the file
	// Format: https://raw.githubusercontent.com/{owner}/{repo}/{branch}/{path}
	// But we need to use the GitHub API's download URL instead

	// The versionInfo.File contains the path like "tools/commands/go-code-reviewer/v1-0-2.zip"
	// We need to get the download URL from GitHub

	fmt.Printf("Downloading %s (%s)...\n", toolName, formatBytes(versionInfo.Size))

	// Download file with progress bar
	data, err := ins.githubClient.DownloadFile(
		ins.buildDownloadURL(versionInfo.File),
		versionInfo.Size,
		true, // Show progress
	)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Write to destination file
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write downloaded file: %w", err)
	}

	return nil
}

// buildDownloadURL constructs the raw GitHub content URL
func (ins *InstallerService) buildDownloadURL(filePath string) string {
	// Get owner and repo from config
	owner := "nghiadoan-work" // Default from registry
	repo := "claude-tools-registry"
	branch := ins.config.Registry.Branch
	if branch == "" {
		branch = "main"
	}

	// Parse owner/repo from registry URL if available
	// Format: https://github.com/owner/repo
	if ins.config.Registry.URL != "" {
		o, r, err := ParseRepoURL(ins.config.Registry.URL)
		if err == nil {
			owner = o
			repo = r
		}
	}

	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s",
		owner, repo, branch, filePath)
}

// getInstallPath returns the installation directory for a tool
func (ins *InstallerService) getInstallPath(toolName string, toolType models.ToolType) string {
	// Format: .claude/{type}s/{name}/
	// Examples:
	// - .claude/agents/code-reviewer/
	// - .claude/commands/git-helper/
	// - .claude/skills/test-writer/
	return filepath.Join(ins.baseDir, string(toolType)+"s", toolName)
}

// rollback removes a partially installed tool
func (ins *InstallerService) rollback(toolName string, destDir string) error {
	// Remove the destination directory
	if destDir != "" {
		if err := ins.fsManager.RemoveDir(destDir); err != nil {
			// Log but don't fail - rollback is best effort
			fmt.Printf("Warning: failed to remove directory during rollback: %v\n", err)
		}
	}

	// Try to remove from lock file (might not be there)
	ins.lockFileService.RemoveTool(toolName)

	return nil
}

// formatBytes formats a byte count as a human-readable string
func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

// GetInstalledTools returns a list of all installed tools
func (ins *InstallerService) GetInstalledTools() (map[string]*models.InstalledTool, error) {
	return ins.lockFileService.ListTools()
}

// IsInstalled checks if a tool is installed
func (ins *InstallerService) IsInstalled(toolName string) (bool, error) {
	return ins.lockFileService.IsInstalled(toolName)
}

// GetInstalledVersion returns the installed version of a tool
func (ins *InstallerService) GetInstalledVersion(toolName string) (string, error) {
	tool, err := ins.lockFileService.GetTool(toolName)
	if err != nil {
		return "", err
	}
	return tool.Version, nil
}

// ShowProgress displays a progress bar for an operation
func (ins *InstallerService) ShowProgress(description string, total int64) *progressbar.ProgressBar {
	return progressbar.DefaultBytes(
		total,
		description,
	)
}
