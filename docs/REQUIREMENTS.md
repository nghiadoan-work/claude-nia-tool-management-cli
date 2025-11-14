# Claude Code Package Manager - Requirements Document

## Project Overview
A CLI package manager for Claude Code tools (agents, commands, and skills) that pulls from and pushes to a GitHub repository, similar to npm for Node.js packages.

## 1. Core Concept

### Package Manager for Claude Code
- **Pull**: Download and install tools from a remote GitHub registry
- **Push**: Publish local tools to the remote GitHub registry
- **Update**: Keep installed tools up-to-date
- **Browse**: Discover available tools in the registry

### Remote Structure (GitHub Repository)
```
claude-tools-registry/                    (GitHub repo)
├── registry.json                         # Catalog/manifest of all tools
├── agents/
│   ├── code-reviewer/
│   │   ├── code-reviewer.zip            # Zipped agent files
│   │   └── metadata.json                # Version, description, etc.
│   └── debugger/
│       ├── debugger.zip
│       └── metadata.json
├── commands/
│   └── git-helper/
│       ├── git-helper.zip
│       └── metadata.json
└── skills/
    └── api-integration/
        ├── api-integration.zip
        └── metadata.json
```

### Local Structure
```
.claude/                                  # User's project
├── agents/
│   └── code-reviewer/                   # Installed from registry
│       └── [agent files]
├── commands/
│   └── git-helper/
│       └── [command files]
└── skills/
    └── api-integration/
        └── [skill files]

.claude-lock.json                         # Tracks installed tools (like package-lock.json)
```

## 2. Core Features

### 2.1 Remote Registry Management

**Registry File Format (registry.json)**:
```json
{
  "version": "1.0",
  "updated_at": "2025-11-14T17:32:00Z",
  "tools": {
    "agents": [
      {
        "name": "code-reviewer",
        "version": "1.2.0",
        "description": "Automated code review agent",
        "author": "john@example.com",
        "tags": ["code-review", "quality"],
        "file": "agents/code-reviewer/code-reviewer.zip",
        "size": 15420,
        "downloads": 1250,
        "created_at": "2025-10-01T10:00:00Z",
        "updated_at": "2025-11-10T15:30:00Z"
      }
    ],
    "commands": [...],
    "skills": [...]
  }
}
```

**Tool Metadata Format (metadata.json)**:
```json
{
  "name": "code-reviewer",
  "version": "1.2.0",
  "description": "Automated code review agent",
  "type": "agent",
  "author": "john@example.com",
  "tags": ["code-review", "quality"],
  "dependencies": [],
  "changelog": {
    "1.2.0": "Added support for Python",
    "1.1.0": "Initial release"
  }
}
```

### 2.2 Install/Pull Operations

**Commands**:
```bash
# Search remote registry
tool search <query>
tool search --type agent "code review"

# List available tools in registry
tool list --remote
tool list --remote --type command

# Install a tool from registry
tool install <name>
tool install code-reviewer
tool install code-reviewer@1.1.0              # specific version

# Install multiple tools
tool install code-reviewer git-helper api-integration

# Show tool information before installing
tool info <name>
```

**Process**:
1. Fetch registry.json from GitHub
2. Find the tool in registry
3. Download the ZIP file from GitHub
4. Extract to `.claude/<type>/<name>/`
5. Update `.claude-lock.json`

### 2.3 Update Operations

**Commands**:
```bash
# Check for updates
tool outdated

# Update specific tool
tool update <name>
tool update code-reviewer

# Update all tools
tool update --all

# Update to specific version
tool update code-reviewer@1.2.0
```

**Process**:
1. Read `.claude-lock.json`
2. Fetch latest registry.json
3. Compare versions
4. Download and replace updated tools

### 2.4 Publish/Push Operations

**Commands**:
```bash
# Create new tool locally (interactive)
tool create agent my-agent
tool create command my-command --template basic

# Publish to registry
tool publish my-agent
tool publish my-agent --version 1.0.0 --message "Initial release"

# Update existing published tool
tool publish my-agent --version 1.1.0 --message "Bug fixes"
```

**Process**:
1. Validate tool locally
2. Create ZIP file of tool directory
3. Create/update metadata.json
4. Push to GitHub (create PR or direct commit depending on permissions)
5. Update registry.json
6. Commit and push to GitHub

### 2.5 Local Management

**Commands**:
```bash
# List installed tools
tool list
tool list --type agent

# Show installed tool details
tool show <name>

# Remove installed tool
tool remove <name>
tool uninstall <name>

# Verify tool integrity
tool verify <name>
```

### 2.6 Browse & Discovery

**Commands**:
```bash
# Browse popular tools
tool browse
tool browse --sort downloads
tool browse --sort recent

# Search by tags
tool search --tag code-review
tool search --tag testing --type agent

# Show trending tools
tool trending
```

## 3. Lock File (.claude-lock.json)

Tracks installed tools and their versions:

```json
{
  "version": "1.0",
  "updated_at": "2025-11-14T17:32:00Z",
  "registry": "https://github.com/username/claude-tools-registry",
  "tools": {
    "code-reviewer": {
      "version": "1.2.0",
      "type": "agent",
      "installed_at": "2025-11-14T17:32:00Z",
      "source": "registry",
      "integrity": "sha256-abc123..."
    },
    "git-helper": {
      "version": "2.0.1",
      "type": "command",
      "installed_at": "2025-11-14T17:35:00Z",
      "source": "registry",
      "integrity": "sha256-def456..."
    }
  }
}
```

## 4. Configuration

