// Package prune provides utilities for cleaning up environment variable sets.
//
// It identifies and optionally removes:
//   - Orphaned keys: keys that exist in a target but in no other target,
//     suggesting they may be unused or forgotten.
//   - Empty-value keys: keys whose value is an empty string, which are often
//     placeholders that were never filled in.
//
// All operations support a dry-run mode that reports what would be removed
// without mutating the configuration.
package prune
