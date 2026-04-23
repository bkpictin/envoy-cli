// Package search provides functionality to find environment variable keys
// across one or more deployment targets.
package search

import (
	"sort"
	"strings"

	"envoy-cli/internal/config"
)

// Result holds a single match found during a search.
type Result struct {
	Target string
	Key    string
	Value  string
}

// Options controls how a search is performed.
type Options struct {
	// Target restricts the search to a single target. Empty means all targets.
	Target string
	// CaseSensitive controls whether the key/value match is case-sensitive.
	CaseSensitive bool
	// SearchValues includes values in the match in addition to keys.
	SearchValues bool
}

// Keys searches for environment variable keys (and optionally values) that
// contain the given query string, returning all matching results sorted by
// target then key.
func Keys(cfg *config.Config, query string, opts Options) ([]Result, error) {
	if query == "" {
		return nil, nil
	}

	needle := query
	if !opts.CaseSensitive {
		needle = strings.ToLower(query)
	}

	var results []Result

	targets := cfg.Targets
	if opts.Target != "" {
		envs, ok := cfg.Targets[opts.Target]
		if !ok {
			return nil, fmt.Errorf("target %q not found", opts.Target)
		}
		targets = map[string]map[string]string{opts.Target: envs}
	}

	for targetName, envs := range targets {
		for k, v := range envs {
			keyHay := k
			valHay := v
			if !opts.CaseSensitive {
				keyHay = strings.ToLower(k)
				valHay = strings.ToLower(v)
			}
			if strings.Contains(keyHay, needle) || (opts.SearchValues && strings.Contains(valHay, needle)) {
				results = append(results, Result{
					Target: targetName,
					Key:    k,
					Value:  v,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Target != results[j].Target {
			return results[i].Target < results[j].Target
		}
		return results[i].Key < results[j].Key
	})

	return results, nil
}
