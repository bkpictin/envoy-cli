// Package lint implements heuristic analysis of environment variable
// configurations managed by envoy-cli.
//
// Rules currently enforced:
//
//   - Empty values are flagged as warnings.
//   - Keys that are not fully uppercase are flagged as warnings.
//   - Keys that contain spaces are flagged as errors.
//   - Keys whose value is identical across two or more targets are flagged as
//     warnings, suggesting the value could be promoted to a shared default.
package lint
