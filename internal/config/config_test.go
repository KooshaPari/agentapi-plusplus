// Package config_test provides tests for the config package.
package config_test

import (
	"os"
	"testing"

	"github.com/coder/agentapi/internal/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_WithoutFile tests loading config without a config file.
func TestLoadConfig_WithoutFile(t *testing.T) {
	// Clear viper for clean test
	viper.Reset()

	cfg, err := config.LoadConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Check defaults
	assert.Equal(t, 3284, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "/chat", cfg.Server.ChatBasePath)
	assert.False(t, cfg.Server.PrintOpenAPI)
	assert.Equal(t, uint16(80), cfg.Server.TermWidth)
	assert.Equal(t, uint16(1000), cfg.Server.TermHeight)
}

// TestLoadConfig_WithMissingFile tests loading config with a non-existent file.
func TestLoadConfig_WithMissingFile(t *testing.T) {
	// Clear viper for clean test
	viper.Reset()

	cfg, err := config.LoadConfig("/nonexistent/file.yml")
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

// TestLoadConfig_WithValidFile tests loading config with a valid file.
func TestLoadConfig_WithValidFile(t *testing.T) {
	// Clear viper for clean test
	viper.Reset()

	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "agentapi-config-*.yml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
server:
  port: 8080
  host: 0.0.0.0
  chat_base_path: /api/chat
  term_width: 120
  term_height: 2000
agent:
  type: claude
  initial_prompt: "Hello"
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	cfg, err := config.LoadConfig(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "/api/chat", cfg.Server.ChatBasePath)
	assert.Equal(t, uint16(120), cfg.Server.TermWidth)
	assert.Equal(t, uint16(2000), cfg.Server.TermHeight)
	assert.Equal(t, "claude", cfg.Agent.Type)
	assert.Equal(t, "Hello", cfg.Agent.InitialPrompt)
}

// TestLoadConfigWithEnv_EnvironmentOverrides tests that environment variables override config file values.
func TestLoadConfigWithEnv_EnvironmentOverrides(t *testing.T) {
	// Clear viper for clean test
	viper.Reset()

	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "agentapi-config-*.yml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	configContent := `
server:
  port: 8080
  host: 0.0.0.0
`
	_, err = tmpFile.WriteString(configContent)
	require.NoError(t, err)
	tmpFile.Close()

	// Set environment variable
	t.Setenv("AGENTAPI_PORT", "9999")

	cfg, err := config.LoadConfigWithEnv(tmpFile.Name())
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Environment variable should override config file
	assert.Equal(t, 9999, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
}

// TestServerConfig tests ServerConfig struct.
func TestServerConfig(t *testing.T) {
	cfg := config.ServerConfig{
		Port:         8080,
		Host:         "localhost",
		ChatBasePath: "/api/chat",
		TermWidth:    120,
		TermHeight:   1000,
	}

	assert.Equal(t, 8080, cfg.Port)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "/api/chat", cfg.ChatBasePath)
	assert.Equal(t, uint16(120), cfg.TermWidth)
	assert.Equal(t, uint16(1000), cfg.TermHeight)
}

// TestAgentConfig tests AgentConfig struct.
func TestAgentConfig(t *testing.T) {
	cfg := config.AgentConfig{
		Type:          "claude",
		InitialPrompt: "Hello world",
	}

	assert.Equal(t, "claude", cfg.Type)
	assert.Equal(t, "Hello world", cfg.InitialPrompt)
}

// TestAgentAPIConfig tests AgentAPIConfig struct.
func TestAgentAPIConfig(t *testing.T) {
	cfg := config.AgentAPIConfig{
		Server: config.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Agent: config.AgentConfig{
			Type: "claude",
		},
	}

	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "claude", cfg.Agent.Type)
}

// TestLoadConfig_DefaultValues tests that default values are set correctly.
func TestLoadConfig_DefaultValues(t *testing.T) {
	// Clear viper for clean test
	viper.Reset()

	cfg, err := config.LoadConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify all defaults
	assert.Equal(t, 3284, cfg.Server.Port)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, "/chat", cfg.Server.ChatBasePath)
	assert.False(t, cfg.Server.PrintOpenAPI)
	assert.Equal(t, uint16(80), cfg.Server.TermWidth)
	assert.Equal(t, uint16(1000), cfg.Server.TermHeight)

	// Check allowed hosts
	assert.Len(t, cfg.Server.AllowedHosts, 3)
	assert.Contains(t, cfg.Server.AllowedHosts, "localhost")
	assert.Contains(t, cfg.Server.AllowedHosts, "127.0.0.1")

	// Check allowed origins
	assert.Len(t, cfg.Server.AllowedOrigins, 3)
	assert.Contains(t, cfg.Server.AllowedOrigins, "http://localhost:3284")
}

// TestBindEnvVars tests binding environment variables.
func TestBindEnvVars(t *testing.T) {
	// Clear viper for clean test
	viper.Reset()

	err := config.BindEnvVars()
	assert.NoError(t, err)

	// Set some environment variables
	t.Setenv("AGENTAPI_PORT", "5000")
	t.Setenv("AGENTAPI_HOST", "0.0.0.0")

	// Reload viper configuration
	viper.AutomaticEnv()
	port := viper.GetInt("server.port")
	host := viper.GetString("server.host")

	assert.Equal(t, 5000, port)
	assert.Equal(t, "0.0.0.0", host)
}
