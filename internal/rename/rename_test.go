package rename_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/rename"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {"DB_HOST": "localhost", "PORT": "5432"},
			"prod": {"DB_HOST": "prod.db", "TIMEOUT": "30"},
		},
	}
}

func TestKeyInTarget(t *testing.T) {
	cfg := newCfg()
	if err := rename.KeyInTarget(cfg, "dev", "PORT", "APP_PORT"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Targets["dev"]["PORT"]; ok {
		t.Error("old key PORT should be removed")
	}
	if v := cfg.Targets["dev"]["APP_PORT"]; v != "5432" {
		t.Errorf("expected APP_PORT=5432, got %q", v)
	}
	// prod should be untouched
	if _, ok := cfg.Targets["prod"]["PORT"]; ok {
		t.Error("prod should not have PORT key")
	}
}

func TestKeyInTargetMissingTarget(t *testing.T) {
	cfg := newCfg()
	err := rename.KeyInTarget(cfg, "staging", "PORT", "APP_PORT")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestKeyInTargetMissingKey(t *testing.T) {
	cfg := newCfg()
	err := rename.KeyInTarget(cfg, "dev", "MISSING", "NEW_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestKeyInTargetConflict(t *testing.T) {
	cfg := newCfg()
	err := rename.KeyInTarget(cfg, "dev", "PORT", "DB_HOST")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestKeyInAll(t *testing.T) {
	cfg := newCfg()
	if err := rename.KeyInAll(cfg, "DB_HOST", "DATABASE_HOST"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, target := range []string{"dev", "prod"} {
		if _, ok := cfg.Targets[target]["DB_HOST"]; ok {
			t.Errorf("old key DB_HOST should be removed from %s", target)
		}
		if _, ok := cfg.Targets[target]["DATABASE_HOST"]; !ok {
			t.Errorf("new key DATABASE_HOST missing from %s", target)
		}
	}
}

func TestKeyInAllConflictAbortsAll(t *testing.T) {
	cfg := newCfg()
	// TIMEOUT exists only in prod; introduce a conflict in prod by renaming DB_HOST -> TIMEOUT
	err := rename.KeyInAll(cfg, "DB_HOST", "TIMEOUT")
	if err == nil {
		t.Fatal("expected error due to conflict in prod")
	}
	// dev should still have original keys (operation aborted before any mutation)
	if _, ok := cfg.Targets["dev"]["DB_HOST"]; !ok {
		t.Error("dev DB_HOST should be preserved after aborted rename")
	}
}
