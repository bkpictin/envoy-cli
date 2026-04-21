package rollback_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/rollback"
	"envoy-cli/internal/snapshot"
)

// TestMultipleRollbackSteps verifies that rolling back to an older snapshot
// correctly ignores newer snapshots and restores the right state.
func TestMultipleRollbackSteps(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"staging": {"PORT": "8080", "MODE": "debug"},
		},
		Snapshots: map[string][]config.Snapshot{},
	}

	// Save initial state
	if err := snapshot.Create(cfg, "staging", "initial"); err != nil {
		t.Fatal(err)
	}

	// Modify and snapshot again
	cfg.Targets["staging"]["PORT"] = "9090"
	cfg.Targets["staging"]["MODE"] = "release"
	if err := snapshot.Create(cfg, "staging", "after-update"); err != nil {
		t.Fatal(err)
	}

	// Further modify without snapshotting
	cfg.Targets["staging"]["PORT"] = "3000"

	// Roll back to the very first snapshot
	if err := rollback.ToSnapshot(cfg, "staging", "initial"); err != nil {
		t.Fatalf("rollback to initial failed: %v", err)
	}
	if cfg.Targets["staging"]["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", cfg.Targets["staging"]["PORT"])
	}
	if cfg.Targets["staging"]["MODE"] != "debug" {
		t.Errorf("expected MODE=debug, got %s", cfg.Targets["staging"]["MODE"])
	}

	// Snapshots should still be intact after rollback
	names, err := rollback.ListAvailable(cfg, "staging")
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 snapshots to remain, got %d", len(names))
	}
}
