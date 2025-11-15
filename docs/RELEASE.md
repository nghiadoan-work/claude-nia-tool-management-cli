# cntm v1.0.0 - Claude Code Package Manager

**Release Date:** November 15, 2025

We're excited to announce the first stable release of **cntm** (Claude Nia Tool Management CLI) - a package manager for Claude Code tools, similar to npm for Node.js!

## What is cntm?

`cntm` is a command-line tool that helps you manage Claude Code agents, commands, and skills. It provides:

- **Easy Installation** - Install tools from a GitHub registry with one command
- **Automatic Updates** - Keep your tools up-to-date with version management
- **Tool Discovery** - Search and browse available tools
- **Publishing** - Share your own tools with the community
- **Multi-Platform** - Works on macOS, Linux, and Windows

## Key Features

### Installation & Management
```bash
# Initialize your project
cntm init

# Install a tool
cntm install code-reviewer

# Install with version pinning
cntm install git-helper@1.0.0

# Install multiple tools
cntm install agent1 agent2 agent3

# Remove a tool
cntm remove old-tool
```

### Search & Discovery
```bash
# Search for tools
cntm search "code review"

# Browse all tools
cntm browse

# Show trending tools
cntm trending

# View tool details
cntm info code-reviewer
```

### Updates
```bash
# Check for outdated tools
cntm outdated

# Update a specific tool
cntm update code-reviewer

# Update all tools
cntm update --all
```

### Publishing (Structure Ready)
```bash
# Create a new tool
cntm create agent my-agent

# Publish to registry
cntm publish my-agent
```

## Installation

### Quick Install (Unix/Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/claude-nia-tool-management-cli/main/scripts/install.sh | bash
```

### Quick Install (Windows)

```powershell
iwr -useb https://raw.githubusercontent.com/yourusername/claude-nia-tool-management-cli/main/scripts/install.ps1 | iex
```

### Manual Installation

Download the appropriate binary for your platform from the [releases page](https://github.com/yourusername/claude-nia-tool-management-cli/releases):

- **macOS (Intel)**: `cntm-1.0.0-darwin-amd64.tar.gz`
- **macOS (Apple Silicon)**: `cntm-1.0.0-darwin-arm64.tar.gz`
- **Linux (x64)**: `cntm-1.0.0-linux-amd64.tar.gz`
- **Linux (ARM64)**: `cntm-1.0.0-linux-arm64.tar.gz`
- **Windows (x64)**: `cntm-1.0.0-windows-amd64.zip`

Extract and move to your PATH:

```bash
# Unix/Linux/macOS
tar -xzf cntm-*.tar.gz
sudo mv cntm-* /usr/local/bin/cntm
chmod +x /usr/local/bin/cntm

# Verify
cntm version
```

### Build from Source

```bash
git clone https://github.com/yourusername/claude-nia-tool-management-cli.git
cd claude-nia-tool-management-cli
go build -o cntm
./cntm version
```

## What's Included

### Core Functionality
- ✅ Install tools from GitHub registry
- ✅ Update tools with semantic versioning
- ✅ Search and browse available tools
- ✅ Remove/uninstall tools
- ✅ Lock file management (`.claude-lock.json`)
- ✅ Multi-level configuration support
- ✅ Publishing infrastructure (ready for use)

### Security Features
- ✅ SHA256 integrity verification
- ✅ Path traversal protection
- ✅ ZIP bomb protection
- ✅ Secure token storage
- ✅ Safe extraction with validation

### Developer Experience
- ✅ Colored terminal output
- ✅ Progress bars for downloads
- ✅ Interactive confirmations
- ✅ Helpful error messages
- ✅ JSON/YAML output support
- ✅ Table formatting for listings

### Platform Support
- ✅ macOS (Intel and Apple Silicon)
- ✅ Linux (x64 and ARM64)
- ✅ Windows (x64)

## Documentation

Comprehensive documentation is available in the repository:

- **[README.md](../README.md)** - Quick start guide
- **[COMMANDS.md](./COMMANDS.md)** - Complete command reference (500+ lines)
- **[CONFIGURATION.md](./CONFIGURATION.md)** - Configuration guide (600+ lines)
- **[TROUBLESHOOTING.md](./TROUBLESHOOTING.md)** - Common issues and solutions (500+ lines)
- **[ARCHITECTURE.md](./ARCHITECTURE.md)** - System design
- **[CHANGELOG.md](../CHANGELOG.md)** - Full release notes

## Quick Start

```bash
# 1. Install cntm
curl -fsSL https://raw.githubusercontent.com/yourusername/claude-nia-tool-management-cli/main/scripts/install.sh | bash

