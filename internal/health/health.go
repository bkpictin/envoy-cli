// Package health checks the integrity of an envoy configuration,
// reporting missing targets, empty env maps, and snapshot inconsistencies.
package health

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Severity levels for health issues.
const (
	Warn  = "WARN"
	Error = "ERROR"
)

// Issue represents a single health finding.
type Issue struct {
	Level   string
	Target  string
	Message string
}

// Report holds all issues found during a health check.
type Report struct {
	Issues []Issue
}

// OK returns true when no issues were found.
func (r *Report) OK() bool { return len(r.Issues) == 0 }

// Check runs all health checks against cfg and returns a Report.
func Check(cfg *config.Config) Report {
	var issues []Issue

	if len(cfg.Targets) == 0 {
		issues = append(issues, Issue{Level: Warn, Message: "no targets defined"})
		return Report{Issues: issues}
	}

	for _, name := range sortedTargets(cfg) {
		envs := cfg.Targets[name]
		if len(envs) == 0 {
			issues = append(issues, Issue{Level: Warn, Target: name, Message: "target has no environment variables"})
		}
		for k, v := range envs {
			if v == "" {
				issues = append(issues, Issue{Level: Warn, Target: name, Message: fmt.Sprintf("key %q has an empty value", k)})
			}
		}
	}

	for snapName, snap := range cfg.Snapshots {
		if _, ok := cfg.Targets[snap.Target]; !ok {
			issues = append(issues, Issue{
				Level:   Error,
				Target:  snap.Target,
				Message: fmt.Sprintf("snapshot %q references non-existent target %q", snapName, snap.Target),
			})
		}
	}

	return Report{Issues: issues}
}

func sortedTargets(cfg *config.Config) []string {
	keys := make([]string, 0, len(cfg.Targets))
	for k := range cfg.Targets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
