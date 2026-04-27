// Package flatten provides utilities for collapsing all targets into a
// single merged key-value map, with configurable conflict resolution.
package flatten

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Result holds the flattened output and any conflicts detected.
type Result struct {
	Envs      map[string]string
	Conflicts []Conflict
}

// Conflict describes a key that existed in more than one target with
// differing values.
type Conflict struct {
	Key     string
	Targets []string // targets that held this key
	Values  []string // corresponding values (parallel slice)
}

// All merges every target in cfg into a single map. When two targets share a
// key with different values a Conflict entry is recorded. The last target in
// sorted order wins unless overwrite is false, in which case the first value
// is kept.
func All(cfg *config.Config, overwrite bool) Result {
	type entry struct {
		value  string
		target string
	}

	seen := map[string]entry{}
	conflictMap := map[string]*Conflict{}

	targets := make([]string, 0, len(cfg.Targets))
	for t := range cfg.Targets {
		targets = append(targets, t)
	}
	sort.Strings(targets)

	for _, t := range targets {
		for k, v := range cfg.Targets[t] {
			if prev, exists := seen[k]; exists && prev.value != v {
				if c, ok := conflictMap[k]; ok {
					c.Targets = append(c.Targets, t)
					c.Values = append(c.Values, v)
				} else {
					conflictMap[k] = &Conflict{
						Key:     k,
						Targets: []string{prev.target, t},
						Values:  []string{prev.value, v},
					}
				}
				if overwrite {
					seen[k] = entry{value: v, target: t}
				}
			} else if !exists {
				seen[k] = entry{value: v, target: t}
			}
		}
	}

	envs := make(map[string]string, len(seen))
	for k, e := range seen {
		envs[k] = e.value
	}

	conflicts := make([]Conflict, 0, len(conflictMap))
	for _, c := range conflictMap {
		conflicts = append(conflicts, *c)
	}
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	return Result{Envs: envs, Conflicts: conflicts}
}

// Format returns a human-readable summary of a Result.
func Format(r Result) string {
	if len(r.Conflicts) == 0 {
		return fmt.Sprintf("%d keys, no conflicts", len(r.Envs))
	}
	out := fmt.Sprintf("%d keys, %d conflict(s):\n", len(r.Envs), len(r.Conflicts))
	for _, c := range r.Conflicts {
		out += fmt.Sprintf("  %s: ", c.Key)
		for i, t := range c.Targets {
			out += fmt.Sprintf("%s=%q", t, c.Values[i])
			if i < len(c.Targets)-1 {
				out += ", "
			}
		}
		out += "\n"
	}
	return out
}
