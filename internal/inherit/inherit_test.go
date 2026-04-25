package inherit_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/inherit"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"base": {"HOST": "localhost", "PORT": "5432", "DEBUG": "true"},
			"staging": {"HOST": "staging.example.com"},
			"prod": {},
		},
	}
}

func TestApplySkipsExisting(t *testing.T) {
	cfg := newCfg()
	results, err := inherit.Apply(cfg, "base", []string{"staging"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", results[0].Skipped)
	}
	if cfg.Targets["staging"]["HOST"] != "staging.example.com" {
		t.Errorf("existing key should not be overwritten")
	}
}

func TestApplyOverwrite(t *testing.T) {
	cfg := newCfg()
	_, err := inherit.Apply(cfg, "base", []string{"staging"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["staging"]["HOST"] != "localhost" {
		t.Errorf("expected HOST to be overwritten to 'localhost'")
	}
}

func TestApplyMissingBase(t *testing.T) {
	cfg := newCfg()
	_, err := inherit.Apply(cfg, "missing", []string{"staging"}, false)
	if err == nil {
		t.Fatal("expected error for missing base target")
	}
}

func TestApplyMissingChild(t *testing.T) {
	cfg := newCfg()
	_, err := inherit.Apply(cfg, "base", []string{"ghost"}, false)
	if err == nil {
		t.Fatal("expected error for missing child target")
	}
}

func TestApplySameBaseAndChild(t *testing.T) {
	cfg := newCfg()
	_, err := inherit.Apply(cfg, "base", []string{"base"}, false)
	if err == nil {
		t.Fatal("expected error when base and child are the same")
	}
}

func TestListInherited(t *testing.T) {
	cfg := newCfg()
	// prod has no keys, staging has HOST
	shared, err := inherit.ListInherited(cfg, "base", "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shared) != 1 || shared[0] != "HOST" {
		t.Errorf("expected [HOST], got %v", shared)
	}
}

func TestListInheritedEmpty(t *testing.T) {
	cfg := newCfg()
	shared, err := inherit.ListInherited(cfg, "base", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(shared) != 0 {
		t.Errorf("expected no shared keys, got %v", shared)
	}
}
