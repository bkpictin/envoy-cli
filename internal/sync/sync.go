// Package sync provides functionality to synchronise environment variable
// sets between two targets, pushing only keys that differ or are missing.
package sync

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/config"
)

// Result holds the outcome of a sync operation.
type Result struct {
	Added   []string
	Updated []string
	Skipped []string
}

// Options controls sync behaviour.
type Options struct {
	// Overwrite existing keys in the destination when true.
	Overwrite bool
	// Keys restricts the sync to the listed keys; empty means all keys.
	Keys []string
}

// Targets synchronises env vars from src into dst.
// Keys present in src but absent in dst are always added.
// Keys present in both are updated only when opts.Overwrite is true.
func Targets(cfg *config.Config, src, dst string, opts Options) (Result, error) {
	srcEnvs, ok := cfg.Targets[src]
	if !ok {
		return Result{}, fmt.Errorf("source target %q not found", src)
	}
	dstEnvs, ok := cfg.Targets[dst]
	if !ok {
		return Result{}, fmt.Errorf("destination target %q not found", dst)
	}

	wantKeys := buildKeySet(opts.Keys, srcEnvs)

	var res Result
	for _, key := range wantKeys {
		srcVal := srcEnvs[key]
		dstVal, exists := dstEnvs[key]
		switch {
		case !exists:
			dstEnvs[key] = srcVal
			res.Added = append(res.Added, key)
		case exists && srcVal != dstVal && opts.Overwrite:
			dstEnvs[key] = srcVal
			res.Updated = append(res.Updated, key)
		default:
			res.Skipped = append(res.Skipped, key)
		}
	}
	cfg.Targets[dst] = dstEnvs
	return res, nil
}

// buildKeySet returns the intersection of requested keys and available keys,
// or all available keys when requested is empty.
func buildKeySet(requested []string, envs map[string]string) []string {
	if len(requested) == 0 {
		keys := make([]string, 0, len(envs))
		for k := range envs {
			keys = append(keys, k)
		}
		return keys
	}
	var out []string
	for _, k := range requested {
		if _, ok := envs[k]; ok {
			out = append(out, k)
		}
	}
	return out
}
