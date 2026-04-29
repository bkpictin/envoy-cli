package extract_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/extract"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"APP_HOST": "localhost",
				"APP_PORT": "8080",
				"DB_URL":   "postgres://dev",
			},
			"prod": {
				"APP_HOST": "prod.example.com",
			},
		},
	}
}

func TestFromTargetAllKeys(t *testing.T) {
	cfg := newCfg()
	r, err := extract.FromTarget(cfg, "dev", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Keys))
	}
}

func TestFromTargetSelectedKeys(t *testing.T) {
	cfg := newCfg()
	r, err := extract.FromTarget(cfg, "dev", []string{"APP_HOST", "APP_PORT"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Keys))
	}
	if r.Keys["APP_HOST"] != "localhost" {
		t.Errorf("unexpected value for APP_HOST: %s", r.Keys["APP_HOST"])
	}
}

func TestFromTargetMissingKeyLenient(t *testing.T) {
	cfg := newCfg()
	r, err := extract.FromTarget(cfg, "dev", []string{"APP_HOST", "MISSING"}, false)
	if err != nil {
		t.Fatalf("unexpected error in lenient mode: %v", err)
	}
	if _, ok := r.Keys["MISSING"]; ok {
		t.Error("MISSING key should not be present")
	}
	if len(r.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(r.Keys))
	}
}

func TestFromTargetMissingKeyStrict(t *testing.T) {
	cfg := newCfg()
	_, err := extract.FromTarget(cfg, "dev", []string{"MISSING"}, true)
	if err == nil {
		t.Fatal("expected error in strict mode, got nil")
	}
}

func TestFromTargetMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := extract.FromTarget(cfg, "staging", nil, false)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestIntoTarget(t *testing.T) {
	cfg := newCfg()
	r := extract.Result{Target: "dev", Keys: map[string]string{"APP_HOST": "localhost"}}
	if err := extract.IntoTarget(cfg, "prod", r, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["prod"]["APP_HOST"] != "prod.example.com" {
		t.Error("existing key should not be overwritten when overwrite=false")
	}
}

func TestIntoTargetOverwrite(t *testing.T) {
	cfg := newCfg()
	r := extract.Result{Target: "dev", Keys: map[string]string{"APP_HOST": "new-host"}}
	if err := extract.IntoTarget(cfg, "prod", r, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["prod"]["APP_HOST"] != "new-host" {
		t.Errorf("expected new-host, got %s", cfg.Targets["prod"]["APP_HOST"])
	}
}

func TestFormat(t *testing.T) {
	r := extract.Result{Target: "dev", Keys: map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}}
	out := extract.Format(r)
	if !strings.Contains(out, "2 key(s)") {
		t.Errorf("expected key count in output, got: %s", out)
	}
	if !strings.Contains(out, "APP_HOST=localhost") {
		t.Errorf("expected APP_HOST in output")
	}
}
