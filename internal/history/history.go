// Package history provides functions to query and display the audit log
// history for specific targets or keys within the envoy configuration.
package history

import (
	"fmt"
	"strings"

	"github.com/your-org/envoy-cli/internal/config"
)

// ForTarget returns all audit entries for the given target, in chronological order.
func ForTarget(cfg *config.Config, target string) []config.AuditEntry {
	var entries []config.AuditEntry
	for _, e := range cfg.Audit {
		if e.Target == target {
			entries = append(entries, e)
		}
	}
	return entries
}

// ForKey returns all audit entries for a specific key within a target.
func ForKey(cfg *config.Config, target, key string) []config.AuditEntry {
	var entries []config.AuditEntry
	for _, e := range cfg.Audit {
		if e.Target == target && e.Key == key {
			entries = append(entries, e)
		}
	}
	return entries
}

// Format renders a slice of audit entries as a human-readable string.
func Format(entries []config.AuditEntry) string {
	if len(entries) == 0 {
		return "(no history)"
	}
	var sb strings.Builder
	for _, e := range entries {
		var line string
		switch e.Op {
		case "set":
			if e.OldValue == "" {
				line = fmt.Sprintf("%s  [%s] SET %s = %s",
					e.Timestamp, e.Target, e.Key, truncate(e.NewValue, 40))
			} else {
				line = fmt.Sprintf("%s  [%s] SET %s: %s → %s",
					e.Timestamp, e.Target, e.Key,
					truncate(e.OldValue, 30), truncate(e.NewValue, 30))
			}
		case "delete":
			line = fmt.Sprintf("%s  [%s] DELETE %s (was: %s)",
				e.Timestamp, e.Target, e.Key, truncate(e.OldValue, 30))
		default:
			line = fmt.Sprintf("%s  [%s] %s %s",
				e.Timestamp, e.Target, strings.ToUpper(e.Op), e.Key)
		}
		sb.WriteString(line)
		sb.WriteRune('\n')
	}
	return sb.String()
}

// truncate shortens a string to maxLen characters, appending "..." if needed.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
