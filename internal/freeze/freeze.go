// Package freeze provides functionality to freeze and unfreeze targets,
// preventing any modifications to their environment variables.
package freeze

import (
	"errors"
	"fmt"

	"envoy-cli/internal/config"
)

const frozenKey = "__frozen__"

// Freeze marks a target as frozen, preventing modifications.
func Freeze(cfg *config.Config, target string) error {
	if _, ok := cfg.Targets[target]; !ok {
		return fmt.Errorf("target %q not found", target)
	}
	if cfg.Targets[target] == nil {
		cfg.Targets[target] = map[string]string{}
	}
	if cfg.Targets[target][frozenKey] == "true" {
		return fmt.Errorf("target %q is already frozen", target)
	}
	cfg.Targets[target][frozenKey] = "true"
	return nil
}

// Unfreeze removes the frozen marker from a target.
func Unfreeze(cfg *config.Config, target string) error {
	if _, ok := cfg.Targets[target]; !ok {
		return fmt.Errorf("target %q not found", target)
	}
	if cfg.Targets[target][frozenKey] != "true" {
		return fmt.Errorf("target %q is not frozen", target)
	}
	delete(cfg.Targets[target], frozenKey)
	return nil
}

// IsFrozen reports whether a target is currently frozen.
func IsFrozen(cfg *config.Config, target string) bool {
	envs, ok := cfg.Targets[target]
	if !ok {
		return false
	}
	return envs[frozenKey] == "true"
}

// List returns all frozen target names.
func List(cfg *config.Config) []string {
	var frozen []string
	for name := range cfg.Targets {
		if IsFrozen(cfg, name) {
			frozen = append(frozen, name)
		}
	}
	return frozen
}

// ErrFrozen is returned when attempting to modify a frozen target.
var ErrFrozen = errors.New("target is frozen")

// GuardWrite returns ErrFrozen if the target is frozen, nil otherwise.
func GuardWrite(cfg *config.Config, target string) error {
	if IsFrozen(cfg, target) {
		return fmt.Errorf("%w: %q cannot be modified", ErrFrozen, target)
	}
	return nil
}
