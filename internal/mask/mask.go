// Package mask provides utilities for masking sensitive environment variable
// values in output, preventing accidental exposure of secrets.
package mask

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
	"envoy-cli/internal/encrypt"
)

const (
	defaultMaskChar   = "*"
	defaultVisibleLen = 4
)

// Result holds a key and its masked value.
type Result struct {
	Target string
	Key    string
	Masked string
}

// MaskValue replaces most characters of val with asterisks, leaving the last
// visibleLen characters visible. Encrypted values are fully masked.
func MaskValue(val string, visibleLen int) string {
	if encrypt.IsEncrypted(val) {
		return strings.Repeat(defaultMaskChar, 8) + "[encrypted]"
	}
	if len(val) <= visibleLen {
		return strings.Repeat(defaultMaskChar, len(val))
	}
	return strings.Repeat(defaultMaskChar, len(val)-visibleLen) + val[len(val)-visibleLen:]
}

// Target returns masked results for all keys in the given target.
func Target(cfg *config.Config, target string, visibleLen int) ([]Result, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return nil, fmt.Errorf("target %q not found", target)
	}
	results := make([]Result, 0, len(envs))
	for k, v := range envs {
		results = append(results, Result{
			Target: target,
			Key:    k,
			Masked: MaskValue(v, visibleLen),
		})
	}
	return results, nil
}

// All returns masked results for every target in the config.
func All(cfg *config.Config, visibleLen int) []Result {
	var results []Result
	for target := range cfg.Targets {
		res, _ := Target(cfg, target, visibleLen)
		results = append(results, res...)
	}
	return results
}
