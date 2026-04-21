// Package importenv implements importing environment variables from external
// files into an envoy target.
//
// Supported formats:
//   - dotenv  — KEY=VALUE lines, comments and blank lines ignored
//   - shell   — export KEY=VALUE lines
//   - json    — flat JSON object {"KEY": "value"}
//
// Quoted values (double-quotes) are automatically stripped.
package importenv
