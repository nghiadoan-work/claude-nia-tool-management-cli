---
name: go-code-reviewer-apply
description: Execute Go code review tasks from the checklist and check off completed items
---

# Go Code Reviewer - Task Execution

## Usage

Work through the code review task checklist and mark items as completed.

## Command Behavior

When invoked, this command will:

1. **Read checklist** - Load the task file from `go-code-reviewer-progress/{task-name}.md`
2. **Execute review tasks** - Work through unchecked items sequentially
3. **Run analysis tools** - Execute go vet, golangci-lint, tests, etc.
4. **Update progress** - Mark completed items with `[x]`
5. **Document findings** - Record issues, warnings, and suggestions
6. **Update status** - Change status field (Not Started → In Progress → Completed)
7. **Save changes** - Write updates back to the task file

## Workflow

### Initial State
```markdown
**Status**: Not Started

## Checklist

### Code Analysis
- [ ] Run `go vet` on target code
- [ ] Run `golangci-lint run` for comprehensive linting
- [ ] Check for race conditions with `go build -race`
```

### During Execution
```markdown
**Status**: In Progress

## Checklist

### Code Analysis
- [x] Run `go vet` on target code
- [x] Run `golangci-lint run` for comprehensive linting ← Currently working on
- [ ] Check for race conditions with `go build -race`

## Review Summary
- **Issues Found**: 3 warnings from golangci-lint
- **Critical**: None
- **Warnings**:
  - Line 45: exported function missing documentation
  - Line 67: error return value not checked
  - Line 89: variable name too short (use descriptive names)
```

### Completion
```markdown
**Status**: Completed
**Completed**: 2025-11-16

## Checklist

### Code Analysis
- [x] Run `go vet` on target code
- [x] Run `golangci-lint run` for comprehensive linting
- [x] Check for race conditions with `go build -race`
- [x] Verify `go mod tidy` - check for unused dependencies

### Type Safety & Compilation
- [x] Run `go build` to ensure code compiles
- [x] Check type assertions and conversions
- [x] Verify interface implementations
- [x] Review error handling patterns

[All items checked...]

## Review Summary
- **Issues Found**: 12 total
- **Critical**: 2
  - Security: Potential SQL injection in query builder (line 156)
  - Memory leak: Missing Close() call in HTTP client (line 234)
- **Warnings**: 5
  - Missing error checks (3 occurrences)
  - Exported functions missing documentation (2 occurrences)
- **Suggestions**: 5
  - Consider using sync.Pool for frequent allocations (line 45)
  - Refactor parseConfig function - too complex (line 123)
  - Add context timeout to HTTP requests (line 201)
  - Use constants for magic numbers (lines 78, 92, 145)
  - Consider breaking up large function (getUserData - 150 lines)
```

## Status Updates

Update the status field as work progresses:

- **Not Started** - No items completed
- **In Progress** - Some items completed, review ongoing
- **Completed** - All items checked off, review summary documented
- **Blocked** - Cannot proceed (add reason in notes, e.g., missing access to code, dependency issues)

## Execution Process

### 1. Code Analysis Phase

Run automated tools and document results:

```bash
# Navigate to target code
cd /path/to/code

# Run go vet
go vet ./...

# Run golangci-lint (comprehensive)
golangci-lint run

# Check for race conditions
go build -race ./...

# Verify dependencies
go mod tidy
go mod verify
```

Document findings for each tool.

### 2. Type Safety & Compilation Phase

```bash
# Compile the code
go build ./...

# Run tests
go test ./...

# Check test coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Review type safety issues and test results.

### 3. Code Quality Review Phase

Manual review of code quality:

- Read through code files
- Check naming conventions
- Review function complexity
- Verify error handling
- Look for code smells
- Check documentation

Use tools like:
```bash
# Cyclomatic complexity
gocyclo -over 15 .

# Code duplication
dupl -threshold 50 .
```

### 4. Security Review Phase

Focus on security concerns:

- Input validation
- SQL queries (parameterization)
- Authentication/authorization
- Secrets management
- Error messages (information disclosure)

### 5. Document Findings

Add all findings to the "Review Summary" section:

```markdown
## Review Summary
- **Issues Found**: {count}
- **Critical**: {count}
  - {description} (file:line)
  - {description} (file:line)
- **Warnings**: {count}
  - {description} (file:line)
- **Suggestions**: {count}
  - {description} (file:line)
```

## Example Usage

```bash
# User request
"Apply: Work on review-auth-package task"

# Command:
# 1. Opens go-code-reviewer-progress/review-auth-package.md
# 2. Shows current progress (0/25 items completed)
# 3. Starts with first unchecked item: "Run go vet"
# 4. Executes: cd internal/auth && go vet ./...
# 5. Documents results in review summary
# 6. Marks as [x] when done
# 7. Moves to next item
# 8. Continues until all complete or user stops
```

## Implementation Guidelines

1. **Show Progress**: Display current completion percentage and which phase you're in
2. **One at a Time**: Focus on one checklist item at a time
3. **Run Commands**: Actually execute the go tools and analyze output
4. **Document Everything**: Record all findings with file and line numbers
5. **Verify Completion**: Ensure each step is fully done before checking off
6. **Add Notes**: Document any issues, decisions, or follow-up needed
7. **Update Timestamp**: Add "Last Updated" and "Completed" dates

## Handling Issues

### When Critical Issues Found

- Mark them clearly in Review Summary
- Consider stopping review to address critical issues first
- Update priority if needed
- Document remediation needed

### When Blocked

- Update status to "Blocked"
- Document blocking issue in Notes section
- Specify what's needed to unblock
- Consider partial completion and archival

## Output Format

Maintain consistent formatting for findings:

```markdown
- {Severity}: {Brief description} ({file}:{line})
  - Details: {additional context}
  - Recommendation: {how to fix}
```

Example:
```markdown
- **Critical**: SQL injection vulnerability (pkg/database/query.go:156)
  - Details: User input directly concatenated into SQL query
  - Recommendation: Use parameterized queries with $1, $2 placeholders

- **Warning**: Missing error check (internal/auth/token.go:67)
  - Details: Return value of jwt.Parse() not checked
  - Recommendation: Add error handling for token parsing failure

- **Suggestion**: Consider using constant (pkg/config/settings.go:45)
  - Details: Magic number 3600 used for timeout
  - Recommendation: Define const DefaultTimeout = 3600 * time.Second
```

## Review Quality Guidelines

Ensure high-quality reviews:

1. **Be Thorough**: Don't skip steps, even if code looks clean
2. **Be Objective**: Focus on code quality, not personal preferences
3. **Be Constructive**: Provide actionable recommendations
4. **Be Specific**: Include file names, line numbers, and exact issues
5. **Prioritize**: Distinguish critical bugs from style suggestions
6. **Consider Context**: Understand the code's purpose and constraints

## Tools and Commands Reference

Quick reference for common Go review commands:

```bash
# Static analysis
go vet ./...
golangci-lint run
staticcheck ./...

# Testing
go test ./...
go test -race ./...
go test -cover ./...
go test -bench=. ./...

# Build checks
go build ./...
go build -race ./...

# Dependencies
go mod tidy
go mod verify
go mod graph

# Code metrics
gocyclo -over 15 .
gofmt -d .
goimports -d .

# Security scanning
gosec ./...
```
