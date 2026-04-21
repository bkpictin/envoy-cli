package diff

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
)

// Result holds the differences between two targets.
type Result struct {
	OnlyInA  map[string]string
	OnlyInB  map[string]string
	Changed  map[string][2]string // key -> [valueA, valueB]
	Unchanged map[string]string
}

// HasDifferences returns true if there are any differences between the two targets.
func (r *Result) HasDifferences() bool {
	return len(r.OnlyInA) > 0 || len(r.OnlyInB) > 0 || len(r.Changed) > 0
}

// Targets compares environment variables of two targets within a config.
func Targets(cfg *config.Config, targetA, targetB string) (*Result, error) {
	envA, ok := cfg.Targets[targetA]
	if !ok {
		return nil, fmt.Errorf("target %q not found", targetA)
	}
	envB, ok := cfg.Targets[targetB]
	if !ok {
		return nil, fmt.Errorf("target %q not found", targetB)
	}

	result := &Result{
		OnlyInA:   make(map[string]string),
		OnlyInB:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for k, v := range envA {
		if vb, exists := envB[k]; !exists {
			result.OnlyInA[k] = v
		} else if v != vb {
			result.Changed[k] = [2]string{v, vb}
		} else {
			result.Unchanged[k] = v
		}
	}
	for k, v := range envB {
		if _, exists := envA[k]; !exists {
			result.OnlyInB[k] = v
		}
	}
	return result, nil
}

// Format returns a human-readable diff string.
func Format(r *Result, targetA, targetB string) string {
	out := ""
	keys := func(m map[string]string) []string {
		ks := make([]string, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Strings(ks)
		return ks
	}
	for _, k := range keys(r.OnlyInA) {
		out += fmt.Sprintf("- [%s] %s=%s\n", targetA, k, r.OnlyInA[k])
	}
	for _, k := range keys(r.OnlyInB) {
		out += fmt.Sprintf("+ [%s] %s=%s\n", targetB, k, r.OnlyInB[k])
	}
	for _, k := range sortedChanged(r.Changed) {
		v := r.Changed[k]
		out += fmt.Sprintf("~ %s: %s → %s\n", k, v[0], v[1])
	}
	if out == "" {
		out = "No differences found.\n"
	}
	return out
}

func sortedChanged(m map[string][2]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}
