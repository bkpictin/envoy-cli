// Package inherit provides functionality for inheriting environment variables
// from a base target into one or more child targets.
package inherit

import (
	"fmt"

	"envoy-cli/internal/config"
)

// Result holds the outcome of an inherit operation for a single child target.
type Result struct {
	Target  string
	Added   int
	Skipped int
}

// Apply copies all keys from baseTarget into each child target.
// Existing keys in child targets are skipped unless overwrite is true.
func Apply(cfg *config.Config, baseTarget string, children []string, overwrite bool) ([]Result, error) {
	base, ok := cfg.Targets[baseTarget]
	if !ok {
		return nil, fmt.Errorf("base target %q not found", baseTarget)
	}

	var results []Result

	for _, child := range children {
		if child == baseTarget {
			return nil, fmt.Errorf("child target cannot be the same as base target %q", baseTarget)
		}
		dest, ok := cfg.Targets[child]
		if !ok {
			return nil, fmt.Errorf("child target %q not found", child)
		}

		result := Result{Target: child}
		for k, v := range base {
			if _, exists := dest[k]; exists && !overwrite {
				result.Skipped++
				continue
			}
			dest[k] = v
			result.Added++
		}
		cfg.Targets[child] = dest
		results = append(results, result)
	}

	return results, nil
}

// ListInherited returns the keys in childTarget that match keys in baseTarget.
func ListInherited(cfg *config.Config, baseTarget, childTarget string) ([]string, error) {
	base, ok := cfg.Targets[baseTarget]
	if !ok {
		return nil, fmt.Errorf("base target %q not found", baseTarget)
	}
	child, ok := cfg.Targets[childTarget]
	if !ok {
		return nil, fmt.Errorf("child target %q not found", childTarget)
	}

	var shared []string
	for k := range base {
		if _, exists := child[k]; exists {
			shared = append(shared, k)
		}
	}
	return shared, nil
}
