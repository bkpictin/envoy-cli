// Package rollback provides functionality to revert environment variable
// changes by restoring a previous snapshot or undoing the last modification.
package rollback

import (
	"errors"
	"fmt"

	"envoy-cli/internal/config"
	"envoy-cli/internal/snapshot"
)

// ErrNoSnapshots is returned when no snapshots exist for a target.
var ErrNoSnapshots = errors.New("no snapshots available for target")

// ToSnapshot restores a target to a specific named snapshot.
func ToSnapshot(cfg *config.Config, target, name string) error {
	snaps, ok := cfg.Snapshots[target]
	if !ok || len(snaps) == 0 {
		return ErrNoSnapshots
	}
	for _, s := range snaps {
		if s.Name == name {
			return snapshot.Restore(cfg, target, name)
		}
	}
	return fmt.Errorf("snapshot %q not found for target %q", name, target)
}

// ToPrevious restores a target to the most recently created snapshot.
// It returns ErrNoSnapshots if no snapshots exist.
func ToPrevious(cfg *config.Config, target string) (string, error) {
	snaps, ok := cfg.Snapshots[target]
	if !ok || len(snaps) == 0 {
		return "", ErrNoSnapshots
	}
	// snapshots are appended in order; last is most recent
	last := snaps[len(snaps)-1]
	if err := snapshot.Restore(cfg, target, last.Name); err != nil {
		return "", err
	}
	return last.Name, nil
}

// ListAvailable returns the names of all snapshots available for rollback
// on the given target, ordered from oldest to newest.
func ListAvailable(cfg *config.Config, target string) ([]string, error) {
	snaps, ok := cfg.Snapshots[target]
	if !ok || len(snaps) == 0 {
		return nil, ErrNoSnapshots
	}
	names := make([]string, len(snaps))
	for i, s := range snaps {
		names[i] = s.Name
	}
	return names, nil
}
