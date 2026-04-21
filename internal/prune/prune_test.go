package prune_test

import (
	"testing"

	"github.com/envoy-cli/internal/config"
	"github.com/envoy-cli/internal/prune"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"APP_PORT":   "8080",
				"DB_HOST":    "localhost",
				"ONLY_DEV":  "secret",
				"EMPTY_KEY": "",
			},
			"prod": {
				"APP_PORT": "80",
				"DB_HOST":  "db.prod",
			},
		},
	}
}

func TestOrphanedKeysDryRun(t *testing.T) {
	cfg := newCfg()
	results, err := prune.OrphanedKeys(cfg, "dev", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 { // ONLY_DEV and EMPTY_KEY are orphaned
		t.Fatalf("expected 2 orphaned keys, got %d", len(results))
	}
	// dry-run must not mutate
	if _, ok := cfg.Targets["dev"]["ONLY_DEV"]; !ok {
		t.Error("dry-run should not delete ONLY_DEV")
	}
}

func TestOrphanedKeysApply(t *testing.T) {
	cfg := newCfg()
	_, err := prune.OrphanedKeys(cfg, "dev", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Targets["dev"]["ONLY_DEV"]; ok {
		t.Error("ONLY_DEV should have been pruned")
	}
	if _, ok := cfg.Targets["dev"]["APP_PORT"]; !ok {
		t.Error("APP_PORT should be retained")
	}
}

func TestOrphanedKeysMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := prune.OrphanedKeys(cfg, "staging", false)
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestEmptyValues(t *testing.T) {
	cfg := newCfg()
	results, err := prune.EmptyValues(cfg, "dev", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 empty key, got %d", len(results))
	}
	if results[0].Key != "EMPTY_KEY" {
		t.Errorf("expected EMPTY_KEY, got %s", results[0].Key)
	}
	if _, ok := cfg.Targets["dev"]["EMPTY_KEY"]; ok {
		t.Error("EMPTY_KEY should have been deleted")
	}
}

func TestEmptyValuesDryRun(t *testing.T) {
	cfg := newCfg()
	_, err := prune.EmptyValues(cfg, "dev", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Targets["dev"]["EMPTY_KEY"]; !ok {
		t.Error("dry-run should not delete EMPTY_KEY")
	}
}
