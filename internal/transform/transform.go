// Package transform applies value transformations to environment variables
// within a target, such as upper-casing, lower-casing, or base64 encoding.
package transform

import (
	"encoding/base64"
	"fmt"
	"strings"

	"envoy-cli/internal/config"
)

// Kind represents a transformation type.
type Kind string

const (
	Uppercase Kind = "uppercase"
	Lowercase Kind = "lowercase"
	Base64Encode Kind = "base64encode"
	Base64Decode Kind = "base64decode"
	TrimSpace    Kind = "trimspace"
)

// Result holds the outcome of a single key transformation.
type Result struct {
	Key    string
	Before string
	After  string
	Changed bool
}

// Target applies the given transformation to all (or selected) keys in a target.
// If keys is empty, all keys in the target are transformed.
func Target(cfg *config.Config, target string, kind Kind, keys []string, dryRun bool) ([]Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}

	selected := keySet(keys)
	var results []Result

	for k, v := range envs {
		if len(selected) > 0 && !selected[k] {
			continue
		}
		after, err := apply(kind, v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		results = append(results, Result{Key: k, Before: v, After: after, Changed: v != after})
		if !dryRun && v != after {
			cenvs := cfg.Targets[target]
			cenvs[k] = after
			cfg.Targets[target] = cenvs
		}
	}
	return results, nil
}

func apply(kind Kind, value string) (string, error) {
	switch kind {
	case Uppercase:
		return strings.ToUpper(value), nil
	case Lowercase:
		return strings.ToLower(value), nil
	case Base64Encode:
		return base64.StdEncoding.EncodeToString([]byte(value)), nil
	case Base64Decode:
		b, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return "", fmt.Errorf("invalid base64: %w", err)
		}
		return string(b), nil
	case TrimSpace:
		return strings.TrimSpace(value), nil
	default:
		return "", fmt.Errorf("unknown transform kind %q", kind)
	}
}

func keySet(keys []string) map[string]bool {
	m := make(map[string]bool, len(keys))
	for _, k := range keys {
		m[k] = true
	}
	return m
}
