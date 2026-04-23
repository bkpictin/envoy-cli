// Package template provides variable interpolation across environment targets.
// It allows rendering a template string using variables from a named target,
// supporting ${VAR} and $VAR syntax.
package template

import (
	"fmt"
	"regexp"
	"strings"

	"envoy-cli/internal/config"
)

var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// RenderResult holds the output of a template render operation.
type RenderResult struct {
	Output   string
	Missing  []string
}

// Render interpolates the template string using env vars from the given target.
// Missing variables are collected and returned rather than causing an error,
// unless strict is true, in which case an error is returned on first missing var.
func Render(cfg *config.Config, target, tmpl string, strict bool) (RenderResult, error) {
	vars, ok := cfg.Targets[target]
	if !ok {
		return RenderResult{}, fmt.Errorf("target %q not found", target)
	}

	missingSet := map[string]struct{}{}
	var missing []string

	output := varPattern.ReplaceAllStringFunc(tmpl, func(match string) string {
		subs := varPattern.FindStringSubmatch(match)
		key := subs[1]
		if key == "" {
			key = subs[2]
		}
		val, exists := vars[key]
		if !exists {
			if _, seen := missingSet[key]; !seen {
				missingSet[key] = struct{}{}
				missing = append(missing, key)
			}
			return match
		}
		return val
	})

	if strict && len(missing) > 0 {
		return RenderResult{}, fmt.Errorf("missing variables in target %q: %s", target, strings.Join(missing, ", "))
	}

	return RenderResult{Output: output, Missing: missing}, nil
}

// ListVars returns all variable references found in the template string.
func ListVars(tmpl string) []string {
	seen := map[string]struct{}{}
	var vars []string
	for _, subs := range varPattern.FindAllStringSubmatch(tmpl, -1) {
		key := subs[1]
		if key == "" {
			key = subs[2]
		}
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			vars = append(vars, key)
		}
	}
	return vars
}
