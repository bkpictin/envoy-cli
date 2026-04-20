package lint

import (
	"fmt"
	"sort"
	"strings"
)

// Summary returns a human-readable summary string for a slice of issues,
// grouping counts by level.
func Summary(issues []Issue) string {
	if len(issues) == 0 {
		return "lint: no issues found"
	}
	counts := map[string]int{}
	for _, iss := range issues {
		counts[iss.Level]++
	}
	levels := make([]string, 0, len(counts))
	for l := range counts {
		levels = append(levels, l)
	}
	sort.Strings(levels)
	parts := make([]string, 0, len(levels))
	for _, l := range levels {
		parts = append(parts, fmt.Sprintf("%d %s(s)", counts[l], l))
	}
	return "lint: " + strings.Join(parts, ", ")
}

// FilterByLevel returns only the issues matching the given level string.
func FilterByLevel(issues []Issue, level string) []Issue {
	var out []Issue
	for _, iss := range issues {
		if iss.Level == level {
			out = append(out, iss)
		}
	}
	return out
}

// FilterByTarget returns only the issues for the specified target name.
func FilterByTarget(issues []Issue, target string) []Issue {
	var out []Issue
	for _, iss := range issues {
		if iss.Target == target {
			out = append(out, iss)
		}
	}
	return out
}
