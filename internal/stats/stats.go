// Package stats provides summary statistics across targets and keys.
package stats

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// TargetStat holds statistics for a single target.
type TargetStat struct {
	Target    string
	KeyCount  int
	EmptyKeys int
	HasSnaps  int
}

// Summary holds aggregate statistics for the entire config.
type Summary struct {
	TotalTargets  int
	TotalKeys     int
	TotalEmpty    int
	TotalSnapshots int
	Targets       []TargetStat
}

// Collect gathers statistics from the given config.
func Collect(cfg *config.Config) Summary {
	var s Summary
	s.TotalTargets = len(cfg.Targets)

	for name, target := range cfg.Targets {
		stat := TargetStat{
			Target:   name,
			KeyCount: len(target.Envs),
			HasSnaps: len(cfg.Snapshots[name]),
		}
		for _, v := range target.Envs {
			if v == "" {
				stat.EmptyKeys++
			}
		}
		s.TotalKeys += stat.KeyCount
		s.TotalEmpty += stat.EmptyKeys
		s.TotalSnapshots += stat.HasSnaps
		s.Targets = append(s.Targets, stat)
	}

	sort.Slice(s.Targets, func(i, j int) bool {
		return s.Targets[i].Target < s.Targets[j].Target
	})
	return s
}

// Format returns a human-readable summary string.
func Format(s Summary) string {
	out := fmt.Sprintf("Targets: %d  Keys: %d  Empty: %d  Snapshots: %d\n",
		s.TotalTargets, s.TotalKeys, s.TotalEmpty, s.TotalSnapshots)
	out += fmt.Sprintf("%-20s %6s %6s %9s\n", "TARGET", "KEYS", "EMPTY", "SNAPSHOTS")
	out += fmt.Sprintf("%-20s %6s %6s %9s\n", "------", "----", "-----", "---------")
	for _, t := range s.Targets {
		out += fmt.Sprintf("%-20s %6d %6d %9d\n", t.Target, t.KeyCount, t.EmptyKeys, t.HasSnaps)
	}
	return out
}
