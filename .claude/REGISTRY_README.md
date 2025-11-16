# Claude Tools Registry

This registry contains Claude Code tools (agents, commands, and skills) that can be installed using the `cntm` CLI tool.

## Structure

```
.claude/
├── registry.json          # Registry index with latest version of each tool
└── tools/                 # Tool packages
    ├── agents/
    │   └── {agent-name}/
    │       ├── v1-0-0.zip
    │       ├── v1-1-0.zip
    │       └── v2-0-0.zip
    ├── commands/
    │   └── {command-name}/
    │       └── v1-0-0.zip
    └── skills/
        └── {skill-name}/
            ├── v1-0-0.zip
            └── v1-1-0.zip
```

## Versioning System

### ZIP File Naming

Tool versions are stored as versioned ZIP files:
- Version `1.0.0` → `v1-0-0.zip`
- Version `2.1.3` → `v2-1-3.zip`
- Version `1.0.0-beta` → `v1-0-0-beta.zip`

### Directory Organization

Each tool has its own directory containing all versions:
```
tools/skills/github-api/
├── v1-0-0.zip    # Version 1.0.0
├── v1-1-0.zip    # Version 1.1.0
└── v2-0-0.zip    # Version 2.0.0
```

## Registry Index (`registry.json`)

The registry index tracks the **latest version** of each tool:

```json
{
  "version": "1.0",
  "updated_at": "2025-11-16T19:00:00+07:00",
  "tools": {
    "skill": [
      {
        "name": "github-api",
        "version": "1.0.0",
        "description": "GitHub API patterns and best practices",
        "type": "skill",
        "author": "Claude Code",
        "tags": ["github", "api", "development"],
        "file": "tools/skills/github-api/v1-0-0.zip",
        "size": 3167,
        "downloads": 0,
        "created_at": "2025-11-16T02:55:25+07:00",
        "updated_at": "2025-11-16T19:00:00+07:00"
      }
    ]
  }
}
```

## Available Tools

### Skills

#### github-api
- **Version**: 1.0.0
- **Description**: GitHub API patterns and best practices for cntm development
- **Tags**: github, api, development
- **File**: `tools/skills/github-api/v1-0-0.zip`
- **Size**: 3.1 KB

#### go-code-reviewer
- **Version**: 1.0.0
- **Description**: Expert Go code review
- **Tags**: go, review
- **File**: `tools/skills/go-code-reviewer/v1-0-0.zip`
- **Size**: 1.1 KB

## Installing Tools

### Install Latest Version
```bash
cntm install github-api
```

### Install Specific Version
```bash
cntm install github-api@1.0.0
```

## Publishing Tools

When publishing a new version:

1. **Version is specified**: `1.1.0`
2. **ZIP is created**: `tools/skills/{name}/v1-1-0.zip`
3. **Registry is updated**: Points to new version
4. **Old versions remain**: Previous ZIPs stay available

### Example Publishing Workflow

```bash
# Publish version 1.0.0
cntm publish skill my-tool
# Creates: tools/skills/my-tool/v1-0-0.zip

# Later, publish version 1.1.0
cntm publish skill my-tool
# Creates: tools/skills/my-tool/v1-1-0.zip
# Updates: registry.json to point to v1-1-0
# Preserves: v1-0-0.zip still exists
```

## Benefits

✅ **Version History**: All versions preserved in repository
✅ **Easy Rollback**: Install any previous version
✅ **Clean Organization**: Tools grouped by name, then version
✅ **Single Source of Truth**: registry.json tracks latest
✅ **Backward Compatible**: Existing install commands work

## Registry URL

This registry is hosted at:
```
https://github.com/nghiadoan-work/claude-tools-registry
```

## Contributing

To publish a tool to this registry:

1. Create your tool in `.claude/{type}s/{name}/`
2. Add metadata.json with version, description, etc.
3. Run `cntm publish {type} {name}`
4. A PR will be created automatically (if configured)

## License

Tools in this registry are provided by their respective authors. Check individual tool licenses for details.
