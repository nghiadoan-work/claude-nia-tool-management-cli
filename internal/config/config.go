package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nghiadoan-work/claude-nia-tool-management-cli/pkg/models"
	"gopkg.in/yaml.v3"
)

// ConfigService handles configuration loading and management
type ConfigService struct {
	config *models.Config
}

// NewConfigService creates a new ConfigService with the provided config
func NewConfigService(config *models.Config) *ConfigService {
	return &ConfigService{
		config: config,
	}
}

// GetConfig returns the current configuration
func (cs *ConfigService) GetConfig() *models.Config {
	return cs.config
}

// LoadConfig loads configuration with the following precedence:
// 1. Project config (.claude-tools-config.yaml in current directory) - highest priority
// 2. Global config (~/.claude-tools-config.yaml)
// 3. Default config - lowest priority
//
// Project-level config overrides global config for per-project customization.
func LoadConfig(configPath string) (*models.Config, error) {
	// Start with default config
	config := models.NewDefaultConfig()

	// Try loading global config first
	if err := loadGlobalConfig(config); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load global config: %w", err)
	}

	// Try loading project config (overrides global)
	if err := loadProjectConfig(config); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load project config: %w", err)
	}

	// If a specific config path is provided, load it (highest priority)
	if configPath != "" {
		if err := loadConfigFromFile(config, configPath); err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	}

	// Validate final config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// loadGlobalConfig loads config from ~/.claude-tools-config.yaml
func loadGlobalConfig(config *models.Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	globalPath := filepath.Join(homeDir, ".claude-tools-config.yaml")
	return loadConfigFromFile(config, globalPath)
}

// loadProjectConfig loads config from .claude-tools-config.yaml in current directory
func loadProjectConfig(config *models.Config) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	projectPath := filepath.Join(currentDir, ".claude-tools-config.yaml")
	return loadConfigFromFile(config, projectPath)
}

// loadConfigFromFile loads and merges config from a YAML file
func loadConfigFromFile(config *models.Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var fileConfig models.Config
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Merge with existing config (non-empty values override)
	mergeConfig(config, &fileConfig)

	return nil
}

// mergeConfig merges source config into target config (non-empty values from source override target)
func mergeConfig(target, source *models.Config) {
	// Registry config
	if source.Registry.URL != "" {
		target.Registry.URL = source.Registry.URL
	}
	if source.Registry.Branch != "" {
		target.Registry.Branch = source.Registry.Branch
	}
	if source.Registry.AuthToken != "" {
		target.Registry.AuthToken = source.Registry.AuthToken
	}

	// Local config
	if source.Local.DefaultPath != "" {
		target.Local.DefaultPath = source.Local.DefaultPath
	}
	// Only override bool if explicitly set (check if different from default)
	if source.Local.AutoUpdateCheck != target.Local.AutoUpdateCheck {
		target.Local.AutoUpdateCheck = source.Local.AutoUpdateCheck
	}
	if source.Local.UpdateCheckInterval > 0 {
		target.Local.UpdateCheckInterval = source.Local.UpdateCheckInterval
	}

	// Publish config
	if source.Publish.DefaultAuthor != "" {
		target.Publish.DefaultAuthor = source.Publish.DefaultAuthor
	}
	if source.Publish.AutoVersionBump != "" {
		target.Publish.AutoVersionBump = source.Publish.AutoVersionBump
	}
	if source.Publish.CreatePR != target.Publish.CreatePR {
		target.Publish.CreatePR = source.Publish.CreatePR
	}
}

// SaveConfig saves the config to a YAML file
func SaveConfig(config *models.Config, path string) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetGlobalConfigPath returns the path to the global config file
func GetGlobalConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".claude-tools-config.yaml"), nil
}

// GetProjectConfigPath returns the path to the project config file
func GetProjectConfigPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(currentDir, ".claude-tools-config.yaml"), nil
}