# 2. Initialize your project
mkdir my-project && cd my-project
cntm init

# 3. Configure your registry (if using a custom registry)
cntm config set registry.url https://github.com/yourorg/your-registry

# 4. Search for tools
cntm search "code review"

# 5. Install a tool
cntm install code-reviewer

# 6. Use the tool with Claude Code!
# Your tool is now in .claude/agents/code-reviewer/
```

## Configuration

cntm supports multiple configuration levels:

1. **Global**: `~/.claude-tools-config.yaml`
2. **Project**: `./.claude-tools-config.yaml`
3. **Environment variables**: `CLAUDE_REGISTRY_URL`, etc.
4. **Command-line flags**: `--registry`, `--path`, etc.

Example configuration:

```yaml
registry:
  url: "https://github.com/username/claude-tools-registry"
  branch: "main"
  auth_token: ""  # GitHub PAT for private repos

local:
  default_path: ".claude"
  auto_update_check: true
  update_check_interval: 86400
```

## Examples

### Install and Use a Code Review Agent

```bash
# Search for code review tools
cntm search "code review" --type agent

# Install the tool
cntm install code-reviewer

# The agent is now available in .claude/agents/code-reviewer/
# Use it with Claude Code!
```

### Keep Your Tools Updated

```bash
# Check what's outdated
cntm outdated

# Update everything
cntm update --all --yes
```

### Create and Share Your Own Tool

```bash
# Create a new agent
cntm create agent awesome-debugger

# Develop your agent...
# (Edit files in .claude/agents/awesome-debugger/)

# Publish to share (when ready)
cntm publish awesome-debugger --version 1.0.0
```

## Technical Details

### Built With
- **Go 1.21+** - Performance and cross-platform support
- **Cobra** - CLI framework
- **go-github** - GitHub API client
- **Clean Architecture** - Maintainable and testable code

### Test Coverage
- Services: 72%+
- Config: 88%
- Data Layer: 80.1%
- Models: 80.6%
- UI: 63.6%

### Performance
- Registry caching with TTL
- Streaming downloads
- Efficient ZIP operations
- Minimal memory footprint

## Known Limitations

- Publishing requires manual GitHub token setup
- Single registry support (multiple registries planned for v1.1)
- No dependency management yet (planned for v1.2)
- No rollback functionality (planned for v1.3)

## Roadmap

### v1.1 (Q1 2026)
- Multiple registry support
- Private registry authentication
- Registry priorities

### v1.2 (Q2 2026)
- Tool dependencies
- Automatic dependency resolution
- Dependency graph

### v1.3 (Q3 2026)
- Tool ratings and reviews
- Auto-update mechanism
- Installation rollback

### v2.0 (Q4 2026)
- Web-based registry browser
- Community features
- CI/CD integrations

## Support & Contributing

- **Issues**: [GitHub Issues](https://github.com/yourusername/claude-nia-tool-management-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/claude-nia-tool-management-cli/discussions)
- **Contributing**: See [CONTRIBUTING.md](../CONTRIBUTING.md)

## Credits

Built with ❤️ for the Claude Code community.

Special thanks to:
- The Go team for an excellent language and toolchain
- The Cobra, go-github, and other open-source library maintainers
- Early testers and contributors

## License

[Your License Here]

---

**Download now and start managing your Claude Code tools like a pro!**

[📥 Download Latest Release](https://github.com/yourusername/claude-nia-tool-management-cli/releases/tag/v1.0.0)
