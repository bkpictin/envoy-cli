// Package merge implements multi-source environment variable merging for envoy-cli.
//
// It supports two conflict resolution strategies:
//
//   - skip: existing keys in the destination are preserved when a conflict occurs.
//   - overwrite: source values replace destination values on conflict.
//
// Example usage:
//
//	res, err := merge.Targets(cfg, "prod", []string{"dev", "staging"}, merge.StrategySkip)
package merge
