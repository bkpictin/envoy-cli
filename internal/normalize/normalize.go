// Package normalize provides utilities for standardizing environment variable
// values across targets, such as trimming whitespace, converting case, and
// collapsing repeated delimiters.
package normalize

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
)

// Result holds the outcome of a normalization operation for a single key.
type Result struct {
	Target string
	Key    string
	Before string
	After  string
}

// Changed reports whether the value was actually modified.
func (r Result) Changed() bool { return r.Before != r.After }

// Options controls which normalization passes are applied.
type Options struct {
	TrimSpace    bool
	UpperKeys    bool
	LowerValues  bool
	CollapseSpaces bool
}

// Target normalises all values in the named target according to opts.
// When dryRun is true the config is not mutated and results are still returned.
func Target(cfg *config.Config, target string, opts Options, dryRun bool) ([]Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}

	var results []Result
	for k, v := range envs {
		newKey := k
		newVal := v

		if opts.TrimSpace {
			newVal = strings.TrimSpace(newVal)
		}
		if opts.CollapseSpaces {
			newVal = strings.Join(strings.Fields(newVal), " ")
		}
		if opts.LowerValues {
			newVal = strings.ToLower(newVal)
		}
		if opts.UpperKeys {
			newKey = strings.ToUpper(k)
		}

		results = append(results, Result{Target: target, Key: k, Before: v, After: newVal})

		if !dryRun {
			delete(cfg.Targets[target], k)
			cfg.Targets[target][newKey] = newVal
		}
	}
	return results, nil
}

// Format returns a human-readable summary of normalization results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no changes\n"
	}
	var sb strings.Builder
	changed := 0
	for _, r := range results {
		if r.Changed() {
			changed++
			fmt.Fprintf(&sb, "  [%s] %s: %q -> %q\n", r.Target, r.Key, r.Before, r.After)
		}
	}
	if changed == 0 {
		return "no changes\n"
	}
	fmt.Fprintf(&sb, "%d value(s) normalized\n", changed)
	return sb.String()
}
