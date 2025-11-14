package services

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockGitHubClient is a mock implementation of GitHubClient for testing
type mockGitHubClient struct {
	fetchFileFunc func(path string) ([]byte, error)
}

func (m *mockGitHubClient) FetchFile(path string) ([]byte, error) {
	if m.fetchFileFunc != nil {
		return m.fetchFileFunc(path)
	}
	return nil, nil
}

// Helper function to create a test registry
func createTestRegistry() *models.Registry {
	now := time.Now()
	return &models.Registry{
		Version:   "1.0",
		UpdatedAt: now,
		Tools: map[models.ToolType][]*models.ToolInfo{
			models.ToolTypeAgent: {
				{
					Name:        "code-reviewer",
					Version:     "1.0.0",
					Description: "Code review automation agent",
					Type:        models.ToolTypeAgent,
					Author:      "Claude Team",
					Tags:        []string{"code-review", "quality"},
					File:        "agents/code-reviewer.zip",
					Downloads:   150,
					CreatedAt:   now.Add(-30 * 24 * time.Hour),
					UpdatedAt:   now.Add(-5 * 24 * time.Hour),
				},
				{
					Name:        "git-helper",
					Version:     "1.2.0",
					Description: "Git workflow helper",
					Type:        models.ToolTypeAgent,
					Author:      "Community",
					Tags:        []string{"git", "workflow"},
					File:        "agents/git-helper.zip",
					Downloads:   89,
					CreatedAt:   now.Add(-60 * 24 * time.Hour),
					UpdatedAt:   now.Add(-10 * 24 * time.Hour),
				},
			},
			models.ToolTypeCommand: {
				{
					Name:        "test-coverage",
					Version:     "1.0.0",
					Description: "Run tests with coverage",
					Type:        models.ToolTypeCommand,
					Author:      "Testing Team",
					Tags:        []string{"testing", "coverage"},
					File:        "commands/test-coverage.zip",
					Downloads:   245,
					CreatedAt:   now.Add(-20 * 24 * time.Hour),
					UpdatedAt:   now.Add(-2 * 24 * time.Hour),
				},
			},
			models.ToolTypeSkill: {
				{
					Name:        "github-api",
					Version:     "1.0.0",
					Description: "GitHub API patterns",
					Type:        models.ToolTypeSkill,
					Author:      "API Team",
					Tags:        []string{"github", "api"},
					File:        "skills/github-api.zip",
					Downloads:   178,
					CreatedAt:   now.Add(-45 * 24 * time.Hour),
					UpdatedAt:   now.Add(-7 * 24 * time.Hour),
				},
			},
		},
	}
}

func TestNewRegistryService(t *testing.T) {
	mockClient := &mockGitHubClient{}
	service := NewRegistryService(mockClient)

	assert.NotNil(t, service)
	assert.Nil(t, service.registry)
}

func TestFetchRegistry_Success(t *testing.T) {
	registry := createTestRegistry()
	registryJSON, err := json.Marshal(registry)
	require.NoError(t, err)

	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			assert.Equal(t, "registry.json", path)
			return registryJSON, nil
		},
	}

	service := NewRegistryService(mockClient)

	result, err := service.FetchRegistry()
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1.0", result.Version)
	assert.Len(t, result.Tools, 3)

	// Verify it's cached
	assert.NotNil(t, service.registry)
}

func TestFetchRegistry_InvalidJSON(t *testing.T) {
	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			return []byte("invalid json"), nil
		},
	}

	service := NewRegistryService(mockClient)

	_, err := service.FetchRegistry()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse registry")
}

func TestGetRegistry_CachedVsFetch(t *testing.T) {
	registry := createTestRegistry()
	registryJSON, _ := json.Marshal(registry)

	callCount := 0
	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			callCount++
			return registryJSON, nil
		},
	}

	service := NewRegistryService(mockClient)

	// First call should fetch
	result1, err := service.GetRegistry()
	require.NoError(t, err)
	assert.NotNil(t, result1)
	assert.Equal(t, 1, callCount)

	// Second call should use cache
	result2, err := service.GetRegistry()
	require.NoError(t, err)
	assert.NotNil(t, result2)
	assert.Equal(t, 1, callCount) // Still 1, not 2

	// Same pointer (cached)
	assert.Same(t, result1, result2)
}

