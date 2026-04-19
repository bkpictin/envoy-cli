package snapshot

import (
	"testing"

	"github.com/user/envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"prod": {"KEY": "value1"},
		},
	}
}

func TestCreate(t *testing.T) {
	cfg := newCfg()
	snap, err := Create(cfg, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Target != "prod" {
		t.Errorf("expected target prod, got %s", snap.Target)
	}
	if snap.Vars["KEY"] != "value1" {
		t.Errorf("expected KEY=value1")
	}
	if len(cfg.Snapshots["prod"]) != 1 {
		t.Errorf("expected 1 snapshot stored")
	}
}

func TestCreateMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := Create(cfg, "staging")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestList(t *testing.T) {
	cfg := newCfg()
	Create(cfg, "prod")
	cfg.Targets["prod"]["KEY"] = "value2"
	Create(cfg, "prod")

	snaps, err := List(cfg, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snaps) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(snaps))
	}
}

func TestRestore(t *testing.T) {
	cfg := newCfg()
	Create(cfg, "prod")
	cfg.Targets["prod"]["KEY"] = "value2"

	if err := Restore(cfg, "prod", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["prod"]["KEY"] != "value1" {
		t.Errorf("expected KEY restored to value1, got %s", cfg.Targets["prod"]["KEY"])
	}
}

func TestRestoreOutOfRange(t *testing.T) {
	cfg := newCfg()
	Create(cfg, "prod")
	if err := Restore(cfg, "prod", 5); err == nil {
		t.Fatal("expected out-of-range error")
	}
}
