package flatten_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/flatten"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{},
	}
}

func TestFlattenNoConflicts(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"APP": "myapp", "PORT": "3000"}
	cfg.Targets["prod"] = map[string]string{"REGION": "us-east-1"}

	r := flatten.All(cfg, false)
	if len(r.Envs) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Envs))
	}
	if len(r.Conflicts) != 0 {
		t.Fatalf("expected no conflicts, got %d", len(r.Conflicts))
	}
}

func TestFlattenConflictDetected(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"PORT": "3000"}
	cfg.Targets["prod"] = map[string]string{"PORT": "8080"}

	r := flatten.All(cfg, false)
	if len(r.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(r.Conflicts))
	}
	if r.Conflicts[0].Key != "PORT" {
		t.Errorf("expected conflict on PORT, got %s", r.Conflicts[0].Key)
	}
}

func TestFlattenOverwriteKeepsLast(t *testing.T) {
	cfg := newCfg()
	// sorted order: dev < prod, so prod wins
	cfg.Targets["dev"] = map[string]string{"PORT": "3000"}
	cfg.Targets["prod"] = map[string]string{"PORT": "8080"}

	r := flatten.All(cfg, true)
	if r.Envs["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", r.Envs["PORT"])
	}
}

func TestFlattenNoOverwriteKeepsFirst(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"PORT": "3000"}
	cfg.Targets["prod"] = map[string]string{"PORT": "8080"}

	r := flatten.All(cfg, false)
	if r.Envs["PORT"] != "3000" {
		t.Errorf("expected PORT=3000 (first/dev), got %s", r.Envs["PORT"])
	}
}

func TestFlattenEmptyConfig(t *testing.T) {
	cfg := newCfg()
	r := flatten.All(cfg, false)
	if len(r.Envs) != 0 {
		t.Errorf("expected empty envs")
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"PORT": "3000"}
	cfg.Targets["prod"] = map[string]string{"PORT": "8080"}

	r := flatten.All(cfg, false)
	out := flatten.Format(r)
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in format output, got: %s", out)
	}
	if !strings.Contains(out, "conflict") {
		t.Errorf("expected 'conflict' in format output, got: %s", out)
	}
}

func TestFormatNoConflicts(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"] = map[string]string{"APP": "myapp"}

	r := flatten.All(cfg, false)
	out := flatten.Format(r)
	if !strings.Contains(out, "no conflicts") {
		t.Errorf("expected 'no conflicts' in output, got: %s", out)
	}
}
