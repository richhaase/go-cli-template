// Package config handles configuration loading from files,
// environment variables, and command-line flags.
package config

import (
	"os"
	"path/filepath"
)

// Config holds the application configuration.
type Config struct {
	// Add your configuration fields here
	Debug bool `yaml:"debug"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Debug: false,
	}
}

// Load loads configuration with the following precedence:
// 1. Command-line flags (highest)
// 2. Environment variables
// 3. Config file
// 4. Defaults (lowest)
func Load() (*Config, error) {
	cfg := DefaultConfig()

	// Try to load from config file
	if configPath := findConfigFile(); configPath != "" {
		// TODO: Load and parse YAML config file
		_ = configPath
	}

	// Override with environment variables
	if os.Getenv("MYCLI_DEBUG") == "true" {
		cfg.Debug = true
	}

	return cfg, nil
}

// findConfigFile looks for config in standard locations.
func findConfigFile() string {
	// Check current directory
	if _, err := os.Stat(".mycli.yaml"); err == nil {
		return ".mycli.yaml"
	}

	if dir := configDir(); dir != "" {
		configPath := filepath.Join(dir, "mycli", "config.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	return ""
}

// configDir returns the XDG-style config directory on every platform:
// $XDG_CONFIG_HOME if set, else ~/.config. CLI tools live in ~/.config
// by convention, macOS included — os.UserConfigDir's
// ~/Library/Application Support is for GUI apps.
func configDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir
	}
	if home, err := os.UserHomeDir(); err == nil {
		return filepath.Join(home, ".config")
	}
	return ""
}
