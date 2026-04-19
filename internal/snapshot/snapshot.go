package snapshot

import (
	"fmt"
	"time"

	"github.com/user/envoy-cli/internal/config"
)

// Snapshot captures the state of a target's env vars at a point in time.
type Snapshot struct {
	Target    string            `yaml:"target"`
	Timestamp time.Time         `yaml:"timestamp"`
	Vars      map[string]string `yaml:"vars"`
}

// Create saves a snapshot of the given target's current env vars.
func Create(cfg *config.Config, target string) (Snapshot, error) {
	vars, ok := cfg.Targets[target]
	if !ok {
		return Snapshot{}, fmt.Errorf("target %q not found", target)
	}

	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}

	snap := Snapshot{
		Target:    target,
		Timestamp: time.Now().UTC(),
		Vars:      copy,
	}

	if cfg.Snapshots == nil {
		cfg.Snapshots = make(map[string][]Snapshot)
	}
	cfg.Snapshots[target] = append(cfg.Snapshots[target], snap)
	return snap, nil
}

// List returns all snapshots for a given target.
func List(cfg *config.Config, target string) ([]Snapshot, error) {
	if _, ok := cfg.Targets[target]; !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	snaps := cfg.Snapshots[target]
	return snaps, nil
}

// Restore replaces a target's env vars with those from snapshot index i.
func Restore(cfg *config.Config, target string, index int) error {
	snaps, err := List(cfg, target)
	if err != nil {
		return err
	}
	if index < 0 || index >= len(snaps) {
		return fmt.Errorf("snapshot index %d out of range (0-%d)", index, len(snaps)-1)
	}
	copy := make(map[string]string, len(snaps[index].Vars))
	for k, v := range snaps[index].Vars {
		copy[k] = v
	}
	cfg.Targets[target] = copy
	return nil
}
