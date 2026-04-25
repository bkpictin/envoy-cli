package prefix_test

import (
	"sort"
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/prefix"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"prod": {
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
				"APP_NAME": "envoy",
			},
		},
	}
}

func TestAddPrefix(t *testing.T) {
	cfg := newCfg()
	res, err := prefix.Add(cfg, "prod", "PROD_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Changed != 3 {
		t.Errorf("expected 3 changed, got %d", res.Changed)
	}
	if _, ok := cfg.Targets["prod"]["PROD_DB_HOST"]; !ok {
		t.Error("expected PROD_DB_HOST to exist")
	}
}

func TestAddPrefixSkipsAlreadyPrefixed(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["prod"]["PROD_EXISTING"] = "yes"
	res, err := prefix.Add(cfg, "prod", "PROD_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestAddPrefixMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := prefix.Add(cfg, "staging", "STG_", false)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestAddPrefixEmptyPrefix(t *testing.T) {
	cfg := newCfg()
	_, err := prefix.Add(cfg, "prod", "", false)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestRemovePrefix(t *testing.T) {
	cfg := newCfg()
	prefix.Add(cfg, "prod", "PROD_", false) //nolint:errcheck
	res, err := prefix.Remove(cfg, "prod", "PROD_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Changed != 3 {
		t.Errorf("expected 3 changed, got %d", res.Changed)
	}
	if _, ok := cfg.Targets["prod"]["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be restored")
	}
}

func TestRemovePrefixMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := prefix.Remove(cfg, "nope", "X_")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestList(t *testing.T) {
	cfg := newCfg()
	keys, err := prefix.List(cfg, "prod", "DB_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sort.Strings(keys)
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(keys), keys)
	}
}

func TestListMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := prefix.List(cfg, "ghost", "DB_")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}
