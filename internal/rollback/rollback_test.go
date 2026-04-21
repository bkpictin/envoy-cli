package rollback_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/rollback"
	"envoy-cli/internal/snapshot"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"prod": {"DB_HOST": "prod-db", "API_KEY": "secret"},
		},
		Snapshots: map[string][]config.Snapshot{},
	}
}

func TestToPrevious(t *testing.T) {
	cfg := newCfg()
	if err := snapshot.Create(cfg, "prod", "snap1"); err != nil {
		t.Fatal(err)
	}
	cfg.Targets["prod"]["DB_HOST"] = "changed-db"
	name, err := rollback.ToPrevious(cfg, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "snap1" {
		t.Errorf("expected snap1, got %s", name)
	}
	if cfg.Targets["prod"]["DB_HOST"] != "prod-db" {
		t.Errorf("expected prod-db after rollback, got %s", cfg.Targets["prod"]["DB_HOST"])
	}
}

func TestToPreviousNoSnapshots(t *testing.T) {
	cfg := newCfg()
	_, err := rollback.ToPrevious(cfg, "prod")
	if err != rollback.ErrNoSnapshots {
		t.Errorf("expected ErrNoSnapshots, got %v", err)
	}
}

func TestToSnapshot(t *testing.T) {
	cfg := newCfg()
	_ = snapshot.Create(cfg, "prod", "v1")
	_ = snapshot.Create(cfg, "prod", "v2")
	cfg.Targets["prod"]["DB_HOST"] = "new-db"
	if err := rollback.ToSnapshot(cfg, "prod", "v1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["prod"]["DB_HOST"] != "prod-db" {
		t.Errorf("expected prod-db, got %s", cfg.Targets["prod"]["DB_HOST"])
	}
}

func TestToSnapshotMissing(t *testing.T) {
	cfg := newCfg()
	_ = snapshot.Create(cfg, "prod", "v1")
	if err := rollback.ToSnapshot(cfg, "prod", "nonexistent"); err == nil {
		t.Error("expected error for missing snapshot")
	}
}

func TestListAvailable(t *testing.T) {
	cfg := newCfg()
	_ = snapshot.Create(cfg, "prod", "v1")
	_ = snapshot.Create(cfg, "prod", "v2")
	names, err := rollback.ListAvailable(cfg, "prod")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 || names[0] != "v1" || names[1] != "v2" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestListAvailableNoSnapshots(t *testing.T) {
	cfg := newCfg()
	_, err := rollback.ListAvailable(cfg, "prod")
	if err != rollback.ErrNoSnapshots {
		t.Errorf("expected ErrNoSnapshots, got %v", err)
	}
}
