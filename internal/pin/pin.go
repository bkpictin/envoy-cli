// Package pin provides functionality for pinning environment variable keys
// to specific values across one or more targets, preventing accidental overwrite.
package pin

import (
	"fmt"

	"envoy-cli/internal/config"
)

const pinnedPrefix = "__pinned__"

// Pin marks a key as pinned in the given target.
// Pinned keys cannot be overwritten by sync, merge, or copy operations
// unless --force is explicitly supplied.
func Pin(cfg *config.Config, target, key string) error {
	t, ok := cfg.Targets[target]
	if !ok {
		return fmt.Errorf("target %q not found", target)
	}
	if _, exists := t.Vars[key]; !exists {
		return fmt.Errorf("key %q does not exist in target %q", key, target)
	}
	metaKey := pinnedPrefix + key
	t.Vars[metaKey] = "1"
	cfg.Targets[target] = t
	return nil
}

// Unpin removes the pinned marker for a key in the given target.
func Unpin(cfg *config.Config, target, key string) error {
	t, ok := cfg.Targets[target]
	if !ok {
		return fmt.Errorf("target %q not found", target)
	}
	metaKey := pinnedPrefix + key
	if _, exists := t.Vars[metaKey]; !exists {
		return fmt.Errorf("key %q is not pinned in target %q", key, target)
	}
	delete(t.Vars, metaKey)
	cfg.Targets[target] = t
	return nil
}

// IsPinned reports whether a key is pinned in the given target.
func IsPinned(cfg *config.Config, target, key string) bool {
	t, ok := cfg.Targets[target]
	if !ok {
		return false
	}
	_, exists := t.Vars[pinnedPrefix+key]
	return exists
}

// List returns all pinned keys for a target.
func List(cfg *config.Config, target string) ([]string, error) {
	t, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	var keys []string
	for k := range t.Vars {
		if len(k) > len(pinnedPrefix) && k[:len(pinnedPrefix)] == pinnedPrefix {
			keys = append(keys, k[len(pinnedPrefix):])
		}
	}
	return keys, nil
}
