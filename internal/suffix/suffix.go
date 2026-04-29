// Package suffix provides utilities for adding and removing key suffixes
// across one or all targets in an envoy configuration.
package suffix

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/envoy/internal/config"
)

// Result holds the outcome of a suffix operation on a single key.
type Result struct {
	Target  string
	OldKey  string
	NewKey  string
	Skipped bool
	Reason  string
}

// Add appends suffix to every key in the given target that does not already
// end with it. When dryRun is true the config is not mutated.
func Add(cfg *config.Config, target, suffix string, dryRun bool) ([]Result, error) {
	if suffix == "" {
		return nil, fmt.Errorf("suffix must not be empty")
	}
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	var results []Result
	updated := make(map[string]string, len(envs))
	for k, v := range envs {
		if strings.HasSuffix(k, suffix) {
			updated[k] = v
			results = append(results, Result{Target: target, OldKey: k, NewKey: k, Skipped: true, Reason: "already has suffix"})
			continue
		}
		newKey := k + suffix
		if _, exists := envs[newKey]; exists {
			updated[k] = v
			results = append(results, Result{Target: target, OldKey: k, NewKey: newKey, Skipped: true, Reason: "destination key already exists"})
			continue
		}
		updated[newKey] = v
		results = append(results, Result{Target: target, OldKey: k, NewKey: newKey})
	}
	if !dryRun {
		cfg.Targets[target] = updated
	}
	return results, nil
}

// Remove strips suffix from every key in the given target that ends with it.
// When dryRun is true the config is not mutated.
func Remove(cfg *config.Config, target, suffix string, dryRun bool) ([]Result, error) {
	if suffix == "" {
		return nil, fmt.Errorf("suffix must not be empty")
	}
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	var results []Result
	updated := make(map[string]string, len(envs))
	for k, v := range envs {
		if !strings.HasSuffix(k, suffix) {
			updated[k] = v
			results = append(results, Result{Target: target, OldKey: k, NewKey: k, Skipped: true, Reason: "does not have suffix"})
			continue
		}
		newKey := strings.TrimSuffix(k, suffix)
		if _, exists := envs[newKey]; exists {
			updated[k] = v
			results = append(results, Result{Target: target, OldKey: k, NewKey: newKey, Skipped: true, Reason: "destination key already exists"})
			continue
		}
		updated[newKey] = v
		results = append(results, Result{Target: target, OldKey: k, NewKey: newKey})
	}
	if !dryRun {
		cfg.Targets[target] = updated
	}
	return results, nil
}

// List returns all keys in the given target that end with suffix.
func List(cfg *config.Config, target, suffix string) ([]string, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	var keys []string
	for k := range envs {
		if strings.HasSuffix(k, suffix) {
			keys = append(keys, k)
		}
	}
	return keys, nil
}
