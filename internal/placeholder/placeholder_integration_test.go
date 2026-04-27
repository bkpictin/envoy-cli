package placeholder_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/placeholder"
)

func buildMultiTargetCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"DB_URL":  "postgres://todo:todo@localhost/db",
				"SECRET":  "real-value",
			},
			"staging": {
				"DB_URL":  "postgres://user:pass@staging/db",
				"SECRET":  "FIXME",
				"API_KEY": "example-key",
			},
			"prod": {
				"DB_URL":  "postgres://user:s3cr3t@prod/db",
				"SECRET":  "s3cr3t",
				"API_KEY": "live-key-xyz",
			},
		},
	}
}

func TestScanAllTargetsFindsCorrectCount(t *testing.T) {
	cfg := buildMultiTargetCfg()
	results, err := placeholder.Find(cfg, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// dev: DB_URL matches "todo" (twice in value but counted once per key), SECRET clean → 1
	// staging: SECRET matches "fixme", API_KEY matches "example" → 2
	// prod: none → 0
	// total: 3
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestProdTargetIsClean(t *testing.T) {
	cfg := buildMultiTargetCfg()
	results, err := placeholder.Find(cfg, []string{"prod"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected prod to be clean, got %d finding(s)", len(results))
	}
}

func TestCustomPatternAcrossTargets(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"a": {"KEY": "INTERNAL_STUB"},
			"b": {"KEY": "real-value"},
			"c": {"KEY": "another_stub"},
		},
	}
	results, err := placeholder.Find(cfg, nil, []string{"stub"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results for custom pattern, got %d", len(results))
	}
}
