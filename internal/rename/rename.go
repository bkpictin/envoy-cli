// Package rename provides functionality for renaming environment variable keys
// across one or all targets in the configuration.
package rename

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/config"
)

// KeyInTarget renames a key within a single target.
// Returns an error if the target does not exist, the old key is not found,
// or the new key already exists in that target.
func KeyInTarget(cfg *config.Config, target, oldKey, newKey string) error {
	envs, ok := cfg.Targets[target]
	if !ok {
		return fmt.Errorf("target %q not found", target)
	}

	val, exists := envs[oldKey]
	if !exists {
		return fmt.Errorf("key %q not found in target %q", oldKey, target)
	}

	if _, conflict := envs[newKey]; conflict {
		return fmt.Errorf("key %q already exists in target %q", newKey, target)
	}

	delete(envs, oldKey)
	envs[newKey] = val
	cfg.Targets[target] = envs
	return nil
}

// KeyInAll renames a key across every target that contains it.
// Targets that do not have the old key are skipped.
// Returns an error if the new key already exists in any target that has the old key.
func KeyInAll(cfg *config.Config, oldKey, newKey string) error {
	for target, envs := range cfg.Targets {
		_, hasOld := envs[oldKey]
		if !hasOld {
			continue
		}
		if _, conflict := envs[newKey]; conflict {
			return fmt.Errorf("key %q already exists in target %q; rename aborted", newKey, target)
		}
	}

	for target, envs := range cfg.Targets {
		val, hasOld := envs[oldKey]
		if !hasOld {
			continue
		}
		delete(envs, oldKey)
		envs[newKey] = val
		cfg.Targets[target] = envs
	}
	return nil
}
