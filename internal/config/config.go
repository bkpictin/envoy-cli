package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all targets and their environment variable maps.
type Config struct {
	Targets map[string]map[string]string `json:"targets"`
}

// DefaultPath returns the default config file path.
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".envoy.json"
	}
	return filepath.Join(home, ".envoy.json")
}

// Init creates a new empty config file at path if it does not exist.
func Init(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	cfg := &Config{Targets: make(map[string]map[string]string)}
	return Save(cfg, path)
}

// Load reads and parses the config from path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Targets == nil {
		cfg.Targets = make(map[string]map[string]string)
	}
	return &cfg, nil
}

// Save serialises cfg to path with indented JSON.
func Save(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
