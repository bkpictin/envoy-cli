package copy

import (
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev":  {"FOO": "bar", "SHARED": "dev-val"},
			"prod": {"SHARED": "prod-val"},
		},
	}
}

func TestCopyEnvs(t *testing.T) {
	cfg := newCfg()
	if err := CopyEnvs(cfg, "dev", "prod", false); err != nil {
		t.Fatal(err)
	}
	if cfg.Targets["prod"]["FOO"] != "bar" {
		t.Error("expected FOO to be copied to prod")
	}
	// SHARED should not be overwritten
	if cfg.Targets["prod"]["SHARED"] != "prod-val" {
		t.Error("expected SHARED to remain prod-val")
	}
}

func TestCopyEnvsOverwrite(t *testing.T) {
	cfg := newCfg()
	if err := CopyEnvs(cfg, "dev", "prod", true); err != nil {
		t.Fatal(err)
	}
	if cfg.Targets["prod"]["SHARED"] != "dev-val" {
		t.Error("expected SHARED to be overwritten with dev-val")
	}
}

func TestCopyMissingSource(t *testing.T) {
	cfg := newCfg()
	if err := CopyEnvs(cfg, "staging", "prod", false); err == nil {
		t.Error("expected error for missing source")
	}
}

func TestCopyMissingDest(t *testing.T) {
	cfg := newCfg()
	if err := CopyEnvs(cfg, "dev", "staging", false); err == nil {
		t.Error("expected error for missing destination")
	}
}

func TestMergeEnvs(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"base":   {"A": "1", "B": "2"},
			"extra":  {"C": "3"},
			"merged": {},
		},
	}
	if err := MergeEnvs(cfg, "merged", false, "base", "extra"); err != nil {
		t.Fatal(err)
	}
	if cfg.Targets["merged"]["A"] != "1" || cfg.Targets["merged"]["C"] != "3" {
		t.Error("expected merged to contain keys from base and extra")
	}
}
