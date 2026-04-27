package placeholder

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"API_KEY":    "todo",
				"DB_PASS":    "real-secret",
				"SMTP_PASS":  "CHANGEME",
				"APP_SECRET": "<your-secret-here>",
			},
			"prod": {
				"API_KEY": "sk-live-abc123",
				"DB_PASS": "hunter2",
			},
		},
	}
}

func TestFindPlaceholders(t *testing.T) {
	cfg := newCfg()
	results, err := Find(cfg, []string{"dev"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestFindNoPlaceholders(t *testing.T) {
	cfg := newCfg()
	results, err := Find(cfg, []string{"prod"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestFindMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := Find(cfg, []string{"staging"}, nil)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestFindAllTargets(t *testing.T) {
	cfg := newCfg()
	results, err := Find(cfg, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// dev has 3 placeholders, prod has 0
	if len(results) != 3 {
		t.Fatalf("expected 3 results across all targets, got %d", len(results))
	}
}

func TestFindExtraPattern(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"qa": {"TOKEN": "MYPLACEHOLDER"},
		},
	}
	results, err := Find(cfg, []string{"qa"}, []string{"myplaceholder"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	results, _ := Find(cfg, []string{"dev"}, nil)
	out := Format(results)
	if !strings.Contains(out, "placeholder value") {
		t.Errorf("expected summary header in output, got: %s", out)
	}
}

func TestFormatEmpty(t *testing.T) {
	out := Format(nil)
	if !strings.Contains(out, "No placeholder") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
