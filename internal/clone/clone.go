// Package clone provides functionality to duplicate an entire target
// including all its environment variables into a new target.
package clone

import (
	"fmt"

	"envoy-cli/internal/config"
)

// Target duplicates all environment variables from src into a newly created
// dest target. If dest already exists the operation is aborted unless
// overwrite is true.
func Target(cfg *config.Config, src, dest string, overwrite bool) error {
	srcEnvs, ok := cfg.Targets[src]
	if !ok {
		return fmt.Errorf("source target %q does not exist", src)
	}

	if _, exists := cfg.Targets[dest]; exists && !overwrite {
		return fmt.Errorf("destination target %q already exists; use --overwrite to replace it", dest)
	}

	copy := make(map[string]string, len(srcEnvs))
	for k, v := range srcEnvs {
		copy[k] = v
	}

	if cfg.Targets == nil {
		cfg.Targets = make(map[string]map[string]string)
	}
	cfg.Targets[dest] = copy
	return nil
}

// WithFilter duplicates only the keys that satisfy the predicate fn.
func WithFilter(cfg *config.Config, src, dest string, overwrite bool, fn func(key string) bool) error {
	srcEnvs, ok := cfg.Targets[src]
	if !ok {
		return fmt.Errorf("source target %q does not exist", src)
	}

	if _, exists := cfg.Targets[dest]; exists && !overwrite {
		return fmt.Errorf("destination target %q already exists; use --overwrite to replace it", dest)
	}

	filtered := make(map[string]string)
	for k, v := range srcEnvs {
		if fn(k) {
			filtered[k] = v
		}
	}

	if cfg.Targets == nil {
		cfg.Targets = make(map[string]map[string]string)
	}
	cfg.Targets[dest] = filtered
	return nil
}
