package services

import (
	"errors"
	"testing"
	"time"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRegistryServiceInterface is a mock for testing
type MockRegistryServiceInterface struct {
	mock.Mock
}

func (m *MockRegistryServiceInterface) GetTool(name string, toolType models.ToolType) (*models.ToolInfo, error) {
	args := m.Called(name, toolType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ToolInfo), args.Error(1)
}

func (m *MockRegistryServiceInterface) GetRegistry() (*models.Registry, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Registry), args.Error(1)
}

// MockLockFileServiceInterface is a mock for testing
type MockLockFileServiceInterface struct {
	mock.Mock
}

func (m *MockLockFileServiceInterface) GetTool(name string) (*models.InstalledTool, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.InstalledTool), args.Error(1)
}

func (m *MockLockFileServiceInterface) AddTool(name string, tool *models.InstalledTool) error {
	args := m.Called(name, tool)
	return args.Error(0)
}

func (m *MockLockFileServiceInterface) RemoveTool(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockLockFileServiceInterface) ListTools() (map[string]*models.InstalledTool, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]*models.InstalledTool), args.Error(1)
}

func (m *MockLockFileServiceInterface) IsInstalled(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func (m *MockLockFileServiceInterface) GetRegistry() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockLockFileServiceInterface) SetRegistry(url string) error {
	args := m.Called(url)
	return args.Error(0)
}

// MockInstallerService is a mock for testing
type MockInstallerService struct {
	mock.Mock
}

func (m *MockInstallerService) InstallWithVersion(toolName, version string) error {
	args := m.Called(toolName, version)
	return args.Error(0)
}

func TestNewUpdaterService(t *testing.T) {
	tests := []struct {
		name             string
		registryService  RegistryServiceInterface
		lockFileService  LockFileServiceInterface
		installerService *InstallerService
		wantErr          bool
		errMsg           string
	}{
		{
			name:             "success",
			registryService:  &MockRegistryServiceInterface{},
			lockFileService:  &MockLockFileServiceInterface{},
			installerService: &InstallerService{},
			wantErr:          false,
		},
		{
			name:             "nil registry service",
			registryService:  nil,
			lockFileService:  &MockLockFileServiceInterface{},
			installerService: &InstallerService{},
			wantErr:          true,
			errMsg:           "registry service cannot be nil",
		},
		{
			name:             "nil lock file service",
			registryService:  &MockRegistryServiceInterface{},
			lockFileService:  nil,
			installerService: &InstallerService{},
			wantErr:          true,
			errMsg:           "lock file service cannot be nil",
		},
		{
			name:             "nil installer service",
			registryService:  &MockRegistryServiceInterface{},
			lockFileService:  &MockLockFileServiceInterface{},
			installerService: nil,
			wantErr:          true,
			errMsg:           "installer service cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, err := NewUpdaterService(tt.registryService, tt.lockFileService, tt.installerService)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, svc)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, svc)
			}
		})
	}
}

func TestUpdaterService_CompareVersions(t *testing.T) {
	// Create service
	svc := &UpdaterService{}

	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{
			name: "v1 < v2",
			v1:   "1.0.0",
			v2:   "2.0.0",
			want: -1,
		},
		{
			name: "v1 = v2",
			v1:   "1.0.0",
			v2:   "1.0.0",
			want: 0,
		},
		{
			name: "v1 > v2",
			v1:   "2.0.0",
			v2:   "1.0.0",
			want: 1,
		},
		{
			name: "with v prefix",
			v1:   "v1.0.0",
			v2:   "v2.0.0",
			want: -1,
		},
		{
			name: "mixed prefix",
			v1:   "1.0.0",
			v2:   "v2.0.0",
			want: -1,
		},
		{
			name: "patch version",
			v1:   "1.0.0",
			v2:   "1.0.1",
			want: -1,
		},
		{
			name: "minor version",
			v1:   "1.1.0",
			v2:   "1.2.0",
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.CompareVersions(tt.v1, tt.v2)
			assert.Equal(t, tt.want, got, "CompareVersions(%s, %s) = %d, want %d", tt.v1, tt.v2, got, tt.want)
		})
	}
}

