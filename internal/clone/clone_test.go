package clone_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/clone"
	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"staging": {"DB_HOST": "localhost", "PORT": "5432", "SECRET": "abc"},
		},
	}
}

func TestCloneTarget(t *testing.T) {
	cfg := newCfg()
	if err := clone.Target(cfg, "staging", "production", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Targets["production"]) != 3 {
		t.Errorf("expected 3 keys, got %d", len(cfg.Targets["production"]))
	}
	if cfg.Targets["production"]["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost")
	}
}

func TestCloneTargetMissingSource(t *testing.T) {
	cfg := newCfg()
	err := clone.Target(cfg, "nonexistent", "production", false)
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}

func TestCloneTargetDestExists(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["production"] = map[string]string{"X": "1"}
	err := clone.Target(cfg, "staging", "production", false)
	if err == nil {
		t.Fatal("expected error when dest exists without overwrite")
	}
}

func TestCloneTargetOverwrite(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["production"] = map[string]string{"X": "1"}
	if err := clone.Target(cfg, "staging", "production", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := cfg.Targets["production"]["DB_HOST"]; !ok {
		t.Error("expected DB_HOST after overwrite")
	}
}

func TestCloneWithFilter(t *testing.T) {
	cfg := newCfg()
	err := clone.WithFilter(cfg, "staging", "production", false, func(key string) bool {
		return strings.HasPrefix(key, "DB_")
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Targets["production"]) != 1 {
		t.Errorf("expected 1 key, got %d", len(cfg.Targets["production"]))
	}
	if cfg.Targets["production"]["DB_HOST"] != "localhost" {
		t.Error("expected DB_HOST=localhost")
	}
}
