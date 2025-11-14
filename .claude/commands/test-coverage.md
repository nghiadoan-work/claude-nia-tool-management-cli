---
description: Run comprehensive test coverage analysis for the cntm project
allowed-tools: Bash
---

Run comprehensive test coverage for the cntm Go project.

## Current status
- Test results: !`go test ./... -v`
- Coverage: !`go test -cover ./...`

## Your tasks

1. Run tests with full coverage report: `go test -race -coverprofile=coverage.out -covermode=atomic ./...`
2. Generate HTML coverage: `go tool cover -html=coverage.out -o coverage.html`
3. Analyze coverage by package: `go tool cover -func=coverage.out`

## Report

Provide summary:
- Total coverage percentage (target: 80%+)
- Packages with low coverage (<80%)
- Specific recommendations for improving test coverage
- Critical untested code paths (especially security-related)

Focus on: error handling, security checks (ZIP validation, path traversal), and GitHub API interactions.
