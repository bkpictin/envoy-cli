package group_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/group"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {"DB_HOST": "db.prod", "DB_PORT": "5432", "API_KEY": "secret"},
		},
	}
}

func TestCreateAndList(t *testing.T) {
	cfg := newCfg()
	if err := group.Create(cfg, "production", "database"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names, err := group.ListGroups(cfg, "production")
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(names) != 1 || names[0] != "database" {
		t.Fatalf("expected [database], got %v", names)
	}
}

func TestCreateDuplicate(t *testing.T) {
	cfg := newCfg()
	_ = group.Create(cfg, "production", "database")
	err := group.Create(cfg, "production", "database")
	if err == nil {
		t.Fatal("expected error for duplicate group")
	}
}

func TestCreateMissingTarget(t *testing.T) {
	cfg := newCfg()
	err := group.Create(cfg, "staging", "database")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestAddAndGetKeys(t *testing.T) {
	cfg := newCfg()
	_ = group.Create(cfg, "production", "database")
	_ = group.AddKey(cfg, "production", "database", "DB_HOST")
	_ = group.AddKey(cfg, "production", "database", "DB_PORT")

	keys, err := group.GetKeys(cfg, "production", "database")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PORT" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestAddKeyIdempotent(t *testing.T) {
	cfg := newCfg()
	_ = group.Create(cfg, "production", "database")
	_ = group.AddKey(cfg, "production", "database", "DB_HOST")
	_ = group.AddKey(cfg, "production", "database", "DB_HOST")
	keys, _ := group.GetKeys(cfg, "production", "database")
	if len(keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(keys))
	}
}

func TestRemoveKey(t *testing.T) {
	cfg := newCfg()
	_ = group.Create(cfg, "production", "database")
	_ = group.AddKey(cfg, "production", "database", "DB_HOST")
	_ = group.AddKey(cfg, "production", "database", "DB_PORT")
	if err := group.RemoveKey(cfg, "production", "database", "DB_HOST"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys, _ := group.GetKeys(cfg, "production", "database")
	if len(keys) != 1 || keys[0] != "DB_PORT" {
		t.Fatalf("expected [DB_PORT], got %v", keys)
	}
}

func TestRemoveMissingKey(t *testing.T) {
	cfg := newCfg()
	_ = group.Create(cfg, "production", "database")
	err := group.RemoveKey(cfg, "production", "database", "NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestDeleteGroup(t *testing.T) {
	cfg := newCfg()
	_ = group.Create(cfg, "production", "database")
	if err := group.Delete(cfg, "production", "database"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names, _ := group.ListGroups(cfg, "production")
	if len(names) != 0 {
		t.Fatalf("expected empty list after delete, got %v", names)
	}
}

func TestDeleteMissingGroup(t *testing.T) {
	cfg := newCfg()
	err := group.Delete(cfg, "production", "nope")
	if err == nil {
		t.Fatal("expected error for missing group")
	}
}
