package intersect_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/intersect"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"APP_HOST": "localhost",
				"APP_PORT": "8080",
				"DEBUG":    "true",
			},
			"staging": {
				"APP_HOST": "staging.example.com",
				"APP_PORT": "443",
				"LOG_LEVEL": "info",
			},
			"prod": {
				"APP_HOST": "prod.example.com",
				"APP_PORT": "443",
				"LOG_LEVEL": "warn",
				"SENTRY_DSN": "https://x@sentry.io/1",
			},
		},
	}
}

func TestIntersectTwoTargets(t *testing.T) {
	cfg := newCfg()
	r, err := intersect.Targets(cfg, []string{"dev", "staging"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 common keys, got %d: %v", len(r.Keys), r.Keys)
	}
	for _, k := range r.Keys {
		if k != "APP_HOST" && k != "APP_PORT" {
			t.Errorf("unexpected key in intersection: %s", k)
		}
	}
}

func TestIntersectThreeTargets(t *testing.T) {
	cfg := newCfg()
	r, err := intersect.Targets(cfg, []string{"dev", "staging", "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 common keys, got %d: %v", len(r.Keys), r.Keys)
	}
}

func TestIntersectNoCommonKeys(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["isolated"] = map[string]string{"UNIQUE_KEY": "value"}
	r, err := intersect.Targets(cfg, []string{"dev", "isolated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 0 {
		t.Errorf("expected no common keys, got %v", r.Keys)
	}
}

func TestIntersectMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := intersect.Targets(cfg, []string{"dev", "nonexistent"})
	if err == nil {
		t.Fatal("expected error for missing target, got nil")
	}
}

func TestIntersectTooFewTargets(t *testing.T) {
	cfg := newCfg()
	_, err := intersect.Targets(cfg, []string{"dev"})
	if err == nil {
		t.Fatal("expected error for fewer than two targets, got nil")
	}
}

func TestFormat(t *testing.T) {
	r := intersect.Result{Targets: []string{"dev", "prod"}, Keys: []string{"APP_HOST", "APP_PORT"}}
	out := intersect.Format(r)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}

func TestFormatEmpty(t *testing.T) {
	r := intersect.Result{Targets: []string{"dev", "prod"}, Keys: []string{}}
	out := intersect.Format(r)
	if out == "" {
		t.Error("expected non-empty format output for empty result")
	}
}