func TestUpdaterService_CheckOutdated(t *testing.T) {
	tests := []struct {
		name           string
		installedTools map[string]*models.InstalledTool
		registry       *models.Registry
		listToolsErr   error
		getRegistryErr error
		wantOutdated   int
		wantErr        bool
	}{
		{
			name: "one outdated tool",
			installedTools: map[string]*models.InstalledTool{
				"code-reviewer": {
					Version:     "1.0.0",
					Type:        models.ToolTypeAgent,
					InstalledAt: time.Now(),
				},
			},
			registry: &models.Registry{
				Version: "1.0",
				Tools: map[models.ToolType][]*models.ToolInfo{
					models.ToolTypeAgent: {
						{
							Name:    "code-reviewer",
							Version: "2.0.0",
							Type:    models.ToolTypeAgent,
						},
					},
				},
			},
			wantOutdated: 1,
			wantErr:      false,
		},
		{
			name: "all up to date",
			installedTools: map[string]*models.InstalledTool{
				"code-reviewer": {
					Version:     "2.0.0",
					Type:        models.ToolTypeAgent,
					InstalledAt: time.Now(),
				},
			},
			registry: &models.Registry{
				Version: "1.0",
				Tools: map[models.ToolType][]*models.ToolInfo{
					models.ToolTypeAgent: {
						{
							Name:    "code-reviewer",
							Version: "2.0.0",
							Type:    models.ToolTypeAgent,
						},
					},
				},
			},
			wantOutdated: 0,
			wantErr:      false,
		},
		{
			name:           "no installed tools",
			installedTools: map[string]*models.InstalledTool{},
			registry: &models.Registry{
				Version: "1.0",
				Tools:   map[models.ToolType][]*models.ToolInfo{},
			},
			wantOutdated: 0,
			wantErr:      false,
		},
		{
			name:           "list tools error",
			installedTools: nil,
			listToolsErr:   errors.New("failed to read lock file"),
			wantErr:        true,
		},
		{
			name: "get registry error",
			installedTools: map[string]*models.InstalledTool{
				"code-reviewer": {
					Version: "1.0.0",
					Type:    models.ToolTypeAgent,
				},
			},
			getRegistryErr: errors.New("network error"),
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistry := new(MockRegistryServiceInterface)
			mockLockFile := new(MockLockFileServiceInterface)
			mockInstaller := &InstallerService{}

			// Set up mocks
			mockLockFile.On("ListTools").Return(tt.installedTools, tt.listToolsErr)
			if tt.listToolsErr == nil && len(tt.installedTools) > 0 {
				mockRegistry.On("GetRegistry").Return(tt.registry, tt.getRegistryErr)
			}

			svc, err := NewUpdaterService(mockRegistry, mockLockFile, mockInstaller)
			assert.NoError(t, err)

			outdated, err := svc.CheckOutdated()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, outdated, tt.wantOutdated)
			}

			mockRegistry.AssertExpectations(t)
			mockLockFile.AssertExpectations(t)
		})
	}
}

func TestUpdaterService_Update(t *testing.T) {
	tests := []struct {
		name          string
		toolName      string
		installedTool *models.InstalledTool
		latestTool    *models.ToolInfo
		getToolErr    error
		getLatestErr  error
		wantSkipped   bool
		wantSuccess   bool
		wantErr       bool
	}{
		{
			name:     "already up to date",
			toolName: "code-reviewer",
			installedTool: &models.InstalledTool{
				Version: "2.0.0",
				Type:    models.ToolTypeAgent,
			},
			latestTool: &models.ToolInfo{
				Name:    "code-reviewer",
				Version: "2.0.0",
				Type:    models.ToolTypeAgent,
			},
			wantSuccess: true,
			wantSkipped: true,
			wantErr:     false,
		},
		{
			name:        "tool not installed",
			toolName:    "missing-tool",
			getToolErr:  errors.New("tool not found"),
			wantSuccess: false,
			wantErr:     true,
		},
		{
			name:     "tool not in registry",
			toolName: "code-reviewer",
			installedTool: &models.InstalledTool{
				Version: "1.0.0",
				Type:    models.ToolTypeAgent,
			},
			getLatestErr: errors.New("tool not found in registry"),
			wantSuccess:  false,
			wantErr:      true,
		},
		{
			name:     "empty tool name",
			toolName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistry := new(MockRegistryServiceInterface)
			mockLockFile := new(MockLockFileServiceInterface)
			// Create a dummy installer - we won't actually call it in these tests
			realInstaller := &InstallerService{}

			// Set up mocks
			if tt.toolName != "" {
				mockLockFile.On("GetTool", tt.toolName).Return(tt.installedTool, tt.getToolErr)

				if tt.getToolErr == nil && tt.installedTool != nil {
					mockRegistry.On("GetTool", tt.toolName, tt.installedTool.Type).
						Return(tt.latestTool, tt.getLatestErr)
				}
			}

			svc, err := NewUpdaterService(mockRegistry, mockLockFile, realInstaller)
			assert.NoError(t, err)

			result, err := svc.Update(tt.toolName)

			if tt.wantErr {
				assert.Error(t, err)
				if result != nil {
					assert.False(t, result.Success)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.wantSuccess, result.Success)
				assert.Equal(t, tt.wantSkipped, result.Skipped)
			}

			mockRegistry.AssertExpectations(t)
			mockLockFile.AssertExpectations(t)
		})
	}
}

