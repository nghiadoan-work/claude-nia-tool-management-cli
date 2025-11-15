# Configuration Guide

Complete guide to configuring the Claude Nia Tool Management CLI (cntm).

## Table of Contents

- [Configuration Files](#configuration-files)
- [Configuration Precedence](#configuration-precedence)
- [Configuration Options](#configuration-options)
- [Environment Variables](#environment-variables)
- [Examples](#examples)

---

## Configuration Files

cntm uses YAML configuration files stored in specific locations:

### File Locations

1. **Global Config**: `~/.claude-config.yaml`
   - User-wide settings
   - Applied to all projects

2. **Project Config**: `.claude-config.yaml` (in project root)
   - Project-specific settings
   - Overrides global config

3. **Custom Config**: Specified with `--config` flag
   - Highest priority
   - Useful for testing different registries

### Basic Structure

```yaml
registry:
  url: "https://github.com/your-org/claude-tools-registry"
  branch: "main"
  authToken: ""  # Optional, use CNTM_GITHUB_TOKEN instead

local:
  defaultPath: ".claude"

cache:
  enabled: true
  ttl: 3600  # seconds
```

---

## Configuration Precedence

Configuration is merged in the following order (highest priority first):

1. **Command-line flags** (e.g., `--path /custom`)
2. **Environment variables** (e.g., `CNTM_GITHUB_TOKEN`)
3. **Project config** (`.claude-config.yaml` in current directory)
4. **Global config** (`~/.claude-config.yaml`)
5. **Default values** (built into cntm)

### Example

```bash
# This command uses:
# 1. --path flag for installation path
# 2. CNTM_GITHUB_TOKEN for authentication
# 3. .claude-config.yaml for registry URL
# 4. Default cache TTL (3600s)

export CNTM_GITHUB_TOKEN=ghp_xxx
cntm install --path /custom code-reviewer
```

---

## Configuration Options

### Registry Section

Controls access to the tool registry.

```yaml
registry:
  # GitHub repository URL (required)
  url: "https://github.com/your-org/claude-tools-registry"

  # Branch name (default: "main")
  branch: "main"

  # GitHub Personal Access Token (optional)
  # Recommended to use CNTM_GITHUB_TOKEN env var instead
  authToken: ""

  # Registry path within repo (default: "registry.json")
  registryFile: "registry.json"
```

#### Registry URL Format

Must be a valid GitHub repository URL:

```
https://github.com/owner/repo
```

Examples:
- `https://github.com/claude-tools/registry`
- `https://github.com/myorg/private-tools`

### Local Section

Controls local tool management.

```yaml
local:
  # Base path for tool installation (default: ".claude")
  defaultPath: ".claude"

  # Lock file name (default: ".claude-lock.json")
  lockFile: ".claude-lock.json"
```

#### Installation Paths

Tools are installed under `defaultPath/<type>/<name>/`:

```
.claude/
├── agents/
│   ├── code-reviewer/
│   └── git-helper/
├── commands/
│   └── deploy-tool/
└── skills/
    └── debug-helper/
```

### Cache Section

Controls registry caching behavior.

```yaml
cache:
  # Enable/disable caching (default: true)
  enabled: true

  # Cache TTL in seconds (default: 3600 = 1 hour)
  ttl: 3600

  # Cache directory (default: ".claude/.cache")
  cacheDir: ".claude/.cache"
```

#### Cache Benefits

- Faster searches and listings
- Reduced API calls (avoid rate limiting)
- Offline browsing of cached data

---

## Environment Variables

### CNTM_GITHUB_TOKEN

**Recommended** way to provide GitHub authentication.

```bash
export CNTM_GITHUB_TOKEN=ghp_xxxxxxxxxxxxx
```

**Why use this over config file?**
- Keeps tokens out of version control
- Easier to rotate tokens
- Works with CI/CD secret management

**Getting a Token:**

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Generate new token (classic)
3. Scopes needed:
   - `repo` (for private registries)
   - `public_repo` (for public registries)
4. Copy token and export it

### CNTM_CONFIG

Override config file location.

```bash
export CNTM_CONFIG=/path/to/custom-config.yaml
cntm install tool-name
```

### CNTM_REGISTRY_URL

Override registry URL without editing config file.

```bash
export CNTM_REGISTRY_URL=https://github.com/test-org/test-registry
cntm search test-tool
```

### CNTM_CACHE_ENABLED

Enable/disable cache.

```bash
export CNTM_CACHE_ENABLED=false
cntm list --remote
```

---

## Examples

### Example 1: Basic Setup

**Global config** (`~/.claude-config.yaml`):

```yaml
registry:
  url: "https://github.com/claude-tools/registry"
  branch: "main"

cache:
  enabled: true
  ttl: 3600
```

**Usage:**

```bash
# Set token once
export CNTM_GITHUB_TOKEN=ghp_xxx

# All projects use global config
cd ~/project1
cntm install code-reviewer

cd ~/project2
cntm install git-helper
```

### Example 2: Project-Specific Registry

**Project config** (`.claude-config.yaml`):

```yaml
registry:
  url: "https://github.com/mycompany/private-tools"
  branch: "production"

local:
  defaultPath: "tools"
```

This project uses a private registry while others use the global default.

### Example 3: Multiple Registries

**Development** (`dev-config.yaml`):

```yaml
registry:
  url: "https://github.com/myorg/tools-dev"
  branch: "develop"
```

**Production** (`prod-config.yaml`):

```yaml
registry:
  url: "https://github.com/myorg/tools-prod"
  branch: "main"
```

**Usage:**

```bash
# Test against dev registry
cntm --config dev-config.yaml search new-tool

# Install from production registry
cntm --config prod-config.yaml install stable-tool
```

### Example 4: CI/CD Setup

**GitHub Actions** (`.github/workflows/tools.yml`):

```yaml
name: Install Tools

on: [push]

jobs:
  install:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Install cntm
        run: |
          wget https://github.com/.../cntm-linux-amd64
          chmod +x cntm-linux-amd64
          sudo mv cntm-linux-amd64 /usr/local/bin/cntm

      - name: Install tools
        env:
          CNTM_GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        run: |
          cntm init
          cntm install --yes code-reviewer
          cntm install --yes test-runner
```

### Example 5: Offline Mode

Disable cache refresh for offline work:

```yaml
cache:
  enabled: true
  ttl: 86400  # 24 hours
```

```bash
# Fetch once while online
cntm list --remote

# Work offline
export CNTM_CACHE_ENABLED=true
cntm search tool  # Uses cached data
```

### Example 6: Custom Installation Path

**Monorepo setup** (`.claude-config.yaml`):

```yaml
local:
  defaultPath: ".tools/claude"

cache:
  cacheDir: ".tools/claude/.cache"
```

Result:

```
.tools/
└── claude/
    ├── agents/
    ├── commands/
    ├── skills/
    ├── .cache/
    └── .claude-lock.json
```

---

## Configuration Templates

### Minimal Config

```yaml
registry:
  url: "https://github.com/claude-tools/registry"
```

### Full Config

```yaml
registry:
  url: "https://github.com/your-org/claude-tools-registry"
  branch: "main"
  authToken: ""  # Use CNTM_GITHUB_TOKEN instead
  registryFile: "registry.json"

local:
  defaultPath: ".claude"
  lockFile: ".claude-lock.json"

cache:
  enabled: true
  ttl: 3600
  cacheDir: ".claude/.cache"
```

---

## Best Practices

### Security

1. **Never commit tokens to git**
   ```bash
   # .gitignore
   .claude-config.yaml  # If it contains tokens
   ```

2. **Use environment variables for tokens**
   ```bash
   export CNTM_GITHUB_TOKEN=ghp_xxx
   ```

3. **Use read-only tokens when possible**
   - For consuming tools: `public_repo` scope
   - For publishing: `repo` scope

### Performance

1. **Adjust cache TTL based on usage**
   - Active development: 600 (10 minutes)
   - Stable environment: 86400 (24 hours)

2. **Use project configs for monorepos**
   ```yaml
   local:
     defaultPath: "packages/tools"
   ```

### Collaboration

1. **Commit project config (without tokens)**
   ```yaml
   # .claude-config.yaml - Safe to commit
   registry:
     url: "https://github.com/team/registry"
     branch: "main"

   local:
     defaultPath: ".claude"
   ```

2. **Document token setup in README**
   ```markdown
   ## Setup

   1. Get GitHub token: https://github.com/settings/tokens
   2. Export token: `export CNTM_GITHUB_TOKEN=ghp_xxx`
   3. Install tools: `cntm install --all`
   ```

---

## Troubleshooting

### "Config file not found"

Create default config:

```bash
mkdir -p ~/.config
cat > ~/.claude-config.yaml <<EOF
registry:
  url: "https://github.com/claude-tools/registry"
EOF
```

### "Invalid registry URL"

Ensure URL is in correct format:

```yaml
# ✓ Correct
registry:
  url: "https://github.com/owner/repo"

# ✗ Wrong
registry:
  url: "github.com/owner/repo"
  url: "https://github.com/owner/repo.git"
```

### "Authentication failed"

1. Check token is set:
   ```bash
   echo $CNTM_GITHUB_TOKEN
   ```

2. Verify token has correct scopes
3. Token may have expired - generate new one

### Cache Issues

Clear cache:

```bash
rm -rf .claude/.cache
cntm list --remote  # Rebuild cache
```

---

## Reference

See also:

- [Commands Reference](COMMANDS.md)
- [Troubleshooting Guide](TROUBLESHOOTING.md)
- [Publishing Guide](PUBLISHING.md)
