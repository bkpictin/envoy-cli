// Package unique identifies keys that exist in exactly one target,
// making it easy to spot environment variables that have not been
// propagated to other deployment targets.
package unique

import (
	"fmt"
	"sort"
	"strings"

	"envoy-cli/internal/config"
)

// Result holds a key that is unique to a single target.
type Result struct {
	Target string
	Key    string
	Value  string
}

// Find returns all keys that appear in exactly one target across the
// entire configuration. If targets is non-empty only those targets are
// considered; otherwise every target is scanned.
func Find(cfg *config.Config, targets []string) ([]Result, error) {
	scope := targets
	if len(scope) == 0 {
		for t := range cfg.Targets {
			scope = append(scope, t)
		}
	}

	for _, t := range scope {
		if _, ok := cfg.Targets[t]; !ok {
			return nil, fmt.Errorf("target %q not found", t)
		}
	}

	// Count how many targets each key appears in (within the scope).
	keyCount := make(map[string]int)
	for _, t := range scope {
		for k := range cfg.Targets[t] {
			keyCount[k]++
		}
	}

	var results []Result
	for _, t := range scope {
		envs := cfg.Targets[t]
		keys := make([]string, 0, len(envs))
		for k := range envs {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if keyCount[k] == 1 {
				results = append(results, Result{Target: t, Key: k, Value: envs[k]})
			}
		}
	}

	return results, nil
}

// Format renders a slice of Result values as a human-readable string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no unique keys found"
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "[%s] %s=%s\n", r.Target, r.Key, r.Value)
	}
	return strings.TrimRight(sb.String(), "\n")
}
