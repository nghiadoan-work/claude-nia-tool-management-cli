# Command Reference

Complete reference for all cntm commands.

## Table of Contents

- [Global Flags](#global-flags)
- [init](#init)
- [search](#search)
- [list](#list)
- [info](#info)
- [browse](#browse)
- [install](#install)
- [update](#update)
- [outdated](#outdated)
- [remove](#remove)
- [create](#create)
- [publish](#publish)

---

## Global Flags

These flags are available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Config file path | `.claude-config.yaml` |
| `--help` | `-h` | Show help | - |

---

## init

Initialize a new Claude Code tools project.

### Usage

```bash
cntm init [flags]
```

### Description

Creates the `.claude` directory structure and initializes the lock file. This should be the first command you run in a new project.

### Structure Created

```
.claude/
├── agents/
├── commands/
├── skills/
└── .claude-lock.json
```

### Flags

None

### Examples

```bash
# Initialize current directory
cntm init

# Initialize and then install tools
cntm init
cntm install code-reviewer
```

### Exit Codes

- `0`: Success
- `1`: Directory already exists (with `--force`, reinitializes)

---

## search

Search for tools in the registry.

### Usage

```bash
cntm search <query> [flags]
```

### Description

Search for tools in the remote registry by name, description, tags, or author. Returns a list of matching tools with their details.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--type` | `-t` | string | Filter by type (agent, command, skill) | all |
| `--json` | | bool | Output in JSON format | false |
| `--limit` | `-l` | int | Limit number of results | 20 |

### Examples

```bash
# Search for code review tools
cntm search code-review

# Search for agents only
cntm search git --type agent

# Search with JSON output
cntm search test --json

# Limit results
cntm search test --limit 5
```

### Output

**Table Format (default)**:
```
NAME            VERSION  TYPE     DESCRIPTION
code-reviewer   1.2.0    agent    Automated code review assistant
git-helper      2.1.5    command  Git workflow automation
```

**JSON Format** (`--json`):
```json
{
  "tools": [
    {
      "name": "code-reviewer",
      "version": "1.2.0",
      "type": "agent",
      "description": "Automated code review assistant",
      "author": "john.doe",
      "tags": ["code-review", "automation"]
    }
  ],
  "count": 1
}
```

---

## list

List installed or available tools.

### Usage

```bash
cntm list [flags]
```

### Description

Show tools that are installed locally or available in the registry.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--remote` | `-r` | bool | List remote registry tools | false |
| `--type` | `-t` | string | Filter by type (agent, command, skill) | all |
| `--json` | | bool | Output in JSON format | false |

### Examples

```bash
# List installed tools
cntm list

# List remote tools
cntm list --remote

# List installed agents
cntm list --type agent

# List with JSON output
cntm list --json
```

### Output

Shows installed version, type, and installation path for local tools. For remote tools, shows available version and description.

---

## info

Show detailed information about a tool.

### Usage

```bash
cntm info <tool-name> [flags]
```

### Description

Display comprehensive information about a specific tool including metadata, version history, dependencies, and installation status.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--json` | | bool | Output in JSON format | false |

### Examples

```bash
# Show tool info
cntm info code-reviewer

# Show info as JSON
cntm info code-reviewer --json
```

### Output

Displays:
- Name and current version
- Type (agent/command/skill)
- Description
- Author
- Tags
- Installation status
- Available versions
- Download URL
- Size and checksum

---

## browse

Browse and discover tools in the registry.

### Usage

```bash
cntm browse [flags]
cntm trending [flags]
```

### Description

Discover tools by browsing the registry with various sorting and filtering options. The `trending` alias shows popular tools by download count.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--sort` | `-s` | string | Sort by (downloads, updated, name) | downloads |
| `--type` | `-t` | string | Filter by type (agent, command, skill) | all |
| `--limit` | `-l` | int | Limit number of results | 10 |
| `--json` | | bool | Output in JSON format | false |

### Examples

```bash
# Browse top tools by downloads
cntm browse

# Browse recently updated tools
cntm browse --sort updated

# Show top 20 trending tools
cntm trending --limit 20

# Browse agents sorted by name
cntm browse --type agent --sort name
```

---

## install

Install tools from the registry.

### Usage

```bash
cntm install <tool-name>[@version] [tool-name2] [...] [flags]
```

### Description

Download and install one or more tools from the registry. Tools are installed to `.claude/<type>/<name>/`.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--force` | `-f` | bool | Force reinstall if already installed | false |
| `--path` | | string | Custom installation path | `.claude` |

### Examples

```bash
# Install latest version
cntm install code-reviewer

# Install specific version
cntm install code-reviewer@1.0.0

# Install multiple tools
cntm install agent1 agent2 agent3

# Force reinstall
cntm install --force code-reviewer

# Install to custom path
cntm install --path /custom code-reviewer
```

### Process

1. Validates tool exists in registry
2. Downloads ZIP file with progress bar
3. Verifies integrity (SHA256)
4. Extracts to installation directory
5. Updates lock file
6. Shows success message

### Error Handling

- **Tool not found**: Suggests running `cntm search <name>`
- **Already installed**: Shows installed version, suggests `--force`
- **Network error**: Suggests checking connection
- **Integrity check failed**: Suggests trying again

---

## update

Update tools to the latest version.

### Usage

```bash
cntm update [tool-name] [flags]
cntm update --all [flags]
```

### Description

Update one or all tools to the latest version available in the registry.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--all` | | bool | Update all outdated tools | false |
| `--yes` | `-y` | bool | Skip confirmation prompts | false |

### Examples

```bash
# Update specific tool
cntm update code-reviewer

# Update all tools
cntm update --all

# Update without confirmation
cntm update --all --yes

# Update specific tool without confirmation
cntm update code-reviewer --yes
```

### Process

1. Checks for available updates
2. Shows version change (current → latest)
3. Prompts for confirmation (unless `--yes`)
4. Downloads and installs new version
5. Updates lock file
6. Shows summary

---

## outdated

Check for outdated tools.

### Usage

```bash
cntm outdated [flags]
```

### Description

Display a list of installed tools that have newer versions available in the registry.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--json` | | bool | Output in JSON format | false |

### Examples

```bash
# Check for outdated tools
cntm outdated

# JSON output
cntm outdated --json
```

### Output

**Table Format**:
```
NAME            INSTALLED  LATEST   STATUS
code-reviewer   1.0.0      1.2.0    outdated
git-helper      2.1.5      2.1.5    up-to-date
```

**JSON Format**:
```json
{
  "outdated": [
    {
      "name": "code-reviewer",
      "current_version": "1.0.0",
      "latest_version": "1.2.0"
    }
  ],
  "up_to_date": ["git-helper"]
}
```

---

## remove

Remove installed tools.

### Usage

```bash
cntm remove <tool-name> [tool-name2] [...] [flags]
```

### Aliases

- `uninstall`
- `rm`

### Description

Remove one or more tools from the local installation. Deletes the tool directory and updates the lock file.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--yes` | `-y` | bool | Skip confirmation prompts | false |

### Examples

```bash
# Remove with confirmation
cntm remove code-reviewer

# Remove multiple tools
cntm remove tool1 tool2 tool3

# Remove without confirmation
cntm remove --yes old-agent

# Using aliases
cntm uninstall code-reviewer
cntm rm code-reviewer
```

### Process

1. Validates tool is installed
2. Shows what will be removed
3. Prompts for confirmation (unless `--yes`)
4. Removes tool directory
5. Updates lock file
6. Shows summary

---

## create

Create a new tool locally.

### Usage

```bash
cntm create <type> <name> [flags]
```

### Description

Create a new tool directory with the proper structure and metadata template. Types can be `agent`, `command`, or `skill`.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--author` | `-a` | string | Tool author name | git user.name |
| `--description` | `-d` | string | Tool description | Interactive prompt |
| `--version` | `-v` | string | Initial version | 1.0.0 |

### Examples

```bash
# Create with interactive prompts
cntm create agent my-agent

# Create with flags
cntm create agent my-agent \
  --author "John Doe" \
  --description "My custom agent" \
  --version "0.1.0"

# Create command
cntm create command my-command
```

### Structure Created

```
my-agent/
├── metadata.json
├── README.md
└── agent.yaml
```

---

## publish

Publish a tool to the registry.

### Usage

```bash
cntm publish <tool-name> [flags]
```

### Description

Package and publish a tool to the GitHub registry via pull request. The tool must exist in the local directory.

### Flags

| Flag | Short | Type | Description | Default |
|------|-------|------|-------------|---------|
| `--version` | `-v` | string | Version to publish | From metadata.json |
| `--message` | `-m` | string | Commit message | Auto-generated |
| `--force` | `-f` | bool | Overwrite existing version | false |

### Examples

```bash
# Publish with metadata version
cntm publish my-agent

# Publish specific version
cntm publish my-agent --version 1.1.0

# Publish with custom message
cntm publish my-agent --message "Add new features"

# Force publish (overwrite)
cntm publish my-agent --force
```

### Process

1. Validates tool structure and metadata
2. Creates ZIP archive
3. Calculates SHA256 checksum
4. Creates branch in registry repo
5. Uploads tool ZIP and metadata
6. Updates registry.json
7. Creates pull request
8. Returns PR URL

### Requirements

- GitHub token in config or `CNTM_GITHUB_TOKEN` env var
- Proper tool structure with metadata.json
- Valid semantic version

---

## Exit Codes

All commands use standard exit codes:

- `0`: Success
- `1`: Error (with descriptive message)

---

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `CNTM_GITHUB_TOKEN` | GitHub personal access token | `ghp_xxx...` |
| `CNTM_CONFIG` | Config file path | `/path/to/config.yaml` |
| `CNTM_REGISTRY_URL` | Registry repository URL | `https://github.com/user/registry` |

---

## Tips and Best Practices

### Performance

- Use `--json` for programmatic access
- Registry data is cached for 1 hour
- Use `--yes` in CI/CD pipelines

### Workflows

```bash
# Initial setup
cntm init
cntm install code-reviewer git-helper

# Keep tools updated
cntm outdated
cntm update --all

# Clean up
cntm remove old-tool
```

### Troubleshooting

See [TROUBLESHOOTING.md](TROUBLESHOOTING.md) for common issues and solutions.
