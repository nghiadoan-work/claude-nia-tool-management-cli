# Skill Template Guide

This guide explains how to create effective skills for Claude Code. Skills provide specialized knowledge, patterns, and best practices for specific domains or technologies.

## Overview

A skill is a knowledge artifact that provides:

1. **Domain Expertise** - Specialized knowledge in a specific area
2. **Implementation Patterns** - Proven approaches and code examples
3. **Best Practices** - Guidelines and recommendations
4. **Reference Materials** - Documentation and resources
5. **Common Pitfalls** - What to avoid and how to troubleshoot

## File Structure

### Skill Directory
```
.claude/skills/{skill-name}/
├── SKILL.md              # Main skill definition
├── examples/             # Code examples and patterns
│   ├── example-1.go
│   ├── example-2.go
│   └── README.md
└── reference/            # Reference documentation
    ├── api-docs.md
    ├── architecture.md
    └── resources.md
```

---

## SKILL.md Structure

**IMPORTANT**: Every SKILL.md file MUST begin with YAML frontmatter at the very top of the file.

### Frontmatter (Required - First Lines of File)

The skill file must start with YAML frontmatter containing name and description:

```markdown
---
name: {skill-name}
description: Brief description of what the skill provides
---
```

**Guidelines:**
- `name`: Use kebab-case, match the directory name
- `description`: One-line summary (50-80 characters)
- Keep it concise and descriptive

**Examples:**

```yaml
---
name: go-error-handling
description: Go error handling patterns and best practices
---
```

```yaml
---
name: react-hooks
description: React Hooks patterns, performance optimization, and common pitfalls
---
```

---

## Skill Sections

### 1. Quick Start

**Purpose**: Provide immediate value with a brief overview and quick usage guide.

**Template:**
```markdown
## Quick Start

This skill provides [brief description of capabilities].

Use this skill when you need to:
- [Use case 1]
- [Use case 2]
- [Use case 3]

**Quick Example:**
[Simple, concrete example that demonstrates the core value]
```

**Example:**
```markdown
## Quick Start

This skill provides comprehensive GitHub API integration patterns for Go applications.

Use this skill when you need to:
- Integrate with GitHub's REST or GraphQL APIs
- Implement OAuth authentication flows
- Handle rate limiting and pagination
- Build tools that interact with repositories, issues, or pull requests

**Quick Example:**
See `./examples/basic-client.go` for a minimal GitHub API client setup.
```

---

### 2. Implementation Workflow

**Purpose**: Guide users through implementing the skill's concepts step-by-step.

**Template:**
```markdown
## Implementation Workflow

### Step 1: [Initial Setup/Preparation]
[What to do first and why]

### Step 2: [Core Implementation]
[Main implementation steps]

### Step 3: [Configuration/Integration]
[How to configure and integrate]

### Step 4: [Testing/Validation]
[How to verify it works]

### Step 5: [Optimization/Production] (Optional)
[Production considerations]
```

**Guidelines:**
- 3-7 steps typically work best
- Each step should be actionable
- Include why, not just how
- Reference examples where applicable
- Keep it logical and sequential

**Example:**
```markdown
## Implementation Workflow

### Step 1: Install Dependencies
Install the GitHub API client library:
```bash
go get github.com/google/go-github/v56/github
go get golang.org/x/oauth2
```

### Step 2: Set Up Authentication
Create an authenticated client (see `./examples/auth-setup.go`):
- Personal Access Token for simple use cases
- OAuth App for user authentication
- GitHub App for organization-wide tools

### Step 3: Implement API Calls
Use typed methods for API operations (see `./examples/api-operations.go`):
- Repository operations
- Issue and PR management
- User and organization queries

### Step 4: Handle Rate Limits
Implement rate limit handling (see `./examples/rate-limiting.go`):
- Check rate limit status
- Implement exponential backoff
- Use conditional requests with ETags

### Step 5: Add Error Handling
Robust error handling (see `./reference/error-handling.md`):
- Check for API errors
- Handle network failures
- Retry transient errors
```

---

### 3. Knowledge Areas

**Purpose**: List the key concepts and topics the skill covers.

**Template:**
```markdown
## Knowledge Areas

[Brief intro about what knowledge this skill provides, reference to ./reference folder]

- **[Area 1]**: [Brief description]
- **[Area 2]**: [Brief description]
- **[Area 3]**: [Brief description]
- **[Area 4]**: [Brief description]

Detailed documentation available in `./reference/` folder.
```

