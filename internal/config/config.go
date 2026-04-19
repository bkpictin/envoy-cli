package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const filename = ".envoy.yaml"

// DefaultPath returns the default config file path in the current directory.
func DefaultPath() string {
	return filepath.Join(".", filename)
}

// Config holds all targets, their env vars, and snapshots.
type Config struct {
	Targets   map[string]map[string]string       `yaml:"targets"`
	Snapshots map[string][]interface{}            `yaml:"snapshots,omitempty"`
}

// Init creates a new empty config file at path.
func Init(path string) error {
	cfg := &Config{
		Targets: make(map[string]map[string]string),
	}
	return Save(path, cfg)
}

// Load reads and parses the config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Targets == nil {
		cfg.Targets = make(map[string]map[string]string)
	}
	return &cfg, nil
}

// Save writes the config to disk at path.
func Save(path string, cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
