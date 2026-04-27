// Package placeholder provides functionality to detect and report
// environment variable keys that contain placeholder or stub values
// such as "TODO", "FIXME", "CHANGEME", "<value>", etc.
package placeholder

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
)

// Result holds a single placeholder finding.
type Result struct {
	Target  string
	Key     string
	Value   string
	Pattern string
}

// defaultPatterns are case-insensitive substrings considered placeholders.
var defaultPatterns = []string{
	"todo",
	"fixme",
	"changeme",
	"replace_me",
	"<value>",
	"<your",
	"example",
	"placeholder",
	"xxx",
}

// Find scans the given targets (or all targets if empty) for placeholder values.
// Extra patterns can be supplied via extraPatterns and are merged with defaults.
func Find(cfg *config.Config, targets []string, extraPatterns []string) ([]Result, error) {
	patterns := append(defaultPatterns, extraPatterns...)

	scopeTargets := targets
	if len(scopeTargets) == 0 {
		for t := range cfg.Targets {
			scopeTargets = append(scopeTargets, t)
		}
	}

	var results []Result
	for _, t := range scopeTargets {
		envs, ok := cfg.Targets[t]
		if !ok {
			return nil, fmt.Errorf("target %q not found", t)
		}
		for k, v := range envs {
			for _, p := range patterns {
				if strings.Contains(strings.ToLower(v), strings.ToLower(p)) {
					results = append(results, Result{
						Target:  t,
						Key:     k,
						Value:   v,
						Pattern: p,
					})
					break
				}
			}
		}
	}
	return results, nil
}

// Format returns a human-readable summary of placeholder findings.
func Format(results []Result) string {
	if len(results) == 0 {
		return "No placeholder values detected.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Found %d placeholder value(s):\n", len(results))
	for _, r := range results {
		fmt.Fprintf(&sb, "  [%s] %s = %q  (matched: %q)\n", r.Target, r.Key, r.Value, r.Pattern)
	}
	return sb.String()
}