**Guidelines:**
- List 4-8 major knowledge areas
- Use bold for area names
- Provide brief 1-sentence descriptions
- Link to detailed docs in ./reference/
- Organize from fundamental to advanced

**Example:**
```markdown
## Knowledge Areas

For detailed reference materials, see the `./reference/` folder.

- **Authentication Methods**: OAuth tokens, Personal Access Tokens, GitHub Apps, and when to use each
- **API Clients**: REST vs GraphQL clients, client configuration, and connection management
- **Rate Limiting**: Understanding rate limits, checking status, and implementing backoff strategies
- **Pagination**: Handling paginated responses, cursor-based vs offset-based pagination
- **Webhooks**: Setting up webhooks, validating signatures, and processing events
- **Error Handling**: API error types, retry logic, and graceful degradation
- **Testing**: Mocking GitHub API responses, integration testing strategies
- **Security**: Token storage, scope management, and secret scanning prevention

Detailed documentation available in `./reference/` folder.
```

---

### 4. Best Practices

**Purpose**: Provide actionable recommendations and guidelines.

**Template:**
```markdown
## Best Practices

1. **[Practice 1 Title]**: [Description and rationale]
2. **[Practice 2 Title]**: [Description and rationale]
3. **[Practice 3 Title]**: [Description and rationale]
4. **[Practice 4 Title]**: [Description and rationale]
5. **[Practice 5 Title]**: [Description and rationale]
```

**Guidelines:**
- 5-10 practices work well
- Make them specific and actionable
- Explain the "why" behind each
- Use bold for the practice title
- Order by importance or workflow sequence

**Example:**
```markdown
## Best Practices

1. **Use Typed Clients**: Always use the typed GitHub client methods rather than raw HTTP calls for better type safety and automatic marshaling

2. **Implement Rate Limit Checks**: Check `client.RateLimits()` before making bulk operations to avoid hitting rate limits

3. **Use Conditional Requests**: Leverage ETags for resources you poll frequently to save rate limit quota

4. **Set User-Agent Header**: Always set a descriptive User-Agent to help GitHub identify your application

5. **Handle Pagination Properly**: Use `ListOptions` with proper page size (100 max) and always check for next page

6. **Store Tokens Securely**: Never hardcode tokens, use environment variables or secure credential storage

7. **Implement Exponential Backoff**: When rate limited or getting 5xx errors, use exponential backoff with jitter

8. **Use Context for Cancellation**: Pass context.Context to all API calls for proper timeout and cancellation handling

9. **Validate Webhook Signatures**: Always verify webhook signatures using HMAC before processing events

10. **Log API Quotas**: Monitor and log rate limit remaining/reset times for proactive management
```

---

### 5. Patterns

**Purpose**: Provide reusable code patterns with explanations.

**Template:**
```markdown
## Patterns

### Pattern 1: [Pattern Name]

[Description of the pattern, when to use it, and why]

**Implementation:**
See `./examples/[example-file]` for complete implementation.

```[language]
// Brief code snippet showing the key parts
```

**Key Points:**
- [Important aspect 1]
- [Important aspect 2]

---

### Pattern 2: [Pattern Name]

[Description of the pattern, when to use it, and why]

**Implementation:**
See `./examples/[example-file]` for complete implementation.

```[language]
// Brief code snippet showing the key parts
```

**Key Points:**
- [Important aspect 1]
- [Important aspect 2]
```

**Guidelines:**
- 3-6 patterns typically sufficient
- Each pattern should solve a specific problem
- Show enough code to understand, link to full example
- Explain when and why to use the pattern
- Include gotchas or important considerations
- Reference example files in ./examples/ folder

