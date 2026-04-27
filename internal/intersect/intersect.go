// Package intersect finds keys that are common across multiple targets.
package intersect

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Result holds the intersection result for a set of targets.
type Result struct {
	Targets []string
	Keys    []string // keys present in ALL specified targets
}

// Targets returns the keys that exist in every one of the given targets.
// Returns an error if fewer than two targets are provided or any target is missing.
func Targets(cfg *config.Config, targets []string) (Result, error) {
	if len(targets) < 2 {
		return Result{}, fmt.Errorf("intersect requires at least two targets")
	}

	for _, t := range targets {
		if _, ok := cfg.Targets[t]; !ok {
			return Result{}, fmt.Errorf("target %q not found", t)
		}
	}

	// seed with keys from the first target
	seed := cfg.Targets[targets[0]]
	candidates := make(map[string]struct{}, len(seed))
	for k := range seed {
		candidates[k] = struct{}{}
	}

	// intersect with each subsequent target
	for _, t := range targets[1:] {
		envs := cfg.Targets[t]
		for k := range candidates {
			if _, exists := envs[k]; !exists {
				delete(candidates, k)
			}
		}
	}

	keys := make([]string, 0, len(candidates))
	for k := range candidates {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return Result{Targets: targets, Keys: keys}, nil
}

// Format returns a human-readable summary of an intersection result.
func Format(r Result) string {
	if len(r.Keys) == 0 {
		return fmt.Sprintf("no common keys across targets: %v\n", r.Targets)
	}
	out := fmt.Sprintf("common keys across %v (%d):\n", r.Targets, len(r.Keys))
	for _, k := range r.Keys {
		out += fmt.Sprintf("  %s\n", k)
	}
	return out
}
