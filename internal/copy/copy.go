package copy

import (
	"fmt"

	"github.com/envoy-cli/envoy-cli/internal/config"
)

// CopyEnvs copies all environment variables from src target to dst target.
// If overwrite is false, existing keys in dst are preserved.
func CopyEnvs(cfg *config.Config, src, dst string, overwrite bool) error {
	if _, ok := cfg.Targets[src]; !ok {
		return fmt.Errorf("source target %q not found", src)
	}
	if _, ok := cfg.Targets[dst]; !ok {
		return fmt.Errorf("destination target %q not found", dst)
	}

	for k, v := range cfg.Targets[src] {
		if _, exists := cfg.Targets[dst][k]; exists && !overwrite {
			continue
		}
		cfg.Targets[dst][k] = v
	}
	return nil
}

// MergeEnvs merges all environment variables from multiple sources into dst.
// Sources are applied in order; later sources overwrite earlier ones when overwrite is true.
func MergeEnvs(cfg *config.Config, dst string, overwrite bool, sources ...string) error {
	for _, src := range sources {
		if err := CopyEnvs(cfg, src, dst, overwrite); err != nil {
			return err
		}
	}
	return nil
}
