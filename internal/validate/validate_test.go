package validate_test

import (
	"testing"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/validate"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {"FOO": "bar"},
		},
	}
}

func TestKeyFormat(t *testing.T) {
	valid := []string{"FOO", "FOO_BAR", "_PRIVATE", "A1"}
	for _, k := range valid {
		if err := validate.KeyFormat(k); err != nil {
			t.Errorf("expected %q to be valid: %v", k, err)
		}
	}
	invalid := []string{"1FOO", "foo-bar", "", "FOO BAR"}
	for _, k := range invalid {
		if err := validate.KeyFormat(k); err == nil {
			t.Errorf("expected %q to be invalid", k)
		}
	}
}

func TestTargetExists(t *testing.T) {
	cfg := newCfg()
	if err := validate.TargetExists(cfg, "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validate.TargetExists(cfg, "staging"); err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestNoReservedKeys(t *testing.T) {
	if err := validate.NoReservedKeys("__ENVOY_TARGET__"); err == nil {
		t.Fatal("expected error for reserved key")
	}
	if err := validate.NoReservedKeys("MY_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAll(t *testing.T) {
	cfg := newCfg()
	if err := validate.All(cfg, "production", "VALID_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := validate.All(cfg, "staging", "VALID_KEY"); err == nil {
		t.Fatal("expected error for missing target")
	}
	if err := validate.All(cfg, "production", "bad-key"); err == nil {
		t.Fatal("expected error for bad key format")
	}
	if err := validate.All(cfg, "production", "__ENVOY_VERSION__"); err == nil {
		t.Fatal("expected error for reserved key")
	}
}
