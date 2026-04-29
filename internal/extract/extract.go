// Package extract provides functionality for extracting a subset of keys
// from one or more targets into a new target or standalone map.
package extract

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Result holds the extracted key-value pairs and metadata.
type Result struct {
	Target string
	Keys   map[string]string
}

// FromTarget extracts the given keys from a source target.
// If keys is empty, all keys are extracted.
// Returns an error if the target does not exist or a requested key is missing.
func FromTarget(cfg *config.Config, target string, keys []string, strict bool) (Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", target)
	}

	out := make(map[string]string)

	if len(keys) == 0 {
		for k, v := range envs {
			out[k] = v
		}
		return Result{Target: target, Keys: out}, nil
	}

	for _, k := range keys {
		v, exists := envs[k]
		if !exists {
			if strict {
				return Result{}, fmt.Errorf("key %q not found in target %q", k, target)
			}
			continue
		}
		out[k] = v
	}

	return Result{Target: target, Keys: out}, nil
}

// IntoTarget writes the extracted keys into a destination target.
// If overwrite is false, existing keys in dest are preserved.
func IntoTarget(cfg *config.Config, dest string, result Result, overwrite bool) error {
	if _, ok := cfg.Targets[dest]; !ok {
		return fmt.Errorf("destination target %q not found", dest)
	}
	for k, v := range result.Keys {
		if _, exists := cfg.Targets[dest][k]; exists && !overwrite {
			continue
		}
		cfg.Targets[dest][k] = v
	}
	return nil
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Keys) == 0 {
		return fmt.Sprintf("[%s] no keys extracted\n", r.Target)
	}
	sorted := make([]string, 0, len(r.Keys))
	for k := range r.Keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	out := fmt.Sprintf("[%s] %d key(s) extracted:\n", r.Target, len(r.Keys))
	for _, k := range sorted {
		out += fmt.Sprintf("  %s=%s\n", k, r.Keys[k])
	}
	return out
}
