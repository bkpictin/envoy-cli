// Package prefix provides utilities for adding, removing, and listing
// key prefixes across one or all targets in an envoy configuration.
package prefix

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/envoy/internal/config"
)

// Result holds the outcome of a prefix operation for a single target.
type Result struct {
	Target  string
	Changed int
	Skipped int
}

// Add prepends prefix to every key in target. If overwrite is false, keys that
// would conflict with an already-prefixed key are skipped.
func Add(cfg *config.Config, target, pfx string, overwrite bool) (Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", target)
	}
	if pfx == "" {
		return Result{}, fmt.Errorf("prefix must not be empty")
	}

	res := Result{Target: target}
	updated := make(map[string]string, len(envs))

	for k, v := range envs {
		if strings.HasPrefix(k, pfx) {
			updated[k] = v
			res.Skipped++
			continue
		}
		newKey := pfx + k
		if _, exists := envs[newKey]; exists && !overwrite {
			updated[k] = v
			res.Skipped++
			continue
		}
		updated[newKey] = v
		res.Changed++
	}

	cfg.Targets[target] = updated
	return res, nil
}

// Remove strips prefix from every key in target that carries it.
func Remove(cfg *config.Config, target, pfx string) (Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", target)
	}
	if pfx == "" {
		return Result{}, fmt.Errorf("prefix must not be empty")
	}

	res := Result{Target: target}
	updated := make(map[string]string, len(envs))

	for k, v := range envs {
		if strings.HasPrefix(k, pfx) {
			updated[strings.TrimPrefix(k, pfx)] = v
			res.Changed++
		} else {
			updated[k] = v
			res.Skipped++
		}
	}

	cfg.Targets[target] = updated
	return res, nil
}

// List returns every key in target that starts with pfx.
func List(cfg *config.Config, target, pfx string) ([]string, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	var keys []string
	for k := range envs {
		if strings.HasPrefix(k, pfx) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}
