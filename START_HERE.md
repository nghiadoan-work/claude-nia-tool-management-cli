# START HERE - Project Overview

## What This Project Is

This is a **Package Manager for Claude Code Tools** - think npm/pip for Claude agents, commands, and skills.

## The Correct Approach (Package Manager)

Based on your clarification, you want a tool that:

✅ **Pulls** tools from a remote GitHub registry
✅ **Installs** them to your local `.claude/` directory
✅ **Updates** tools to newer versions
✅ **Publishes** your own tools to share with others

Like this:

```bash
# Search for tools in remote registry
tool search "code review"

# Install from registry (downloads ZIP, extracts to .claude/)
tool install code-reviewer

# Check for updates
tool outdated

# Update to latest
tool update code-reviewer

# Create and publish your own
tool create agent my-agent
tool publish my-agent
```

## Document Map

### 📚 Documentation

All documentation is now in the `docs/` folder:

1. **[docs/REQUIREMENTS.md](docs/REQUIREMENTS.md)** ⭐ START HERE
   - Complete feature requirements
   - How the package manager works
   - GitHub registry structure
   - All commands explained

2. **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)**
   - System design and components
   - Service layer breakdown
   - Data flow examples
   - API interfaces

3. **[docs/ROADMAP.md](docs/ROADMAP.md)**
   - 10-week implementation plan
   - Phase-by-phase milestones
   - What to build first
   - Testing strategy

4. **[docs/SETUP.md](docs/SETUP.md)**
   - Development setup guide
   - Prerequisites and installation
   - IDE configuration
   - Troubleshooting

5. **[README.md](README.md)** (Root)
   - Quick start guide
   - Command reference
   - Configuration examples
   - User documentation

## How It Works

### Remote Registry (GitHub Repository)

```
github.com/yourusername/claude-tools-registry/
├── registry.json                    # Catalog of all tools
├── agents/
│   └── code-reviewer/
│       ├── code-reviewer.zip       # Zipped agent files
│       └── metadata.json           # Version, description
├── commands/
│   └── git-helper/
│       ├── git-helper.zip
│       └── metadata.json
└── skills/
    └── api-integration/
        ├── api-integration.zip
        └── metadata.json
```

### Installation Flow

```
1. User runs: tool install code-reviewer

2. CLI fetches registry.json from GitHub

3. Finds code-reviewer in catalog

4. Downloads code-reviewer.zip from GitHub

5. Extracts to .claude/agents/code-reviewer/

6. Updates .claude-lock.json with installed version

7. Done! Tool is ready to use
```

### Publishing Flow

```
1. User creates tool: tool create agent my-agent

2. User develops tool in .claude/agents/my-agent/

3. User publishes: tool publish my-agent

4. CLI zips the tool directory

5. Creates metadata.json with version info

6. Pushes to GitHub (creates Pull Request)

7. After PR is merged, tool is available for others!
```

## Quick Start

### 1. Install Go

```bash
brew install go  # macOS
```

### 2. Initialize Project

```bash
cd /Volumes/ex-macmini-a/claude_projects/agent_skill_cli_go

# Initialize Go module
go mod init github.com/yourusername/claude-tools-cli

# Install dependencies
go get github.com/spf13/cobra@latest
go get github.com/google/go-github/v56/github
go get golang.org/x/oauth2
go get github.com/schollz/progressbar/v3
go get gopkg.in/yaml.v3
```

### 3. Create Test Registry

Create a new GitHub repository to act as your registry:

```bash
# On GitHub, create: claude-tools-registry

# Clone and initialize
git clone https://github.com/yourusername/claude-tools-registry
cd claude-tools-registry

mkdir -p agents commands skills

# Create registry.json
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

git add .
git commit -m "Initialize registry"
git push
```

### 4. Start Development

Follow the roadmap in **[docs/ROADMAP.md](docs/ROADMAP.md)**:

**Week 1**: Build models and config
**Week 2-3**: GitHub client and registry service
**Week 4-5**: Installation system
**Week 6**: Update system
**Week 7-8**: Publishing system
**Week 9-10**: Polish and docs

## Key Features

### Core Commands

| Command | Description |
|---------|-------------|
| `tool search <query>` | Search remote registry |
| `tool install <name>` | Download and install tool |
| `tool update <name>` | Update to latest version |
| `tool publish <name>` | Publish your tool |
| `tool list` | Show installed tools |
| `tool list --remote` | Show available tools |
| `tool outdated` | Check for updates |

### Files Created

| File | Purpose |
|------|---------|
| `.claude/` | Installed tools directory |
| `.claude-lock.json` | Tracks installed versions |
| `~/.claude-tools-config.yaml` | Global config |

## Architecture Summary

```
CLI Layer (cmd/)
    ↓
Service Layer (internal/services/)
    ├── RegistryService   - Fetch and search registry
    ├── InstallerService  - Download and install
    ├── UpdaterService    - Check for updates
    ├── PublisherService  - Publish tools
    └── GitHubClient      - GitHub API wrapper
    ↓
Data Layer (internal/data/)
    ├── FSManager         - ZIP/file operations
    └── CacheManager      - Cache registry locally
    ↓
GitHub Registry + Local .claude/
```

## Next Steps

1. ✅ **Read [docs/REQUIREMENTS.md](docs/REQUIREMENTS.md)** - Understand what you're building
2. ✅ **Read [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Understand how to build it
3. ✅ **Follow [docs/ROADMAP.md](docs/ROADMAP.md)** - Build it step by step
4. ✅ **Refer to [README.md](README.md)** - For user-facing documentation

## Questions?

Common questions answered:

**Q: Where are tools stored?**
A: Remote: GitHub registry repo. Local: `.claude/` directory.

**Q: How do I share a tool?**
A: Run `tool publish my-tool` - creates PR to registry.

**Q: Can I use private repos?**
A: Yes! Just configure your GitHub token.

**Q: What if I want multiple registries?**
A: That's a v2 feature. v1 supports one registry.

**Q: Do I need authentication?**
A: For public registries: No (for install). Yes (for publish).
For private registries: Yes (GitHub PAT).

## Comparison

| Feature | v1 (Wrong) | v2 (Correct) |
|---------|------------|--------------|
| Remote registry | ❌ Local only | ✅ GitHub-based |
| Pull/Install | ❌ No | ✅ Yes |
| Push/Publish | ❌ No | ✅ Yes |
| Discovery | ❌ Local search | ✅ Registry search |
| Sharing | ❌ Manual files | ✅ Automated via PR |
| Updates | ❌ Manual | ✅ Version-aware |

## Ready to Build?

Start with **Phase 1** in [docs/ROADMAP.md](docs/ROADMAP.md):
1. Set up Go project
2. Create models
3. Implement config service
4. Then move to GitHub integration!

Good luck! 🚀
