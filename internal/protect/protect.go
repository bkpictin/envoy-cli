// Package protect provides functionality for marking environment variable
// keys as protected, preventing accidental modification or deletion.
package protect

import (
	"errors"
	"fmt"

	"envoy-cli/internal/config"
)

var (
	ErrAlreadyProtected = errors.New("key is already protected")
	ErrNotProtected     = errors.New("key is not protected")
	ErrTargetNotFound   = errors.New("target not found")
	ErrKeyNotFound      = errors.New("key not found in target")
)

// Protect marks a key in the given target as protected.
func Protect(cfg *config.Config, target, key string) error {
	t, ok := cfg.Targets[target]
	if !ok {
		return fmt.Errorf("%w: %s", ErrTargetNotFound, target)
	}
	if _, exists := t.Vars[key]; !exists {
		return fmt.Errorf("%w: %s", ErrKeyNotFound, key)
	}
	if cfg.Protected == nil {
		cfg.Protected = make(map[string]map[string]bool)
	}
	if cfg.Protected[target] == nil {
		cfg.Protected[target] = make(map[string]bool)
	}
	if cfg.Protected[target][key] {
		return fmt.Errorf("%w: %s", ErrAlreadyProtected, key)
	}
	cfg.Protected[target][key] = true
	return nil
}

// Unprotect removes the protected status from a key in the given target.
func Unprotect(cfg *config.Config, target, key string) error {
	if cfg.Protected == nil || !cfg.Protected[target][key] {
		return fmt.Errorf("%w: %s", ErrNotProtected, key)
	}
	delete(cfg.Protected[target], key)
	if len(cfg.Protected[target]) == 0 {
		delete(cfg.Protected, target)
	}
	return nil
}

// IsProtected reports whether a key in the given target is protected.
func IsProtected(cfg *config.Config, target, key string) bool {
	if cfg.Protected == nil {
		return false
	}
	return cfg.Protected[target][key]
}

// List returns all protected keys for the given target.
func List(cfg *config.Config, target string) ([]string, error) {
	if _, ok := cfg.Targets[target]; !ok {
		return nil, fmt.Errorf("%w: %s", ErrTargetNotFound, target)
	}
	var keys []string
	if cfg.Protected != nil {
		for k, v := range cfg.Protected[target] {
			if v {
				keys = append(keys, k)
			}
		}
	}
	return keys, nil
}
