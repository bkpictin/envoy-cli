package env

import (
	"errors"
	"fmt"
	"os"

	"github.com/envoy-cli/envoy-cli/internal/config"
)

// Set adds or updates an environment variable in the given target.
func Set(cfg *config.Config, target, key, value string) error {
	if target == "" {
		return errors.New("target must not be empty")
	}
	if key == "" {
		return errors.New("key must not be empty")
	}
	if cfg.Targets == nil {
		cfg.Targets = make(map[string]map[string]string)
	}
	if cfg.Targets[target] == nil {
		cfg.Targets[target] = make(map[string]string)
	}
	cfg.Targets[target][key] = value
	return nil
}

// Get retrieves an environment variable from the given target.
func Get(cfg *config.Config, target, key string) (string, error) {
	if cfg.Targets == nil {
		return "", fmt.Errorf("target %q not found", target)
	}
	vars, ok := cfg.Targets[target]
	if !ok {
		return "", fmt.Errorf("target %q not found", target)
	}
	val, ok := vars[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in target %q", key, target)
	}
	return val, nil
}

// Delete removes an environment variable from the given target.
func Delete(cfg *config.Config, target, key string) error {
	if cfg.Targets == nil {
		return fmt.Errorf("target %q not found", target)
	}
	if _, ok := cfg.Targets[target]; !ok {
		return fmt.Errorf("target %q not found", target)
	}
	delete(cfg.Targets[target], key)
	return nil
}

// Export writes all variables for a target to the current process environment.
func Export(cfg *config.Config, target string) error {
	if cfg.Targets == nil {
		return fmt.Errorf("target %q not found", target)
	}
	vars, ok := cfg.Targets[target]
	if !ok {
		return fmt.Errorf("target %q not found", target)
	}
	for k, v := range vars {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("failed to export %q: %w", k, err)
		}
	}
	return nil
}
