package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ankitpokhrel/jira-cli/internal/config"
)

func TestLoadFromEnvPath(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "test-config.yml")

	yamlContent := []byte("server: https://jira.example.com\nlogin: user@example.com\ntoken: secret123\nproject: PROJ\n")
	require.NoError(t, os.WriteFile(cfgPath, yamlContent, 0600))

	t.Setenv(config.EnvConfigPath, cfgPath)

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "https://jira.example.com", cfg.Server)
	assert.Equal(t, "user@example.com", cfg.Login)
	assert.Equal(t, "secret123", cfg.Token)
	assert.Equal(t, "PROJ", cfg.Project)
}

func TestLoadMissingFileReturnsEmpty(t *testing.T) {
	t.Setenv(config.EnvConfigPath, "/nonexistent/path/.jira-cli.yml")

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Empty(t, cfg.Server)
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "saved-config.yml")

	original := &config.Config{
		Server:   "https://jira.example.com",
		Login:    "admin@example.com",
		Token:    "tok123",
		Project:  "DEMO",
		AuthType: "bearer",
	}

	require.NoError(t, config.Save(original, cfgPath))

	t.Setenv(config.EnvConfigPath, cfgPath)
	loaded, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, original.Server, loaded.Server)
	assert.Equal(t, original.Login, loaded.Login)
	assert.Equal(t, original.Token, loaded.Token)
}

func TestValidate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     config.Config
		wantErr string
	}{
		{"missing server", config.Config{Login: "a", Token: "b"}, "server is required"},
		{"missing login", config.Config{Server: "s", Token: "b"}, "login is required"},
		{"missing token", config.Config{Server: "s", Login: "a"}, "token is required"},
		{"valid", config.Config{Server: "s", Login: "a", Token: "b"}, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cfg.Validate()
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErr)
			}
		})
	}
}
