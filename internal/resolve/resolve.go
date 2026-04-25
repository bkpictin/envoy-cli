// Package resolve provides variable interpolation within env values,
// expanding references like ${OTHER_KEY} using values from the same target.
package resolve

import (
	"fmt"
	"regexp"
	"strings"

	"envoy-cli/internal/config"
)

var refRe = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}`)

// Result holds the resolved value and any warnings produced during expansion.
type Result struct {
	Key      string
	Original string
	Resolved string
	Warnings []string
}

// Target resolves all variable references within a single target's env map.
// References that cannot be resolved are left as-is and a warning is recorded.
func Target(cfg *config.Config, target string) ([]Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}

	var results []Result
	for key, val := range envs {
		resolved, warnings := interpolate(val, envs)
		results = append(results, Result{
			Key:      key,
			Original: val,
			Resolved: resolved,
			Warnings: warnings,
		})
	}
	return results, nil
}

// Value resolves a single value string against the provided env map.
func Value(raw string, envs map[string]string) (string, []string) {
	return interpolate(raw, envs)
}

func interpolate(raw string, envs map[string]string) (string, []string) {
	var warnings []string
	resolved := refRe.ReplaceAllStringFunc(raw, func(match string) string {
		inner := strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}")
		if v, ok := envs[inner]; ok {
			return v
		}
		warnings = append(warnings, fmt.Sprintf("unresolved reference: %s", match))
		return match
	})
	return resolved, warnings
}
