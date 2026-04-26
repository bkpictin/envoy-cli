// Package count provides utilities for counting keys and targets
// across an envoy configuration, with optional filtering by target or prefix.
package count

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envoy-cli/envoy/internal/config"
)

// Result holds the count summary for a single target.
type Result struct {
	Target string
	Total  int
	Empty  int
	Filled int
}

// Summary holds count results across all queried targets.
type Summary struct {
	Results    []Result
	TotalKeys  int
	TotalEmpty int
}

// ByTarget counts keys in a single target, returning an error if the target
// does not exist.
func ByTarget(cfg *config.Config, target string) (Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", target)
	}

	r := Result{Target: target, Total: len(envs)}
	for _, v := range envs {
		if strings.TrimSpace(v) == "" {
			r.Empty++
		}
	}
	r.Filled = r.Total - r.Empty
	return r, nil
}

// All counts keys across every target in the config, sorted by target name.
func All(cfg *config.Config) Summary {
	names := make([]string, 0, len(cfg.Targets))
	for t := range cfg.Targets {
		names = append(names, t)
	}
	sort.Strings(names)

	var s Summary
	for _, t := range names {
		r, _ := ByTarget(cfg, t)
		s.Results = append(s.Results, r)
		s.TotalKeys += r.Total
		s.TotalEmpty += r.Empty
	}
	return s
}

// Format renders a Summary as a human-readable string.
func Format(s Summary) string {
	if len(s.Results) == 0 {
		return "no targets found\n"
	}
	var sb strings.Builder
	for _, r := range s.Results {
		sb.WriteString(fmt.Sprintf("%-20s total=%-4d filled=%-4d empty=%d\n",
			r.Target, r.Total, r.Filled, r.Empty))
	}
	sb.WriteString(fmt.Sprintf("\ntotals: keys=%d empty=%d\n", s.TotalKeys, s.TotalEmpty))
	return sb.String()
}
