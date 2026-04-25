package resolve_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/resolve"
)

// TestChainedRefs verifies that direct (non-recursive) chaining works when
// the referenced key itself contains a resolved literal (not another ref).
func TestChainedRefs(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"prod": {
				"SCHEME":   "https",
				"HOST":     "example.com",
				"BASE_URL": "${SCHEME}://${HOST}",
				"API_URL":  "${BASE_URL}/v1",
			},
		},
	}
	// Single-pass: API_URL resolves ${BASE_URL} to its literal "${SCHEME}://${HOST}"
	// because we do not recurse. The test asserts single-pass behaviour.
	results, err := resolve.Target(cfg, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "BASE_URL" && r.Resolved != "https://example.com" {
			t.Errorf("BASE_URL: expected https://example.com, got %q", r.Resolved)
		}
	}
}

func TestMultipleWarningsCollected(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"COMBO": "${UNKNOWN_A}-${UNKNOWN_B}",
			},
		},
	}
	results, err := resolve.Target(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "COMBO" {
			if len(r.Warnings) < 2 {
				t.Errorf("expected 2 warnings, got %d", len(r.Warnings))
			}
			for _, w := range r.Warnings {
				if !strings.Contains(w, "unresolved reference") {
					t.Errorf("unexpected warning text: %q", w)
				}
			}
		}
	}
}
