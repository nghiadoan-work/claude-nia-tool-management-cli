package models

import (
	"fmt"
	"time"
)

// ToolType represents the type of Claude Code tool
type ToolType string

const (
	ToolTypeAgent   ToolType = "agent"
	ToolTypeCommand ToolType = "command"
	ToolTypeSkill   ToolType = "skill"
)

// Validate checks if the ToolType is valid
func (t ToolType) Validate() error {
	switch t {
	case ToolTypeAgent, ToolTypeCommand, ToolTypeSkill:
		return nil
	default:
		return fmt.Errorf("invalid tool type: %s", t)
	}
}

// ToolInfo represents a tool in the registry
type ToolInfo struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Type        ToolType  `json:"type"`
	Author      string    `json:"author"`
	Tags        []string  `json:"tags"`
	File        string    `json:"file"`       // Path in repo to ZIP file
	Size        int64     `json:"size"`       // Size in bytes
	Downloads   int       `json:"downloads"`  // Download count
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate checks if ToolInfo is valid
func (t *ToolInfo) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if t.Version == "" {
		return fmt.Errorf("tool version cannot be empty")
	}
	if err := t.Type.Validate(); err != nil {
		return err
	}
	if t.File == "" {
		return fmt.Errorf("tool file path cannot be empty")
	}
	return nil
}

// Registry represents the registry.json structure from GitHub
type Registry struct {
	Version   string                    `json:"version"`
	UpdatedAt time.Time                 `json:"updated_at"`
	Tools     map[ToolType][]*ToolInfo `json:"tools"`
}

// Validate checks if Registry is valid
func (r *Registry) Validate() error {
	if r.Version == "" {
		return fmt.Errorf("registry version cannot be empty")
	}
	if r.Tools == nil {
		return fmt.Errorf("registry tools cannot be nil")
	}

	// Validate all tools
	for toolType, tools := range r.Tools {
		if err := toolType.Validate(); err != nil {
			return fmt.Errorf("invalid tool type in registry: %w", err)
		}
		for _, tool := range tools {
			if err := tool.Validate(); err != nil {
				return fmt.Errorf("invalid tool %s: %w", tool.Name, err)
			}
		}
	}

	return nil
}

// GetTool finds a tool by name and type in the registry
func (r *Registry) GetTool(name string, toolType ToolType) (*ToolInfo, error) {
	tools, ok := r.Tools[toolType]
	if !ok {
		return nil, fmt.Errorf("no tools of type %s in registry", toolType)
	}

	for _, tool := range tools {
		if tool.Name == name {
			return tool, nil
		}
	}

	return nil, fmt.Errorf("tool %s not found in registry", name)
}

// InstalledTool represents a tool installed locally
type InstalledTool struct {
	Version     string    `json:"version"`
	Type        ToolType  `json:"type"`
	InstalledAt time.Time `json:"installed_at"`
	Source      string    `json:"source"`    // "registry" or URL
	Integrity   string    `json:"integrity"` // SHA256 hash
}

// Validate checks if InstalledTool is valid
func (i *InstalledTool) Validate() error {
	if i.Version == "" {
		return fmt.Errorf("installed tool version cannot be empty")
	}
	if err := i.Type.Validate(); err != nil {
		return err
	}
	if i.Source == "" {
		return fmt.Errorf("installed tool source cannot be empty")
	}
	return nil
}

// LockFile represents the .claude-lock.json structure
type LockFile struct {
	Version   string                    `json:"version"`
	UpdatedAt time.Time                 `json:"updated_at"`
	Registry  string                    `json:"registry"`
	Tools     map[string]*InstalledTool `json:"tools"` // Key: tool name
}

// Validate checks if LockFile is valid
func (l *LockFile) Validate() error {
	if l.Version == "" {
		return fmt.Errorf("lock file version cannot be empty")
	}
	if l.Registry == "" {
		return fmt.Errorf("lock file registry cannot be empty")
	}
	if l.Tools == nil {
		return fmt.Errorf("lock file tools cannot be nil")
	}

	// Validate all installed tools
	for name, tool := range l.Tools {
		if name == "" {
			return fmt.Errorf("installed tool name cannot be empty")
		}
		if err := tool.Validate(); err != nil {
			return fmt.Errorf("invalid installed tool %s: %w", name, err)
		}
	}

	return nil
}

// AddTool adds a tool to the lock file
func (l *LockFile) AddTool(name string, tool *InstalledTool) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}
	if err := tool.Validate(); err != nil {
		return fmt.Errorf("invalid tool: %w", err)
	}

	if l.Tools == nil {
		l.Tools = make(map[string]*InstalledTool)
	}

	l.Tools[name] = tool
	l.UpdatedAt = time.Now()

	return nil
}