func TestUpdaterService_IsOutdated(t *testing.T) {
	tests := []struct {
		name          string
		toolName      string
		installedTool *models.InstalledTool
		latestTool    *models.ToolInfo
		getToolErr    error
		getLatestErr  error
		wantOutdated  bool
		wantErr       bool
	}{
		{
			name:     "tool is outdated",
			toolName: "code-reviewer",
			installedTool: &models.InstalledTool{
				Version: "1.0.0",
				Type:    models.ToolTypeAgent,
			},
			latestTool: &models.ToolInfo{
				Name:    "code-reviewer",
				Version: "2.0.0",
				Type:    models.ToolTypeAgent,
			},
			wantOutdated: true,
			wantErr:      false,
		},
		{
			name:     "tool is up to date",
			toolName: "code-reviewer",
			installedTool: &models.InstalledTool{
				Version: "2.0.0",
				Type:    models.ToolTypeAgent,
			},
			latestTool: &models.ToolInfo{
				Name:    "code-reviewer",
				Version: "2.0.0",
				Type:    models.ToolTypeAgent,
			},
			wantOutdated: false,
			wantErr:      false,
		},
		{
			name:         "empty tool name",
			toolName:     "",
			wantOutdated: false,
			wantErr:      true,
		},
		{
			name:       "tool not installed",
			toolName:   "missing-tool",
			getToolErr: errors.New("tool not found"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistry := new(MockRegistryServiceInterface)
			mockLockFile := new(MockLockFileServiceInterface)
			mockInstaller := &InstallerService{}

			// Set up mocks
			if tt.toolName != "" {
				mockLockFile.On("GetTool", tt.toolName).Return(tt.installedTool, tt.getToolErr)

				if tt.getToolErr == nil && tt.installedTool != nil {
					mockRegistry.On("GetTool", tt.toolName, tt.installedTool.Type).
						Return(tt.latestTool, tt.getLatestErr)
				}
			}

			svc, err := NewUpdaterService(mockRegistry, mockLockFile, mockInstaller)
			assert.NoError(t, err)

			outdated, err := svc.IsOutdated(tt.toolName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantOutdated, outdated)
			}

			mockRegistry.AssertExpectations(t)
			mockLockFile.AssertExpectations(t)
		})
	}
}

func TestUpdaterService_GetOutdatedCount(t *testing.T) {
	tests := []struct {
		name           string
		installedTools map[string]*models.InstalledTool
		registry       *models.Registry
		wantCount      int
		wantErr        bool
	}{
		{
			name: "two outdated tools",
			installedTools: map[string]*models.InstalledTool{
				"code-reviewer": {
					Version: "1.0.0",
					Type:    models.ToolTypeAgent,
				},
				"git-helper": {
					Version: "1.0.0",
					Type:    models.ToolTypeAgent,
				},
			},
			registry: &models.Registry{
				Version: "1.0",
				Tools: map[models.ToolType][]*models.ToolInfo{
					models.ToolTypeAgent: {
						{Name: "code-reviewer", Version: "2.0.0", Type: models.ToolTypeAgent},
						{Name: "git-helper", Version: "2.0.0", Type: models.ToolTypeAgent},
					},
				},
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "no outdated tools",
			installedTools: map[string]*models.InstalledTool{
				"code-reviewer": {
					Version: "2.0.0",
					Type:    models.ToolTypeAgent,
				},
			},
			registry: &models.Registry{
				Version: "1.0",
				Tools: map[models.ToolType][]*models.ToolInfo{
					models.ToolTypeAgent: {
						{Name: "code-reviewer", Version: "2.0.0", Type: models.ToolTypeAgent},
					},
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistry := new(MockRegistryServiceInterface)
			mockLockFile := new(MockLockFileServiceInterface)
			mockInstaller := &InstallerService{}

			mockLockFile.On("ListTools").Return(tt.installedTools, nil)
			if len(tt.installedTools) > 0 {
				mockRegistry.On("GetRegistry").Return(tt.registry, nil)
			}

			svc, err := NewUpdaterService(mockRegistry, mockLockFile, mockInstaller)
			assert.NoError(t, err)

			count, err := svc.GetOutdatedCount()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, count)
			}

			mockRegistry.AssertExpectations(t)
			mockLockFile.AssertExpectations(t)
		})
	}
}