**Config File (~/.claude-tools-config.yaml)**:
```yaml
# Remote registry configuration
registry:
  url: "https://github.com/username/claude-tools-registry"
  branch: "main"
  auth_token: ""  # GitHub token for private repos and publishing

# Local settings
default_path: ".claude"
auto_update_check: true
update_check_interval: 86400  # seconds (24 hours)

# Publishing settings
publish:
  default_author: "your@email.com"
  auto_version_bump: "patch"  # patch, minor, major
  create_pr: true             # create PR instead of direct commit
```

## 5. GitHub Integration

### 5.1 Authentication
- Support GitHub personal access tokens
- Store in config or environment variable
- Required for:
  - Pulling from private repos
  - Publishing tools
  - Higher rate limits

### 5.2 Publishing Workflow

**Option 1: Direct Commit (if user has write access)**:
1. Clone/pull latest registry repo
2. Add tool ZIP and metadata
3. Update registry.json
4. Commit and push

**Option 2: Pull Request (recommended)**:
1. Fork registry repo (if needed)
2. Create branch
3. Add tool files
4. Update registry.json
5. Create PR to main registry
6. Wait for approval

### 5.3 GitHub API Usage
- Use GitHub REST API for:
  - Fetching registry.json
  - Downloading ZIP files
  - Creating commits/PRs
  - Checking rate limits

## 6. CLI Command Structure

```bash
tool [global-flags] <command> [args] [flags]

Commands:
  search <query>              Search registry for tools
  list                        List installed tools
  list --remote               List available tools in registry
  info <name>                 Show tool information
  install <name>              Install tool from registry
  update <name>               Update tool to latest version
  update --all                Update all tools
  remove <name>               Remove installed tool
  publish <name>              Publish tool to registry
  create <type> <name>        Create new tool locally
  outdated                    Check for outdated tools
  browse                      Browse available tools
  trending                    Show trending tools
  init                        Initialize .claude directory
  config                      Manage configuration

Global Flags:
  --path, -p <path>          Custom .claude directory path
  --registry <url>           Override registry URL
  --verbose, -v              Verbose output
  --help, -h                 Display help
  --version                  Display version
```

## 7. Features Breakdown

### Phase 1: Core Pull/Install (MVP)
- Fetch and parse registry.json from GitHub
- Search registry
- List available tools
- Install tool (download ZIP, extract)
- List installed tools
- Basic configuration

### Phase 2: Updates & Lock File
- Generate .claude-lock.json
- Check for outdated tools
- Update tools
- Verify tool integrity

### Phase 3: Publishing
- Create tool locally
- Zip tool files
- Generate metadata
- Push to GitHub (PR workflow)
- Update registry

### Phase 4: Enhanced Discovery
- Browse interface
- Trending tools
- Tag-based search
- Tool ratings/reviews (optional)

## 8. Technical Requirements

### 8.1 Go Dependencies
- **GitHub API**: `github.com/google/go-github/v56/github`
- **OAuth2**: `golang.org/x/oauth2` for auth
- **CLI Framework**: `github.com/spf13/cobra`
- **ZIP handling**: Standard library `archive/zip`
- **HTTP client**: Standard library with retry logic
- **JSON**: Standard library `encoding/json`
- **YAML**: `gopkg.in/yaml.v3`
- **Progress bars**: `github.com/schollz/progressbar/v3`
- **Table output**: `github.com/olekukonko/tablewriter`

### 8.2 File Operations
- ZIP/UNZIP tool directories
- Atomic file operations
- Integrity checking (SHA256)
- Safe extraction (prevent zip bombs)

### 8.3 Error Handling
- Network errors (retry with backoff)
- GitHub rate limiting
- Authentication errors
- Conflict resolution (tool already exists)
- Validation errors

## 9. Security Considerations

1. **ZIP Bomb Protection**: Limit extraction size and file count
2. **Path Traversal**: Validate ZIP contents don't escape target directory
3. **Integrity**: Verify SHA256 checksums
4. **Auth Token**: Never log or print tokens
5. **HTTPS Only**: All GitHub communication over HTTPS
6. **Validate JSON**: Sanitize registry.json content

## 10. User Experience

### 10.1 Installation Experience
```bash
$ tool install code-reviewer
Fetching registry from GitHub...
Found: code-reviewer v1.2.0
Downloading... ████████████████ 100% (15.4 KB)
Extracting to .claude/agents/code-reviewer/
✓ Successfully installed code-reviewer@1.2.0
```

### 10.2 Publishing Experience
```bash
$ tool create agent my-agent
Creating new agent: my-agent
? Description: My custom agent
? Author: john@example.com
? Tags (comma-separated): automation, custom
✓ Created .claude/agents/my-agent/

$ tool publish my-agent
Preparing my-agent for publication...
Creating ZIP archive...
Generating metadata...
? Version (current: none, suggest: 1.0.0): 1.0.0
? Changelog entry: Initial release
Pushing to GitHub...
✓ Created pull request: https://github.com/.../pull/123
Your tool will be available after PR approval!
```

### 10.3 Update Experience
```bash
$ tool outdated
Checking for updates...

Package         Current    Latest
code-reviewer   1.1.0      1.2.0
git-helper      2.0.0      2.0.1

$ tool update --all
Updating 2 tools...
Updating code-reviewer... ✓ 1.1.0 → 1.2.0
Updating git-helper... ✓ 2.0.0 → 2.0.1
```

## 11. Future Enhancements (Post-v1)

- Multiple registry support
- Private registries
- Tool dependencies
- Automatic updates
- Tool usage analytics
- Rating system
- Comments/reviews
- Web UI for registry
- CI/CD integration
- Webhooks for updates

## 12. Success Criteria

- Can fetch and search remote registry
- Can install tools from GitHub
- Can update installed tools
- Can publish new tools via PR
- Lock file tracks installations correctly
- Handles network errors gracefully
- Works with both public and private repos
- Clear, helpful error messages
- Fast operations (<3s for install on good connection)
