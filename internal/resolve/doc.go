// Package resolve implements variable interpolation for envoy-cli.
//
// It expands ${KEY} references within an environment variable's value using
// other keys defined in the same target. References to unknown keys are
// preserved verbatim and surfaced as warnings so callers can decide how to
// handle missing dependencies.
package resolve
