## Claude Tools for CNTM Development

Specialized tools for building the claude-nia-tool-management-cli (cntm) package manager.

## What's Included

```
.claude/
├── agents/
│   └── cntm-developer.md         # AI pair programmer for cntm
├── commands/
│   ├── test-coverage.md          # Run tests with coverage analysis
│   ├── build-test.md             # Build and verify the CLI
│   ├── review-architecture.md    # Check architecture compliance
│   └── new-feature.md            # TDD implementation guide
└── skills/
    └── github-api/
        └── SKILL.md              # GitHub API reference patterns
```

## Quick Start

### 1. Use the Subagent

The **cntm-developer** subagent automatically activates when you're working on cntm. It knows:
- Complete project architecture (from docs/)
- Go best practices and CLI patterns
- Security requirements (ZIP bombs, path traversal, tokens)
- Testing strategies (TDD, table-driven tests, mocking)

Just work on cntm code and it will provide expert guidance!

### 2. Run Slash Commands

Execute commands in Claude Code:

```
/test-coverage              # Analyze test coverage
/build-test                 # Build and test cntm binary
/review-architecture        # Check code compliance
/new-feature InstallerService  # Guide for implementing features
```

### 3. Reference Skills

The **github-api** skill provides code patterns and best practices. Claude will reference it automatically when needed, or you can mention it explicitly.

## Commands Reference

### /test-coverage
Runs comprehensive test analysis:
- Executes tests with race detection
- Generates coverage reports
- Identifies areas needing tests (<80%)
- Recommends security-critical code to test

**When to use**: Before committing, during code review

### /build-test
Builds and verifies the CLI:
- Compiles cntm binary
- Tests all commands
- Runs linter (if available)
- Checks integration tests
- Reports binary size and status

**When to use**: After implementing features, before releases

### /review-architecture
Checks architecture compliance:
- Verifies layer separation (cmd → services → data)
- Validates service patterns (interfaces, DI)
- Security audit (path validation, ZIP bombs, tokens)
- Testing coverage review

**When to use**: During development, code reviews

### /new-feature [name]
Step-by-step TDD implementation guide:
1. Design interface
2. Write tests FIRST
3. Implement service
4. Create CLI command
5. Verify and test

**When to use**: Starting any new feature

**Example**: `/new-feature InstallerService`

## Skills Reference

### github-api
Provides patterns for GitHub integration:
- Authentication setup
- Rate limit handling
- Common operations (download, fork, PR)
- Error handling with CLIError
- Retry logic with exponential backoff
- Testing with mocks

**When referenced**: Implementing GitHubClient, debugging API issues

## Development Workflow

### Starting a Feature

1. **Check roadmap**: Review `docs/ROADMAP.md`
2. **Run guide**: `/new-feature FeatureName`
3. **Implement with TDD**: Tests first, then code
4. **Verify**: Use `/test-coverage` and `/review-architecture`

### Before Committing

```
/test-coverage          # Ensure >80% coverage
/review-architecture    # Check compliance
/build-test            # Verify everything works
```

### Getting Help

Just ask! The cntm-developer subagent will:
- Implement services following project patterns
- Write tests with proper mocking
- Review code for security issues
- Explain architecture decisions
- Provide code examples

**Example interactions**:
```
"Implement RegistryService following our architecture"
"Review this code for security issues"
"How should I handle GitHub rate limiting?"
"Write tests for InstallerService"
```

## Best Practices

### Commands
- Use for **routine checks** and **step-by-step guides**
- Commands provide structure and checklists
- They tell Claude what to analyze and report

### Subagent
- Use for **implementation** and **code review**
- Automatically active when working on cntm
- Provides expert Go and CLI development guidance
- Understands full project context

### Skills
- Use for **reference patterns**
- Claude auto-activates when relevant
- Provides code examples and best practices

## File Locations

**Project-level** (checked into git):
- `.claude/agents/` - Subagents for team
- `.claude/commands/` - Slash commands for team
- `.claude/skills/` - Shared skills for team

**Personal** (your machine only):
- `~/.claude/agents/` - Your personal subagents
- `~/.claude/commands/` - Your personal commands
- `~/.claude/skills/` - Your personal skills

## Tips

### Effective Subagent Usage

✅ **Do**:
- Let it activate automatically
- Ask for complete implementations
- Request code reviews
- Get security guidance

❌ **Don't**:
- Micromanage individual lines
- Ask for quick fixes without context

### Effective Command Usage

✅ **Do**:
- Run before commits
- Use for routine checks
- Follow the guides

❌ **Don't**:
- Expect commands to write code
- Skip verification steps

## Quality Standards

All code must meet:
- ✓ **Architecture**: Follows docs/ARCHITECTURE.md
- ✓ **Testing**: >80% coverage, table-driven
- ✓ **Security**: Path validation, integrity checks, no token leaks
- ✓ **UX**: Progress bars, clear errors, helpful hints
- ✓ **Quality**: Formatted, linted, documented

## Resources

- **Project Docs**: `docs/` directory
  - REQUIREMENTS.md - What to build
  - ARCHITECTURE.md - How to build it
  - ROADMAP.md - When to build it

- **Claude Code Docs**: https://code.claude.com/docs

## Summary

These tools work together:

```
cntm-developer (Subagent)
    ↓ provides expert implementation guidance
/new-feature (Command)
    ↓ guides TDD approach step-by-step
github-api (Skill)
    ↓ provides reference patterns
/test-coverage (Command)
    ↓ verifies quality
/review-architecture (Command)
    ↓ ensures compliance
/build-test (Command)
    ↓ confirms it works

Result: High-quality, secure, well-tested code! 🚀
```

Focus on building features - let these tools ensure you're following best practices!
