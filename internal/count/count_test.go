package count_test

import (
	"strings"
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/count"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"DB_HOST": "prod.db.example.com",
				"DB_PASS": "",
				"API_KEY": "secret",
			},
			"staging": {
				"DB_HOST": "staging.db.example.com",
				"DB_PASS": "",
			},
			"development": {
				"DB_HOST": "localhost",
			},
		},
	}
}

func TestByTarget(t *testing.T) {
	cfg := newCfg()
	r, err := count.ByTarget(cfg, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Total != 3 {
		t.Errorf("expected total 3, got %d", r.Total)
	}
	if r.Empty != 1 {
		t.Errorf("expected empty 1, got %d", r.Empty)
	}
	if r.Filled != 2 {
		t.Errorf("expected filled 2, got %d", r.Filled)
	}
}

func TestByTargetMissing(t *testing.T) {
	cfg := newCfg()
	_, err := count.ByTarget(cfg, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestAll(t *testing.T) {
	cfg := newCfg()
	s := count.All(cfg)
	if len(s.Results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(s.Results))
	}
	if s.TotalKeys != 6 {
		t.Errorf("expected 6 total keys, got %d", s.TotalKeys)
	}
	if s.TotalEmpty != 2 {
		t.Errorf("expected 2 empty keys, got %d", s.TotalEmpty)
	}
	// results should be sorted alphabetically
	if s.Results[0].Target != "development" {
		t.Errorf("expected first target 'development', got %q", s.Results[0].Target)
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	s := count.All(cfg)
	out := count.Format(s)
	if !strings.Contains(out, "production") {
		t.Error("expected 'production' in output")
	}
	if !strings.Contains(out, "totals:") {
		t.Error("expected 'totals:' summary line in output")
	}
}

func TestFormatEmpty(t *testing.T) {
	cfg := &config.Config{Targets: map[string]map[string]string{}}
	s := count.All(cfg)
	out := count.Format(s)
	if !strings.Contains(out, "no targets found") {
		t.Errorf("expected 'no targets found', got %q", out)
	}
}
