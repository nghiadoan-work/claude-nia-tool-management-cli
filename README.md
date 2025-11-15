# Claude Nia Tool Management CLI (cntm)

A package manager CLI for Claude Code tools - like npm for Claude agents, commands, and skills. Pull from and push to a GitHub-based registry.

**Project:** `claude-nia-tool-management-cli`
**CLI Command:** `cntm`

## Overview

This tool helps you:
- **Install** agents/commands/skills from a shared registry
- **Update** your installed tools to the latest versions
- **Publish** your own tools to share with others
- **Browse** and discover available tools

## Quick Start

### Installation

#### Quick Install (Recommended)

**Unix/Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/yourusername/claude-nia-tool-management-cli/main/scripts/install.sh | bash
```

**Windows (PowerShell):**
```powershell
iwr -useb https://raw.githubusercontent.com/yourusername/claude-nia-tool-management-cli/main/scripts/install.ps1 | iex
```

#### Download Pre-built Binary

Download the latest release for your platform from [Releases](https://github.com/yourusername/claude-nia-tool-management-cli/releases):

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
```

#### Build from Source

```bash
# Clone and build
git clone https://github.com/yourusername/claude-nia-tool-management-cli.git
cd claude-nia-tool-management-cli
go build -o cntm

# Verify installation
./cntm version
```

Requirements:
- Go 1.21 or later
- Git

### Basic Usage

```bash
# Initialize your project
cntm init

# Search for tools
cntm search "code review"

# Install a tool
cntm install code-reviewer

# List installed tools
cntm list

# Check for updates
cntm outdated

# Update a tool
cntm update code-reviewer

# Create and publish your own tool
cntm create agent my-agent
cntm publish my-agent
```

## How It Works

### Remote Registry (GitHub)

Tools are stored in a GitHub repository with this structure:

```
claude-tools-registry/
├── registry.json           # Catalog of all available tools
├── agents/
│   └── code-reviewer/
│       ├── code-reviewer.zip
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

### Local Installation

Tools are installed to your project's `.claude/` directory:

```
your-project/
├── .claude/
│   ├── agents/
│   │   └── code-reviewer/      # Installed from registry
│   ├── commands/
│   └── skills/
└── .claude-lock.json           # Tracks installed tools
```

## Command Reference

### Discovery & Search

```bash
# Search registry for tools
cntm search <query>
cntm search "code review" --type agent

# List available tools in registry
cntm list --remote
cntm list --remote --type command --sort downloads

# Show tool information
cntm info <name>
cntm info code-reviewer

# Browse popular tools
cntm browse
cntm browse --sort downloads

# Show trending tools
cntm trending
```

### Installation

```bash
# Install a tool
cntm install <name>
cntm install code-reviewer

# Install specific version
cntm install code-reviewer@1.2.0

# Install multiple tools
cntm install agent1 agent2 agent3

# List installed tools
cntm list
cntm list --type agent
```

### Updates

```bash
# Check for outdated tools
cntm outdated

# Update specific tool
cntm update <name>
cntm update code-reviewer

# Update to specific version
cntm update code-reviewer@1.3.0

# Update all tools
cntm update --all
```

### Publishing

```bash
# Create new tool locally
cntm create <type> <name>
cntm create agent my-agent
cntm create command my-command
cntm create skill my-skill

# Publish to registry (creates PR)
cntm publish <name>
cntm publish my-agent --version 1.0.0 --message "Initial release"

# Update existing published tool
cntm publish my-agent --version 1.1.0 --message "Bug fixes"
```

### Management

```bash
# Initialize .claude directory
cntm init

# Remove installed tool
cntm remove <name>
cntm uninstall <name>

# Verify tool integrity
cntm verify <name>

# Show configuration
cntm config

# Set configuration value
cntm config set registry.url https://github.com/user/repo
```

## Global Flags

```bash
--path, -p <path>      Custom .claude directory path
--registry <url>       Override registry URL
--verbose, -v          Verbose output
--output <format>      Output format (table, json, yaml)
--help, -h             Display help
--version              Display version
```

## Configuration

Configuration file: `~/.claude-tools-config.yaml`

```yaml
# Remote registry configuration
registry:
  url: "https://github.com/username/claude-tools-registry"
  branch: "main"
  auth_token: ""  # GitHub PAT for private repos and publishing

# Local settings
local:
  default_path: ".claude"
  auto_update_check: true
  update_check_interval: 86400  # 24 hours in seconds

# Publishing settings
publish:
  default_author: "your@email.com"
  auto_version_bump: "patch"  # patch, minor, major
  create_pr: true             # create PR instead of direct commit
