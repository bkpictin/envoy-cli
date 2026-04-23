package search_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/search"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"DATABASE_URL": "postgres://prod",
				"API_KEY":      "secret-prod",
				"LOG_LEVEL":    "error",
			},
			"staging": {
				"DATABASE_URL": "postgres://staging",
				"API_KEY":      "secret-staging",
				"DEBUG":        "true",
			},
		},
	}
}

func TestSearchByKey(t *testing.T) {
	cfg := newCfg()
	results, err := search.Keys(cfg, "API", search.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Key != "API_KEY" {
			t.Errorf("unexpected key %q", r.Key)
		}
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	cfg := newCfg()
	results, err := search.Keys(cfg, "api_key", search.Options{CaseSensitive: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchCaseSensitiveNoMatch(t *testing.T) {
	cfg := newCfg()
	results, err := search.Keys(cfg, "api_key", search.Options{CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchByValue(t *testing.T) {
	cfg := newCfg()
	results, err := search.Keys(cfg, "postgres", search.Options{SearchValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchSingleTarget(t *testing.T) {
	cfg := newCfg()
	results, err := search.Keys(cfg, "DATABASE", search.Options{Target: "production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Target != "production" {
		t.Errorf("expected target production, got %q", results[0].Target)
	}
}

func TestSearchMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := search.Keys(cfg, "KEY", search.Options{Target: "nonexistent"})
	if err == nil {
		t.Fatal("expected error for missing target, got nil")
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	cfg := newCfg()
	results, err := search.Keys(cfg, "", search.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results != nil {
		t.Fatalf("expected nil results for empty query, got %v", results)
	}
}
