// Package trim provides utilities for removing leading/trailing whitespace
// from environment variable values across one or more targets.
package trim

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
)

// Result holds the outcome of a trim operation for a single key.
type Result struct {
	Target string
	Key    string
	Before string
	After  string
}

// Changed reports whether the value was actually modified.
func (r Result) Changed() bool {
	return r.Before != r.After
}

// Target trims whitespace from all values in the named target.
// If dryRun is true the config is not saved and results are returned for
// inspection only.
func Target(cfg *config.Config, target string, dryRun bool) ([]Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}

	var results []Result
	for k, v := range envs {
		trimmed := strings.TrimSpace(v)
		results = append(results, Result{Target: target, Key: k, Before: v, After: trimmed})
		if !dryRun {
			envs[k] = trimmed
		}
	}
	return results, nil
}

// All trims whitespace from every value in every target.
// If dryRun is true the config is not saved.
func All(cfg *config.Config, dryRun bool) ([]Result, error) {
	var all []Result
	for target := range cfg.Targets {
		res, err := Target(cfg, target, dryRun)
		if err != nil {
			return nil, err
		}
		all = append(all, res...)
	}
	return all, nil
}

// Format returns a human-readable summary of trim results.
func Format(results []Result) string {
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if r.Changed() {
			changed++
			fmt.Fprintf(&sb, "[%s] %s: %q -> %q\n", r.Target, r.Key, r.Before, r.After)
		}
	}
	if changed == 0 {
		sb.WriteString("no values required trimming\n")
	}
	return sb.String()
}
