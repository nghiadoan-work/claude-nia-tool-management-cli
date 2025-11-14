# Architecture Document - Claude Code Package Manager

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────┐
│                         CLI Layer                            │
│  (Cobra Commands - search, install, update, publish, etc.)  │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    Service Layer                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Registry   │  │  Installer   │  │  Publisher   │      │
│  │   Service    │  │   Service    │  │   Service    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Updater    │  │   Lock File  │  │   GitHub     │      │
│  │   Service    │  │   Service    │  │   Client     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                   Data Access Layer                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Local FS   │  │   GitHub     │  │    Cache     │      │
│  │   Manager    │  │   API        │  │   Manager    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│              External Resources                              │
│  GitHub Registry Repo  |  Local .claude/  |  Cache           │
└─────────────────────────────────────────────────────────────┘
```

## Component Breakdown

### 1. CLI Layer (cmd/)

**Commands**:
- `root.go`: Root command and global configuration
- `search.go`: Search registry for tools
- `list.go`: List installed/available tools
- `info.go`: Show tool information
- `install.go`: Install tools from registry
- `update.go`: Update installed tools
- `remove.go`: Remove installed tools
- `publish.go`: Publish tools to registry
- `create.go`: Create new tool locally
- `outdated.go`: Check for outdated tools
- `browse.go`: Browse available tools
- `init.go`: Initialize .claude directory
- `config.go`: Manage configuration

**Responsibilities**:
- Parse command-line arguments
- Validate user input
- Call appropriate services
- Format and display output
- Handle errors and exit codes

### 2. Service Layer (internal/services/)

#### 2.1 Registry Service

```go
// RegistryService manages the remote tool registry
type RegistryService interface {
    // Fetch and parse registry from GitHub
    FetchRegistry() (*Registry, error)

    // Search registry for tools
    Search(query string, filter *SearchFilter) ([]*ToolInfo, error)

    // Get specific tool info from registry
    GetToolInfo(name string, toolType ToolType) (*ToolInfo, error)

    // List all tools in registry
    ListTools(filter *ListFilter) ([]*ToolInfo, error)

    // Update local registry cache
    RefreshCache() error
}

type Registry struct {
    Version   string                    `json:"version"`
    UpdatedAt time.Time                 `json:"updated_at"`
    Tools     map[ToolType][]*ToolInfo  `json:"tools"`
}

