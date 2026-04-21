// Package compare provides functionality for comparing environment variable
// sets between two targets or between a target and a snapshot.
package compare

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Result holds the comparison outcome between two env sets.
type Result struct {
	OnlyInA    map[string]string // keys present only in source
	OnlyInB    map[string]string // keys present only in dest
	Different  map[string]Pair   // keys present in both but with different values
	Identical  map[string]string // keys with identical values in both
}

// Pair holds the two differing values for a key.
type Pair struct {
	A string
	B string
}

// Targets compares the env vars of two targets and returns a Result.
func Targets(cfg *config.Config, a, b string) (Result, error) {
	envA, ok := cfg.Targets[a]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", a)
	}
	envB, ok := cfg.Targets[b]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", b)
	}
	return compare(envA, envB), nil
}

// SnapshotVsTarget compares a named snapshot against a live target.
func SnapshotVsTarget(cfg *config.Config, snapshotName, target string) (Result, error) {
	snap, ok := cfg.Snapshots[snapshotName]
	if !ok {
		return Result{}, fmt.Errorf("snapshot %q not found", snapshotName)
	}
	envT, ok := cfg.Targets[target]
	if !ok {
		return Result{}, fmt.Errorf("target %q not found", target)
	}
	return compare(snap, envT), nil
}

func compare(a, b map[string]string) Result {
	r := Result{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Different: make(map[string]Pair),
		Identical: make(map[string]string),
	}
	for k, v := range a {
		if bv, exists := b[k]; !exists {
			r.OnlyInA[k] = v
		} else if v != bv {
			r.Different[k] = Pair{A: v, B: bv}
		} else {
			r.Identical[k] = v
		}
	}
	for k, v := range b {
		if _, exists := a[k]; !exists {
			r.OnlyInB[k] = v
		}
	}
	return r
}

// Summary returns a human-readable one-line summary of the result.
func Summary(r Result) string {
	keys := make([]string, 0)
	for k := range r.Different {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return fmt.Sprintf("+%d -%d ~%d =%d",
		len(r.OnlyInA), len(r.OnlyInB), len(r.Different), len(r.Identical))
}
