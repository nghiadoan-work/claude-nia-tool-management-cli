---
description: Build the cntm CLI binary and verify all commands work correctly
allowed-tools: Bash
---

Build and test the cntm CLI.

## Current state
- Existing binary: !`ls -lh cntm 2>/dev/null || echo "No binary found"`
- Build status: !`go build -o cntm . && echo "Build successful" || echo "Build failed"`

## Your tasks

1. **Build the CLI**: `go build -o cntm -ldflags="-s -w" .`
2. **Verify commands**: Test each command's help text works
3. **Check binary size**: Report the size (should be <20MB)
4. **Run lint** (if golangci-lint available): `golangci-lint run`
5. **Integration tests** (if they exist): `go test -tags=integration ./tests/integration/...`

## Report

Provide:
- Build status (success/failure)
- Binary size
- Which commands are implemented vs. planned
- Any lint warnings or errors
- Integration test results (if applicable)
- Overall status: Ready for testing / In development / Has issues

List any problems found and suggest fixes.
