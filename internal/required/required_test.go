package required

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
)

func newCfg(targets map[string]map[string]string) *config.Config {
	cfg := &config.Config{Targets: make(map[string]map[string]string)}
	for t, envs := range targets {
		cfg.Targets[t] = envs
	}
	return cfg
}

func TestNoMissingKeys(t *testing.T) {
	cfg := newCfg(map[string]map[string]string{
		"dev":  {"DB_URL": "dev-db", "PORT": "3000"},
		"prod": {"DB_URL": "prod-db", "PORT": "443"},
	})
	results, err := Check(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if len(r.Missing) != 0 {
			t.Errorf("target %s should have no missing keys, got %v", r.Target, r.Missing)
		}
	}
}

func TestMissingKeyDetected(t *testing.T) {
	cfg := newCfg(map[string]map[string]string{
		"dev":     {"DB_URL": "dev-db", "PORT": "3000", "SECRET": "x"},
		"staging": {"DB_URL": "stg-db", "PORT": "8080"},
		"prod":    {"DB_URL": "prod-db", "PORT": "443", "SECRET": "y"},
	})
	results, err := Check(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range results {
		if r.Target == "staging" {
			found = true
			if len(r.Missing) != 1 || r.Missing[0] != "SECRET" {
				t.Errorf("expected staging to be missing SECRET, got %v", r.Missing)
			}
		}
	}
	if !found {
		t.Error("staging result not found")
	}
}

func TestEmptyConfig(t *testing.T) {
	cfg := newCfg(map[string]map[string]string{})
	results, err := Check(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for empty config, got %d", len(results))
	}
}

func TestSingleTarget(t *testing.T) {
	cfg := newCfg(map[string]map[string]string{
		"dev": {"DB_URL": "dev-db"},
	})
	results, err := Check(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || len(results[0].Missing) != 0 {
		t.Errorf("single target should have no missing keys: %+v", results)
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg(map[string]map[string]string{
		"dev":  {"DB_URL": "x", "PORT": "3000"},
		"prod": {"DB_URL": "y"},
	})
	results, _ := Check(cfg)
	out := Format(results)
	if !strings.Contains(out, "[MISS]") {
		t.Errorf("expected [MISS] marker in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got:\n%s", out)
	}
}

func TestFormatNoTargets(t *testing.T) {
	out := Format(nil)
	if !strings.Contains(out, "no targets") {
		t.Errorf("expected 'no targets' message, got: %s", out)
	}
}