**Example:**
```markdown
## Patterns

### Pattern 1: Authenticated Client with Token

Create a GitHub client with Personal Access Token authentication.

**Implementation:**
See `./examples/auth-client.go` for complete implementation.

```go
func NewGitHubClient(token string) *github.Client {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)
    return github.NewClient(tc)
}
```

**Key Points:**
- Use `oauth2.StaticTokenSource` for simple token auth
- Client is safe for concurrent use
- Pass context to all API calls

---

### Pattern 2: Rate Limit Aware Requests

Check and handle rate limits before making requests.

**Implementation:**
See `./examples/rate-limit-handler.go` for complete implementation.

```go
func (c *Client) CheckRateLimit(ctx context.Context) error {
    limits, _, err := c.client.RateLimits(ctx)
    if err != nil {
        return err
    }

    if limits.Core.Remaining < 100 {
        return fmt.Errorf("rate limit low: %d remaining, resets at %v",
            limits.Core.Remaining, limits.Core.Reset)
    }
    return nil
}
```

**Key Points:**
- Check before bulk operations
- Consider implementing automatic retry after reset
- Different limits for Core, Search, GraphQL

---

### Pattern 3: Paginated List Fetching

Fetch all items from a paginated API endpoint.

**Implementation:**
See `./examples/pagination.go` for complete implementation.

```go
func FetchAllRepositories(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
    var allRepos []*github.Repository
    opts := &github.RepositoryListByOrgOptions{
        ListOptions: github.ListOptions{PerPage: 100},
    }

    for {
        repos, resp, err := client.Repositories.ListByOrg(ctx, org, opts)
        if err != nil {
            return nil, err
        }
        allRepos = append(allRepos, repos...)
        if resp.NextPage == 0 {
            break
        }
        opts.Page = resp.NextPage
    }

    return allRepos, nil
}
```

**Key Points:**
- Always use max PerPage (100) for efficiency
- Check `resp.NextPage == 0` for completion
- Consider memory implications for large result sets
```

---

### 6. Common Pitfalls

**Purpose**: Warn users about frequent mistakes and issues.

**Template:**
```markdown
## Common Pitfalls

- **[Pitfall 1]**: [What goes wrong] → [How to avoid it]
- **[Pitfall 2]**: [What goes wrong] → [How to avoid it]
- **[Pitfall 3]**: [What goes wrong] → [How to avoid it]
- **[Pitfall 4]**: [What goes wrong] → [How to avoid it]
```

**Guidelines:**
- 4-8 pitfalls are typical
- Describe both the problem and solution
- Use → to separate problem from solution
- Based on real common mistakes
- Order by frequency or severity

**Example:**
```markdown
## Common Pitfalls

- **Hardcoded Tokens**: Never commit tokens to git → Use environment variables or secret management

- **Ignoring Rate Limits**: Hitting rate limits causes 403 errors → Check limits proactively and implement backoff

- **Not Handling Pagination**: Only getting first page of results → Always loop through all pages using NextPage

- **Missing Context**: API calls hang without timeout → Always pass context.Context with timeout

- **Incorrect Scopes**: Token lacks required permissions → Verify token has necessary OAuth scopes for operation

- **Not Validating Webhooks**: Processing unverified webhook events → Always verify HMAC signature before processing

- **Concurrent API Calls Without Limits**: Exhausting rate limits quickly → Implement semaphore or rate limiter for concurrent calls

- **Using String Pointers Incorrectly**: Nil pointer dereference on optional fields → Use `github.String()` helper or check for nil

- **Forgetting Error Response Details**: Missing API error details → Check `github.ErrorResponse` for detailed error info
```

---

### 7. Resources

**Purpose**: Link to external documentation, tools, and learning materials.

**Template:**
```markdown
## Resources

For detailed reference materials, see the `./reference/` folder.

**Official Documentation:**
- [Resource name](url) - Description
- [Resource name](url) - Description

**Tools & Libraries:**
- [Tool name](url) - Description
- [Tool name](url) - Description

**Community Resources:**
- [Resource name](url) - Description
- [Resource name](url) - Description

**Related Skills:**
- [skill-name] - Description
```

**Guidelines:**
- Categorize resources logically
- Include brief description for each
- Keep links current and official
- Reference internal ./reference/ folder
- Link to related skills if applicable

**Example:**
```markdown
## Resources

For detailed reference materials, see the `./reference/` folder:
- `./reference/api-reference.md` - Comprehensive API endpoint documentation
- `./reference/authentication.md` - Detailed auth setup and security guidelines
- `./reference/webhooks.md` - Webhook event types and payload schemas

**Official Documentation:**
- [GitHub REST API Docs](https://docs.github.com/en/rest) - Official REST API reference
- [GitHub GraphQL API](https://docs.github.com/en/graphql) - GraphQL API documentation
- [GitHub OAuth Apps](https://docs.github.com/en/developers/apps/building-oauth-apps) - OAuth app setup

**Tools & Libraries:**
- [go-github](https://github.com/google/go-github) - Official Go client library
- [oauth2](https://pkg.go.dev/golang.org/x/oauth2) - Go OAuth2 library
- [GitHub CLI](https://cli.github.com/) - Official command-line tool

**Community Resources:**
- [GitHub API Best Practices](https://github.blog/2021-04-05-api-best-practices/) - Official blog post
- [go-github Examples](https://github.com/google/go-github/tree/master/example) - Official examples

**Related Skills:**
- `oauth-patterns` - OAuth 2.0 implementation patterns
- `rest-api-design` - RESTful API client design patterns
```

