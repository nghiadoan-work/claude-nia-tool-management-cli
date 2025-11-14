# Setup Guide

## Prerequisites

### Install Go

This project requires Go 1.21 or later.

#### macOS
```bash
# Using Homebrew
brew install go

# Or download from https://go.dev/dl/
```

#### Linux
```bash
# Using package manager (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# Or download from https://go.dev/dl/
```

#### Windows
Download and install from https://go.dev/dl/

#### Verify Installation
```bash
go version
# Should output: go version go1.21.x ...
```

## Project Setup

### 1. Initialize Go Module

```bash
cd agent_skill_cli_go
go mod init github.com/yourusername/claude-tools-cli
```

Replace `yourusername` with your actual GitHub username or organization.

### 2. Install Dependencies

```bash
# CLI framework
go get github.com/spf13/cobra@latest

# Interactive UI
go get github.com/charmbracelet/bubbletea@latest
# Alternative: go get github.com/AlecAivazis/survey/v2@latest

# YAML support
go get gopkg.in/yaml.v3

# UUID generation
go get github.com/google/uuid

# Testing utilities
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/mock

# Optional: File watching (for future features)
# go get github.com/fsnotify/fsnotify
```

### 3. Create Project Structure

```bash
# Create directory structure
mkdir -p cmd
mkdir -p internal/{claude,templates,versioning,search,bundle,repository,config,ui}
mkdir -p pkg/models
mkdir -p templates/{agents,commands,skills}
mkdir -p tests/{integration,fixtures}
```

### 4. Initialize Cobra CLI

```bash
# Initialize cobra (this creates cmd/root.go)
cobra-cli init

# Add subcommands
cobra-cli add agent
cobra-cli add command
cobra-cli add skill
cobra-cli add template
cobra-cli add search
cobra-cli add import
cobra-cli add export
cobra-cli add version
cobra-cli add interactive
```

**Note**: If you don't have cobra-cli installed:
```bash
go install github.com/spf13/cobra-cli@latest
```

### 5. Create main.go

```bash
cat > main.go << 'EOF'
package main

import "github.com/yourusername/claude-tools-cli/cmd"

func main() {
    cmd.Execute()
}
EOF
```

### 6. Build and Test

```bash
# Build the project
go build -o tool

# Test the CLI
./tool --help

# Run tests
go test ./...
```

## Project Structure After Setup

```
agent_skill_cli_go/
├── cmd/
│   ├── root.go
│   ├── agent.go
│   ├── command.go
│   ├── skill.go
│   ├── template.go
│   ├── search.go
│   ├── import.go
│   ├── export.go
│   ├── version.go
│   └── interactive.go
├── internal/
│   ├── claude/
│   ├── templates/
│   ├── versioning/
│   ├── search/
│   ├── bundle/
│   ├── repository/
│   ├── config/
│   └── ui/
├── pkg/
│   └── models/
│       └── models.go
├── templates/
│   ├── agents/
│   ├── commands/
│   └── skills/
├── tests/
│   ├── integration/
│   └── fixtures/
├── main.go
├── go.mod
├── go.sum
├── README.md
├── REQUIREMENTS.md
├── ARCHITECTURE.md
├── ROADMAP.md
└── SETUP.md
```

## Development Workflow

### 1. Start with Phase 1 (Core CRUD)

Follow the roadmap in ROADMAP.md, starting with:
- Implement models in `pkg/models/`
- Implement repository in `internal/repository/`
- Implement agent service in `internal/claude/`
- Wire up CLI commands in `cmd/agent.go`

### 2. Run Tests Frequently

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/repository

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 3. Build and Test Locally

```bash
# Build
go build -o tool

# Test commands
./tool agent list
./tool agent add test-agent
./tool agent get test-agent
```

### 4. Use Air for Hot Reload (Optional)

Install air for automatic rebuilding during development:

```bash
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

Create `.air.toml`:
```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/tool ."
bin = "tmp/tool"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]
```

## IDE Setup

### VS Code

Recommended extensions:
- Go (by Go Team at Google)
- Go Test Explorer
- GoDoc

Create `.vscode/settings.json`:
```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.formatTool": "goimports",
  "editor.formatOnSave": true
}
```

### GoLand

GoLand has built-in Go support. Just open the project directory.

## Code Quality Tools

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

### Formatting

```bash
# Format all files
go fmt ./...

# Or use goimports (better)
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

### Static Analysis

```bash
go vet ./...
```

## Debugging

### Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug the CLI
dlv debug -- agent list

# Debug tests
dlv test ./internal/repository
```

## Next Steps

1. Install Go if not already installed
2. Run the setup commands above
3. Start implementing Phase 1 from ROADMAP.md
4. Reference ARCHITECTURE.md for design guidance
5. Reference REQUIREMENTS.md for feature specifications

## Troubleshooting

### "go: command not found"
- Go is not installed or not in PATH
- Install Go and add to PATH in your shell profile

### "cannot find package"
- Run `go mod tidy` to download dependencies
- Ensure go.mod exists with correct module path

### "permission denied"
- Make the binary executable: `chmod +x tool`

### Import path issues
- Update import paths in all files to match your module name
- Run `go mod tidy` after fixing imports

## Additional Resources

- [Go Documentation](https://go.dev/doc/)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Bubble Tea Documentation](https://github.com/charmbracelet/bubbletea)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
