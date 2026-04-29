// Package filter provides utilities for filtering environment variable entries
// across targets based on key patterns, value patterns, or tag membership.
package filter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/yourusername/envoy-cli/internal/config"
)

// Result holds a single matched entry from a filter operation.
type Result struct {
	Target string
	Key    string
	Value  string
}

// Options controls how filtering is applied.
type Options struct {
	// KeyPattern is a glob-style or regex pattern matched against key names.
	KeyPattern string
	// ValuePattern is a regex pattern matched against values.
	ValuePattern string
	// Tags restricts results to keys that carry all of the given tags.
	Tags []string
	// Targets restricts the search to specific target names. Empty means all.
	Targets []string
	// CaseSensitive controls whether pattern matching is case-sensitive.
	CaseSensitive bool
}

// ByPattern returns all key/value pairs across the requested targets whose
// key or value matches the supplied Options. At least one of KeyPattern or
// ValuePattern must be non-empty, otherwise an error is returned.
func ByPattern(cfg *config.Config, opts Options) ([]Result, error) {
	if opts.KeyPattern == "" && opts.ValuePattern == "" {
		return nil, fmt.Errorf("filter: at least one of KeyPattern or ValuePattern must be set")
	}

	var keyRe, valRe *regexp.Regexp
	var err error

	if opts.KeyPattern != "" {
		pattern := globToRegex(opts.KeyPattern)
		if !opts.CaseSensitive {
			pattern = "(?i)" + pattern
		}
		keyRe, err = regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid key pattern: %w", err)
		}
	}

	if opts.ValuePattern != "" {
		pattern := opts.ValuePattern
		if !opts.CaseSensitive {
			pattern = "(?i)" + pattern
		}
		valRe, err = regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("filter: invalid value pattern: %w", err)
		}
	}

	targetSet := toSet(opts.Targets)
	tagSet := toSet(opts.Tags)

	var results []Result

	for _, target := range cfg.Targets {
		if len(targetSet) > 0 && !targetSet[target.Name] {
			continue
		}

		for _, kv := range target.Envs {
			// Key pattern check.
			if keyRe != nil && !keyRe.MatchString(kv.Key) {
				continue
			}
			// Value pattern check.
			if valRe != nil && !valRe.MatchString(kv.Value) {
				continue
			}
			// Tag membership check.
			if len(tagSet) > 0 && !hasAllTags(kv.Tags, tagSet) {
				continue
			}
			results = append(results, Result{
				Target: target.Name,
				Key:    kv.Key,
				Value:  kv.Value,
			})
		}
	}

	return results, nil
}

// Format renders a slice of Results as a human-readable table string.
func Format(results []Result) string {
	if len(results) == 0 {
		return "no matches found"
	}
	var sb strings.Builder
	for _, r := range results {
		fmt.Fprintf(&sb, "[%s] %s = %s\n", r.Target, r.Key, r.Value)
	}
	return strings.TrimRight(sb.String(), "\n")
}

// globToRegex converts a simple glob pattern (supporting * and ?) into a
// full regular expression string anchored at both ends.
func globToRegex(glob string) string {
	var sb strings.Builder
	sb.WriteString("^")
	for _, ch := range glob {
		switch ch {
		case '*':
			sb.WriteString(".*")
		case '?':
			sb.WriteString(".")
		case '.', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			sb.WriteRune('\\')
			sb.WriteRune(ch)
		default:
			sb.WriteRune(ch)
		}
	}
	sb.WriteString("$")
	return sb.String()
}

func toSet(items []string) map[string]bool {
	s := make(map[string]bool, len(items))
	for _, v := range items {
		s[v] = true
	}
	return s
}

func hasAllTags(kvTags []string, required map[string]bool) bool {
	present := toSet(kvTags)
	for tag := range required {
		if !present[tag] {
			return false
		}
	}
	return true
}
