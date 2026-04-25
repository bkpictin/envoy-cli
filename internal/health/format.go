package health

import (
	"fmt"
	"strings"
)

// Format renders a Report as a human-readable string.
func Format(r Report) string {
	if r.OK() {
		return "✔  all checks passed — configuration looks healthy"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d issue(s) found:\n", len(r.Issues)))

	for _, issue := range r.Issues {
		var line string
		if issue.Target != "" {
			line = fmt.Sprintf("  [%s] (%s) %s", issue.Level, issue.Target, issue.Message)
		} else {
			line = fmt.Sprintf("  [%s] %s", issue.Level, issue.Message)
		}
		sb.WriteString(line + "\n")
	}

	return strings.TrimRight(sb.String(), "\n")
}

// FilterByLevel returns only issues matching the given level (e.g. "ERROR").
func FilterByLevel(r Report, level string) []Issue {
	var out []Issue
	for _, i := range r.Issues {
		if strings.EqualFold(i.Level, level) {
			out = append(out, i)
		}
	}
	return out
}
