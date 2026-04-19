package target_test

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/config"
	"github.com/envoy-cli/envoy-cli/internal/target"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {"KEY": "val"},
		},
	}
}

func TestList(t *testing.T) {
	cfg := newCfg()
	list := target.List(cfg)
	if len(list) != 1 || list[0] != "production" {
		t.Fatalf("expected [production], got %v", list)
	}
}

func TestAdd(t *testing.T) {
	cfg := newCfg()
	if err := target.Add(cfg, "staging"); err != nil {
		t.Fatal(err)
	}
	if _, ok := cfg.Targets["staging"]; !ok {
		t.Fatal("staging target not created")
	}
}

func TestAddDuplicate(t *testing.T) {
	cfg := newCfg()
	if err := target.Add(cfg, "production"); err == nil {
		t.Fatal("expected error for duplicate target")
	}
}

func TestRemove(t *testing.T) {
	cfg := newCfg()
	if err := target.Remove(cfg, "production"); err != nil {
		t.Fatal(err)
	}
	if _, ok := cfg.Targets["production"]; ok {
		t.Fatal("target should have been removed")
	}
}

func TestRemoveMissing(t *testing.T) {
	cfg := newCfg()
	if err := target.Remove(cfg, "ghost"); err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestRename(t *testing.T) {
	cfg := newCfg()
	if err := target.Rename(cfg, "production", "prod"); err != nil {
		t.Fatal(err)
	}
	if _, ok := cfg.Targets["prod"]; !ok {
		t.Fatal("renamed target not found")
	}
	if cfg.Targets["prod"]["KEY"] != "val" {
		t.Fatal("variables not preserved after rename")
	}
}