// RemoveTool removes a tool from the lock file
func (l *LockFile) RemoveTool(name string) error {
	if name == "" {
		return fmt.Errorf("tool name cannot be empty")
	}

	if _, exists := l.Tools[name]; !exists {
		return fmt.Errorf("tool %s not found in lock file", name)
	}

	delete(l.Tools, name)
	l.UpdatedAt = time.Now()

	return nil
}

// GetTool retrieves a tool from the lock file
func (l *LockFile) GetTool(name string) (*InstalledTool, error) {
	if name == "" {
		return nil, fmt.Errorf("tool name cannot be empty")
	}

	tool, exists := l.Tools[name]
	if !exists {
		return nil, fmt.Errorf("tool %s not found in lock file", name)
	}

	return tool, nil
}

// ToolMetadata represents additional metadata for a tool
type ToolMetadata struct {
	Author       string            `json:"author,omitempty" yaml:"author,omitempty"`
	Tags         []string          `json:"tags,omitempty" yaml:"tags,omitempty"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	Version      string            `json:"version,omitempty" yaml:"version,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Changelog    map[string]string `json:"changelog,omitempty" yaml:"changelog,omitempty"`
	Custom       map[string]string `json:"custom,omitempty" yaml:"custom,omitempty"`
}

// SearchFilter represents filter criteria for searching tools
type SearchFilter struct {
	Query         string   `json:"query"`
	Type          ToolType `json:"type,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Author        string   `json:"author,omitempty"`
	MinDownloads  int      `json:"min_downloads,omitempty"`
	Regex         bool     `json:"regex"`
	CaseSensitive bool     `json:"case_sensitive"`
}

// Validate checks if SearchFilter is valid
func (s *SearchFilter) Validate() error {
	if s.Query == "" {
		return fmt.Errorf("search query cannot be empty")
	}
	if s.Type != "" {
		if err := s.Type.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// SortField represents the field to sort by
type SortField string

const (
	SortByName      SortField = "name"
	SortByCreated   SortField = "created"
	SortByUpdated   SortField = "updated"
	SortByDownloads SortField = "downloads"
)

// ListFilter represents filter criteria for listing tools
type ListFilter struct {
	Type     ToolType  `json:"type,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Author   string    `json:"author,omitempty"`
	SortBy   SortField `json:"sort_by,omitempty"`
	SortDesc bool      `json:"sort_desc"`
	Limit    int       `json:"limit,omitempty"`
}

// Validate checks if ListFilter is valid
func (l *ListFilter) Validate() error {
	if l.Type != "" {
		if err := l.Type.Validate(); err != nil {
			return err
		}
	}
	if l.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	return nil
}

// Config represents the application configuration
type Config struct {
	Registry RegistryConfig `yaml:"registry"`
	Local    LocalConfig    `yaml:"local"`
	Publish  PublishConfig  `yaml:"publish"`
}

// RegistryConfig represents registry-specific configuration
type RegistryConfig struct {
	URL       string `yaml:"url"`
	Branch    string `yaml:"branch"`
	AuthToken string `yaml:"auth_token"`
}

// LocalConfig represents local configuration
type LocalConfig struct {
	DefaultPath         string `yaml:"default_path"`
	AutoUpdateCheck     bool   `yaml:"auto_update_check"`
	UpdateCheckInterval int    `yaml:"update_check_interval"` // seconds
}

// PublishConfig represents publishing configuration
type PublishConfig struct {
	DefaultAuthor   string `yaml:"default_author"`
	AutoVersionBump string `yaml:"auto_version_bump"` // patch, minor, major
	CreatePR        bool   `yaml:"create_pr"`
}

// Validate checks if Config is valid
func (c *Config) Validate() error {
	if c.Registry.URL == "" {
		return fmt.Errorf("registry URL cannot be empty")
	}
	if c.Registry.Branch == "" {
		return fmt.Errorf("registry branch cannot be empty")
	}
	if c.Local.DefaultPath == "" {
		return fmt.Errorf("default path cannot be empty")
	}
	if c.Local.UpdateCheckInterval < 0 {
		return fmt.Errorf("update check interval cannot be negative")
	}
	return nil
}

// NewDefaultConfig creates a new Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		Registry: RegistryConfig{
			URL:    "https://github.com/nghiadt/claude-tools-registry",
			Branch: "main",
		},
		Local: LocalConfig{
			DefaultPath:         ".claude",
			AutoUpdateCheck:     true,
			UpdateCheckInterval: 86400, // 24 hours
		},
		Publish: PublishConfig{
			AutoVersionBump: "patch",
			CreatePR:        true,
		},
	}
}
