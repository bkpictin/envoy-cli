package normalize_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/normalize"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"DATABASE_URL": "  postgres://localhost/prod  ",
				"log_level":    "DEBUG",
				"APP_NAME":     "my   app",
			},
		},
	}
}

func TestTrimSpace(t *testing.T) {
	cfg := newCfg()
	opts := normalize.Options{TrimSpace: true}
	results, err := normalize.Target(cfg, "production", opts, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "DATABASE_URL" && r.After != "postgres://localhost/prod" {
			t.Errorf("expected trimmed value, got %q", r.After)
		}
	}
	if v := cfg.Targets["production"]["DATABASE_URL"]; v != "postgres://localhost/prod" {
		t.Errorf("config not mutated: got %q", v)
	}
}

func TestDryRun(t *testing.T) {
	cfg := newCfg()
	opts := normalize.Options{TrimSpace: true}
	_, err := normalize.Target(cfg, "production", opts, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// config must remain unchanged
	if v := cfg.Targets["production"]["DATABASE_URL"]; v != "  postgres://localhost/prod  " {
		t.Errorf("dry-run mutated config: got %q", v)
	}
}

func TestCollapseSpaces(t *testing.T) {
	cfg := newCfg()
	opts := normalize.Options{CollapseSpaces: true}
	results, err := normalize.Target(cfg, "production", opts, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "APP_NAME" && r.After != "my app" {
			t.Errorf("expected collapsed spaces, got %q", r.After)
		}
	}
}

func TestMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := normalize.Target(cfg, "staging", normalize.Options{}, false)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	opts := normalize.Options{TrimSpace: true}
	results, _ := normalize.Target(cfg, "production", opts, true)
	out := normalize.Format(results)
	if !strings.Contains(out, "normalized") && !strings.Contains(out, "no changes") {
		t.Errorf("unexpected format output: %q", out)
	}
}

func TestFormatNoChanges(t *testing.T) {
	out := normalize.Format(nil)
	if out != "no changes\n" {
		t.Errorf("expected 'no changes', got %q", out)
	}
}
