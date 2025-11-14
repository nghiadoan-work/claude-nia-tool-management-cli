package services

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/nghiadt/claude-nia-tool-management-cli/pkg/models"
)

// GitHubClientInterface defines the methods needed from GitHubClient
type GitHubClientInterface interface {
	FetchFile(path string) ([]byte, error)
}

// CacheManagerInterface defines the methods needed from CacheManager
type CacheManagerInterface interface {
	GetRegistry() (*models.Registry, error)
	SetRegistry(registry *models.Registry) error
	IsValid() bool
	Invalidate() error
}

// RegistryService manages tool registry operations
type RegistryService struct {
	githubClient GitHubClientInterface
	cacheManager CacheManagerInterface
	registry     *models.Registry
	useCache     bool
}

// NewRegistryService creates a new RegistryService with cache support
func NewRegistryService(githubClient GitHubClientInterface, cacheManager CacheManagerInterface) *RegistryService {
	return &RegistryService{
		githubClient: githubClient,
		cacheManager: cacheManager,
		useCache:     cacheManager != nil,
	}
}

// NewRegistryServiceWithoutCache creates a new RegistryService without cache support
func NewRegistryServiceWithoutCache(githubClient GitHubClientInterface) *RegistryService {
	return &RegistryService{
		githubClient: githubClient,
		cacheManager: nil,
		useCache:     false,
	}
}

// FetchRegistry fetches and parses the registry.json from GitHub
func (rs *RegistryService) FetchRegistry() (*models.Registry, error) {
	// Fetch registry.json from GitHub
	data, err := rs.githubClient.FetchFile("registry.json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch registry: %w", err)
	}

	// Parse JSON
	var registry models.Registry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("failed to parse registry: %w", err)
	}

	// Validate registry
	if err := registry.Validate(); err != nil {
		return nil, fmt.Errorf("invalid registry: %w", err)
	}

	// Cache the registry in memory
	rs.registry = &registry

	// Cache to disk if cache manager is available
	if rs.useCache && rs.cacheManager != nil {
		if err := rs.cacheManager.SetRegistry(&registry); err != nil {
			// Log warning but don't fail - cache is not critical
			// In production, this would use a proper logger
			_ = err
		}
	}

	return &registry, nil
}

// GetRegistry returns the cached registry or fetches it if not available
func (rs *RegistryService) GetRegistry() (*models.Registry, error) {
	// First check in-memory cache
	if rs.registry != nil {
		return rs.registry, nil
	}

	// Then check disk cache if available
	if rs.useCache && rs.cacheManager != nil && rs.cacheManager.IsValid() {
		registry, err := rs.cacheManager.GetRegistry()
		if err == nil {
			// Update in-memory cache
			rs.registry = registry
			return registry, nil
		}
		// If cache read fails, continue to fetch from GitHub
	}

	// No cache available or cache invalid, fetch from GitHub
	return rs.FetchRegistry()
}

// RefreshRegistry forces a refresh of the registry from GitHub
func (rs *RegistryService) RefreshRegistry() (*models.Registry, error) {
	// Invalidate disk cache if available
	if rs.useCache && rs.cacheManager != nil {
		_ = rs.cacheManager.Invalidate() // Ignore errors
	}

	// Clear in-memory cache
	rs.registry = nil

	// Fetch fresh data from GitHub
	return rs.FetchRegistry()
}

// InvalidateCache invalidates the disk cache
func (rs *RegistryService) InvalidateCache() error {
	if rs.useCache && rs.cacheManager != nil {
		return rs.cacheManager.Invalidate()
	}
	return nil
}

// GetTool finds a specific tool by name and type
func (rs *RegistryService) GetTool(name string, toolType models.ToolType) (*models.ToolInfo, error) {
	registry, err := rs.GetRegistry()
	if err != nil {
		return nil, err
	}

	return registry.GetTool(name, toolType)
}

// SearchTools searches for tools matching the filter criteria
func (rs *RegistryService) SearchTools(filter *models.SearchFilter) ([]*models.ToolInfo, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid search filter: %w", err)
	}

	registry, err := rs.GetRegistry()
	if err != nil {
		return nil, err
	}

	var results []*models.ToolInfo
	var pattern *regexp.Regexp

	// Compile regex if needed
	if filter.Regex {
		flags := ""
		if !filter.CaseSensitive {
			flags = "(?i)"
		}
		pattern, err = regexp.Compile(flags + filter.Query)
		if err != nil {
			return nil, fmt.Errorf("invalid regex pattern: %w", err)
		}
	}

	// Search through all tools
	for toolType, tools := range registry.Tools {
		// Skip if type filter is set and doesn't match
		if filter.Type != "" && toolType != filter.Type {
			continue
		}

		for _, tool := range tools {
			if rs.matchesTool(tool, filter, pattern) {
				results = append(results, tool)
			}
		}
	}

	return results, nil
}

