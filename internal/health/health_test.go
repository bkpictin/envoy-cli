package health_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/health"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets:   make(map[string]map[string]string),
		Snapshots: make(map[string]config.Snapshot),
	}
}

func TestHealthyConfig(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["prod"] = map[string]string{"DB_URL": "postgres://localhost/prod"}
	cfg.Targets["staging"] = map[string]string{"DB_URL": "postgres://localhost/staging"}

	r := health.Check(cfg)
	if !r.OK() {
		t.Fatalf("expected healthy config, got issues: %v", r.Issues)
	}
}

func TestNoTargets(t *testing.T) {
	cfg := newCfg()
	r := health.Check(cfg)
	if r.OK() {
		t.Fatal("expected warn for no targets")
	}
	if len(r.Issues) != 1 || !strings.Contains(r.Issues[0].Message, "no targets") {
		t.Fatalf("unexpected issues: %v", r.Issues)
	}
}

func TestEmptyTarget(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{}
	r := health.Check(cfg)
	if r.OK() {
		t.Fatal("expected warn for empty target")
	}
}

func TestEmptyValue(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"API_KEY": ""}
	r := health.Check(cfg)
	if r.OK() {
		t.Fatal("expected warn for empty value")
	}
	if !strings.Contains(r.Issues[0].Message, "API_KEY") {
		t.Fatalf("expected key name in message, got: %s", r.Issues[0].Message)
	}
}

func TestOrphanedSnapshot(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["prod"] = map[string]string{"X": "1"}
	cfg.Snapshots["snap1"] = config.Snapshot{Target: "ghost", Envs: map[string]string{}}
	r := health.Check(cfg)
	if r.OK() {
		t.Fatal("expected error for orphaned snapshot")
	}
	errs := health.FilterByLevel(r, health.Error)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	r := health.Check(cfg)
	out := health.Format(r)
	if !strings.Contains(out, "[WARN]") {
		t.Fatalf("expected WARN in output, got: %s", out)
	}

	cfg2 := newCfg()
	cfg2.Targets["prod"] = map[string]string{"K": "v"}
	r2 := health.Check(cfg2)
	out2 := health.Format(r2)
	if !strings.Contains(out2, "healthy") {
		t.Fatalf("expected healthy message, got: %s", out2)
	}
}
