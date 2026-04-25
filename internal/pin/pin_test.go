package pin_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/pin"
)

func newCfg() *config.Config {
	cfg := &config.Config{
		Targets: map[string]config.Target{
			"dev": {Vars: map[string]string{"API_KEY": "abc", "DEBUG": "true"}},
			"prod": {Vars: map[string]string{"API_KEY": "xyz"}},
		},
	}
	return cfg
}

func TestPinAndIsPinned(t *testing.T) {
	cfg := newCfg()
	if err := pin.Pin(cfg, "dev", "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !pin.IsPinned(cfg, "dev", "API_KEY") {
		t.Error("expected API_KEY to be pinned")
	}
	if pin.IsPinned(cfg, "dev", "DEBUG") {
		t.Error("DEBUG should not be pinned")
	}
}

func TestPinMissingTarget(t *testing.T) {
	cfg := newCfg()
	if err := pin.Pin(cfg, "staging", "API_KEY"); err == nil {
		t.Error("expected error for missing target")
	}
}

func TestPinMissingKey(t *testing.T) {
	cfg := newCfg()
	if err := pin.Pin(cfg, "dev", "NONEXISTENT"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestUnpin(t *testing.T) {
	cfg := newCfg()
	_ = pin.Pin(cfg, "dev", "API_KEY")
	if err := pin.Unpin(cfg, "dev", "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pin.IsPinned(cfg, "dev", "API_KEY") {
		t.Error("expected API_KEY to be unpinned")
	}
}

func TestUnpinNotPinned(t *testing.T) {
	cfg := newCfg()
	if err := pin.Unpin(cfg, "dev", "DEBUG"); err == nil {
		t.Error("expected error when unpinning a non-pinned key")
	}
}

func TestList(t *testing.T) {
	cfg := newCfg()
	_ = pin.Pin(cfg, "dev", "API_KEY")
	_ = pin.Pin(cfg, "dev", "DEBUG")
	keys, err := pin.List(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 pinned keys, got %d", len(keys))
	}
}

func TestListMissingTarget(t *testing.T) {
	cfg := newCfg()
	if _, err := pin.List(cfg, "ghost"); err == nil {
		t.Error("expected error for missing target")
	}
}
