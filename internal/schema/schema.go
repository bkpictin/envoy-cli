// Package schema provides validation of environment variable sets
// against a declared schema (required keys, allowed keys, types).
package schema

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
)

// Rule describes a single key constraint in a schema.
type Rule struct {
	Required    bool
	Description string
}

// Schema maps key names to their rules.
type Schema map[string]Rule

// Result holds a single schema violation.
type Result struct {
	Target  string
	Key     string
	Message string
	Level   string // "error" | "warning"
}

// Validate checks the given target's env vars against the provided schema.
// Missing required keys are errors; keys not in the schema are warnings.
func Validate(cfg *config.Config, target string, s Schema) ([]Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}

	var results []Result

	// Check required keys.
	for key, rule := range s {
		if !rule.Required {
			continue
		}
		if _, exists := envs[key]; !exists {
			results = append(results, Result{
				Target:  target,
				Key:     key,
				Message: fmt.Sprintf("required key %q is missing", key),
				Level:   "error",
			})
		}
	}

	// Check for undeclared keys.
	for key := range envs {
		if _, declared := s[key]; !declared {
			results = append(results, Result{
				Target:  target,
				Key:     key,
				Message: fmt.Sprintf("key %q is not declared in schema", key),
				Level:   "warning",
			})
		}
	}

	return results, nil
}

// Format renders results as a human-readable string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "schema validation passed with no issues"
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "[%s] %s: %s\n", strings.ToUpper(r.Level), r.Target, r.Message)
	}
	return strings.TrimRight(sb.String(), "\n")
}
