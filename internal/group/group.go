// Package group provides functionality for organizing environment variable
// keys into named logical groups within a target.
package group

import (
	"errors"
	"fmt"
	"sort"

	"github.com/envoy-cli/envoy/internal/config"
)

// ErrGroupNotFound is returned when a referenced group does not exist.
var ErrGroupNotFound = errors.New("group not found")

// ErrGroupExists is returned when creating a group that already exists.
var ErrGroupExists = errors.New("group already exists")

// ErrKeyNotInGroup is returned when a key is not a member of the group.
var ErrKeyNotInGroup = errors.New("key not in group")

// Create adds a new empty group to the given target.
func Create(cfg *config.Config, target, group string) error {
	if _, ok := cfg.Targets[target]; !ok {
		return fmt.Errorf("target %q not found", target)
	}
	if cfg.Groups == nil {
		cfg.Groups = map[string]map[string][]string{}
	}
	if cfg.Groups[target] == nil {
		cfg.Groups[target] = map[string][]string{}
	}
	if _, exists := cfg.Groups[target][group]; exists {
		return fmt.Errorf("%w: %q", ErrGroupExists, group)
	}
	cfg.Groups[target][group] = []string{}
	return nil
}

// Delete removes a group from the given target.
func Delete(cfg *config.Config, target, group string) error {
	if cfg.Groups == nil || cfg.Groups[target] == nil {
		return fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	if _, exists := cfg.Groups[target][group]; !exists {
		return fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	delete(cfg.Groups[target], group)
	return nil
}

// AddKey associates a key with a group in the given target.
func AddKey(cfg *config.Config, target, group, key string) error {
	if cfg.Groups == nil || cfg.Groups[target] == nil {
		return fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	keys, exists := cfg.Groups[target][group]
	if !exists {
		return fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	for _, k := range keys {
		if k == key {
			return nil // already a member
		}
	}
	cfg.Groups[target][group] = append(keys, key)
	return nil
}

// RemoveKey disassociates a key from a group in the given target.
func RemoveKey(cfg *config.Config, target, group, key string) error {
	if cfg.Groups == nil || cfg.Groups[target] == nil {
		return fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	keys, exists := cfg.Groups[target][group]
	if !exists {
		return fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	next := keys[:0]
	found := false
	for _, k := range keys {
		if k == key {
			found = true
			continue
		}
		next = append(next, k)
	}
	if !found {
		return fmt.Errorf("%w: %q", ErrKeyNotInGroup, key)
	}
	cfg.Groups[target][group] = next
	return nil
}

// ListGroups returns all group names for a target in sorted order.
func ListGroups(cfg *config.Config, target string) ([]string, error) {
	if _, ok := cfg.Targets[target]; !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	if cfg.Groups == nil || cfg.Groups[target] == nil {
		return []string{}, nil
	}
	names := make([]string, 0, len(cfg.Groups[target]))
	for name := range cfg.Groups[target] {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// GetKeys returns the keys belonging to a group, sorted.
func GetKeys(cfg *config.Config, target, group string) ([]string, error) {
	if cfg.Groups == nil || cfg.Groups[target] == nil {
		return nil, fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	keys, exists := cfg.Groups[target][group]
	if !exists {
		return nil, fmt.Errorf("%w: %q", ErrGroupNotFound, group)
	}
	out := make([]string, len(keys))
	copy(out, keys)
	sort.Strings(out)
	return out, nil
}
