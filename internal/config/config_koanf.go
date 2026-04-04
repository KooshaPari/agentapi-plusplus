// Package config provides configuration management for AgentAPI using koanf.
// Migration from viper: replaced spf13/viper with knadh/koanf/v2 for better
// type safety and concurrent map access handling.
package config

import (
	"os"

	"github.com/knadh/koanf/v2"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/sops"
	"github.com/mitchellh/copystructure"
)

// ServerConfig represents the server configuration for AgentAPI.
type ServerConfig struct {
	Port           int      `koanf:"port"`
	Host           string   `koanf:"host"`
	ChatBasePath   string   `koanf:"chat_base_path"`
	AllowedHosts   []string `koanf:"allowed_hosts"`
	AllowedOrigins []string `koanf:"allowed_origins"`
	TermWidth      uint16   `koanf:"term_width"`
	TermHeight     uint16   `koanf:"term_height"`
	PrintOpenAPI   bool     `koanf:"print_openapi"`
}

// AgentAPIConfig represents the complete configuration for AgentAPI.
type AgentAPIConfig struct {
	Server ServerConfig `koanf:"server"`
	Agent  AgentConfig  `koanf:"agent"`
}

// AgentConfig represents agent-related configuration.
type AgentConfig struct {
	Type          string `koanf:"type"`
	InitialPrompt string `koanf:"initial_prompt"`
}

// defaults returns the default configuration map using koanf's confmap provider.
func defaults() *koanf.Koanf {
	k := koanf.NewWithConfmapProvider(".", confmap.Provider(map[string]any{
		"server.port":             3284,
		"server.host":             "localhost",
		"server.chat_base_path":   "/chat",
		"server.allowed_hosts":    []string{"localhost", "127.0.0.1", "[::1]"},
		"server.allowed_origins":  []string{"http://localhost:3284", "http://localhost:3000", "http://localhost:3001"},
		"server.term_width":       uint16(80),
		"server.term_height":      uint16(1000),
		"server.print_openapi":    false,
		"agent.type":              "",
		"agent.initial_prompt":    "",
	}, "."))
	return k
}

// LoadConfig loads the configuration from a file and environment variables.
// Uses koanf for type-safe configuration loading.
func LoadConfig(filePath string) (*AgentAPIConfig, error) {
	// Start with defaults
	k := defaults()

	// Load from file if provided and exists
	if filePath != "" {
		if _, statErr := os.Stat(filePath); statErr == nil {
			// Detect file type and load accordingly
			if err := k.Load(file.Provider(filePath, getFileExtension(filePath)), getParser(filePath), koanf.WithMergeFunc(file.Merge)); err != nil {
				return nil, err
			}
		}
	}

	// Load from environment variables (AGENTAPI_* prefix)
	// These take precedence over file values
	if err := k.Load(env.Provider("AGENTAPI_", ".", func(s string) string {
		// Convert AGENTAPI_SERVER_PORT -> server.port
		// Convert AGENTAPI_ALLOWED_HOSTS -> server.allowed_hosts
		return s
	}), koanf.WithMergeFunc(env.Merge), koanf.EnvOverride()); err != nil {
		return nil, err
	}

	// Unmarshal into config struct
	var cfg AgentAPIConfig
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// LoadConfigWithEnv loads configuration from environment variables and a config file.
// Environment variables take precedence over config file values.
func LoadConfigWithEnv(filePath string) (*AgentAPIConfig, error) {
	// Start with defaults
	k := defaults()

	// Load from file first
	if filePath != "" {
		if _, statErr := os.Stat(filePath); statErr == nil {
			if err := k.Load(file.Provider(filePath, getFileExtension(filePath)), getParser(filePath), koanf.WithMergeFunc(file.Merge)); err != nil {
				return nil, err
			}
		}
	}

	// Then load from environment to override
	if err := k.Load(env.Provider("AGENTAPI_", ".", nil), koanf.WithMergeFunc(env.Merge), koanf.EnvOverride()); err != nil {
		return nil, err
	}

	// Unmarshal
	var cfg AgentAPIConfig
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// BindEnvVars binds specific environment variables to configuration keys.
// With koanf, environment variables are loaded with the env.Provider.
// This function is kept for backwards compatibility.
func BindEnvVars() error {
	// No-op: koanf handles env vars automatically via env.Provider
	// Kept for API compatibility
	return nil
}

// getFileExtension returns the file extension for parsing
func getFileExtension(path string) string {
	switch {
	case contains(path, ".yaml") || contains(path, ".yml"):
		return "yaml"
	case contains(path, ".json"):
		return "json"
	case contains(path, ".toml"):
		return "toml"
	case contains(path, ".sops"):
		return "sops"
	default:
		return "yaml" // default
	}
}

// getParser returns the appropriate parser for the config file
func getParser(path string) koanf.Parser {
	ext := getFileExtension(path)
	switch ext {
	case "yaml":
		return koanf.YAML()
	case "json":
		return koanf.JSON()
	case "toml":
		return koanf.TOML()
	case "sops":
		return sops.Provider()
	default:
		return koanf.YAML()
	}
}

// contains is a simple string contains check
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Copy creates a deep copy of the config (needed for copystructure interface)
func (c *AgentAPIConfig) Copy() (*AgentAPIConfig, error) {
	data, err := copystructure.Copy(*c)
	if err != nil {
		return nil, err
	}
	result := data.(AgentAPIConfig)
	return &result, nil
}