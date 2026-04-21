// Package prune removes stale or unused keys across targets.
package prune

import (
	"fmt"

	"github.com/envoy-cli/internal/config"
)

// Result holds information about a single pruned key.
type Result struct {
	Target string
	Key    string
}

// OrphanedKeys removes keys from targetName that do not exist in any other
// target. If dryRun is true the keys are identified but not deleted.
func OrphanedKeys(cfg *config.Config, targetName string, dryRun bool) ([]Result, error) {
	envs, ok := cfg.Targets[targetName]
	if !ok {
		return nil, fmt.Errorf("target %q not found", targetName)
	}

	// Build a set of keys that appear in at least one other target.
	shared := map[string]bool{}
	for name, other := range cfg.Targets {
		if name == targetName {
			continue
		}
		for k := range other {
			shared[k] = true
		}
	}

	var results []Result
	for k := range envs {
		if !shared[k] {
			results = append(results, Result{Target: targetName, Key: k})
			if !dryRun {
				delete(cfg.Targets[targetName], k)
			}
		}
	}
	return results, nil
}

// EmptyValues removes keys whose value is an empty string from targetName.
// If dryRun is true the keys are identified but not deleted.
func EmptyValues(cfg *config.Config, targetName string, dryRun bool) ([]Result, error) {
	envs, ok := cfg.Targets[targetName]
	if !ok {
		return nil, fmt.Errorf("target %q not found", targetName)
	}

	var results []Result
	for k, v := range envs {
		if v == "" {
			results = append(results, Result{Target: targetName, Key: k})
			if !dryRun {
				delete(cfg.Targets[targetName], k)
			}
		}
	}
	return results, nil
}
