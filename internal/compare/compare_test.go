package compare

import (
	"testing"

	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"staging": {
				"APP_URL":  "https://staging.example.com",
				"DB_HOST":  "db-staging",
				"SHARED":   "same",
			},
			"production": {
				"APP_URL":  "https://prod.example.com",
				"DB_HOST":  "db-prod",
				"SHARED":   "same",
				"PROD_KEY": "secret",
			},
		},
		Snapshots: map[string]map[string]string{
			"snap-1": {
				"APP_URL": "https://staging.example.com",
				"DB_HOST": "db-old",
			},
		},
	}
}

func TestTargetsCompare(t *testing.T) {
	cfg := newCfg()
	r, err := Targets(cfg, "staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.OnlyInA) != 0 {
		t.Errorf("expected 0 only-in-A, got %d", len(r.OnlyInA))
	}
	if len(r.OnlyInB) != 1 {
		t.Errorf("expected 1 only-in-B (PROD_KEY), got %d", len(r.OnlyInB))
	}
	if len(r.Different) != 2 {
		t.Errorf("expected 2 different keys, got %d", len(r.Different))
	}
	if len(r.Identical) != 1 {
		t.Errorf("expected 1 identical key (SHARED), got %d", len(r.Identical))
	}
}

func TestTargetsMissingSource(t *testing.T) {
	cfg := newCfg()
	_, err := Targets(cfg, "nope", "production")
	if err == nil {
		t.Fatal("expected error for missing source target")
	}
}

func TestTargetsMissingDest(t *testing.T) {
	cfg := newCfg()
	_, err := Targets(cfg, "staging", "nope")
	if err == nil {
		t.Fatal("expected error for missing dest target")
	}
}

func TestSnapshotVsTarget(t *testing.T) {
	cfg := newCfg()
	r, err := SnapshotVsTarget(cfg, "snap-1", "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// snap-1 has APP_URL same, DB_HOST different; staging has SHARED extra
	if len(r.Identical) != 1 {
		t.Errorf("expected 1 identical key, got %d", len(r.Identical))
	}
	if len(r.Different) != 1 {
		t.Errorf("expected 1 different key (DB_HOST), got %d", len(r.Different))
	}
	if len(r.OnlyInB) != 1 {
		t.Errorf("expected 1 only-in-target (SHARED), got %d", len(r.OnlyInB))
	}
}

func TestSummary(t *testing.T) {
	r := Result{
		OnlyInA:   map[string]string{"A": "1"},
		OnlyInB:   map[string]string{"B": "2", "C": "3"},
		Different: map[string]Pair{"D": {"x", "y"}},
		Identical: map[string]string{},
	}
	s := Summary(r)
	if s != "+1 -2 ~1 =0" {
		t.Errorf("unexpected summary: %q", s)
	}
}