type ToolInfo struct {
    Name        string    `json:"name"`
    Version     string    `json:"version"`
    Description string    `json:"description"`
    Type        ToolType  `json:"type"`
    Author      string    `json:"author"`
    Tags        []string  `json:"tags"`
    File        string    `json:"file"`  // Path in repo
    Size        int64     `json:"size"`
    Downloads   int       `json:"downloads"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type ToolType string

const (
    ToolTypeAgent   ToolType = "agent"
    ToolTypeCommand ToolType = "command"
    ToolTypeSkill   ToolType = "skill"
)
```

#### 2.2 Installer Service

```go
// InstallerService handles tool installation
type InstallerService interface {
    // Install a tool from registry
    Install(name string, version string) error

    // Install multiple tools
    InstallMultiple(tools []string) error

    // Verify tool after installation
    Verify(name string) error

    // Get installation status
    IsInstalled(name string) (bool, error)
}

type Installation struct {
    Tool        *ToolInfo
    TempPath    string  // Temporary download location
    TargetPath  string  // Final installation path
    Progress    *ProgressTracker
}

// Installation process:
// 1. Fetch tool info from registry
// 2. Download ZIP from GitHub
// 3. Verify integrity (checksum)
// 4. Extract to .claude/<type>/<name>/
// 5. Update lock file
// 6. Clean up temp files
```

#### 2.3 Publisher Service

```go
// PublisherService handles tool publishing
type PublisherService interface {
    // Publish a tool to registry
    Publish(name string, opts *PublishOptions) (*PublishResult, error)

    // Create local tool from template
    CreateLocal(toolType ToolType, name string, opts *CreateOptions) error

    // Validate tool before publishing
    Validate(name string) error

    // Generate metadata for tool
    GenerateMetadata(name string) (*ToolMetadata, error)
}

type PublishOptions struct {
    Version       string
    ChangelogMsg  string
    CreatePR      bool
    Force         bool  // Force publish even if version exists
}

type PublishResult struct {
    Success   bool
    PRUrl     string  // If PR was created
    CommitSHA string  // If direct commit
    Message   string
}

type ToolMetadata struct {
    Name         string            `json:"name"`
    Version      string            `json:"version"`
    Description  string            `json:"description"`
    Type         ToolType          `json:"type"`
    Author       string            `json:"author"`
    Tags         []string          `json:"tags"`
    Dependencies []string          `json:"dependencies"`
    Changelog    map[string]string `json:"changelog"`
}

// Publishing process:
// 1. Validate tool locally
// 2. Create ZIP of tool directory
// 3. Generate/update metadata.json
// 4. Calculate integrity hash
// 5. Fork/clone registry repo (if PR workflow)
// 6. Add files and update registry.json
// 7. Commit and push
// 8. Create PR (or direct commit)
```

#### 2.4 Updater Service

```go
// UpdaterService handles tool updates
type UpdaterService interface {
    // Check for outdated tools
    CheckOutdated() ([]*OutdatedTool, error)

    // Update specific tool
    Update(name string, version string) error

    // Update all outdated tools
    UpdateAll() error

    // Get available versions for tool
    GetVersions(name string) ([]string, error)
}

type OutdatedTool struct {
    Name           string
    Type           ToolType
    CurrentVersion string
    LatestVersion  string
    ChangesSummary string
}

// Update process:
// 1. Read current version from lock file
// 2. Fetch latest version from registry
// 3. Compare versions
// 4. Download new version
// 5. Backup old version (optional)
// 6. Replace with new version
// 7. Update lock file
```

#### 2.5 Lock File Service

```go
// LockFileService manages .claude-lock.json
type LockFileService interface {
    // Read lock file
    Read() (*LockFile, error)

    // Write lock file
    Write(lockFile *LockFile) error

    // Add installed tool to lock file
    AddTool(tool *InstalledTool) error

    // Remove tool from lock file
    RemoveTool(name string) error

    // Update tool version in lock file
    UpdateTool(name string, version string) error

    // Get installed tool info
    GetTool(name string) (*InstalledTool, error)
}

type LockFile struct {
    Version   string                    `json:"version"`
    UpdatedAt time.Time                 `json:"updated_at"`
    Registry  string                    `json:"registry"`
    Tools     map[string]*InstalledTool `json:"tools"`
}

type InstalledTool struct {
    Version     string    `json:"version"`
    Type        ToolType  `json:"type"`
    InstalledAt time.Time `json:"installed_at"`
    Source      string    `json:"source"`      // "registry" or URL
    Integrity   string    `json:"integrity"`   // SHA256 hash
}
```

#### 2.6 GitHub Client Service

```go
// GitHubClient handles all GitHub API interactions
type GitHubClient interface {
    // Fetch file from repo
    FetchFile(path string) ([]byte, error)

    // Download file (for large files like ZIPs)
    Download(path string, dest string) error

    // Create commit in repo
    CreateCommit(files map[string][]byte, message string) (string, error)

    // Create pull request
    CreatePR(title, body, branch string) (string, error)

    // Fork repository (if needed)
    ForkRepo() error

    // Check rate limit
    GetRateLimit() (*RateLimit, error)
}

type GitHubClientImpl struct {
    client       *github.Client
    owner        string
    repo         string
    branch       string
    authToken    string
}

// Uses go-github library for API calls
// Implements retry logic with exponential backoff
// Handles rate limiting
```

### 3. Data Access Layer (internal/data/)

#### 3.1 Local File System Manager

```go
// FSManager handles local file operations
type FSManager interface {
    // Extract ZIP to directory
    ExtractZIP(zipPath, destDir string) error

    // Create ZIP from directory
    CreateZIP(srcDir, zipPath string) error

    // Get tool directory path
    GetToolPath(toolType ToolType, name string) string

    // Check if tool exists locally
    ToolExists(toolType ToolType, name string) (bool, error)

    // Remove tool directory
    RemoveTool(toolType ToolType, name string) error

    // Calculate file/directory hash
    CalculateHash(path string) (string, error)

    // Ensure .claude directory structure exists
    EnsureDirectories() error
}

// Security checks:
// - Validate ZIP doesn't escape target directory
// - Check for zip bombs (size limits)
// - Atomic operations where possible
```

#### 3.2 Cache Manager

```go
// CacheManager handles local caching
type CacheManager interface {
    // Get cached registry
    GetRegistry() (*Registry, error)

    // Set cached registry
    SetRegistry(registry *Registry, ttl time.Duration) error

    // Clear cache
    Clear() error

    // Check if cache is valid
    IsValid() bool
}

type CacheManagerImpl struct {
    cacheDir string
    ttl      time.Duration
}

// Cache structure:
// ~/.cache/claude-tools/
// ├── registry.json
// └── .cache-metadata
```

### 4. Models (pkg/models/)

All the data structures defined in service interfaces above, plus:

```go
// SearchFilter for searching tools
type SearchFilter struct {
    Query         string
    Type          ToolType
    Tags          []string
    Author        string
    MinDownloads  int
}

// ListFilter for listing tools
type ListFilter struct {
    Type      ToolType
    SortBy    SortField
    SortDesc  bool
    Limit     int
}

type SortField string

const (
    SortByName      SortField = "name"
    SortByDownloads SortField = "downloads"
    SortByUpdated   SortField = "updated"
    SortByCreated   SortField = "created"
)

// ProgressTracker for showing progress
type ProgressTracker struct {
    Total     int64
    Current   int64
    StartTime time.Time
}

// Config structure
type Config struct {
    Registry RegistryConfig `yaml:"registry"`
    Local    LocalConfig    `yaml:"local"`
    Publish  PublishConfig  `yaml:"publish"`
}

type RegistryConfig struct {
    URL       string `yaml:"url"`
    Branch    string `yaml:"branch"`
    AuthToken string `yaml:"auth_token"`
}

type LocalConfig struct {
    DefaultPath          string `yaml:"default_path"`
    AutoUpdateCheck      bool   `yaml:"auto_update_check"`
    UpdateCheckInterval  int    `yaml:"update_check_interval"`
}

type PublishConfig struct {
    DefaultAuthor    string `yaml:"default_author"`
    AutoVersionBump  string `yaml:"auto_version_bump"`
    CreatePR         bool   `yaml:"create_pr"`
}
```

## Data Flow Examples

### Example 1: Installing a Tool

```
User: tool install code-reviewer

1. CLI (cmd/install.go)
   └─> Parse arguments
   └─> Call InstallerService.Install("code-reviewer", "latest")

2. InstallerService
   └─> Call RegistryService.GetToolInfo("code-reviewer")
   └─> Call GitHubClient.Download(toolInfo.File, tempPath)
   └─> Call FSManager.ExtractZIP(tempPath, targetPath)
   └─> Calculate integrity hash
   └─> Call LockFileService.AddTool(installedTool)

3. Output
   └─> Show progress bar during download
   └─> Display success message with version
```

### Example 2: Publishing a Tool

```
User: tool publish my-agent --version 1.0.0

1. CLI (cmd/publish.go)
   └─> Parse arguments
   └─> Call PublisherService.Publish("my-agent", opts)

2. PublisherService
   └─> Call PublisherService.Validate("my-agent")
   └─> Call PublisherService.GenerateMetadata("my-agent")
   └─> Call FSManager.CreateZIP(toolPath, tempZipPath)
   └─> Call FSManager.CalculateHash(tempZipPath)
   └─> Call GitHubClient.ForkRepo() [if needed]
   └─> Create branch
   └─> Upload ZIP and metadata
   └─> Update registry.json
   └─> Call GitHubClient.CreatePR()

3. Output
   └─> Show progress during each step
   └─> Display PR URL
```

### Example 3: Checking for Updates

```
User: tool outdated

1. CLI (cmd/outdated.go)
   └─> Call UpdaterService.CheckOutdated()

2. UpdaterService
   └─> Call LockFileService.Read()
   └─> Call RegistryService.FetchRegistry()
   └─> Compare versions for each installed tool

3. Output
   └─> Display table of outdated tools
```

## Configuration Management

### Config File Locations
1. Global: `~/.claude-tools-config.yaml`
2. Project: `./.claude-tools-config.yaml`
3. Environment variables (override config)
4. Priority: ENV > Project > Global > Defaults

### Environment Variables
```bash
CLAUDE_REGISTRY_URL=https://github.com/user/repo
CLAUDE_REGISTRY_TOKEN=ghp_xxx
CLAUDE_REGISTRY_BRANCH=main
CLAUDE_DEFAULT_PATH=.claude
```

## Error Handling

### Error Types

```go
type ErrorType int

const (
    ErrorTypeNetwork ErrorType = iota
    ErrorTypeAuth
    ErrorTypeNotFound
    ErrorTypeAlreadyExists
    ErrorTypeValidation
    ErrorTypeIntegrity
    ErrorTypeRateLimit
    ErrorTypeInternal
)

type CLIError struct {
    Type    ErrorType
    Message string
    Err     error
    Hint    string  // Helpful hint for user
}
```

### Retry Logic

- Network errors: Retry with exponential backoff (max 3 retries)
- Rate limiting: Wait and retry based on GitHub headers
- Other errors: Fail immediately with clear message

## Security Considerations

1. **ZIP Bomb Protection**:
   - Limit uncompressed size (e.g., 100MB)
   - Limit file count (e.g., 1000 files)

2. **Path Traversal Prevention**:
   - Validate all ZIP entry paths
   - Reject paths with `..` or absolute paths

3. **Integrity Verification**:
   - SHA256 hash for all downloads
   - Compare with registry hash

4. **Token Security**:
   - Never log tokens
   - Store securely in config file (chmod 600)
   - Support env variables for CI/CD

5. **HTTPS Only**:
   - All GitHub API calls over HTTPS
   - Reject HTTP redirects

## Performance Optimizations

1. **Registry Caching**:
   - Cache registry.json locally
   - TTL configurable (default: 1 hour)
   - Auto-refresh in background

2. **Parallel Downloads**:
   - Download multiple tools in parallel
   - Connection pooling

3. **Resume Downloads**:
   - Support resuming interrupted downloads
   - Use Range requests

4. **Progress Indication**:
   - Show progress bars for downloads
   - Estimate time remaining

## Testing Strategy

### Unit Tests
- Each service tested independently
- Mock GitHub client
- Mock file system operations
- Test error scenarios

### Integration Tests
- Test with real GitHub API (rate-limited)
- Test actual ZIP operations
- Test lock file operations

### End-to-End Tests
- Full workflow tests
- Test with test registry repo

## Project Structure

```
agent_skill_cli_go/
├── cmd/
│   ├── root.go
│   ├── search.go
│   ├── list.go
│   ├── info.go
│   ├── install.go
│   ├── update.go
│   ├── remove.go
│   ├── publish.go
│   ├── create.go
│   ├── outdated.go
│   ├── browse.go
│   ├── init.go
│   └── config.go
├── internal/
│   ├── services/
│   │   ├── registry.go
│   │   ├── installer.go
│   │   ├── publisher.go
│   │   ├── updater.go
│   │   ├── lockfile.go
│   │   └── github.go
│   ├── data/
│   │   ├── fs.go
│   │   └── cache.go
│   ├── config/
│   │   └── config.go
│   └── ui/
│       └── formatter.go
├── pkg/
│   └── models/
│       └── models.go
├── main.go
├── go.mod
├── go.sum
└── README.md
```

This architecture provides clean separation of concerns, is testable, and supports the package manager workflow effectively.
