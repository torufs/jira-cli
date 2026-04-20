package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFile = ".jira-cli.yml"
	EnvConfigPath     = "JIRA_CLI_CONFIG"
)

// Config holds the application configuration.
type Config struct {
	Server   string `yaml:"server"`
	Login    string `yaml:"login"`
	Project  string `yaml:"project"`
	Board    string `yaml:"board"`
	AuthType string `yaml:"auth_type"`
	Token    string `yaml:"token"`
}

// Load reads configuration from a YAML file.
// It checks the environment variable JIRA_CLI_CONFIG first,
// then falls back to ~/.jira-cli.yml.
func Load() (*Config, error) {
	path := os.Getenv(EnvConfigPath)
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(home, DefaultConfigFile)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save writes the configuration to the given file path.
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// Validate checks that required fields are present.
func (c *Config) Validate() error {
	if c.Server == "" {
		return errors.New("server is required")
	}
	if c.Login == "" {
		return errors.New("login is required")
	}
	if c.Token == "" {
		return errors.New("token is required")
	}
	return nil
}
