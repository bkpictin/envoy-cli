package reorder_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/reorder"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"ZEBRA": "z",
				"APPLE": "a",
				"MANGO": "m",
			},
		},
	}
}

func TestAlphabetical(t *testing.T) {
	cfg := newCfg()
	res, err := reorder.Alphabetical(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, k := range res.After {
		if k != expected[i] {
			t.Errorf("After[%d] = %q, want %q", i, k, expected[i])
		}
	}
}

func TestAlphabeticalMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := reorder.Alphabetical(cfg, "staging")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestCustomOrder(t *testing.T) {
	cfg := newCfg()
	order := []string{"MANGO", "ZEBRA", "APPLE"}
	res, err := reorder.Custom(cfg, "production", order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, k := range res.After {
		if k != order[i] {
			t.Errorf("After[%d] = %q, want %q", i, k, order[i])
		}
	}
}

func TestCustomOrderPartial(t *testing.T) {
	cfg := newCfg()
	// Only specify two keys; APPLE should be appended alphabetically.
	order := []string{"ZEBRA", "MANGO"}
	res, err := reorder.Custom(cfg, "production", order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.After[0] != "ZEBRA" || res.After[1] != "MANGO" || res.After[2] != "APPLE" {
		t.Errorf("unexpected order: %v", res.After)
	}
}

func TestCustomOrderMissingKey(t *testing.T) {
	cfg := newCfg()
	_, err := reorder.Custom(cfg, "production", []string{"NONEXISTENT"})
	if err == nil {
		t.Fatal("expected error for missing key in order list")
	}
}

func TestCustomOrderMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := reorder.Custom(cfg, "staging", []string{"KEY"})
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}