---

## Examples Directory

### Purpose

The `./examples/` directory contains working code examples that demonstrate the patterns and concepts in the skill.

### Structure

```
examples/
├── README.md                    # Overview of all examples
├── basic-client.go             # Simple, minimal example
├── auth-setup.go               # Authentication patterns
├── rate-limiting.go            # Rate limit handling
├── pagination.go               # Pagination example
├── error-handling.go           # Error handling patterns
└── complete-integration.go     # Full-featured example
```

### Example README.md

```markdown
# Examples

This directory contains working code examples for the [skill-name] skill.

## Examples Index

- **basic-client.go** - Minimal client setup, good starting point
- **auth-setup.go** - Different authentication methods
- **rate-limiting.go** - Rate limit checking and handling
- **pagination.go** - Fetching all pages from paginated endpoints
- **error-handling.go** - Robust error handling patterns
- **complete-integration.go** - Full-featured real-world example

## Running Examples

```bash
# Set your GitHub token
export GITHUB_TOKEN="your-token-here"

# Run an example
go run basic-client.go
```

## Prerequisites

- Go 1.21 or later
- GitHub Personal Access Token
- Dependencies: `go mod download`
```

### Guidelines for Examples

1. **Self-Contained**: Each example should run independently
2. **Well-Commented**: Explain what the code does and why
3. **Realistic**: Use real-world scenarios
4. **Progressive**: Order from simple to complex
5. **Runnable**: Include instructions to execute
6. **Error Handling**: Show proper error handling
7. **Best Practices**: Follow the patterns from SKILL.md

---

## Reference Directory

### Purpose

The `./reference/` directory contains detailed documentation, specifications, and reference materials.

### Structure

```
reference/
├── api-reference.md        # API documentation
├── architecture.md         # Architectural patterns
├── configuration.md        # Configuration options
├── error-codes.md         # Error reference
├── security.md            # Security guidelines
└── troubleshooting.md     # Common issues and solutions
```

### Reference Document Template

```markdown
# [Topic Name]

## Overview

[Brief description of what this reference covers]

## [Section 1]

[Detailed information]

## [Section 2]

[Detailed information]

## Quick Reference

[Tables, lists, or quick-lookup information]

## See Also

- [Related reference doc]
- [Related example]
```

### Guidelines for Reference Materials

1. **Comprehensive**: Cover topics in depth
2. **Searchable**: Use clear headings and structure
3. **Tables**: Use tables for comparisons and quick reference
4. **Cross-References**: Link to related docs and examples
5. **Keep Updated**: Maintain accuracy with latest versions
6. **Separate Concerns**: One topic per file

---

## Complete Skill Example

Here's a complete minimal skill structure. **Note: The YAML frontmatter MUST be the first lines of the file.**

