// Package rollback provides utilities for reverting environment variable
// changes on a target by restoring a previously created snapshot.
//
// Usage:
//
//	// Roll back to the most recent snapshot
//	name, err := rollback.ToPrevious(cfg, "production")
//
//	// Roll back to a specific named snapshot
//	err := rollback.ToSnapshot(cfg, "production", "pre-deploy-v2")
//
//	// List all available rollback points
//	names, err := rollback.ListAvailable(cfg, "production")
//
// Rollback is non-destructive to snapshot history; the snapshot remains
// available for future restores after a rollback is performed.
package rollback