// ListTools lists tools with optional filtering
func (rs *RegistryService) ListTools(filter *models.ListFilter) ([]*models.ToolInfo, error) {
	if err := filter.Validate(); err != nil {
		return nil, fmt.Errorf("invalid list filter: %w", err)
	}

	registry, err := rs.GetRegistry()
	if err != nil {
		return nil, err
	}

	var results []*models.ToolInfo

	// Collect tools matching the filter
	for toolType, tools := range registry.Tools {
		// Skip if type filter is set and doesn't match
		if filter.Type != "" && toolType != filter.Type {
			continue
		}

		for _, tool := range tools {
			// Check tag filter
			if len(filter.Tags) > 0 && !hasAnyTag(tool.Tags, filter.Tags) {
				continue
			}

			// Check author filter
			if filter.Author != "" && !strings.EqualFold(tool.Author, filter.Author) {
				continue
			}

			results = append(results, tool)
		}
	}

	// Sort results
	if filter.SortBy != "" {
		sortTools(results, filter.SortBy, filter.SortDesc)
	}

	// Apply limit
	if filter.Limit > 0 && len(results) > filter.Limit {
		results = results[:filter.Limit]
	}

	return results, nil
}

// GetToolsByType returns all tools of a specific type
func (rs *RegistryService) GetToolsByType(toolType models.ToolType) ([]*models.ToolInfo, error) {
	if err := toolType.Validate(); err != nil {
		return nil, err
	}

	registry, err := rs.GetRegistry()
	if err != nil {
		return nil, err
	}

	tools, ok := registry.Tools[toolType]
	if !ok {
		return []*models.ToolInfo{}, nil
	}

	return tools, nil
}

// matchesTool checks if a tool matches the search criteria
func (rs *RegistryService) matchesTool(tool *models.ToolInfo, filter *models.SearchFilter, pattern *regexp.Regexp) bool {
	// Match against name, description, tags, and author
	searchTargets := []string{
		tool.Name,
		tool.Description,
		tool.Author,
		strings.Join(tool.Tags, " "),
	}

	query := filter.Query
	if !filter.CaseSensitive {
		query = strings.ToLower(query)
	}

	for _, target := range searchTargets {
		if !filter.CaseSensitive {
			target = strings.ToLower(target)
		}

		var matches bool
		if filter.Regex {
			matches = pattern.MatchString(target)
		} else {
			matches = strings.Contains(target, query)
		}

		if matches {
			// Apply additional filters
			if len(filter.Tags) > 0 && !hasAnyTag(tool.Tags, filter.Tags) {
				continue
			}
			if filter.Author != "" && !strings.EqualFold(tool.Author, filter.Author) {
				continue
			}
			if filter.MinDownloads > 0 && tool.Downloads < filter.MinDownloads {
				continue
			}
			return true
		}
	}

	return false
}

// hasAnyTag checks if slice1 contains any element from slice2
func hasAnyTag(slice1, slice2 []string) bool {
	for _, s1 := range slice1 {
		for _, s2 := range slice2 {
			if strings.EqualFold(s1, s2) {
				return true
			}
		}
	}
	return false
}

// sortTools sorts tools by the specified field
func sortTools(tools []*models.ToolInfo, sortBy models.SortField, desc bool) {
	// Simple bubble sort (for small lists, this is fine)
	n := len(tools)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			shouldSwap := false

			switch sortBy {
			case models.SortByName:
				if desc {
					shouldSwap = tools[j].Name < tools[j+1].Name
				} else {
					shouldSwap = tools[j].Name > tools[j+1].Name
				}
			case models.SortByCreated:
				if desc {
					shouldSwap = tools[j].CreatedAt.Before(tools[j+1].CreatedAt)
				} else {
					shouldSwap = tools[j].CreatedAt.After(tools[j+1].CreatedAt)
				}
			case models.SortByUpdated:
				if desc {
					shouldSwap = tools[j].UpdatedAt.Before(tools[j+1].UpdatedAt)
				} else {
					shouldSwap = tools[j].UpdatedAt.After(tools[j+1].UpdatedAt)
				}
			case models.SortByDownloads:
				if desc {
					shouldSwap = tools[j].Downloads < tools[j+1].Downloads
				} else {
					shouldSwap = tools[j].Downloads > tools[j+1].Downloads
				}
			}

			if shouldSwap {
				tools[j], tools[j+1] = tools[j+1], tools[j]
			}
		}
	}
}
