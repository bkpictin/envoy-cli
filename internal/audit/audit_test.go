package audit

import (
	"strings"
	"testing"
	"time"

	"github.com/envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{},
	}
}

func TestLog(t *testing.T) {
	cfg := newCfg()
	Log(cfg, "set", "production", "API_KEY", "")
	if len(cfg.Audit) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(cfg.Audit))
	}
	e := cfg.Audit[0]
	if e.Action != "set" || e.Target != "production" || e.Key != "API_KEY" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestListAll(t *testing.T) {
	cfg := newCfg()
	Log(cfg, "set", "staging", "FOO", "")
	Log(cfg, "delete", "production", "BAR", "")
	entries := List(cfg, "")
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestListFiltered(t *testing.T) {
	cfg := newCfg()
	Log(cfg, "set", "staging", "FOO", "")
	Log(cfg, "delete", "production", "BAR", "")
	entries := List(cfg, "staging")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Target != "staging" {
		t.Errorf("expected staging, got %s", entries[0].Target)
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	Log(cfg, "set", "production", "SECRET", "updated via CLI")
	entries := List(cfg, "")
	out := Format(entries)
	if !strings.Contains(out, "set") || !strings.Contains(out, "production") {
		t.Errorf("unexpected format output: %s", out)
	}
	if !strings.Contains(out, "SECRET") {
		t.Errorf("expected key in output: %s", out)
	}
}

func TestFormatEmpty(t *testing.T) {
	out := Format(nil)
	if out != "no audit entries found" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestTimestampUTC(t *testing.T) {
	cfg := newCfg()
	Log(cfg, "set", "dev", "X", "")
	if cfg.Audit[0].Timestamp.Location() != time.UTC {
		t.Error("timestamp should be UTC")
	}
}
