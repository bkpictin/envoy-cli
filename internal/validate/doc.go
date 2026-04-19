// Package validate provides helpers for validating environment variable keys
// and target names before they are persisted to the envoy configuration.
//
// Rules enforced:
//   - Keys must match the POSIX pattern [A-Z_][A-Z0-9_]* (case-insensitive input is upper-cased before comparison).
//   - Keys must not collide with envoy internal metadata prefixes.
//   - Target names must already exist in the loaded Config.
package validate
