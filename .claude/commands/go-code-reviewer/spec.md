---
name: go-code-reviewer-spec
description: Specify and scope Go code review tasks into a detailed checklist
---

# Go Code Reviewer - Task Specification

## Usage

Create a detailed task specification and checklist for Go code review tasks.

## Command Behavior

When invoked, this command will:

1. **Create working directory** - Initialize `go-code-reviewer-progress/` folder (if not exists)
2. **Analyze the request** - Understand the code review requirements (files, packages, or entire project)
3. **Break down into steps** - Create a logical sequence of review subtasks
4. **Create checklist** - Generate a markdown checklist file
5. **Save to progress folder** - Store in `go-code-reviewer-progress/{task-name}.md`

## Checklist Format

Each code review task file will contain:

```markdown
# Task: {Task Name}

**Created**: {date}
**Status**: Not Started
**Priority**: {High/Medium/Low}
**Scope**: {files/package/project}

## Overview
{Brief description of what needs to be reviewed}

## Files/Packages to Review
- {file/package 1}
- {file/package 2}
- ...

## Checklist

### Code Analysis
- [ ] Run `go vet` on target code
- [ ] Run `golangci-lint run` for comprehensive linting
- [ ] Check for race conditions with `go build -race`
- [ ] Verify `go mod tidy` - check for unused dependencies

### Type Safety & Compilation
- [ ] Run `go build` to ensure code compiles
- [ ] Check type assertions and conversions
- [ ] Verify interface implementations
- [ ] Review error handling patterns

### Code Quality Review
- [ ] Review function complexity and readability
- [ ] Check for proper error handling (no ignored errors)
- [ ] Verify naming conventions (Go standards)
- [ ] Review documentation and comments
- [ ] Check for code duplication
- [ ] Verify proper use of goroutines and channels
- [ ] Review resource cleanup (defer, Close())

### Testing & Coverage
- [ ] Run `go test ./...` - all tests pass
- [ ] Check test coverage with `go test -cover ./...`
- [ ] Review test quality and completeness
- [ ] Check for missing edge case tests

### Performance & Best Practices
- [ ] Review memory allocations and efficiency
- [ ] Check for potential memory leaks
- [ ] Verify proper use of pointers vs values
- [ ] Review concurrent code for deadlocks
- [ ] Check for proper context usage

### Security Review
- [ ] Check for SQL injection vulnerabilities
- [ ] Review input validation
- [ ] Check for hardcoded credentials
- [ ] Review error messages (no sensitive info leaks)

## Notes
{Any additional context, concerns, or focus areas}

## Review Summary
{To be filled during apply phase}
- **Issues Found**:
- **Critical**:
- **Warnings**:
- **Suggestions**:
```

## Example Usage

```bash
# User request
"Spec: Review authentication package for security and best practices"

# Command generates
go-code-reviewer-progress/
├── review-auth-package.md
└── archived/  (created if doesn't exist)
```

## Review Scope Options

When creating a spec, determine the appropriate scope:

1. **Single File Review**: Focus on specific file(s)
   - Quick targeted reviews
   - Bug fix validation
   - Specific feature review

2. **Package Review**: Entire package analysis
   - Module-level architecture review
   - API consistency check
   - Package-level best practices

3. **Project Review**: Full codebase analysis
   - Overall architecture review
   - Cross-package consistency
   - Project-wide standards compliance

## Implementation Guidelines

1. **Be Specific**: Clearly identify what code needs review (files, packages, or project scope)
2. **Logical Order**: Organize checks from compilation → linting → quality → security
3. **Reasonable Scope**: Don't try to review entire large projects in one task
4. **Include Context**: Note any specific concerns, recent changes, or focus areas
5. **Set Priority**: Mark urgent reviews (security, production bugs) as High priority

## Checklist Customization

Adapt the checklist based on review type:

- **Security Review**: Add more security-focused checks
- **Performance Review**: Focus on benchmarks and profiling
- **Refactoring Review**: Emphasize design patterns and architecture
- **Bug Fix Review**: Focus on the specific bug area and regression tests
- **Feature Review**: Check feature completeness and integration

## Tools Used

The review process leverages these Go tools:
- `go build` - Compilation check
- `go vet` - Static analysis
- `go test` - Testing
- `golangci-lint` - Comprehensive linting
- `go mod` - Dependency management
- `go build -race` - Race condition detection

Ensure these tools are installed before starting reviews.
