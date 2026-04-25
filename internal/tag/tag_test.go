package tag_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/tag"
)

func newCfg() *config.Config {
	cfg := &config.Config{
		Targets: map[string]config.Target{
			"prod": {
				Vars: map[string]string{"DB_URL": "postgres://", "API_KEY": "secret"},
				Tags: map[string][]string{},
			},
		},
	}
	return cfg
}

func TestAddAndListForKey(t *testing.T) {
	cfg := newCfg()
	if err := tag.Add(cfg, "prod", "DB_URL", "sensitive"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, err := tag.ListForKey(cfg, "prod", "DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 || tags[0] != "sensitive" {
		t.Fatalf("expected [sensitive], got %v", tags)
	}
}

func TestAddDuplicateTag(t *testing.T) {
	cfg := newCfg()
	_ = tag.Add(cfg, "prod", "DB_URL", "sensitive")
	_ = tag.Add(cfg, "prod", "DB_URL", "sensitive")
	tags, _ := tag.ListForKey(cfg, "prod", "DB_URL")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
}

func TestRemoveTag(t *testing.T) {
	cfg := newCfg()
	_ = tag.Add(cfg, "prod", "API_KEY", "secret")
	if err := tag.Remove(cfg, "prod", "API_KEY", "secret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, _ := tag.ListForKey(cfg, "prod", "API_KEY")
	if len(tags) != 0 {
		t.Fatalf("expected no tags, got %v", tags)
	}
}

func TestRemoveMissingTag(t *testing.T) {
	cfg := newCfg()
	err := tag.Remove(cfg, "prod", "DB_URL", "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing tag")
	}
}

func TestListByTag(t *testing.T) {
	cfg := newCfg()
	_ = tag.Add(cfg, "prod", "DB_URL", "sensitive")
	_ = tag.Add(cfg, "prod", "API_KEY", "sensitive")
	keys, err := tag.ListByTag(cfg, "prod", "sensitive")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestAddMissingTarget(t *testing.T) {
	cfg := newCfg()
	err := tag.Add(cfg, "staging", "DB_URL", "sensitive")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestAddMissingKey(t *testing.T) {
	cfg := newCfg()
	err := tag.Add(cfg, "prod", "MISSING_KEY", "sensitive")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}
