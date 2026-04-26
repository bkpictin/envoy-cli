package dedupe_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/dedupe"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{},
	}
}

func TestNoDuplicates(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"HOST": "localhost", "PORT": "3000"}
	cfg.Targets["prod"] = map[string]string{"HOST": "example.com", "PORT": "443"}

	r := dedupe.FindCrossTarget(cfg)
	if len(r.Matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(r.Matches))
	}
}

func TestFindsDuplicates(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"DB_PASS": "secret", "PORT": "3000"}
	cfg.Targets["staging"] = map[string]string{"DB_PASS": "secret", "PORT": "8080"}
	cfg.Targets["prod"] = map[string]string{"DB_PASS": "secret", "PORT": "443"}

	r := dedupe.FindCrossTarget(cfg)
	if len(r.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d: %+v", len(r.Matches), r.Matches)
	}
	m := r.Matches[0]
	if m.Key != "DB_PASS" {
		t.Errorf("expected key DB_PASS, got %s", m.Key)
	}
	if m.Value != "secret" {
		t.Errorf("expected value secret, got %s", m.Value)
	}
	if len(m.Targets) != 3 {
		t.Errorf("expected 3 targets, got %v", m.Targets)
	}
}

func TestMultipleDuplicatePairs(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["a"] = map[string]string{"X": "1", "Y": "same"}
	cfg.Targets["b"] = map[string]string{"X": "1", "Y": "same"}

	r := dedupe.FindCrossTarget(cfg)
	if len(r.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(r.Matches))
	}
}

func TestFormatNoDuplicates(t *testing.T) {
	r := dedupe.Result{}
	out := dedupe.Format(r)
	if !strings.Contains(out, "No duplicate") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatWithDuplicates(t *testing.T) {
	r := dedupe.Result{
		Matches: []dedupe.Match{
			{Key: "API_KEY", Value: "abc123", Targets: []string{"dev", "staging"}},
		},
	}
	out := dedupe.Format(r)
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output: %s", out)
	}
	if !strings.Contains(out, "1 duplicate") {
		t.Errorf("expected count in output: %s", out)
	}
}

func TestEmptyConfig(t *testing.T) {
	cfg := newCfg()
	r := dedupe.FindCrossTarget(cfg)
	if len(r.Matches) != 0 {
		t.Fatalf("expected 0 matches for empty config")
	}
}
