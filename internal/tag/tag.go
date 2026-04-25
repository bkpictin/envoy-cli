// Package tag provides functionality for tagging environment variables
// across targets with arbitrary labels for grouping and filtering.
package tag

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Add attaches a tag to a key within a target.
func Add(cfg *config.Config, target, key, tag string) error {
	t, ok := cfg.Targets[target]
	if !ok {
		return fmt.Errorf("target %q not found", target)
	}
	if _, ok := t.Vars[key]; !ok {
		return fmt.Errorf("key %q not found in target %q", key, target)
	}
	if t.Tags == nil {
		t.Tags = make(map[string][]string)
	}
	for _, existing := range t.Tags[key] {
		if existing == tag {
			return nil // already tagged
		}
	}
	t.Tags[key] = append(t.Tags[key], tag)
	sort.Strings(t.Tags[key])
	cfg.Targets[target] = t
	return nil
}

// Remove detaches a tag from a key within a target.
func Remove(cfg *config.Config, target, key, tag string) error {
	t, ok := cfg.Targets[target]
	if !ok {
		return fmt.Errorf("target %q not found", target)
	}
	if t.Tags == nil {
		return fmt.Errorf("tag %q not found on key %q", tag, key)
	}
	tags := t.Tags[key]
	updated := tags[:0]
	for _, existing := range tags {
		if existing != tag {
			updated = append(updated, existing)
		}
	}
	if len(updated) == len(tags) {
		return fmt.Errorf("tag %q not found on key %q", tag, key)
	}
	t.Tags[key] = updated
	cfg.Targets[target] = t
	return nil
}

// ListByTag returns all keys in a target that carry the given tag.
func ListByTag(cfg *config.Config, target, tag string) ([]string, error) {
	t, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	var keys []string
	for key, tags := range t.Tags {
		for _, tg := range tags {
			if tg == tag {
				keys = append(keys, key)
				break
			}
		}
	}
	sort.Strings(keys)
	return keys, nil
}

// ListForKey returns all tags attached to a key within a target.
func ListForKey(cfg *config.Config, target, key string) ([]string, error) {
	t, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	if _, ok := t.Vars[key]; !ok {
		return nil, fmt.Errorf("key %q not found in target %q", key, target)
	}
	tags := make([]string, len(t.Tags[key]))
	copy(tags, t.Tags[key])
	return tags, nil
}
