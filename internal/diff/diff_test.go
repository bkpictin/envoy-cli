package diff

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"staging": {
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
				"DEBUG":   "true",
			},
			"production": {
				"DB_HOST": "prod.db.internal",
				"DB_PORT": "5432",
				"LOG_LEVEL": "warn",
			},
		},
	}
}

func TestDiffTargets(t *testing.T) {
	cfg := newCfg()
	r, err := Targets(cfg, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := r.OnlyInA["DEBUG"]; !ok {
		t.Error("expected DEBUG to be only in staging")
	}
	if _, ok := r.OnlyInB["LOG_LEVEL"]; !ok {
		t.Error("expected LOG_LEVEL to be only in production")
	}
	if v, ok := r.Changed["DB_HOST"]; !ok || v[0] != "localhost" || v[1] != "prod.db.internal" {
		t.Error("expected DB_HOST to be changed")
	}
	if _, ok := r.Unchanged["DB_PORT"]; !ok {
		t.Error("expected DB_PORT to be unchanged")
	}
}

func TestDiffMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := Targets(cfg, "staging", "missing")
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	r, _ := Targets(cfg, "staging", "production")
	out := Format(r, "staging", "production")
	if !strings.Contains(out, "DEBUG") {
		t.Error("expected DEBUG in output")
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in output")
	}
}

func TestFormatNoDiff(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"a": {"KEY": "val"},
			"b": {"KEY": "val"},
		},
	}
	r, _ := Targets(cfg, "a", "b")
	out := Format(r, "a", "b")
	if !strings.Contains(out, "No differences") {
		t.Error("expected no-differences message")
	}
}
