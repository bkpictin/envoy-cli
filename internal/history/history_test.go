package history_test

import (
	"testing"

	"github.com/your-org/envoy-cli/internal/config"
	"github.com/your-org/envoy-cli/internal/history"
)

func newCfg() *config.Config {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"API_URL":  "https://api.prod.example.com",
				"LOG_LEVEL": "error",
			},
			"staging": {
				"API_URL":  "https://api.staging.example.com",
				"LOG_LEVEL": "debug",
			},
		},
		Audit: []config.AuditEntry{
			{Target: "production", Key: "API_URL", OldValue: "", NewValue: "https://api.prod.example.com", Op: "set", Timestamp: "2024-01-01T10:00:00Z"},
			{Target: "production", Key: "LOG_LEVEL", OldValue: "", NewValue: "warn", Op: "set", Timestamp: "2024-01-01T10:01:00Z"},
			{Target: "production", Key: "LOG_LEVEL", OldValue: "warn", NewValue: "error", Op: "set", Timestamp: "2024-01-02T09:00:00Z"},
			{Target: "staging", Key: "API_URL", OldValue: "", NewValue: "https://api.staging.example.com", Op: "set", Timestamp: "2024-01-03T08:00:00Z"},
			{Target: "production", Key: "DEBUG", OldValue: "true", NewValue: "", Op: "delete", Timestamp: "2024-01-04T07:00:00Z"},
		},
	}
	return cfg
}

func TestForTarget(t *testing.T) {
	cfg := newCfg()
	entries := history.ForTarget(cfg, "production")
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries for production, got %d", len(entries))
	}
}

func TestForTargetMissing(t *testing.T) {
	cfg := newCfg()
	entries := history.ForTarget(cfg, "nonexistent")
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries for nonexistent target, got %d", len(entries))
	}
}

func TestForKey(t *testing.T) {
	cfg := newCfg()
	entries := history.ForKey(cfg, "production", "LOG_LEVEL")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for LOG_LEVEL in production, got %d", len(entries))
	}
	if entries[0].NewValue != "warn" {
		t.Errorf("expected first entry NewValue=warn, got %s", entries[0].NewValue)
	}
	if entries[1].NewValue != "error" {
		t.Errorf("expected second entry NewValue=error, got %s", entries[1].NewValue)
	}
}

func TestForKeyMissing(t *testing.T) {
	cfg := newCfg()
	entries := history.ForKey(cfg, "production", "NONEXISTENT")
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries for missing key, got %d", len(entries))
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	entries := history.ForKey(cfg, "production", "LOG_LEVEL")
	out := history.Format(entries)
	if len(out) == 0 {
		t.Fatal("expected non-empty formatted output")
	}
}