```

### Environment Variables

Override config with environment variables:

```bash
export CLAUDE_REGISTRY_URL=https://github.com/user/repo
export CLAUDE_REGISTRY_TOKEN=ghp_xxxxx
export CLAUDE_REGISTRY_BRANCH=main
export CLAUDE_DEFAULT_PATH=.claude
```

## Lock File

`.claude-lock.json` tracks your installed tools:

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
    }
  }
}
```

## Publishing Guide

### 1. Create a Tool Locally

```bash
# Interactive creation
$ cntm create agent my-agent
? Description: My custom code review agent
? Author: john@example.com
? Tags (comma-separated): code-review, automation
✓ Created .claude/agents/my-agent/
```

### 2. Develop Your Tool

Edit the files in `.claude/agents/my-agent/` to build your agent.

### 3. Publish to Registry

```bash
$ cntm publish my-agent
Preparing my-agent for publication...
Creating ZIP archive...
Generating metadata...
? Version (current: none, suggest: 1.0.0): 1.0.0
? Changelog entry: Initial release of my agent
Pushing to GitHub...
✓ Created pull request: https://github.com/.../pull/123

Your tool will be available after PR approval!
```

### 4. Update Your Published Tool

```bash
# Make changes to your tool
# Then publish update

$ cntm publish my-agent --version 1.1.0 --message "Added new features"
```

## Setting Up Your Own Registry

### 1. Create GitHub Repository

```bash
# Create new repo on GitHub
# Clone it locally
git clone https://github.com/yourusername/claude-tools-registry
cd claude-tools-registry
```

### 2. Initialize Registry Structure

```bash
# Create directory structure
mkdir -p agents commands skills

# Create initial registry.json
cat > registry.json << 'EOF'
{
  "version": "1.0",
  "updated_at": "2025-11-14T00:00:00Z",
  "tools": {
    "agents": [],
    "commands": [],
    "skills": []
  }
}
EOF

# Commit and push
git add .
git commit -m "Initialize registry"
git push
```

### 3. Configure CLI to Use Your Registry

```bash
cntm config set registry.url https://github.com/yourusername/claude-tools-registry
cntm config set registry.auth_token YOUR_GITHUB_TOKEN
```

### 4. Publish Your First Tool

```bash
cntm create agent my-first-agent
cntm publish my-first-agent
```

## Authentication

For publishing and private registries, you need a GitHub Personal Access Token (PAT):

### 1. Create GitHub PAT

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Generate new token (classic)
3. Select scopes:
   - `repo` (for private repos)
   - `public_repo` (for public repos)
4. Copy the token

### 2. Configure Token

```bash
# Option 1: Config file
cntm config set registry.auth_token ghp_xxxxx

# Option 2: Environment variable
export CLAUDE_REGISTRY_TOKEN=ghp_xxxxx
```

## Development

### Prerequisites

- Go 1.21+
- Git
- GitHub account

### Build from Source

```bash
# Clone repository
git clone <repo-url>
cd claude-nia-tool-management-cli

# Install dependencies
go mod download

# Build
go build -o cntm

# Run tests
go test ./...
```

### Project Structure

```
claude-nia-tool-management-cli/
├── cmd/                  # CLI commands
├── internal/
│   ├── services/        # Business logic
│   ├── data/            # File system & GitHub access
│   └── config/          # Configuration management
├── pkg/
│   └── models/          # Data models
├── docs/                # Documentation
├── tests/               # Integration tests
├── main.go
└── go.mod
```

## Documentation

- [Requirements](docs/REQUIREMENTS.md) - Detailed requirements
- [Architecture](docs/ARCHITECTURE.md) - System design
- [Roadmap](docs/ROADMAP.md) - Development plan
- [Setup Guide](docs/SETUP.md) - Development setup

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

### Create and Share Your Own Tool

```bash
# Create new agent
cntm create agent awesome-debugger

# Develop your agent...
# (Edit files in .claude/agents/awesome-debugger/)

# Publish to share
cntm publish awesome-debugger --version 1.0.0

# Others can now install it
cntm install awesome-debugger
```

### Keep Tools Updated

```bash
# Check what's outdated
cntm outdated

# Update everything
cntm update --all
```

## Troubleshooting

### "GitHub rate limit exceeded"

Use authentication:
```bash
cntm config set registry.auth_token YOUR_TOKEN
```

### "Tool not found in registry"

Make sure registry URL is correct:
```bash
cntm config
cntm search --remote  # Check available tools
```

### "Permission denied when publishing"

Ensure your GitHub token has correct permissions (repo scope).

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

[Your License]

## Support

For issues and feature requests, open an issue on GitHub.

---

**Built for the Claude Code community** 🚀
# claude-nia-tool-management-cli
