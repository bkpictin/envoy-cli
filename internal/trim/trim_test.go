package trim_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/trim"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"KEY_A": "  hello  ",
				"KEY_B": "no-space",
				"KEY_C": "\ttabbed\t",
			},
			"prod": {
				"KEY_X": " value ",
			},
		},
	}
}

func TestTrimTarget(t *testing.T) {
	cfg := newCfg()
	results, err := trim.Target(cfg, "dev", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["dev"]["KEY_A"] != "hello" {
		t.Errorf("expected KEY_A to be trimmed, got %q", cfg.Targets["dev"]["KEY_A"])
	}
	if cfg.Targets["dev"]["KEY_B"] != "no-space" {
		t.Errorf("KEY_B should be unchanged, got %q", cfg.Targets["dev"]["KEY_B"])
	}
	changed := 0
	for _, r := range results {
		if r.Changed() {
			changed++
		}
	}
	if changed != 2 {
		t.Errorf("expected 2 changed results, got %d", changed)
	}
}

func TestTrimTargetDryRun(t *testing.T) {
	cfg := newCfg()
	_, err := trim.Target(cfg, "dev", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["dev"]["KEY_A"] != "  hello  " {
		t.Errorf("dry run should not modify values, got %q", cfg.Targets["dev"]["KEY_A"])
	}
}

func TestTrimMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := trim.Target(cfg, "staging", false)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestTrimAll(t *testing.T) {
	cfg := newCfg()
	results, err := trim.All(cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["prod"]["KEY_X"] != "value" {
		t.Errorf("expected KEY_X trimmed, got %q", cfg.Targets["prod"]["KEY_X"])
	}
	if len(results) == 0 {
		t.Error("expected non-empty results")
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	results, _ := trim.All(cfg, true)
	out := trim.Format(results)
	if !strings.Contains(out, "KEY_A") {
		t.Errorf("expected KEY_A in format output, got: %s", out)
	}
}

func TestFormatNoChanges(t *testing.T) {
	out := trim.Format([]trim.Result{
		{Target: "dev", Key: "K", Before: "v", After: "v"},
	})
	if !strings.Contains(out, "no values required trimming") {
		t.Errorf("expected no-change message, got: %s", out)
	}
}
