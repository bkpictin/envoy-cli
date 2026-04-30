package cascade_test

import (
	"testing"

	"envoy-cli/internal/cascade"
	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"base": {"DB_HOST": "localhost", "DB_PORT": "5432", "LOG_LEVEL": "debug"},
			"staging": {"LOG_LEVEL": "info"},
			"prod": {},
		},
	}
}

func TestCascadeApplyNoOverwrite(t *testing.T) {
	cfg := newCfg()
	results, err := cascade.Apply(cfg, "base", []string{"staging", "prod"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	// staging already had LOG_LEVEL, so it should be skipped
	if cfg.Targets["staging"]["LOG_LEVEL"] != "info" {
		t.Errorf("staging LOG_LEVEL should remain 'info', got %q", cfg.Targets["staging"]["LOG_LEVEL"])
	}
	if cfg.Targets["staging"]["DB_HOST"] != "localhost" {
		t.Errorf("staging DB_HOST should be 'localhost', got %q", cfg.Targets["staging"]["DB_HOST"])
	}
	if results[0].Skipped != 1 {
		t.Errorf("expected 1 skipped for staging, got %d", results[0].Skipped)
	}
	// prod had nothing, all keys should be applied
	if cfg.Targets["prod"]["LOG_LEVEL"] != "info" {
		t.Errorf("prod LOG_LEVEL should be 'info' (propagated from staging), got %q", cfg.Targets["prod"]["LOG_LEVEL"])
	}
}

func TestCascadeApplyWithOverwrite(t *testing.T) {
	cfg := newCfg()
	_, err := cascade.Apply(cfg, "base", []string{"staging"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["staging"]["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL to be overwritten to 'debug', got %q", cfg.Targets["staging"]["LOG_LEVEL"])
	}
}

func TestCascadeMissingBase(t *testing.T) {
	cfg := newCfg()
	_, err := cascade.Apply(cfg, "nonexistent", []string{"staging"}, false)
	if err == nil {
		t.Fatal("expected error for missing base target")
	}
}

func TestCascadeMissingChainTarget(t *testing.T) {
	cfg := newCfg()
	_, err := cascade.Apply(cfg, "base", []string{"staging", "ghost"}, false)
	if err == nil {
		t.Fatal("expected error for missing chain target")
	}
}

func TestCascadeEmptyChain(t *testing.T) {
	cfg := newCfg()
	results, err := cascade.Apply(cfg, "base", []string{}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty chain, got %d", len(results))
	}
}
