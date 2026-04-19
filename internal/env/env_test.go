package env_test

import (
	"os"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/config"
	"github.com/envoy-cli/envoy-cli/internal/env"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: make(map[string]map[string]string),
	}
}

func TestSetAndGet(t *testing.T) {
	cfg := newCfg()
	if err := env.Set(cfg, "production", "DB_HOST", "localhost"); err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	val, err := env.Get(cfg, "production", "DB_HOST")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "localhost" {
		t.Errorf("expected localhost, got %q", val)
	}
}

func TestGetMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := env.Get(cfg, "staging", "KEY")
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestDelete(t *testing.T) {
	cfg := newCfg()
	_ = env.Set(cfg, "dev", "API_KEY", "secret")
	if err := env.Delete(cfg, "dev", "API_KEY"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err := env.Get(cfg, "dev", "API_KEY")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestExport(t *testing.T) {
	cfg := newCfg()
	_ = env.Set(cfg, "local", "ENVOY_TEST_VAR", "hello")
	if err := env.Export(cfg, "local"); err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if got := os.Getenv("ENVOY_TEST_VAR"); got != "hello" {
		t.Errorf("expected hello, got %q", got)
	}
	os.Unsetenv("ENVOY_TEST_VAR")
}

func TestSetEmptyKey(t *testing.T) {
	cfg := newCfg()
	if err := env.Set(cfg, "prod", "", "value"); err == nil {
		t.Error("expected error for empty key")
	}
}
