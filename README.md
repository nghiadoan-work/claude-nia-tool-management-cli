# cntm - Claude Tool Manager

A package manager for Claude Code tools (agents, commands, and skills). Like npm for Claude.

## Installation

### Using npx (No Installation Required) ⭐ Recommended

Run cntm directly from GitHub without installing:

```bash
npx github:nghiadoan-work/claude-nia-tool-management-cli init
npx github:nghiadoan-work/claude-nia-tool-management-cli search "code review"
npx github:nghiadoan-work/claude-nia-tool-management-cli install code-reviewer
```

**Create a shortcut alias** to avoid typing the long command:

For **zsh** (macOS default):
```bash
echo 'alias cntm="npx github:nghiadoan-work/claude-nia-tool-management-cli"' >> ~/.zshrc
source ~/.zshrc
```

For **bash**:
```bash
echo 'alias cntm="npx github:nghiadoan-work/claude-nia-tool-management-cli"' >> ~/.bashrc
source ~/.bashrc
```

After setting up the alias, you can use it like:
```bash
cntm init
cntm search "code review"
cntm install code-reviewer
```

### Install via npm

Install globally from the repository:

```bash
git clone https://github.com/nghiadoan-work/claude-nia-tool-management-cli.git
cd claude-nia-tool-management-cli
npm install -g .

# Then use directly
cntm init
cntm search "code review"
cntm install code-reviewer
```

### Build from Source (Go Only)

If you prefer to build with Go without npm:

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

# Create a new tool (agent, command, or skill)
cntm create

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

### Project Setup
- `cntm init` - Initialize .claude directory structure

### Tool Creation
- `cntm create` - Create a new tool (interactive)
- `cntm create --type agent --name "My Agent"` - Create an agent
- `cntm create --type command --name "My Command"` - Create a command
- `cntm create --type skill --name "My Skill"` - Create a skill

### Tool Management
- `cntm search <query>` - Search for tools in registry
- `cntm install <name>` - Install a tool from registry
- `cntm update --all` - Update all installed tools
- `cntm remove <name>` - Remove an installed tool

### Publishing
- `cntm publish <name>` - Publish your tool to registry

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
