---
description: Guide me through implementing a new feature in cntm using TDD and following architecture patterns
argument-hint: [feature-name]
allowed-tools: Read, Write, Edit, Bash
---

Guide me through implementing: **$ARGUMENTS**

## Architecture reference
First, review the architecture: @docs/ARCHITECTURE.md and roadmap: @docs/ROADMAP.md

## TDD Implementation Steps

### Step 1: Design Interface
Help me design the interface for $ARGUMENTS. Consider:
- What methods are needed?
- What parameters and return values?
- How does it fit into the existing architecture?

### Step 2: Write Tests FIRST
Guide me to write table-driven tests before implementation:
- Test success cases
- Test error cases
- Mock external dependencies
- Aim for >80% coverage

### Step 3: Implement Service
Help me implement the service following cntm patterns:
- Dependency injection via constructor
- Context support for cancellation
- Error wrapping with fmt.Errorf and %w
- CLIError for user-facing errors with hints
- Security checks where applicable

### Step 4: Create CLI Command
Guide me to create the Cobra command:
- Clear usage and examples
- Appropriate flags
- Progress indication for long operations
- Good error messages

### Step 5: Verify
Help me verify:
- Tests pass: `go test ./...`
- Coverage: `go test -cover ./...`
- No race conditions: `go test -race ./...`
- Code formatted: `go fmt ./...`
- Lint clean: `golangci-lint run`

## Code patterns to follow

Show me examples using cntm patterns:
- Service implementation with DI
- Table-driven test structure
- Error handling with CLIError and hints
- CLI command with progress bars
- Context cancellation support

## Security checklist
For features involving file operations or external data:
- Path validation (prevent traversal)
- Size limits (prevent ZIP bombs)
- Integrity verification (SHA256)
- No token/secret logging

Guide me through each step, providing code examples and checking my implementation.
