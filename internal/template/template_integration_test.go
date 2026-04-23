package template_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/template"
)

// TestRenderMultipleTargets verifies that the same template produces
// different output for different targets.
func TestRenderMultipleTargets(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"dev":  {"API_URL": "http://localhost:3000"},
			"prod": {"API_URL": "https://api.example.com"},
		},
	}
	tmpl := "endpoint=${API_URL}/health"

	for _, tc := range []struct{ target, want string }{
		{"dev", "endpoint=http://localhost:3000/health"},
		{"prod", "endpoint=https://api.example.com/health"},
	} {
		res, err := template.Render(cfg, tc.target, tmpl, true)
		if err != nil {
			t.Fatalf("target %s: unexpected error: %v", tc.target, err)
		}
		if res.Output != tc.want {
			t.Errorf("target %s: got %q, want %q", tc.target, res.Output, tc.want)
		}
	}
}

// TestRenderDeduplicatesMissing ensures that a variable referenced
// multiple times only appears once in the Missing slice.
func TestRenderDeduplicatesMissing(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"staging": {},
		},
	}
	res, err := template.Render(cfg, "staging", "${FOO} and ${FOO} again", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 {
		t.Errorf("expected 1 unique missing var, got %d: %v", len(res.Missing), res.Missing)
	}
	if !strings.Contains(res.Output, "${FOO}") {
		t.Errorf("expected placeholder preserved, got %q", res.Output)
	}
}

// TestListVarsEmpty confirms an empty slice is returned for a plain string.
func TestListVarsEmpty(t *testing.T) {
	vars := template.ListVars("no variables here")
	if len(vars) != 0 {
		t.Errorf("expected empty, got %v", vars)
	}
}
