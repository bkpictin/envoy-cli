package unique

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"SHARED_KEY": "dev-val",
				"DEV_ONLY":   "secret",
			},
			"staging": {
				"SHARED_KEY":     "staging-val",
				"STAGING_ONLY":   "hello",
			},
			"prod": {
				"SHARED_KEY": "prod-val",
				"PROD_ONLY":  "world",
			},
		},
	}
}

func TestFindUniqueKeys(t *testing.T) {
	cfg := newCfg()
	results, err := Find(cfg, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect DEV_ONLY, STAGING_ONLY, PROD_ONLY — 3 unique keys.
	if len(results) != 3 {
		t.Fatalf("expected 3 unique keys, got %d", len(results))
	}
	keys := make(map[string]string)
	for _, r := range results {
		keys[r.Key] = r.Target
	}
	for _, expected := range []string{"DEV_ONLY", "STAGING_ONLY", "PROD_ONLY"} {
		if _, ok := keys[expected]; !ok {
			t.Errorf("expected key %q in results", expected)
		}
	}
}

func TestFindUniqueScopedTargets(t *testing.T) {
	cfg := newCfg()
	// Only scan dev and staging — PROD_ONLY should not appear.
	results, err := Find(cfg, []string{"dev", "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "PROD_ONLY" {
			t.Error("PROD_ONLY should not appear when prod is excluded from scope")
		}
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 unique keys in scoped search, got %d", len(results))
	}
}

func TestFindUniqueMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := Find(cfg, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestFormat(t *testing.T) {
	results := []Result{
		{Target: "dev", Key: "DEV_ONLY", Value: "secret"},
	}
	out := Format(results)
	if !strings.Contains(out, "[dev]") || !strings.Contains(out, "DEV_ONLY") {
		t.Errorf("unexpected format output: %q", out)
	}
}

func TestFormatEmpty(t *testing.T) {
	out := Format(nil)
	if out != "no unique keys found" {
		t.Errorf("expected empty message, got %q", out)
	}
}