```markdown
---
name: error-handling-go
description: Go error handling patterns and best practices
---

# Go Error Handling Skill

Expert guidance on error handling in Go, from basic patterns to advanced error wrapping and custom error types.

## Quick Start

This skill provides comprehensive error handling patterns for Go applications.

Use this skill when you need to:
- Design error handling strategies
- Implement custom error types
- Use error wrapping and unwrapping
- Handle errors in concurrent code

**Quick Example:**
See `./examples/basic-errors.go` for fundamental error handling patterns.

## Implementation Workflow

### Step 1: Choose Error Strategy
Decide between sentinel errors, error types, or wrapped errors based on your use case.

### Step 2: Define Custom Errors
Create meaningful error types for your domain (see `./examples/custom-errors.go`).

### Step 3: Implement Error Wrapping
Use `fmt.Errorf` with `%w` to wrap errors and preserve context.

### Step 4: Add Error Handling
Check errors at every level and provide appropriate context.

### Step 5: Log and Monitor
Implement structured logging for errors in production.

## Knowledge Areas

For detailed reference materials, see the `./reference/` folder.

- **Error Types**: Sentinel errors, custom types, wrapped errors
- **Error Wrapping**: Using %w, errors.Is(), errors.As()
- **Stack Traces**: Capturing and logging call stacks
- **Concurrent Errors**: Handling errors in goroutines and channels
- **Testing Errors**: Error assertion and testing strategies

## Best Practices

1. **Always Check Errors**: Never ignore returned errors
2. **Add Context**: Wrap errors with additional context using %w
3. **Use Sentinel Errors for Behavior**: Define package-level error variables for expected errors
4. **Create Typed Errors for Data**: Use custom error types when you need to attach data
5. **Don't Panic**: Reserve panic for truly exceptional circumstances

## Patterns

### Pattern 1: Error Wrapping

Add context while preserving the original error.

**Implementation:**
See `./examples/error-wrapping.go` for complete implementation.

```go
func ProcessFile(path string) error {
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to process file %s: %w", path, err)
    }
    // ... process data
    return nil
}
```

**Key Points:**
- Use `%w` to wrap errors
- Add contextual information
- Preserve error chain for `errors.Is()` and `errors.As()`

## Common Pitfalls

- **Ignoring Errors**: Using `_` to ignore errors → Always handle or explicitly log
- **Double Wrapping**: Wrapping already wrapped errors → Check if error needs additional context
- **Logging and Returning**: Both logging and returning same error → Choose one or the other
- **Missing Context**: Just returning raw errors → Add context about what failed

## Resources

For detailed reference materials, see the `./reference/` folder.

**Official Documentation:**
- [Go Blog: Error Handling](https://go.dev/blog/error-handling-and-go) - Official guidance
- [errors package](https://pkg.go.dev/errors) - Standard library reference

**Related Skills:**
- `go-testing` - Testing error handling code
- `go-logging` - Structured logging for errors
```

---

## Best Practices for Skill Creation

### Content Quality

1. **Be Specific**: Focus on one domain or technology
2. **Be Practical**: Emphasize real-world usage over theory
3. **Be Complete**: Cover the topic comprehensively
4. **Be Current**: Keep examples and references up-to-date
5. **Be Concise**: Clear and direct writing

### Organization

1. **Logical Flow**: Quick Start → Workflow → Deep Knowledge → Resources
2. **Progressive Complexity**: Simple examples first, advanced last
3. **Cross-Reference**: Link between SKILL.md, examples, and reference
4. **Consistent Naming**: Use kebab-case for files and directories
5. **Clear Structure**: Use headings effectively

### Examples and References

1. **Working Code**: All examples must run without errors
2. **Documented**: Comment code thoroughly
3. **Realistic**: Use real-world scenarios
4. **Separated**: Keep examples and reference docs separate
5. **Indexed**: Provide README.md in examples/ and reference/

### Maintenance

1. **Version Aware**: Note which versions the skill applies to
2. **Update Regularly**: Keep pace with technology changes
3. **Track Dependencies**: Document required libraries/tools
4. **Community Input**: Accept feedback and contributions
5. **Test Examples**: Verify examples work with current versions

---

## Skill Creation Checklist

Before publishing a skill, verify:

- [ ] YAML frontmatter at the very top of SKILL.md with name and description (first lines of file)
- [ ] Quick Start section with clear use cases
- [ ] Implementation Workflow with 3-7 actionable steps
- [ ] Knowledge Areas with reference to ./reference/ folder
- [ ] 5-10 Best Practices with rationale
- [ ] 3-6 Patterns with code examples
- [ ] Common Pitfalls with solutions
- [ ] Resources section with ./reference/ folder mentioned
- [ ] examples/ directory with working code
- [ ] examples/README.md with index and instructions
- [ ] reference/ directory with detailed documentation
- [ ] All links and file references are valid
- [ ] Code examples are tested and work
- [ ] Grammar and spelling checked
- [ ] Consistent formatting throughout

---

## Summary

A well-crafted skill provides:

- **Immediate Value**: Quick start gets users productive fast
- **Deep Knowledge**: Reference materials for comprehensive understanding
- **Practical Examples**: Working code demonstrating patterns
- **Guidance**: Best practices and pitfall warnings
- **Resources**: Links to additional learning materials

Structure your skills to be both reference material (for looking up specifics) and learning resource (for understanding concepts).
