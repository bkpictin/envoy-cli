package merge

import (
	"testing"

	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"APP_PORT": "3000",
				"LOG_LEVEL": "debug",
			},
			"staging": {
				"APP_PORT": "8080",
				"DB_URL": "postgres://staging/db",
			},
			"prod": {
				"APP_PORT": "80",
			},
		},
	}
}

func TestMergeSkip(t *testing.T) {
	cfg := newCfg()
	res, err := Targets(cfg, "prod", []string{"dev", "staging"}, StrategySkip)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["prod"]["APP_PORT"] != "80" {
		t.Errorf("expected APP_PORT to remain 80, got %s", cfg.Targets["prod"]["APP_PORT"])
	}
	if cfg.Targets["prod"]["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %s", cfg.Targets["prod"]["LOG_LEVEL"])
	}
	if res.Merged != 2 {
		t.Errorf("expected 2 merged, got %d", res.Merged)
	}
	if res.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", res.Skipped)
	}
}

func TestMergeOverwrite(t *testing.T) {
	cfg := newCfg()
	res, err := Targets(cfg, "prod", []string{"staging"}, StrategyOverwrite)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Targets["prod"]["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT=8080 after overwrite, got %s", cfg.Targets["prod"]["APP_PORT"])
	}
	if res.Overwrote != 1 {
		t.Errorf("expected 1 overwritten, got %d", res.Overwrote)
	}
}

func TestMergeMissingDest(t *testing.T) {
	cfg := newCfg()
	_, err := Targets(cfg, "nonexistent", []string{"dev"}, StrategySkip)
	if err == nil {
		t.Error("expected error for missing destination target")
	}
}

func TestMergeMissingSource(t *testing.T) {
	cfg := newCfg()
	_, err := Targets(cfg, "prod", []string{"ghost"}, StrategySkip)
	if err == nil {
		t.Error("expected error for missing source target")
	}
}
