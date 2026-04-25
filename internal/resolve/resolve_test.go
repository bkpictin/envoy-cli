package resolve_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/resolve"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{},
	}
}

func TestResolveSimpleRef(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["prod"] = map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	results, err := resolve.Target(cfg, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "API_URL" && r.Resolved != "https://example.com/api" {
			t.Errorf("expected resolved API_URL = https://example.com/api, got %q", r.Resolved)
		}
	}
}

func TestResolveMissingRef(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["staging"] = map[string]string{
		"SERVICE_URL": "${MISSING_KEY}/svc",
	}
	results, err := resolve.Target(cfg, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "SERVICE_URL" {
			if r.Resolved != "${MISSING_KEY}/svc" {
				t.Errorf("expected original value preserved, got %q", r.Resolved)
			}
			if len(r.Warnings) == 0 {
				t.Error("expected a warning for unresolved reference")
			}
		}
	}
}

func TestResolveMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := resolve.Target(cfg, "ghost")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestResolveNoRefs(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{
		"PLAIN": "just-a-value",
	}
	results, err := resolve.Target(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Key == "PLAIN" {
			if r.Resolved != "just-a-value" {
				t.Errorf("unexpected change to plain value: %q", r.Resolved)
			}
			if len(r.Warnings) != 0 {
				t.Errorf("unexpected warnings for plain value")
			}
		}
	}
}

func TestValueHelper(t *testing.T) {
	envs := map[string]string{"HOST": "localhost", "PORT": "5432"}
	resolved, warnings := resolve.Value("${HOST}:${PORT}", envs)
	if resolved != "localhost:5432" {
		t.Errorf("expected localhost:5432, got %q", resolved)
	}
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
}
