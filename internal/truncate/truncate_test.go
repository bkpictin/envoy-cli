package truncate_test

import (
	"strings"
	"testing"

	"github.com/user/envoy-cli/internal/config"
	"github.com/user/envoy-cli/internal/truncate"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"SHORT": "hi",
				"LONG": "this-is-a-very-long-value-exceeding-limit",
				"EXACT": "12345",
			},
			"staging": {
				"KEY": "another-long-staging-value",
			},
		},
	}
}

func TestTruncateTarget(t *testing.T) {
	cfg := newCfg()
	res, err := truncate.Target(cfg, "production", 5, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changed := 0
	for _, r := range res {
		if r.Changed {
			changed++
			if len(r.Truncated) > 5 {
				t.Errorf("truncated value too long: %s", r.Truncated)
			}
			if cfg.Targets["production"][r.Key] != r.Truncated {
				t.Errorf("config not updated for key %s", r.Key)
			}
		}
	}
	if changed == 0 {
		t.Error("expected at least one changed result")
	}
}

func TestTruncateDryRun(t *testing.T) {
	cfg := newCfg()
	orig := cfg.Targets["production"]["LONG"]
	_, err := truncate.Target(cfg, "production", 5, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["production"]["LONG"] != orig {
		t.Error("dry run must not modify config")
	}
}

func TestTruncateMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := truncate.Target(cfg, "nonexistent", 10, false)
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestTruncateInvalidMaxLen(t *testing.T) {
	cfg := newCfg()
	_, err := truncate.Target(cfg, "production", 0, false)
	if err == nil {
		t.Error("expected error for maxLen=0")
	}
}

func TestTruncateAll(t *testing.T) {
	cfg := newCfg()
	res, err := truncate.All(cfg, 5, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) == 0 {
		t.Error("expected results from all targets")
	}
	targets := map[string]bool{}
	for _, r := range res {
		targets[r.Target] = true
	}
	if !targets["production"] || !targets["staging"] {
		t.Error("expected results for both targets")
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	res, _ := truncate.Target(cfg, "production", 5, false)
	out := truncate.Format(res)
	if !strings.Contains(out, "truncated") && !strings.Contains(out, "within limit") {
		t.Errorf("unexpected format output: %s", out)
	}
}

func TestFormatNoChanges(t *testing.T) {
	cfg := newCfg()
	res, _ := truncate.Target(cfg, "production", 1000, false)
	out := truncate.Format(res)
	if !strings.Contains(out, "within limit") {
		t.Errorf("expected 'within limit' message, got: %s", out)
	}
}
