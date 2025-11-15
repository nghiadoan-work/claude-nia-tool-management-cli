# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-15

### Added

#### Core Features
- **Installation System** - Install Claude Code tools from GitHub registry
  - Install single or multiple tools
  - Version pinning support (e.g., `tool@1.2.0`)
  - Integrity verification with SHA256 checksums
  - Progress tracking for downloads
  - Automatic rollback on installation failure
  - Force reinstall with `--force` flag

- **Update System** - Keep tools up-to-date
  - Check for outdated tools with `cntm outdated`
  - Update specific tools or all at once
  - Semantic version comparison
  - Interactive confirmation prompts
  - JSON output support

- **Search & Discovery**
  - Search tools by name, description, tags, and author
  - List all available tools in registry
  - View detailed tool information
  - Browse trending/popular tools
  - Filter by tool type (agent, command, skill)
  - Sort by downloads or last updated
  - Table and JSON output formats

- **Publishing** (Structure Ready)
  - Create new tools locally with templates
  - Publish tools to GitHub registry via pull requests
  - Automatic metadata generation
  - Version management

- **Management Commands**
  - Initialize `.claude` directory structure
  - Remove/uninstall tools
  - List installed tools
  - Version command with build information

#### GitHub Integration
- **GitHub Client** with comprehensive features:
  - Fetch files from repositories
  - Download large files with progress indication
  - Personal Access Token (PAT) authentication
  - Rate limit handling and warnings
  - Retry logic with exponential backoff
  - Context-based cancellation support

- **Registry Service**:
  - Fetch and parse registry.json from GitHub
  - Search functionality across all tool metadata
  - Cache management with TTL-based expiration
  - Version comparison and updates

#### Security Features
- **Path Traversal Prevention** - Validates all ZIP file paths
- **ZIP Bomb Protection** - Limits file size and count
- **Integrity Verification** - SHA256 checksum validation
- **Safe Extraction** - Validates destination paths
- **Token Security** - Secure configuration storage

#### Configuration
- **Multi-level Configuration**:
  - Global config: `~/.claude-tools-config.yaml`
  - Project config: `.claude-tools-config.yaml`
  - Environment variables
  - Command-line flags
  - Precedence: ENV > Project > Global > Defaults

- **Lock File Management**:
  - Track installed tools in `.claude-lock.json`
  - Version tracking
  - Integrity hashes
  - Installation timestamps
  - Atomic operations

#### User Experience
- **Rich CLI Interface**:
  - Colored output for better readability
  - Progress bars for downloads
  - Spinners for operations
  - Interactive confirmation prompts
  - Clear error messages with helpful hints
  - Table formatting for listings
  - JSON/YAML output support

- **UI Utilities Package**:
  - Success/Error/Warning/Info messages
  - Progress bars and spinners
  - Confirmation prompts
  - Relative time formatting
  - Consistent styling across all commands

#### Testing
- **Comprehensive Test Suite**:
  - Unit tests for all services (72%+ coverage)
  - Config tests (88% coverage)
  - Data layer tests (80.1% coverage)
  - Models tests (80.6% coverage)
  - UI utilities tests (63.6% coverage)
  - Table-driven test patterns
  - Mock implementations

#### Documentation
- **Complete Documentation Set**:
  - README.md - Quick start and overview
  - COMMANDS.md - Detailed command reference (500+ lines)
  - CONFIGURATION.md - Configuration guide (600+ lines)
  - TROUBLESHOOTING.md - Common issues and solutions (500+ lines)
  - ARCHITECTURE.md - System design
  - REQUIREMENTS.md - Detailed specifications
  - ROADMAP.md - Development phases
  - SETUP.md - Development setup

#### Build & Release
- **Multi-Platform Support**:
  - macOS (amd64, arm64 / Apple Silicon)
  - Linux (amd64, arm64)
  - Windows (amd64)

- **Installation Scripts**:
  - Unix shell script (install.sh)
  - Windows PowerShell script (install.ps1)
  - Automatic platform detection
  - Checksum verification
  - PATH management

- **Build System**:
  - Multi-platform build script
  - Version injection via ldflags
  - Binary size optimization
  - Automatic archive creation
  - SHA256 checksum generation

### Commands

#### Discovery Commands
- `cntm search <query>` - Search for tools
- `cntm list` - List installed tools
- `cntm list --remote` - List available tools in registry
- `cntm info <name>` - Show detailed tool information
- `cntm browse` - Browse all tools
- `cntm trending` - Show trending tools

#### Installation Commands
- `cntm init` - Initialize project
- `cntm install <name>` - Install tool(s)
- `cntm remove <name>` - Remove tool(s)
- `cntm update <name>` - Update tool(s)
- `cntm update --all` - Update all tools
- `cntm outdated` - Check for updates

#### Publishing Commands
- `cntm create <type> <name>` - Create new tool
- `cntm publish <name>` - Publish tool to registry

#### Utility Commands
- `cntm version` - Show version information
- `cntm config` - Manage configuration

### Technical Details

#### Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/google/go-github/v56/github` - GitHub API client
- `golang.org/x/oauth2` - OAuth2 authentication
- `github.com/schollz/progressbar/v3` - Progress indication
- `github.com/olekukonko/tablewriter` - Table formatting
- `gopkg.in/yaml.v3` - YAML parsing
- `github.com/stretchr/testify` - Testing utilities
- `github.com/fatih/color` - Colored terminal output
- `github.com/briandowns/spinner` - Loading spinners

#### Architecture
- **Clean Architecture** with clear layer separation:
  - CMD Layer - User interface only
  - Service Layer - Business logic
  - Data Layer - File system and GitHub access
  - Models - Pure data structures

- **Design Patterns**:
  - Dependency injection
  - Interface-based design
  - Repository pattern
  - Service pattern
  - Error wrapping with context

#### Code Quality
- Go 1.21+ compatibility
- `go fmt` formatted
- `golangci-lint` compliant
- Comprehensive error handling
- Context-based cancellation
- Race condition prevention

### Performance
- Registry caching with TTL
- Streaming downloads for large files
- Parallel operations where possible
- Efficient ZIP operations
- Minimal memory footprint

### Security
- No credentials in logs or output
- Secure token storage
- Path validation and sanitization
- Resource limits (ZIP bombs, file size)
- Checksum verification

## [Unreleased]

### Planned for v1.1
- Multiple registry support
- Private registry authentication
- Registry priorities and fallbacks
- Dependency management

### Planned for v1.2
- Tool dependencies
- Automatic dependency resolution
- Dependency graph visualization

### Planned for v1.3
- Tool ratings and reviews
- Usage analytics
- Auto-update mechanism
- Installation rollback
- Tool aliases

### Planned for v2.0
- Web-based registry browser
- Online tool editor
- Community features
- CI/CD integrations

---

## Version History

- **v1.0.0** (2025-11-15) - Initial release
  - Core functionality complete
  - Multi-platform support
  - Comprehensive documentation
  - Production-ready

[1.0.0]: https://github.com/yourusername/claude-nia-tool-management-cli/releases/tag/v1.0.0
