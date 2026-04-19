package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envoy-cli/envoy-cli/internal/config"
)

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "envoy.json")

	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"dev": {"FOO": "bar"},
		},
	}
	if err := config.Save(cfg, path); err != nil {
		t.Fatal(err)
	}

	loaded, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Targets["dev"]["FOO"] != "bar" {
		t.Fatalf("expected bar, got %s", loaded.Targets["dev"]["FOO"])
	}
}

func TestInit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "envoy.json")

	if err := config.Init(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal("config file not created")
	}

	// calling Init again should not overwrite
	cfg, _ := config.Load(path)
	cfg.Targets["prod"] = map[string]string{"X": "1"}
	_ = config.Save(cfg, path)

	_ = config.Init(path)
	reloaded, _ := config.Load(path)
	if _, ok := reloaded.Targets["prod"]; !ok {
		t.Fatal("Init overwrote existing config")
	}
}
