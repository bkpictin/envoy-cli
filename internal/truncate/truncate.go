// Package truncate provides utilities for truncating environment variable
// values to a maximum length across one or all targets.
package truncate

import (
	"fmt"

	"github.com/user/envoy-cli/internal/config"
)

// Result holds the outcome of a single truncation operation.
type Result struct {
	Target   string
	Key      string
	Original string
	Truncated string
	Changed  bool
}

// Target truncates all values in the named target that exceed maxLen.
// If dryRun is true the config is not modified.
func Target(cfg *config.Config, target string, maxLen int, dryRun bool) ([]Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	if maxLen <= 0 {
		return nil, fmt.Errorf("maxLen must be greater than zero")
	}

	var results []Result
	for k, v := range envs {
		if len(v) > maxLen {
			trunc := v[:maxLen]
			results = append(results, Result{
				Target:    target,
				Key:       k,
				Original:  v,
				Truncated: trunc,
				Changed:   true,
			})
			if !dryRun {
				cfg.Targets[target][k] = trunc
			}
		} else {
			results = append(results, Result{
				Target:    target,
				Key:       k,
				Original:  v,
				Truncated: v,
				Changed:   false,
			})
		}
	}
	return results, nil
}

// All truncates values across every target.
func All(cfg *config.Config, maxLen int, dryRun bool) ([]Result, error) {
	if maxLen <= 0 {
		return nil, fmt.Errorf("maxLen must be greater than zero")
	}
	var all []Result
	for t := range cfg.Targets {
		res, err := Target(cfg, t, maxLen, dryRun)
		if err != nil {
			return nil, err
		}
		all = append(all, res...)
	}
	return all, nil
}

// Format returns a human-readable summary of truncation results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no values to truncate\n"
	}
	out := ""
	changed := 0
	for _, r := range results {
		if r.Changed {
			out += fmt.Sprintf("[%s] %s: %q → %q\n", r.Target, r.Key, r.Original, r.Truncated)
			changed++
		}
	}
	if changed == 0 {
		return "all values within limit\n"
	}
	out += fmt.Sprintf("%d value(s) truncated\n", changed)
	return out
}
