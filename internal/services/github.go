package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v56/github"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/oauth2"
)

// GitHubClient handles interactions with GitHub API
type GitHubClient struct {
	client    *github.Client
	owner     string
	repo      string
	branch    string
	ctx       context.Context
	authToken string
}

// GitHubClientConfig holds configuration for GitHubClient
type GitHubClientConfig struct {
	Owner     string
	Repo      string
	Branch    string
	AuthToken string
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(config GitHubClientConfig) *GitHubClient {
	ctx := context.Background()

	// Try to get auth token from various sources if not provided
	authToken := config.AuthToken
	if authToken == "" {
		authToken = GetGitHubToken()
	}

	var client *github.Client
	if authToken != "" {
		// Authenticated client (5000 req/hr)
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: authToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		// Unauthenticated client (60 req/hr)
		client = github.NewClient(nil)
	}

	return &GitHubClient{
		client:    client,
		owner:     config.Owner,
		repo:      config.Repo,
		branch:    config.Branch,
		ctx:       ctx,
		authToken: authToken,
	}
}

// FetchFile fetches a file from the GitHub repository
func (gc *GitHubClient) FetchFile(path string) ([]byte, error) {
	var content []byte
	var err error

	// Retry with exponential backoff
	err = gc.retryWithBackoff(func() error {
		fileContent, _, resp, fetchErr := gc.client.Repositories.GetContents(
			gc.ctx,
			gc.owner,
			gc.repo,
			path,
			&github.RepositoryContentGetOptions{Ref: gc.branch},
		)

		if fetchErr != nil {
			// Check for rate limit
			if resp != nil && resp.StatusCode == http.StatusForbidden {
				if gc.isRateLimited(resp) {
					return &RateLimitError{RetryAfter: gc.getRateLimitReset(resp)}
				}
			}
			return fetchErr
		}

		if fileContent == nil {
			return fmt.Errorf("file not found: %s", path)
		}

		contentStr, fetchErr := fileContent.GetContent()
		if fetchErr != nil {
			return fetchErr
		}
		content = []byte(contentStr)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to fetch file %s: %w", path, err)
	}

	return content, nil
}

// DownloadFile downloads a file from a URL with progress bar
func (gc *GitHubClient) DownloadFile(url string, size int64, showProgress bool) ([]byte, error) {
	var data []byte
	var err error

	err = gc.retryWithBackoff(func() error {
		req, reqErr := http.NewRequestWithContext(gc.ctx, "GET", url, nil)
		if reqErr != nil {
			return reqErr
		}

		// Add auth token if available
		if gc.authToken != "" {
			req.Header.Set("Authorization", "token "+gc.authToken)
		}

		client := &http.Client{
			Timeout: 10 * time.Minute,
		}

		resp, respErr := client.Do(req)
		if respErr != nil {
			return respErr
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusForbidden && gc.isRateLimitedHTTP(resp) {
				return &RateLimitError{RetryAfter: gc.getRateLimitResetHTTP(resp)}
			}
			return fmt.Errorf("HTTP error: %s", resp.Status)
		}

		var reader io.Reader = resp.Body

		// Add progress bar if requested and size is known
		if showProgress && size > 0 {
			bar := progressbar.DefaultBytes(
				size,
				"Downloading",
			)
			reader = io.TeeReader(resp.Body, bar)
		}

		data, respErr = io.ReadAll(reader)
		return respErr
	})

	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}

	return data, nil
}

// GetRateLimit returns current rate limit information
func (gc *GitHubClient) GetRateLimit() (*github.RateLimits, error) {
	limits, _, err := gc.client.RateLimits(gc.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate limits: %w", err)
	}
	return limits, nil
}

// retryWithBackoff retries a function with exponential backoff
func (gc *GitHubClient) retryWithBackoff(fn func() error) error {
	maxRetries := 3
	backoff := 1 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// Check if it's a rate limit error
		if rateLimitErr, ok := err.(*RateLimitError); ok {
			waitTime := rateLimitErr.RetryAfter
			if waitTime > 0 {
				time.Sleep(waitTime)
				continue
			}
		}

		// For other errors, use exponential backoff
		if i < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// isRateLimited checks if the response indicates rate limiting
func (gc *GitHubClient) isRateLimited(resp *github.Response) bool {
	return resp.Rate.Remaining == 0
}

// getRateLimitReset returns the duration until rate limit reset
func (gc *GitHubClient) getRateLimitReset(resp *github.Response) time.Duration {
	resetTime := resp.Rate.Reset.Time
	now := time.Now()
	if resetTime.After(now) {
		return resetTime.Sub(now)
	}
	return 0
}

// isRateLimitedHTTP checks if HTTP response indicates rate limiting
func (gc *GitHubClient) isRateLimitedHTTP(resp *http.Response) bool {
	return resp.Header.Get("X-RateLimit-Remaining") == "0"
}

// getRateLimitResetHTTP returns the duration until rate limit reset from HTTP response
func (gc *GitHubClient) getRateLimitResetHTTP(resp *http.Response) time.Duration {
	resetHeader := resp.Header.Get("X-RateLimit-Reset")
	if resetHeader != "" {
		var resetUnix int64
		fmt.Sscanf(resetHeader, "%d", &resetUnix)
		resetTime := time.Unix(resetUnix, 0)
		now := time.Now()
		if resetTime.After(now) {
			return resetTime.Sub(now)
		}
	}
	return 0
}

// RateLimitError represents a rate limit error
type RateLimitError struct {
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded, retry after %v", e.RetryAfter)
}

// ParseRepoURL parses a GitHub repository URL into owner and repo
// Supports formats:
// - https://github.com/owner/repo
// - https://github.com/owner/repo.git
// - github.com/owner/repo
func ParseRepoURL(url string) (owner, repo string, err error) {
	// Remove common prefixes
	url = removePrefix(url, "https://")
	url = removePrefix(url, "http://")
	url = removePrefix(url, "github.com/")
	url = removeSuffix(url, ".git")

	// Split by /
	parts := splitPath(url)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid GitHub URL format: %s", url)
	}

	return parts[0], parts[1], nil
}

// Helper functions
func removePrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}

func removeSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	for i := 0; i < len(path); i++ {
		if path[i] == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(path[i])
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// ForkRepository forks a repository to the authenticated user's account
func (gc *GitHubClient) ForkRepository(owner, repo string) (*github.Repository, error) {
	fork, _, err := gc.client.Repositories.CreateFork(gc.ctx, owner, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fork repository: %w", err)
	}

	// Wait for fork to be ready (GitHub needs time to prepare the fork)
	time.Sleep(3 * time.Second)
	return fork, nil
}

// GetAuthenticatedUser returns the authenticated user's login
func (gc *GitHubClient) GetAuthenticatedUser() (string, error) {
	user, _, err := gc.client.Users.Get(gc.ctx, "")
	if err != nil {
		return "", fmt.Errorf("failed to get authenticated user: %w", err)
	}
	return user.GetLogin(), nil
}

// GetDefaultBranch gets the default branch of a repository
func (gc *GitHubClient) GetDefaultBranch(owner, repo string) (string, error) {
	repository, _, err := gc.client.Repositories.Get(gc.ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get repository: %w", err)
	}
	return repository.GetDefaultBranch(), nil
}

// CreateBranch creates a new branch from a base branch
func (gc *GitHubClient) CreateBranch(owner, repo, newBranch, baseBranch string) error {
	// Get the base branch reference
	baseRef, _, err := gc.client.Git.GetRef(gc.ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		return fmt.Errorf("failed to get base branch: %w", err)
	}

	// Create new branch reference
	newRef := &github.Reference{
		Ref: github.String("refs/heads/" + newBranch),
		Object: &github.GitObject{
			SHA: baseRef.Object.SHA,
		},
	}

	_, _, err = gc.client.Git.CreateRef(gc.ctx, owner, repo, newRef)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	return nil
}

// UploadFile uploads a file to a repository
func (gc *GitHubClient) UploadFile(owner, repo, path, branch string, content []byte, message string) error {
	// Check if file exists
	_, _, resp, _ := gc.client.Repositories.GetContents(
		gc.ctx, owner, repo, path,
		&github.RepositoryContentGetOptions{Ref: branch},
	)

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		Branch:  github.String(branch),
	}

	// If file exists (status 200), we need its SHA for update
	if resp != nil && resp.StatusCode == http.StatusOK {
		fileContent, _, _, err := gc.client.Repositories.GetContents(
			gc.ctx, owner, repo, path,
			&github.RepositoryContentGetOptions{Ref: branch},
		)
		if err == nil && fileContent != nil {
			opts.SHA = fileContent.SHA
		}
	}

	_, _, err := gc.client.Repositories.CreateFile(gc.ctx, owner, repo, path, opts)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}

// CreatePullRequest creates a pull request
func (gc *GitHubClient) CreatePullRequest(owner, repo, title, body, head, base string) (*github.PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(title),
		Head:  github.String(head),
		Base:  github.String(base),
		Body:  github.String(body),
	}

	pr, _, err := gc.client.PullRequests.Create(gc.ctx, owner, repo, newPR)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return pr, nil
}

// GetGitHubToken attempts to get a GitHub token from various sources
func GetGitHubToken() string {
	// 1. Check GITHUB_TOKEN environment variable
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}

	// 2. Check GH_TOKEN environment variable (used by gh CLI)
	if token := os.Getenv("GH_TOKEN"); token != "" {
		return token
	}

	// 3. Try to get token from gh CLI
	if token := getTokenFromGHCLI(); token != "" {
		return token
	}

	return ""
}

// getTokenFromGHCLI attempts to get the GitHub token from gh CLI
func getTokenFromGHCLI() string {
	// Check if gh CLI is installed
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		// gh CLI not installed or not authenticated
		return ""
	}

	// Get the token
	cmd = exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	token := strings.TrimSpace(string(output))
	return token
}
