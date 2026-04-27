// Package placeholder detects stub or placeholder environment variable values
// across one or more deployment targets.
//
// It scans each key's value against a built-in list of common placeholder
// patterns (e.g. "TODO", "CHANGEME", "<value>") and any user-supplied
// patterns, reporting every match with its target, key, value, and the
// pattern that triggered the match.
//
// Typical usage:
//
//	results, err := placeholder.Find(cfg, []string{"production"}, nil)
//	fmt.Print(placeholder.Format(results))
package placeholder
