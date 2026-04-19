package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const defaultConfigFile = ".envoy.yaml"

// Config represents the top-level envoy configuration.
type Config struct {
	Version  string              `yaml:"version"`
	Targets  map[string]Target   `yaml:"targets"`
}

// Target represents a deployment target with its own env vars.
type Target struct {
	Description string            `yaml:"description,omitempty"`
	Env         map[string]string `yaml:"env"`
}

// Load reads the config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to the given path.
func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directories: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// DefaultPath returns the default config file path in the current directory.
func DefaultPath() string {
	return defaultConfigFile
}

// Init creates a new empty config file at path if it does not exist.
func Init(path string) error {
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config file already exists: %s", path)
	}
	cfg := &Config{
		Version: "1",
		Targets: map[string]Target{},
	}
	return Save(path, cfg)
}
