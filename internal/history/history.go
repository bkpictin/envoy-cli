// Package history provides functionality for tracking and displaying
// the change history of environment variables across targets.
package history

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"envoy-cli/internal/config"
)

// Entry represents a single historical change to an environment variable.
type Entry struct {
	Timestamp time.Time
	Target    string
	Key       string
	OldValue  string
	NewValue  string
	Operation string // "set", "delete", "import", "restore"
}

// ForTarget returns all audit log entries for a given target, converted into
// history entries ordered from oldest to newest.
func ForTarget(cfg *config.Config, target string) ([]Entry, error) {
	if _, ok := cfg.Targets[target]; !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}

	var entries []Entry
	for _, record := range cfg.AuditLog {
		if record.Target != target {
			continue
		}
		entries = append(entries, Entry{
			Timestamp: record.Timestamp,
			Target:    record.Target,
			Key:       record.Key,
			OldValue:  record.OldValue,
			NewValue:  record.NewValue,
			Operation: record.Operation,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	return entries, nil
}

// ForKey returns all audit log entries for a specific key within a target.
func ForKey(cfg *config.Config, target, key string) ([]Entry, error) {
	all, err := ForTarget(cfg, target)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, e := range all {
		if strings.EqualFold(e.Key, key) {
			entries = append(entries, e)
		}
	}
	return entries, nil
}

// Format renders a slice of history entries as a human-readable table string.
func Format(entries []Entry) string {
	if len(entries) == 0 {
		return "(no history found)"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-22s %-10s %-24s %-20s %-20s\n",
		"TIMESTAMP", "OPERATION", "KEY", "OLD VALUE", "NEW VALUE"))
	sb.WriteString(strings.Repeat("-", 100) + "\n")

	for _, e := range entries {
		oldVal := truncate(e.OldValue, 18)
		newVal := truncate(e.NewValue, 18)
		sb.WriteString(fmt.Sprintf("%-22s %-10s %-24s %-20s %-20s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			e.Operation,
			e.Key,
			oldVal,
			newVal,
		))
	}
	return sb.String()
}

// truncate shortens a string to maxLen, appending "…" if it was shortened.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}
