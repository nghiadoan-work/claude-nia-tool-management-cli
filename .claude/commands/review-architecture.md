---
description: Review codebase for compliance with cntm architecture patterns from docs/ARCHITECTURE.md
allowed-tools: Read, Grep, Glob, Bash
---

Review the cntm codebase for architecture compliance.

## Reference
Read @docs/ARCHITECTURE.md to understand the expected architecture.

## Project structure
Current structure: !`find . -type d -name "cmd" -o -name "internal" -o -name "pkg" | head -20`

## Review checklist

### 1. Layer Separation
- Check if `cmd/` contains only CLI code (no business logic)
- Verify business logic is in `internal/services/`
- Confirm data access is in `internal/data/`
- Check models are in `pkg/models/`

Search for violations:
- Services in cmd: !`grep -r "type.*Service struct" cmd/ 2>/dev/null || echo "None found"`
- Direct file I/O in cmd: !`grep -r "os\\.Open\\|ioutil" cmd/ 2>/dev/null || echo "None found"`

### 2. Service Patterns
- Verify services define interfaces
- Check dependency injection (constructor pattern)
- Confirm error handling uses CLIError type

### 3. Security Checks
- Path validation (prevent directory traversal)
- ZIP bomb protection (size/count limits)
- SHA256 integrity verification
- No token logging

Search for security patterns:
- Path validation: !`grep -r "filepath\\.Clean\\|HasPrefix" internal/ 2>/dev/null | head -5`
- Hash verification: !`grep -r "sha256\\|SHA256" internal/ 2>/dev/null | head -5`

### 4. Testing
- Count test files: !`find . -name "*_test.go" | wc -l`
- Quick coverage: !`go test -cover ./... 2>/dev/null | grep coverage`

## Report

Provide assessment:
- **Layer Separation**: PASS/FAIL with violations found
- **Service Patterns**: PASS/NEEDS IMPROVEMENT with details
- **Security**: PASS/CRITICAL ISSUES with specifics
- **Testing**: Coverage percentage and missing tests
- **Overall**: GOOD/NEEDS WORK/CRITICAL ISSUES

List priority actions to fix issues found.
