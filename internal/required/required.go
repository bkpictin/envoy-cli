// Package required identifies environment variable keys that are present in
// all targets and flags any target that is missing one of those keys.
package required

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Result holds the outcome of a required-key check for a single target.
type Result struct {
	Target  string
	Missing []string // keys that are required but absent from this target
}

// Check computes the set of keys that appear in every target (the
// "required" set) and returns, for each target, which of those keys are
// missing.  Targets whose missing list is empty are considered healthy.
//
// If fewer than two targets exist the required set is empty and every
// Result will have a nil Missing slice.
func Check(cfg *config.Config) ([]Result, error) {
	if len(cfg.Targets) == 0 {
		return nil, nil
	}

	// Build per-target key sets and collect all target names.
	keySets := make(map[string]map[string]struct{}, len(cfg.Targets))
	names := make([]string, 0, len(cfg.Targets))
	for name, envs := range cfg.Targets {
		names = append(names, name)
		set := make(map[string]struct{}, len(envs))
		for k := range envs {
			set[k] = struct{}{}
		}
		keySets[name] = set
	}
	sort.Strings(names)

	// Intersect all key sets to find the required keys.
	required := make(map[string]struct{})
	for k := range keySets[names[0]] {
		required[k] = struct{}{}
	}
	for _, name := range names[1:] {
		for k := range required {
			if _, ok := keySets[name][k]; !ok {
				delete(required, k)
			}
		}
	}

	// Build results.
	results := make([]Result, 0, len(names))
	for _, name := range names {
		var missing []string
		for k := range required {
			if _, ok := keySets[name][k]; !ok {
				missing = append(missing, k)
			}
		}
		sort.Strings(missing)
		results = append(results, Result{Target: name, Missing: missing})
	}
	return results, nil
}

// Format renders a human-readable summary of required-key check results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no targets found\n"
	}
	out := ""
	for _, r := range results {
		if len(r.Missing) == 0 {
			out += fmt.Sprintf("[ok]   %s\n", r.Target)
		} else {
			for _, k := range r.Missing {
				out += fmt.Sprintf("[MISS] %s  missing key: %s\n", r.Target, k)
			}
		}
	}
	return out
}
