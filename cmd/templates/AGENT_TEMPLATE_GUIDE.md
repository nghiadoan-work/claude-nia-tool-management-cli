# Agent Template Guide

Guide for creating the `{agent-name}.md` file for Claude Code agents.

## Template Structure

```markdown
---
name: {agent-name}
description: {What the agent does and when to use it}
tools: Read, Write, Edit, Bash, Grep, Glob, WebFetch
model: inherit
---

# {Agent Name}

## Purpose
{1-2 sentence description of what this agent does}

## Instructions
When invoked, you should:

1. **{Action 1}**
   - {Detail}
   - {Detail}

2. **{Action 2}**
   - {Detail}
   - {Detail}

3. **{Action 3}**
   - {Detail}
   - {Detail}

## Guidelines
- {Behavioral rule}
- {Constraint}
- {Best practice}

## Output Format
{How the agent should structure its output}

## Scope
This agent WILL:
- {Capability}
- {Capability}

This agent WILL NOT:
- {Limitation}
- {Limitation}

## Error Handling
- **{Error Type}**: {How to handle}
- **{Error Type}**: {How to handle}
```

---

## Section Breakdown

### Frontmatter (Required)
YAML config at the very top of the file.

- `name`: Agent identifier (kebab-case)
- `description`: What it does and when to use it
- `tools`: Which tools it can access (optional)
- `model`: Which model to use - `inherit`, `sonnet`, `opus`, `haiku` (optional)

### Purpose (Required)
Brief description of what the agent does.

### Instructions (Required)
Step-by-step actions the agent should take. Be specific.

### Guidelines (Required)
Behavioral rules, priorities, and constraints.

### Output Format (Required)
Define the structure of the agent's response.

### Scope (Recommended)
What the agent WILL and WILL NOT do.

### Error Handling (Recommended)
How to handle common errors and edge cases.

---

## Example

```markdown
---
name: go-linter
description: Analyze Go code for common issues and best practices
tools: Read, Grep, Glob
model: inherit
---

# Go Linter Agent

## Purpose
Analyze Go code for syntax errors, common mistakes, and style violations.

## Instructions
When invoked, you should:

1. **Scan the code** for:
   - Syntax errors
   - Unused variables
   - Missing error checks
   - Style violations

2. **Report findings** with:
   - File and line number
   - Issue description
   - Suggested fix

## Guidelines
- Prioritize errors over warnings
- Include code examples in suggestions
- Reference Go documentation when relevant

## Output Format
**[LEVEL]** `file.go:line` - {Issue}
Fix: {Suggestion}

## Scope
This agent WILL:
- Analyze Go files
- Check syntax and style

This agent WILL NOT:
- Modify files
- Run tests

## Error Handling
- **Invalid syntax**: Report location and skip detailed analysis
- **File not found**: Report error and continue
```

---

## Quick Tips

- **Name**: Use kebab-case, be specific (`go-linter` not `helper`)
- **Instructions**: 3-6 specific steps, use bold headers
- **Tools**: Only list tools you actually need
- **Scope**: Be explicit about limitations
- **Output**: Define a consistent format
