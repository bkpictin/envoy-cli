// Package cascade provides functionality for applying environment variable
// values from a base target down through a chain of derived targets,
// only setting keys that are not already defined in each target.
package cascade

import (
	"fmt"

	"envoy-cli/internal/config"
)

// Result holds the outcome of a cascade operation for a single target.
type Result struct {
	Target  string
	Applied int
	Skipped int
}

// Apply cascades environment variables from the base target through each
// target in the chain in order. Keys already present in a target are
// skipped unless overwrite is true. Each target in the chain also
// inherits any keys propagated to the previous target.
func Apply(cfg *config.Config, base string, chain []string, overwrite bool) ([]Result, error) {
	if _, ok := cfg.Targets[base]; !ok {
		return nil, fmt.Errorf("base target %q not found", base)
	}
	for i, t := range chain {
		if _, ok := cfg.Targets[t]; !ok {
			return nil, fmt.Errorf("chain target %q (index %d) not found", t, i)
		}
	}

	// Work with a running source that accumulates as we walk the chain.
	source := copyMap(cfg.Targets[base])
	results := make([]Result, 0, len(chain))

	for _, t := range chain {
		dest := cfg.Targets[t]
		applied, skipped := 0, 0

		for k, v := range source {
			if _, exists := dest[k]; exists && !overwrite {
				skipped++
				continue
			}
			dest[k] = v
			applied++
		}
		cfg.Targets[t] = dest
		results = append(results, Result{Target: t, Applied: applied, Skipped: skipped})

		// Next iteration sources from the now-updated target.
		source = copyMap(dest)
	}
	return results, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
