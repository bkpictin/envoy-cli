package audit

import (
	"fmt"
	"time"

	"github.com/envoy-cli/internal/config"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Key       string    `json:"key,omitempty"`
	Detail    string    `json:"detail,omitempty"`
}

// Log appends an audit entry to the config.
func Log(cfg *config.Config, action, target, key, detail string) {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Action:    action,
		Target:    target,
		Key:       key,
		Detail:    detail,
	}
	if cfg.Audit == nil {
		cfg.Audit = []config.AuditEntry{}
	}
	cfg.Audit = append(cfg.Audit, config.AuditEntry{
		Timestamp: entry.Timestamp,
		Action:    entry.Action,
		Target:    entry.Target,
		Key:       entry.Key,
		Detail:    entry.Detail,
	})
}

// List returns all audit entries, optionally filtered by target.
func List(cfg *config.Config, target string) []config.AuditEntry {
	var result []config.AuditEntry
	for _, e := range cfg.Audit {
		if target == "" || e.Target == target {
			result = append(result, e)
		}
	}
	return result
}

// Format renders audit entries as human-readable lines.
func Format(entries []config.AuditEntry) string {
	if len(entries) == 0 {
		return "no audit entries found"
	}
	out := ""
	for _, e := range entries {
		line := fmt.Sprintf("%s  %-12s  target=%-15s", e.Timestamp.Format(time.RFC3339), e.Action, e.Target)
		if e.Key != "" {
			line += fmt.Sprintf("  key=%s", e.Key)
		}
		if e.Detail != "" {
			line += fmt.Sprintf("  (%s)", e.Detail)
		}
		out += line + "\n"
	}
	return out
}
