// Package mask provides value-masking utilities for envoy-cli.
//
// It is used to safely display environment variable values in terminal output
// without revealing full secret contents. Encrypted values (prefixed with
// "enc:") are always fully masked and annotated. Plain values show only the
// last N characters (default 4) with the remainder replaced by asterisks.
//
// Example output:
//
//	API_KEY  ************alue
//	DB_PASS  ********[encrypted]
package mask
