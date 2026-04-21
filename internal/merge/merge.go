// Package merge provides functionality for merging environment variables
// from multiple source targets into a single destination target.
package merge

import (
	"fmt"

	"envoy-cli/internal/config"
)

// Strategy defines how conflicts are resolved during a merge.
type Strategy string

const (
	// StrategySkip keeps the existing value in the destination on conflict.
	StrategySkip Strategy = "skip"
	// StrategyOverwrite replaces the destination value with the source value on conflict.
	StrategyOverwrite Strategy = "overwrite"
)

// Result holds the outcome of a merge operation.
type Result struct {
	Merged    int
	Skipped   int
	Overwrote int
}

// Targets merges environment variables from one or more source targets into
// the destination target using the provided conflict resolution strategy.
func Targets(cfg *config.Config, dest string, sources []string, strategy Strategy) (Result, error) {
	if _, ok := cfg.Targets[dest]; !ok {
		return Result{}, fmt.Errorf("destination target %q not found", dest)
	}

	var result Result

	for _, src := range sources {
		envs, ok := cfg.Targets[src]
		if !ok {
			return Result{}, fmt.Errorf("source target %q not found", src)
		}

		for key, val := range envs {
			if _, exists := cfg.Targets[dest][key]; exists {
				if strategy == StrategySkip {
					result.Skipped++
					continue
				}
				cfg.Targets[dest][key] = val
				result.Overwrote++
			} else {
				if cfg.Targets[dest] == nil {
					cfg.Targets[dest] = make(map[string]string)
				}
				cfg.Targets[dest][key] = val
				result.Merged++
			}
		}
	}

	return result, nil
}
