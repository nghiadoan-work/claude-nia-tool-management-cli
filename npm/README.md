# cntm - Claude Nia Tool Management CLI

[![npm version](https://badge.fury.io/js/cntm.svg)](https://www.npmjs.com/package/cntm)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A package manager for Claude Code tools (agents, commands, and skills). Like npm for Node.js, cntm helps you install, update, and publish Claude Code tools.

## Installation

### Using npx (Recommended)

Run cntm without installation:

```bash
npx cntm search "code review"
npx cntm install code-reviewer
```

### Global Installation

Install globally to use as a command:

```bash
npm install -g cntm
cntm --help
```

### Local Installation

Install in your project:

```bash
npm install cntm
npx cntm init
```

### Direct from GitHub

```bash
npx github:nghiadt/claude-nia-tool-management-cli
```

## Quick Start

```bash
# Initialize a new project
cntm init

# Search for tools
cntm search "code review"

# Install a tool
cntm install code-reviewer

# List installed tools
cntm list

# Update tools
cntm update --all

# Browse trending tools
cntm trending
```

## All Commands

| Command | Description |
|---------|-------------|
| `cntm init` | Initialize a new Claude tools project |
| `cntm search <query>` | Search for tools in the registry |
| `cntm browse` | Browse and explore tools |
| `cntm trending` | Show trending tools |
| `cntm info <name>` | Display detailed tool information |
| `cntm list` | List installed tools |
| `cntm install <name>` | Install a tool |
| `cntm outdated` | Check for outdated tools |
| `cntm update [name]` | Update tools |
| `cntm remove <name>` | Remove a tool |
| `cntm create <type> <name>` | Create a new tool |
| `cntm publish <name>` | Publish a tool to the registry |
| `cntm version` | Show version information |

## Programmatic Usage

You can also use cntm programmatically in Node.js:

```javascript
const cntm = require('cntm');

// Get version
const version = await cntm.version();
console.log(`cntm version: ${version}`);

// Search for tools
const results = await cntm.search('code review');
console.log(results);

// Install a tool
await cntm.install('code-reviewer');

// List installed tools
const tools = await cntm.list();
console.log(tools);

// Execute custom commands
const result = await cntm.execute(['outdated', '--json']);
console.log(result.stdout);
```

## Features

- 🔍 **Search & Discovery** - Find tools easily
- 📦 **Install & Update** - Manage tools with version control
- 🚀 **Publish** - Share your tools with others
- 🎨 **Beautiful UI** - Color-coded output with spinners
- 🔒 **Secure** - SHA256 verification and path validation
- 🌐 **Multi-platform** - Works on macOS, Linux, and Windows
- 📚 **Well documented** - Comprehensive guides and examples

## Configuration

Create a `.cntm-config.yaml` file:

```yaml
registry:
  url: https://github.com/nghiadoan-work/claude-tools-registry
  branch: main
  auth_token: ${GITHUB_TOKEN}  # Optional

local:
  default_path: .claude
  cache_ttl: 3600
```

## Environment Variables

- `CNTM_REGISTRY_URL` - Override registry URL
- `CNTM_REGISTRY_BRANCH` - Override registry branch
- `CNTM_AUTH_TOKEN` - GitHub personal access token
- `CNTM_INSTALL_PATH` - Override installation path

## Documentation

- [Command Reference](https://github.com/nghiadt/claude-nia-tool-management-cli/blob/main/docs/COMMANDS.md)
- [Configuration Guide](https://github.com/nghiadt/claude-nia-tool-management-cli/blob/main/docs/CONFIGURATION.md)
- [Publishing Guide](https://github.com/nghiadt/claude-nia-tool-management-cli/blob/main/docs/PUBLISHING.md)
- [Troubleshooting](https://github.com/nghiadt/claude-nia-tool-management-cli/blob/main/docs/TROUBLESHOOTING.md)

## Platform Support

This npm package downloads the appropriate binary for your platform:

- macOS (Intel and Apple Silicon)
- Linux (x64 and ARM64)
- Windows (x64)

## Alternative Installation Methods

### Homebrew (macOS/Linux)

```bash
# Coming soon
brew install cntm
```

### Download Binary

Download pre-built binaries from [GitHub Releases](https://github.com/nghiadt/claude-nia-tool-management-cli/releases).

### Build from Source

```bash
git clone https://github.com/nghiadt/claude-nia-tool-management-cli.git
cd claude-nia-tool-management-cli
go build -o cntm
./cntm version
```

## Contributing

Contributions are welcome! Please read our [Contributing Guide](https://github.com/nghiadt/claude-nia-tool-management-cli/blob/main/CONTRIBUTING.md).

## License

MIT License - see [LICENSE](https://github.com/nghiadt/claude-nia-tool-management-cli/blob/main/LICENSE) file for details.

## Links

- [GitHub Repository](https://github.com/nghiadt/claude-nia-tool-management-cli)
- [Issue Tracker](https://github.com/nghiadt/claude-nia-tool-management-cli/issues)
- [Changelog](https://github.com/nghiadt/claude-nia-tool-management-cli/blob/main/CHANGELOG.md)
- [Registry](https://github.com/nghiadoan-work/claude-tools-registry)

## Support

- 📖 [Documentation](https://github.com/nghiadt/claude-nia-tool-management-cli/tree/main/docs)
- 🐛 [Report Bug](https://github.com/nghiadt/claude-nia-tool-management-cli/issues)
- 💡 [Request Feature](https://github.com/nghiadt/claude-nia-tool-management-cli/issues)

---

Made with ❤️ for the Claude Code community
