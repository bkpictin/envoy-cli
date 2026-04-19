package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/envoy-cli/envoy/internal/config"
)

var keyRegex = regexp.MustCompile(`^[A-Z_][A-Z0-9_]*$`)

// KeyFormat checks that an env var key is valid (POSIX-style).
func KeyFormat(key string) error {
	if !keyRegex.MatchString(strings.ToUpper(key)) {
		return fmt.Errorf("invalid key %q: must match [A-Z_][A-Z0-9_]*", key)
	}
	return nil
}

// TargetExists returns an error if the named target is not present in cfg.
func TargetExists(cfg *config.Config, target string) error {
	if _, ok := cfg.Targets[target]; !ok {
		return fmt.Errorf("target %q does not exist", target)
	}
	return nil
}

// NoReservedKeys rejects keys that clash with envoy internal metadata.
var reservedKeys = map[string]struct{}{
	"__ENVOY_TARGET__": {},
	"__ENVOY_VERSION__": {},
}

func NoReservedKeys(key string) error {
	if _, ok := reservedKeys[strings.ToUpper(key)]; ok {
		return fmt.Errorf("key %q is reserved for internal use", key)
	}
	return nil
}

// All runs all validations for a key before setting it.
func All(cfg *config.Config, target, key string) error {
	if err := TargetExists(cfg, target); err != nil {
		return err
	}
	if err := KeyFormat(key); err != nil {
		return err
	}
	if err := NoReservedKeys(key); err != nil {
		return err
	}
	return nil
}
