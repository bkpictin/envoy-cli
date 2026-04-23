package template_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/template"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"APP_HOST": "prod.example.com",
				"APP_PORT": "443",
				"DB_URL":   "postgres://prod-db/app",
			},
			"staging": {
				"APP_HOST": "staging.example.com",
				"APP_PORT": "8080",
			},
		},
	}
}

func TestRenderBraceStyle(t *testing.T) {
	cfg := newCfg()
	res, err := template.Render(cfg, "production", "https://${APP_HOST}:${APP_PORT}/api", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://prod.example.com:443/api"
	if res.Output != want {
		t.Errorf("got %q, want %q", res.Output, want)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing vars, got %v", res.Missing)
	}
}

func TestRenderBareStyle(t *testing.T) {
	cfg := newCfg()
	res, err := template.Render(cfg, "staging", "host=$APP_HOST port=$APP_PORT", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "host=staging.example.com port=8080"
	if res.Output != want {
		t.Errorf("got %q, want %q", res.Output, want)
	}
}

func TestRenderMissingVarLenient(t *testing.T) {
	cfg := newCfg()
	res, err := template.Render(cfg, "staging", "db=${DB_URL}", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output != "db=${DB_URL}" {
		t.Errorf("expected original placeholder, got %q", res.Output)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "DB_URL" {
		t.Errorf("expected missing [DB_URL], got %v", res.Missing)
	}
}

func TestRenderMissingVarStrict(t *testing.T) {
	cfg := newCfg()
	_, err := template.Render(cfg, "staging", "db=${DB_URL}", true)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestRenderMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := template.Render(cfg, "ghost", "${APP_HOST}", false)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestListVars(t *testing.T) {
	vars := template.ListVars("${APP_HOST}:${APP_PORT} $APP_HOST")
	if len(vars) != 2 {
		t.Errorf("expected 2 unique vars, got %d: %v", len(vars), vars)
	}
}
