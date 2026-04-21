// Package clone implements deep-copy semantics for envoy-cli targets.
//
// It allows users to duplicate an existing target — either in full or
// filtered by key prefix — into a new named target. This is useful when
// bootstrapping a new deployment environment from an existing one.
//
// Usage:
//
//	clone.Target(cfg, "staging", "production", false)
//	clone.WithFilter(cfg, "staging", "production", false, func(k string) bool {
//	    return strings.HasPrefix(k, "APP_")
//	})
package clone
