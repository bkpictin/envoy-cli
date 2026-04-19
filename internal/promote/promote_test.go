package promote

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"staging": {"DB_HOST": "staging-db", "API_KEY": "abc", "PORT": "8080"},
			"production": {"DB_HOST": "prod-db"},
		},
	}
}

func TestPromoteAll(t *testing.T) {
	cfg := newCfg()
	promoted, err := Promote(cfg, "staging", "production", Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// DB_HOST exists and Overwrite=false, so only API_KEY and PORT should be promoted
	if len(promoted) != 2 {
		t.Errorf("expected 2 promoted keys, got %d", len(promoted))
	}
	if cfg.Targets["production"]["DB_HOST"] != "prod-db" {
		t.Error("DB_HOST should not have been overwritten")
	}
}

func TestPromoteWithOverwrite(t *testing.T) {
	cfg := newCfg()
	_, err := Promote(cfg, "staging", "production", Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["production"]["DB_HOST"] != "staging-db" {
		t.Error("DB_HOST should have been overwritten")
	}
}

func TestPromoteFilteredKeys(t *testing.T) {
	cfg := newCfg()
	promoted, err := Promote(cfg, "staging", "production", Options{Overwrite: true, Keys: []string{"PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(promoted) != 1 || promoted[0] != "PORT" {
		t.Errorf("expected only PORT to be promoted, got %v", promoted)
	}
}

func TestPromoteMissingSource(t *testing.T) {
	cfg := newCfg()
	_, err := Promote(cfg, "dev", "production", Options{})
	if err == nil {
		t.Error("expected error for missing source target")
	}
}

func TestPromoteMissingDest(t *testing.T) {
	cfg := newCfg()
	_, err := Promote(cfg, "staging", "dev", Options{})
	if err == nil {
		t.Error("expected error for missing destination target")
	}
}
