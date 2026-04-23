// Package template provides variable interpolation for envoy-cli.
//
// It supports rendering template strings using environment variables
// from a named deployment target. Both ${VAR} and $VAR syntax are
// recognised. Missing variables are either collected as warnings
// (lenient mode) or cause an immediate error (strict mode).
//
// Example usage:
//
//	res, err := template.Render(cfg, "production", "https://${APP_HOST}/api", false)
//	if err != nil { ... }
//	fmt.Println(res.Output)
package template
