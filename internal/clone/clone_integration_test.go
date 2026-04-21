package clone_test

import (
	"testing"

	"envoy-cli/internal/clone"
	"envoy-cli/internal/config"
)

// TestClonePreservesIsolation verifies that mutating the cloned target does
// not affect the original source (deep-copy semantics).
func TestClonePreservesIsolation(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"staging": {"KEY": "original"},
		},
	}

	if err := clone.Target(cfg, "staging", "production", false); err != nil {
		t.Fatalf("clone failed: %v", err)
	}

	// Mutate the clone.
	cfg.Targets["production"]["KEY"] = "mutated"

	if cfg.Targets["staging"]["KEY"] != "original" {
		t.Error("source target was mutated by changes to the clone")
	}
}

// TestCloneEmptySource ensures cloning a target with no keys creates an empty
// destination target without error.
func TestCloneEmptySource(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"empty": {},
		},
	}

	if err := clone.Target(cfg, "empty", "also-empty", false); err != nil {
		t.Fatalf("unexpected error cloning empty target: %v", err)
	}

	if dest, ok := cfg.Targets["also-empty"]; !ok {
		t.Error("destination target was not created")
	} else if len(dest) != 0 {
		t.Errorf("expected 0 keys, got %d", len(dest))
	}
}
