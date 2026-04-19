package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envoy.yaml")

	cfg := &config.Config{
		Version: "1",
		Targets: map[string]config.Target{
			"production": {
				Description: "prod env",
				Env: map[string]string{"APP_ENV": "production", "PORT": "8080"},
			},
		},
	}

	if err := config.Save(path, cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Version != cfg.Version {
		t.Errorf("version mismatch: got %q want %q", loaded.Version, cfg.Version)
	}

	prod, ok := loaded.Targets["production"]
	if !ok {
		t.Fatal("expected production target")
	}
	if prod.Env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV mismatch: got %q", prod.Env["APP_ENV"])
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envoy.yaml")

	if err := config.Init(path); err != nil {
		t.Fatalf("Init() error: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}

	if err := config.Init(path); err == nil {
		t.Error("expected error when calling Init() on existing file")
	}
}
