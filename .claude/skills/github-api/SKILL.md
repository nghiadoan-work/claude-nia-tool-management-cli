---
name: github-api
description: GitHub API patterns and best practices for cntm development. Use when implementing GitHub client, handling authentication, rate limits, or repository operations.
---

# GitHub API Integration for CNTM

Reference for working with GitHub API in the cntm project.

## Authentication

### Setup Personal Access Token

Generate token at: GitHub Settings → Developer settings → Personal access tokens

Required scopes:
- `repo` (for private repos)
- `public_repo` (for public repos)

### In Go Code

```go
import (
    "github.com/google/go-github/v56/github"
    "golang.org/x/oauth2"
)

func NewGitHubClient(token string) *github.Client {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
    tc := oauth2.NewClient(ctx, ts)
    return github.NewClient(tc)
}
```

## Rate Limits

**Limits**:
- Unauthenticated: 60 req/hour
- Authenticated: 5,000 req/hour

**Check Before Operations**:
```go
func (c *GitHubClient) CheckRateLimit(ctx context.Context) error {
    limits, _, err := c.client.RateLimits(ctx)
    if err != nil {
        return err
    }
    if limits.Core.Remaining < 10 {
        return fmt.Errorf("rate limit low: %d remaining", limits.Core.Remaining)
    }
    return nil
}
```

## Common Operations

### Download File from Repository

```go
func (c *GitHubClient) DownloadFile(ctx context.Context, owner, repo, path string) ([]byte, error) {
    fileContent, _, _, err := c.client.Repositories.GetContents(
        ctx, owner, repo, path,
        &github.RepositoryContentGetOptions{Ref: "main"},
    )
    if err != nil {
        return nil, fmt.Errorf("failed to get file: %w", err)
    }

    content, err := fileContent.GetContent()
    if err != nil {
        return nil, fmt.Errorf("failed to decode: %w", err)
    }

    return []byte(content), nil
}
```

**For large files (>1MB)**: Use `raw.githubusercontent.com` directly with HTTP client.

### Create Pull Request

```go
func (c *GitHubClient) CreatePR(ctx context.Context, owner, repo, title, body, head, base string) (*github.PullRequest, error) {
    newPR := &github.NewPullRequest{
        Title: github.String(title),
        Head:  github.String(head),
        Base:  github.String(base),
        Body:  github.String(body),
    }

    pr, _, err := c.client.PullRequests.Create(ctx, owner, repo, newPR)
    if err != nil {
        return nil, fmt.Errorf("failed to create PR: %w", err)
    }

    return pr, nil
}
```

### Fork Repository

```go
func (c *GitHubClient) ForkRepo(ctx context.Context, owner, repo string) (*github.Repository, error) {
    fork, _, err := c.client.Repositories.CreateFork(ctx, owner, repo, nil)
    if err != nil {
        return nil, fmt.Errorf("fork failed: %w", err)
    }

    // Wait for fork to be ready
    time.Sleep(2 * time.Second)
    return fork, nil
}
```

## Error Handling

### HTTP Status Codes

```go
func handleGitHubError(err error) error {
    var ghErr *github.ErrorResponse
    if !errors.As(err, &ghErr) {
        return err
    }

    switch ghErr.Response.StatusCode {
    case 401:
        return &CLIError{
            Type: ErrorTypeAuth,
            Message: "GitHub authentication failed",
            Hint: "Check your GitHub token",
        }
    case 403:
        return &CLIError{
            Type: ErrorTypeRateLimit,
            Message: "Rate limit exceeded",
            Hint: "Wait or use authentication",
        }
    case 404:
        return &CLIError{
            Type: ErrorTypeNotFound,
            Message: "Repository or file not found",
            Hint: "Verify the URL",
        }
    }
    return err
}
```

## Retry Pattern

### Exponential Backoff

```go
func (c *GitHubClient) retryOperation(ctx context.Context, op func() error) error {
    maxAttempts := 3

    for attempt := 0; attempt < maxAttempts; attempt++ {
        err := op()
        if err == nil {
            return nil
        }

        // Don't retry 4xx errors (except 429)
        var ghErr *github.ErrorResponse
        if errors.As(err, &ghErr) {
            status := ghErr.Response.StatusCode
            if status >= 400 && status < 500 && status != 429 {
                return err
            }
        }

        if attempt < maxAttempts-1 {
            backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
            time.Sleep(backoff)
        }
    }

    return fmt.Errorf("operation failed after %d attempts", maxAttempts)
}
```

## Testing with Mocks

```go
type MockGitHubClient struct {
    mock.Mock
}

func (m *MockGitHubClient) DownloadFile(ctx context.Context, owner, repo, path string) ([]byte, error) {
    args := m.Called(ctx, owner, repo, path)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]byte), args.Error(1)
}

// In tests:
mockGH := new(MockGitHubClient)
mockGH.On("DownloadFile", mock.Anything, "owner", "repo", "file").
    Return([]byte("content"), nil)
```

## Best Practices for CNTM

1. **Always check rate limits** before batch operations
2. **Use context** for cancellation support
3. **Implement retry** with exponential backoff
4. **Cache responses** when appropriate (registry.json)
5. **Never log tokens** - security critical
6. **Test with mocks** - don't hit real API in unit tests
7. **Handle errors gracefully** with user-friendly CLIError messages

## Common Issues

**401 Unauthorized**: Token invalid/expired → Generate new token

**403 Rate Limit**: Too many requests → Use authentication or wait for reset

**404 Not Found**: Repo/file doesn't exist → Verify path and permissions

**File >1MB**: API limit → Use raw.githubusercontent.com instead

## Quick Reference

```go
// Essential packages
import (
    "github.com/google/go-github/v56/github"
    "golang.org/x/oauth2"
)

// Create authenticated client
client := github.NewClient(oauth2.NewClient(ctx, tokenSource))

// Common methods
client.Repositories.GetContents(ctx, owner, repo, path, opts)
client.Repositories.CreateFork(ctx, owner, repo, opts)
client.PullRequests.Create(ctx, owner, repo, newPR)
client.RateLimits(ctx)
```

Use these patterns when implementing GitHubClient service in cntm.
