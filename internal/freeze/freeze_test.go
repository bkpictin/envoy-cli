package freeze_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/freeze"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev":  {"FOO": "bar"},
			"prod": {"FOO": "baz"},
		},
	}
}

func TestFreezeAndIsFrozen(t *testing.T) {
	cfg := newCfg()
	if err := freeze.Freeze(cfg, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !freeze.IsFrozen(cfg, "prod") {
		t.Fatal("expected prod to be frozen")
	}
	if freeze.IsFrozen(cfg, "dev") {
		t.Fatal("dev should not be frozen")
	}
}

func TestFreezeAlreadyFrozen(t *testing.T) {
	cfg := newCfg()
	_ = freeze.Freeze(cfg, "prod")
	if err := freeze.Freeze(cfg, "prod"); err == nil {
		t.Fatal("expected error for already-frozen target")
	}
}

func TestFreezeMissingTarget(t *testing.T) {
	cfg := newCfg()
	if err := freeze.Freeze(cfg, "staging"); err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestUnfreeze(t *testing.T) {
	cfg := newCfg()
	_ = freeze.Freeze(cfg, "prod")
	if err := freeze.Unfreeze(cfg, "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if freeze.IsFrozen(cfg, "prod") {
		t.Fatal("expected prod to be unfrozen")
	}
}

func TestUnfreezeNotFrozen(t *testing.T) {
	cfg := newCfg()
	if err := freeze.Unfreeze(cfg, "dev"); err == nil {
		t.Fatal("expected error when unfreezing a non-frozen target")
	}
}

func TestList(t *testing.T) {
	cfg := newCfg()
	_ = freeze.Freeze(cfg, "prod")
	list := freeze.List(cfg)
	if len(list) != 1 || list[0] != "prod" {
		t.Fatalf("expected [prod], got %v", list)
	}
}

func TestGuardWrite(t *testing.T) {
	cfg := newCfg()
	_ = freeze.Freeze(cfg, "prod")
	if err := freeze.GuardWrite(cfg, "prod"); err == nil {
		t.Fatal("expected ErrFrozen for frozen target")
	}
	if err := freeze.GuardWrite(cfg, "dev"); err != nil {
		t.Fatalf("expected no error for unfrozen target, got: %v", err)
	}
}
