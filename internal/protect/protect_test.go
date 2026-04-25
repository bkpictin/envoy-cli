package protect_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/protect"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]config.Target{
			"dev": {Vars: map[string]string{"DB_URL": "postgres://localhost", "API_KEY": "abc"}},
			"prod": {Vars: map[string]string{"DB_URL": "postgres://prod"}},
		},
	}
}

func TestProtectAndIsProtected(t *testing.T) {
	cfg := newCfg()
	if err := protect.Protect(cfg, "dev", "DB_URL"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !protect.IsProtected(cfg, "dev", "DB_URL") {
		t.Error("expected DB_URL to be protected")
	}
	if protect.IsProtected(cfg, "dev", "API_KEY") {
		t.Error("API_KEY should not be protected")
	}
}

func TestProtectAlreadyProtected(t *testing.T) {
	cfg := newCfg()
	_ = protect.Protect(cfg, "dev", "DB_URL")
	err := protect.Protect(cfg, "dev", "DB_URL")
	if err == nil {
		t.Fatal("expected error for duplicate protection")
	}
}

func TestProtectMissingTarget(t *testing.T) {
	cfg := newCfg()
	err := protect.Protect(cfg, "staging", "DB_URL")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestProtectMissingKey(t *testing.T) {
	cfg := newCfg()
	err := protect.Protect(cfg, "dev", "MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestUnprotect(t *testing.T) {
	cfg := newCfg()
	_ = protect.Protect(cfg, "dev", "API_KEY")
	if err := protect.Unprotect(cfg, "dev", "API_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if protect.IsProtected(cfg, "dev", "API_KEY") {
		t.Error("expected API_KEY to be unprotected")
	}
}

func TestUnprotectNotProtected(t *testing.T) {
	cfg := newCfg()
	err := protect.Unprotect(cfg, "dev", "DB_URL")
	if err == nil {
		t.Fatal("expected error when unprotecting a non-protected key")
	}
}

func TestList(t *testing.T) {
	cfg := newCfg()
	_ = protect.Protect(cfg, "dev", "DB_URL")
	_ = protect.Protect(cfg, "dev", "API_KEY")
	keys, err := protect.List(cfg, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 protected keys, got %d", len(keys))
	}
}

func TestListMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := protect.List(cfg, "nope")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}