func TestRefreshRegistry(t *testing.T) {
	registry := createTestRegistry()
	registryJSON, _ := json.Marshal(registry)

	callCount := 0
	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			callCount++
			return registryJSON, nil
		},
	}

	service := NewRegistryService(mockClient)

	// Initial fetch
	_, err := service.GetRegistry()
	require.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Refresh should fetch again
	_, err = service.RefreshRegistry()
	require.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

func TestGetTool(t *testing.T) {
	registry := createTestRegistry()
	registryJSON, _ := json.Marshal(registry)

	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			return registryJSON, nil
		},
	}

	service := NewRegistryService(mockClient)

	tests := []struct {
		name     string
		toolName string
		toolType models.ToolType
		wantErr  bool
	}{
		{"found agent", "code-reviewer", models.ToolTypeAgent, false},
		{"found command", "test-coverage", models.ToolTypeCommand, false},
		{"not found", "non-existent", models.ToolTypeAgent, true},
		{"wrong type", "code-reviewer", models.ToolTypeCommand, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool, err := service.GetTool(tt.toolName, tt.toolType)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tool)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, tool)
				assert.Equal(t, tt.toolName, tool.Name)
			}
		})
	}
}

func TestSearchTools(t *testing.T) {
	registry := createTestRegistry()
	registryJSON, _ := json.Marshal(registry)

	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			return registryJSON, nil
		},
	}

	service := NewRegistryService(mockClient)

	tests := []struct {
		name        string
		filter      *models.SearchFilter
		wantCount   int
		wantNames   []string
		wantErr     bool
	}{
		{
			name:      "search by name",
			filter:    &models.SearchFilter{Query: "code-reviewer"},
			wantCount: 1,
			wantNames: []string{"code-reviewer"},
		},
		{
			name:      "search by tag",
			filter:    &models.SearchFilter{Query: "git"},
			wantCount: 2, // git-helper and github-api
		},
		{
			name:      "search with type filter",
			filter:    &models.SearchFilter{Query: "test", Type: models.ToolTypeCommand},
			wantCount: 1,
			wantNames: []string{"test-coverage"},
		},
		{
			name:      "search case insensitive",
			filter:    &models.SearchFilter{Query: "CODE", CaseSensitive: false},
			wantCount: 1,
			wantNames: []string{"code-reviewer"},
		},
		{
			name:      "search with regex",
			filter:    &models.SearchFilter{Query: "^git-", Regex: true},
			wantCount: 1,
			wantNames: []string{"git-helper"},
		},
		{
			name:      "search with author filter",
			filter:    &models.SearchFilter{Query: "helper", Author: "Community"},
			wantCount: 1,
			wantNames: []string{"git-helper"},
		},
		{
			name:      "search with min downloads",
			filter:    &models.SearchFilter{Query: "test", MinDownloads: 200},
			wantCount: 1,
			wantNames: []string{"test-coverage"},
		},
		{
			name:    "invalid empty query",
			filter:  &models.SearchFilter{Query: ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.SearchTools(tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, results, tt.wantCount)

			if len(tt.wantNames) > 0 {
				for i, name := range tt.wantNames {
					assert.Equal(t, name, results[i].Name)
				}
			}
		})
	}
}

