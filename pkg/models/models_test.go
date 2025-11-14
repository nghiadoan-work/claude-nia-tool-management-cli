package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolType_Validate(t *testing.T) {
	tests := []struct {
		name    string
		toolType ToolType
		wantErr bool
	}{
		{"valid agent", ToolTypeAgent, false},
		{"valid command", ToolTypeCommand, false},
		{"valid skill", ToolTypeSkill, false},
		{"invalid type", ToolType("invalid"), true},
		{"empty type", ToolType(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.toolType.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestToolInfo_Validate(t *testing.T) {
	validTool := &ToolInfo{
		Name:    "test-agent",
		Version: "1.0.0",
		Type:    ToolTypeAgent,
		File:    "agents/test-agent/test-agent.zip",
	}

	tests := []struct {
		name    string
		tool    *ToolInfo
		wantErr bool
	}{
		{"valid tool", validTool, false},
		{"missing name", &ToolInfo{Version: "1.0.0", Type: ToolTypeAgent, File: "test.zip"}, true},
		{"missing version", &ToolInfo{Name: "test", Type: ToolTypeAgent, File: "test.zip"}, true},
		{"invalid type", &ToolInfo{Name: "test", Version: "1.0.0", Type: ToolType("invalid"), File: "test.zip"}, true},
		{"missing file", &ToolInfo{Name: "test", Version: "1.0.0", Type: ToolTypeAgent}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tool.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegistry_Validate(t *testing.T) {
	validRegistry := &Registry{
		Version: "1.0",
		Tools: map[ToolType][]*ToolInfo{
			ToolTypeAgent: {
				{Name: "agent1", Version: "1.0.0", Type: ToolTypeAgent, File: "test.zip"},
			},
		},
	}

	tests := []struct {
		name     string
		registry *Registry
		wantErr  bool
	}{
		{"valid registry", validRegistry, false},
		{"missing version", &Registry{Tools: map[ToolType][]*ToolInfo{}}, true},
		{"nil tools", &Registry{Version: "1.0"}, true},
		{"invalid tool type", &Registry{
			Version: "1.0",
			Tools: map[ToolType][]*ToolInfo{
				ToolType("invalid"): {},
			},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.registry.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegistry_GetTool(t *testing.T) {
	registry := &Registry{
		Version: "1.0",
		Tools: map[ToolType][]*ToolInfo{
			ToolTypeAgent: {
				{Name: "agent1", Version: "1.0.0", Type: ToolTypeAgent, File: "test.zip"},
				{Name: "agent2", Version: "1.0.0", Type: ToolTypeAgent, File: "test.zip"},
			},
		},
	}

	tests := []struct {
		name     string
		toolName string
		toolType ToolType
		wantErr  bool
	}{
		{"found", "agent1", ToolTypeAgent, false},
		{"not found", "agent3", ToolTypeAgent, true},
		{"wrong type", "agent1", ToolTypeCommand, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool, err := registry.GetTool(tt.toolName, tt.toolType)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tool)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tool)
				assert.Equal(t, tt.toolName, tool.Name)
			}
		})
	}
}

func TestInstalledTool_Validate(t *testing.T) {
	validTool := &InstalledTool{
		Version: "1.0.0",
		Type:    ToolTypeAgent,
		Source:  "registry",
	}

	tests := []struct {
		name    string
		tool    *InstalledTool
		wantErr bool
	}{
		{"valid tool", validTool, false},
		{"missing version", &InstalledTool{Type: ToolTypeAgent, Source: "registry"}, true},
		{"invalid type", &InstalledTool{Version: "1.0.0", Type: ToolType("invalid"), Source: "registry"}, true},
		{"missing source", &InstalledTool{Version: "1.0.0", Type: ToolTypeAgent}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tool.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLockFile_AddTool(t *testing.T) {
	lockFile := &LockFile{
		Version:  "1.0",
		Registry: "https://github.com/test/registry",
		Tools:    make(map[string]*InstalledTool),
	}

	tool := &InstalledTool{
		Version: "1.0.0",
		Type:    ToolTypeAgent,
		Source:  "registry",
	}

	err := lockFile.AddTool("test-agent", tool)
	assert.NoError(t, err)
	assert.Len(t, lockFile.Tools, 1)
	assert.Equal(t, tool, lockFile.Tools["test-agent"])
}

func TestLockFile_RemoveTool(t *testing.T) {
	lockFile := &LockFile{
		Version:  "1.0",
		Registry: "https://github.com/test/registry",
		Tools: map[string]*InstalledTool{
			"test-agent": {Version: "1.0.0", Type: ToolTypeAgent, Source: "registry"},
		},
	}

	// Remove existing tool
	err := lockFile.RemoveTool("test-agent")
	assert.NoError(t, err)
	assert.Len(t, lockFile.Tools, 0)

	// Try to remove non-existent tool
	err = lockFile.RemoveTool("non-existent")
	assert.Error(t, err)
}

func TestLockFile_GetTool(t *testing.T) {
	expectedTool := &InstalledTool{Version: "1.0.0", Type: ToolTypeAgent, Source: "registry"}
	lockFile := &LockFile{
		Version:  "1.0",
		Registry: "https://github.com/test/registry",
		Tools: map[string]*InstalledTool{
			"test-agent": expectedTool,
		},
	}

	// Get existing tool
	tool, err := lockFile.GetTool("test-agent")
	assert.NoError(t, err)
	assert.Equal(t, expectedTool, tool)

	// Get non-existent tool
	tool, err = lockFile.GetTool("non-existent")
	assert.Error(t, err)
	assert.Nil(t, tool)
}

func TestSearchFilter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		filter  *SearchFilter
		wantErr bool
	}{
		{"valid filter", &SearchFilter{Query: "test"}, false},
		{"valid with type", &SearchFilter{Query: "test", Type: ToolTypeAgent}, false},
		{"empty query", &SearchFilter{Query: ""}, true},
		{"invalid type", &SearchFilter{Query: "test", Type: ToolType("invalid")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListFilter_Validate(t *testing.T) {
	tests := []struct {
		name    string
		filter  *ListFilter
		wantErr bool
	}{
		{"valid filter", &ListFilter{}, false},
		{"valid with type", &ListFilter{Type: ToolTypeAgent}, false},
		{"valid with limit", &ListFilter{Limit: 10}, false},
		{"invalid type", &ListFilter{Type: ToolType("invalid")}, true},
		{"negative limit", &ListFilter{Limit: -1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.filter.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	validConfig := &Config{
		Registry: RegistryConfig{
			URL:    "https://github.com/test/registry",
			Branch: "main",
		},
		Local: LocalConfig{
			DefaultPath:         ".claude",
			UpdateCheckInterval: 86400,
		},
	}

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{"valid config", validConfig, false},
		{"missing registry URL", &Config{
			Registry: RegistryConfig{Branch: "main"},
			Local:    LocalConfig{DefaultPath: ".claude"},
		}, true},
		{"missing registry branch", &Config{
			Registry: RegistryConfig{URL: "https://github.com/test/registry"},
			Local:    LocalConfig{DefaultPath: ".claude"},
		}, true},
		{"missing default path", &Config{
			Registry: RegistryConfig{URL: "https://github.com/test/registry", Branch: "main"},
			Local:    LocalConfig{},
		}, true},
		{"negative interval", &Config{
			Registry: RegistryConfig{URL: "https://github.com/test/registry", Branch: "main"},
			Local:    LocalConfig{DefaultPath: ".claude", UpdateCheckInterval: -1},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()
	assert.NotNil(t, config)
	assert.NotEmpty(t, config.Registry.URL)
	assert.NotEmpty(t, config.Registry.Branch)
	assert.NotEmpty(t, config.Local.DefaultPath)
	assert.True(t, config.Local.AutoUpdateCheck)
	assert.Greater(t, config.Local.UpdateCheckInterval, 0)
	assert.True(t, config.Publish.CreatePR)

	// Validate default config
	err := config.Validate()
	assert.NoError(t, err)
}
