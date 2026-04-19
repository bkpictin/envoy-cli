package promote

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/config"
)

// Options controls promotion behaviour.
type Options struct {
	Overwrite bool
	Keys      []string // if empty, promote all keys
}

// Promote copies environment variables from src to dst target,
// optionally filtered to a specific set of keys.
func Promote(cfg *config.Config, src, dst string, opts Options) ([]string, error) {
	srcEnvs, ok := cfg.Targets[src]
	if !ok {
		return nil, fmt.Errorf("source target %q not found", src)
	}
	if _, ok := cfg.Targets[dst]; !ok {
		return nil, fmt.Errorf("destination target %q not found", dst)
	}

	filter := map[string]bool{}
	for _, k := range opts.Keys {
		filter[k] = true
	}

	promoted := []string{}
	for k, v := range srcEnvs {
		if len(filter) > 0 && !filter[k] {
			continue
		}
		_, exists := cfg.Targets[dst][k]
		if exists && !opts.Overwrite {
			continue
		}
		if cfg.Targets[dst] == nil {
			cfg.Targets[dst] = map[string]string{}
		}
		cfg.Targets[dst][k] = v
		promoted = append(promoted, k)
	}
	return promoted, nil
}
