package stats_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/stats"
)

func newCfg() *config.Config {
	cfg := &config.Config{
		Targets: map[string]*config.Target{
			"dev": {
				Envs: map[string]string{
					"KEY1": "val1",
					"KEY2": "",
					"KEY3": "val3",
				},
			},
			"prod": {
				Envs: map[string]string{
					"KEY1": "prodval",
				},
			},
		},
		Snapshots: map[string][]config.Snapshot{
			"dev": {{Name: "snap1"}, {Name: "snap2"}},
		},
	}
	return cfg
}

func TestCollect(t *testing.T) {
	cfg := newCfg()
	s := stats.Collect(cfg)

	if s.TotalTargets != 2 {
		t.Errorf("expected 2 targets, got %d", s.TotalTargets)
	}
	if s.TotalKeys != 4 {
		t.Errorf("expected 4 total keys, got %d", s.TotalKeys)
	}
	if s.TotalEmpty != 1 {
		t.Errorf("expected 1 empty key, got %d", s.TotalEmpty)
	}
	if s.TotalSnapshots != 2 {
		t.Errorf("expected 2 snapshots, got %d", s.TotalSnapshots)
	}
}

func TestCollectSortedTargets(t *testing.T) {
	cfg := newCfg()
	s := stats.Collect(cfg)

	if len(s.Targets) < 2 {
		t.Fatal("expected at least 2 target stats")
	}
	if s.Targets[0].Target != "dev" {
		t.Errorf("expected first target 'dev', got %s", s.Targets[0].Target)
	}
	if s.Targets[1].Target != "prod" {
		t.Errorf("expected second target 'prod', got %s", s.Targets[1].Target)
	}
}

func TestFormat(t *testing.T) {
	cfg := newCfg()
	s := stats.Collect(cfg)
	out := stats.Format(s)

	if !strings.Contains(out, "dev") {
		t.Error("expected 'dev' in formatted output")
	}
	if !strings.Contains(out, "prod") {
		t.Error("expected 'prod' in formatted output")
	}
	if !strings.Contains(out, "Targets: 2") {
		t.Error("expected summary line in output")
	}
}
