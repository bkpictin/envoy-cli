package group_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/group"
)

// TestMultipleGroupsPerTarget verifies that a target can hold several
// independent groups without interference.
func TestMultipleGroupsPerTarget(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"DB_HOST": "db.prod",
				"DB_PORT": "5432",
				"REDIS_URL": "redis://prod",
				"API_KEY": "secret",
			},
		},
	}

	_ = group.Create(cfg, "production", "database")
	_ = group.Create(cfg, "production", "cache")
	_ = group.Create(cfg, "production", "auth")

	_ = group.AddKey(cfg, "production", "database", "DB_HOST")
	_ = group.AddKey(cfg, "production", "database", "DB_PORT")
	_ = group.AddKey(cfg, "production", "cache", "REDIS_URL")
	_ = group.AddKey(cfg, "production", "auth", "API_KEY")

	names, err := group.ListGroups(cfg, "production")
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(names) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(names))
	}

	dbKeys, _ := group.GetKeys(cfg, "production", "database")
	if len(dbKeys) != 2 {
		t.Fatalf("expected 2 db keys, got %d", len(dbKeys))
	}

	cacheKeys, _ := group.GetKeys(cfg, "production", "cache")
	if len(cacheKeys) != 1 || cacheKeys[0] != "REDIS_URL" {
		t.Fatalf("unexpected cache keys: %v", cacheKeys)
	}
}

// TestKeyInMultipleGroups verifies that a key may belong to more than one group.
func TestKeyInMultipleGroups(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"staging": {"SHARED_SECRET": "abc123"},
		},
	}

	_ = group.Create(cfg, "staging", "auth")
	_ = group.Create(cfg, "staging", "security")
	_ = group.AddKey(cfg, "staging", "auth", "SHARED_SECRET")
	_ = group.AddKey(cfg, "staging", "security", "SHARED_SECRET")

	for _, g := range []string{"auth", "security"} {
		keys, err := group.GetKeys(cfg, "staging", g)
		if err != nil {
			t.Fatalf("error getting keys for group %q: %v", g, err)
		}
		if len(keys) != 1 || keys[0] != "SHARED_SECRET" {
			t.Fatalf("group %q: expected [SHARED_SECRET], got %v", g, keys)
		}
	}
}