func TestListTools(t *testing.T) {
	registry := createTestRegistry()
	registryJSON, _ := json.Marshal(registry)

	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			return registryJSON, nil
		},
	}

	service := NewRegistryService(mockClient)

	tests := []struct {
		name      string
		filter    *models.ListFilter
		wantCount int
		checkFunc func(*testing.T, []*models.ToolInfo)
	}{
		{
			name:      "list all",
			filter:    &models.ListFilter{},
			wantCount: 4, // All tools
		},
		{
			name:      "list by type",
			filter:    &models.ListFilter{Type: models.ToolTypeAgent},
			wantCount: 2,
		},
		{
			name:      "list with limit",
			filter:    &models.ListFilter{Limit: 2},
			wantCount: 2,
		},
		{
			name:      "list by author",
			filter:    &models.ListFilter{Author: "Claude Team"},
			wantCount: 1,
		},
		{
			name:      "list by tags",
			filter:    &models.ListFilter{Tags: []string{"git"}},
			wantCount: 1, // git-helper only
		},
		{
			name:      "sort by name",
			filter:    &models.ListFilter{SortBy: models.SortByName},
			wantCount: 4,
			checkFunc: func(t *testing.T, results []*models.ToolInfo) {
				assert.Equal(t, "code-reviewer", results[0].Name)
			},
		},
		{
			name:      "sort by downloads desc",
			filter:    &models.ListFilter{SortBy: models.SortByDownloads, SortDesc: true},
			wantCount: 4,
			checkFunc: func(t *testing.T, results []*models.ToolInfo) {
				assert.Equal(t, "test-coverage", results[0].Name) // 245 downloads
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := service.ListTools(tt.filter)
			require.NoError(t, err)
			assert.Len(t, results, tt.wantCount)

			if tt.checkFunc != nil {
				tt.checkFunc(t, results)
			}
		})
	}
}

func TestGetToolsByType(t *testing.T) {
	registry := createTestRegistry()
	registryJSON, _ := json.Marshal(registry)

	mockClient := &mockGitHubClient{
		fetchFileFunc: func(path string) ([]byte, error) {
			return registryJSON, nil
		},
	}

	service := NewRegistryService(mockClient)

	tests := []struct {
		name      string
		toolType  models.ToolType
		wantCount int
		wantErr   bool
	}{
		{"agents", models.ToolTypeAgent, 2, false},
		{"commands", models.ToolTypeCommand, 1, false},
		{"skills", models.ToolTypeSkill, 1, false},
		{"invalid type", models.ToolType("invalid"), 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tools, err := service.GetToolsByType(tt.toolType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, tools, tt.wantCount)
		})
	}
}

func TestHasAnyTag(t *testing.T) {
	tests := []struct {
		name   string
		slice1 []string
		slice2 []string
		want   bool
	}{
		{"has match", []string{"a", "b", "c"}, []string{"b", "d"}, true},
		{"no match", []string{"a", "b"}, []string{"c", "d"}, false},
		{"case insensitive", []string{"Git", "Test"}, []string{"git"}, true},
		{"empty slice1", []string{}, []string{"a"}, false},
		{"empty slice2", []string{"a"}, []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasAnyTag(tt.slice1, tt.slice2)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSortTools(t *testing.T) {
	now := time.Now()
	tools := []*models.ToolInfo{
		{Name: "zebra", Downloads: 100, CreatedAt: now.Add(-10 * time.Hour), UpdatedAt: now.Add(-1 * time.Hour)},
		{Name: "alpha", Downloads: 200, CreatedAt: now.Add(-5 * time.Hour), UpdatedAt: now.Add(-2 * time.Hour)},
		{Name: "beta", Downloads: 50, CreatedAt: now.Add(-15 * time.Hour), UpdatedAt: now.Add(-3 * time.Hour)},
	}

	t.Run("sort by name asc", func(t *testing.T) {
		toolsCopy := make([]*models.ToolInfo, len(tools))
		copy(toolsCopy, tools)

		sortTools(toolsCopy, models.SortByName, false)
		assert.Equal(t, "alpha", toolsCopy[0].Name)
		assert.Equal(t, "beta", toolsCopy[1].Name)
		assert.Equal(t, "zebra", toolsCopy[2].Name)
	})

	t.Run("sort by downloads desc", func(t *testing.T) {
		toolsCopy := make([]*models.ToolInfo, len(tools))
		copy(toolsCopy, tools)

		sortTools(toolsCopy, models.SortByDownloads, true)
		assert.Equal(t, 200, toolsCopy[0].Downloads)
		assert.Equal(t, 100, toolsCopy[1].Downloads)
		assert.Equal(t, 50, toolsCopy[2].Downloads)
	})
}
