// Package transform provides value transformation operations for environment
// variable entries within a target. Supported transformations include:
//
//   - uppercase  – convert value to upper-case
//   - lowercase  – convert value to lower-case
//   - base64encode – encode value as standard base64
//   - base64decode – decode a base64-encoded value
//   - trimspace  – strip leading and trailing whitespace
//
// All operations support a dry-run mode that reports what would change without
// modifying the underlying config.
package transform
