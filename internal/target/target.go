package target

import (
	"fmt"

	"github.com/envoy-cli/envoy-cli/internal/config"
)

// List returns all target names defined in the config.
func List(cfg *config.Config) []string {
	targets := make([]string, 0, len(cfg.Targets))
	for name := range cfg.Targets {
		targets = append(targets, name)
	}
	return targets
}

// Add creates a new empty target in the config.
func Add(cfg *config.Config, name string) error {
	if _, exists := cfg.Targets[name]; exists {
		return fmt.Errorf("target %q already exists", name)
	}
	if cfg.Targets == nil {
		cfg.Targets = make(map[string]map[string]string)
	}
	cfg.Targets[name] = make(map[string]string)
	return nil
}

// Remove deletes a target and all its variables from the config.
func Remove(cfg *config.Config, name string) error {
	if _, exists := cfg.Targets[name]; !exists {
		return fmt.Errorf("target %q not found", name)
	}
	delete(cfg.Targets, name)
	return nil
}

// Rename renames an existing target, preserving its variables.
func Rename(cfg *config.Config, oldName, newName string) error {
	envs, exists := cfg.Targets[oldName]
	if !exists {
		return fmt.Errorf("target %q not found", oldName)
	}
	if _, exists := cfg.Targets[newName]; exists {
		return fmt.Errorf("target %q already exists", newName)
	}
	cfg.Targets[newName] = envs
	delete(cfg.Targets, oldName)
	return nil
}
