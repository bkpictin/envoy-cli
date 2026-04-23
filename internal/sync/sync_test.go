package sync

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"staging": {
				"APP_PORT": "8080",
				"LOG_LEVEL": "debug",
				"FEATURE_X": "true",
			},
			"production": {
				"APP_PORT": "80",
				"LOG_LEVEL": "info",
			},
		},
	}
}

func TestSyncAddsNewKeys(t *testing.T) {
	cfg := newCfg()
	res, err := Targets(cfg, "staging", "production", Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "FEATURE_X" {
		t.Errorf("expected FEATURE_X added, got %v", res.Added)
	}
	if cfg.Targets["production"]["FEATURE_X"] != "true" {
		t.Error("FEATURE_X should be set in production")
	}
}

func TestSyncSkipsExistingWithoutOverwrite(t *testing.T) {
	cfg := newCfg()
	res, err := Targets(cfg, "staging", "production", Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// APP_PORT and LOG_LEVEL exist in both — should be skipped
	if len(res.Updated) != 0 {
		t.Errorf("expected no updates, got %v", res.Updated)
	}
	if cfg.Targets["production"]["APP_PORT"] != "80" {
		t.Error("APP_PORT should remain 80 in production")
	}
}

func TestSyncOverwritesExisting(t *testing.T) {
	cfg := newCfg()
	res, err := Targets(cfg, "staging", "production", Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) == 0 {
		t.Error("expected at least one update")
	}
	if cfg.Targets["production"]["APP_PORT"] != "8080" {
		t.Error("APP_PORT should be overwritten to 8080")
	}
}

func TestSyncFilteredKeys(t *testing.T) {
	cfg := newCfg()
	res, err := Targets(cfg, "staging", "production", Options{
		Overwrite: true,
		Keys:      []string{"LOG_LEVEL"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Updated) != 1 || res.Updated[0] != "LOG_LEVEL" {
		t.Errorf("expected only LOG_LEVEL updated, got %v", res.Updated)
	}
	// APP_PORT should not have changed
	if cfg.Targets["production"]["APP_PORT"] != "80" {
		t.Error("APP_PORT should be unchanged")
	}
}

func TestSyncMissingSource(t *testing.T) {
	cfg := newCfg()
	_, err := Targets(cfg, "nope", "production", Options{})
	if err == nil {
		t.Error("expected error for missing source")
	}
}

func TestSyncMissingDest(t *testing.T) {
	cfg := newCfg()
	_, err := Targets(cfg, "staging", "nope", Options{})
	if err == nil {
		t.Error("expected error for missing destination")
	}
}
