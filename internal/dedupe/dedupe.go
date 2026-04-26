// Package dedupe identifies and removes duplicate environment variable values
// across targets, helping keep configurations clean and consistent.
package dedupe

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Match represents a duplicate value found across targets.
type Match struct {
	Key     string
	Value   string
	Targets []string
}

// Result holds all duplicate matches found in a config.
type Result struct {
	Matches []Match
}

// FindCrossTarget scans all targets and returns keys whose values are identical
// in two or more targets.
func FindCrossTarget(cfg *config.Config) Result {
	type entry struct {
		targets []string
	}
	// map[key+"="+value] -> targets
	index := map[string]*entry{}

	targetNames := make([]string, 0, len(cfg.Targets))
	for t := range cfg.Targets {
		targetNames = append(targetNames, t)
	}
	sort.Strings(targetNames)

	for _, t := range targetNames {
		for k, v := range cfg.Targets[t] {
			key := fmt.Sprintf("%s=%s", k, v)
			if _, ok := index[key]; !ok {
				index[key] = &entry{}
			}
			index[key].targets = append(index[key].targets, t)
		}
	}

	var matches []Match
	for composite, e := range index {
		if len(e.targets) < 2 {
			continue
		}
		// parse key and value back out
		var k, v string
		fmt.Sscanf(composite, "%s", &k) // fallback
		for i, ch := range composite {
			if ch == '=' {
				k = composite[:i]
				v = composite[i+1:]
				break
			}
		}
		sort.Strings(e.targets)
		matches = append(matches, Match{Key: k, Value: v, Targets: e.targets})
	}

	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Key != matches[j].Key {
			return matches[i].Key < matches[j].Key
		}
		return matches[i].Value < matches[j].Value
	})

	return Result{Matches: matches}
}

// Format returns a human-readable report of duplicate findings.
func Format(r Result) string {
	if len(r.Matches) == 0 {
		return "No duplicate values found across targets.\n"
	}
	out := fmt.Sprintf("Found %d duplicate value(s):\n", len(r.Matches))
	for _, m := range r.Matches {
		out += fmt.Sprintf("  %-24s  targets: %v\n", m.Key, m.Targets)
	}
	return out
}
