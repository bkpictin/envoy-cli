// Package reorder provides utilities for reordering environment variable keys
// within a target, allowing deterministic or custom ordering.
package reorder

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Result holds the outcome of a reorder operation.
type Result struct {
	Target  string
	Before  []string
	After   []string
}

// Alphabetical sorts all keys in the given target alphabetically.
func Alphabetical(cfg *config.Config, target string) (Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", target)
	}

	before := make([]string, 0, len(envs))
	for k := range envs {
		before = append(before, k)
	}
	sort.Strings(before)

	after := make([]string, len(before))
	copy(after, before)
	sort.Strings(after)

	return Result{Target: target, Before: before, After: after}, nil
}

// Custom reorders keys in the target according to the provided order slice.
// Keys not mentioned in order are appended alphabetically at the end.
func Custom(cfg *config.Config, target string, order []string) (Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", target)
	}

	// Validate all keys in order exist in target.
	for _, k := range order {
		if _, exists := envs[k]; !exists {
			return Result{}, fmt.Errorf("key %q not found in target %q", k, target)
		}
	}

	before := make([]string, 0, len(envs))
	for k := range envs {
		before = append(before, k)
	}
	sort.Strings(before)

	seen := make(map[string]bool, len(order))
	after := make([]string, 0, len(envs))
	for _, k := range order {
		if !seen[k] {
			after = append(after, k)
			seen[k] = true
		}
	}

	// Append remaining keys alphabetically.
	remaining := make([]string, 0)
	for _, k := range before {
		if !seen[k] {
			remaining = append(remaining, k)
		}
	}
	sort.Strings(remaining)
	after = append(after, remaining...)

	return Result{Target: target, Before: before, After: after}, nil
}
