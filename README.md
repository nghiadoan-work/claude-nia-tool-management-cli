# cntm - Claude Tool Manager

A package manager for Claude Code tools (agents, commands, and skills). Like npm for Claude.

## Installation

### Install from GitHub via npm

Install directly from the GitHub repository:

```bash
npm install -g github:nghiadoan-work/claude-nia-tool-management-cli

# Then use directly
cntm init
cntm search "code review"
cntm install code-reviewer
```

### Using npx (No Installation Required)

Run cntm directly from GitHub without installing:

```bash
npx github:nghiadoan-work/claude-nia-tool-management-cli init
npx github:nghiadoan-work/claude-nia-tool-management-cli search "code review"
npx github:nghiadoan-work/claude-nia-tool-management-cli install code-reviewer
```

### Build from Source

If you prefer to build from source with Go:

```bash
git clone https://github.com/nghiadoan-work/claude-nia-tool-management-cli.git
cd claude-nia-tool-management-cli
go build -o cntm
./cntm --version
```

## Quick Start

### Usage

```bash
# Initialize project
cntm init

# Search for tools
cntm search "code review"

# Install a tool
cntm install code-reviewer

# Update tools
cntm update --all

# Publish your tool
cntm publish my-agent

# Remove a tool
cntm remove code-reviewer
```

## Configuration

Create `.claude-tools-config.yaml` in your project or `~/.claude-tools-config.yaml` globally:

```yaml
registry:
  url: https://github.com/yourusername/your-registry
  branch: main
  auth_token: your_github_token  # Optional, for private repos

local:
  default_path: .claude
  auto_update_check: true
  update_check_interval: 86400

publish:
  default_author: Your Name
  auto_version_bump: patch
  create_pr: true
```

Project-level config overrides global config.

## Commands

- `cntm init` - Initialize .claude directory
- `cntm search <query>` - Search for tools
- `cntm install <name>` - Install a tool
- `cntm update --all` - Update all tools
- `cntm publish <name>` - Publish your tool
- `cntm remove <name>` - Remove a tool

## Directory Structure

```
your-project/
├── .claude/
│   ├── agents/
│   ├── commands/
│   ├── skills/
│   ├── AGENT_TEMPLATE_GUIDE.md
│   ├── SKILL_TEMPLATE_GUIDE.md
│   ├── COMMAND_TEMPLATE_GUIDE.md
│   └── .claude-lock.json
└── .claude-tools-config.yaml  # Optional project config
```

## License

MIT
