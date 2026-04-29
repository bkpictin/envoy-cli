package suffix_test

import (
	"sort"
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/suffix"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"prod": {
				"DB_HOST": "db.prod",
				"DB_PORT": "5432",
				"API_KEY": "secret",
			},
		},
	}
}

func TestAddSuffix(t *testing.T) {
	cfg := newCfg()
	results, err := suffix.Add(cfg, "prod", "_V2", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Skipped {
			t.Errorf("key %q should not be skipped", r.OldKey)
		}
	}
	if _, ok := cfg.Targets["prod"]["DB_HOST_V2"]; !ok {
		t.Error("expected DB_HOST_V2 to exist")
	}
	if _, ok := cfg.Targets["prod"]["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be removed")
	}
}

func TestAddSuffixDryRun(t *testing.T) {
	cfg := newCfg()
	_, err := suffix.Add(cfg, "prod", "_V2", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Targets["prod"]["DB_HOST"]; !ok {
		t.Error("dry-run should not mutate config")
	}
}

func TestAddSuffixAlreadyHas(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"prod": {"DB_HOST_V2": "db.prod"},
		},
	}
	results, err := suffix.Add(cfg, "prod", "_V2", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Error("expected key to be skipped as it already has the suffix")
	}
}

func TestRemoveSuffix(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"prod": {
				"DB_HOST_OLD": "db.prod",
				"DB_PORT_OLD": "5432",
				"API_KEY": "secret",
			},
		},
	}
	results, err := suffix.Remove(cfg, "prod", "_OLD", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var changed, skipped int
	for _, r := range results {
		if r.Skipped {
			skipped++
		} else {
			changed++
		}
	}
	if changed != 2 {
		t.Errorf("expected 2 changed, got %d", changed)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", skipped)
	}
	if _, ok := cfg.Targets["prod"]["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to exist after remove")
	}
}

func TestRemoveSuffixMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := suffix.Remove(cfg, "staging", "_OLD", false)
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestList(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"prod": {
				"DB_HOST_PROD": "db",
				"DB_PORT_PROD": "5432",
				"API_KEY": "secret",
			},
		},
	}
	keys, err := suffix.List(cfg, "prod", "_PROD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "DB_HOST_PROD" || keys[1] != "DB_PORT_PROD" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestAddEmptySuffix(t *testing.T) {
	cfg := newCfg()
	_, err := suffix.Add(cfg, "prod", "", false)
	if err == nil {
		t.Error("expected error for empty suffix")
	}
}
