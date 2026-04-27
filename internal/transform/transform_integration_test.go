package transform_test

import (
	"encoding/base64"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/transform"
)

// TestRoundtripBase64 encodes then decodes and expects the original value.
func TestRoundtripBase64(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"dev": {"TOKEN": "supersecret"},
		},
	}

	_, err := transform.Target(cfg, "dev", transform.Base64Encode, []string{"TOKEN"}, false)
	if err != nil {
		t.Fatal(err)
	}
	encoded := cfg.Targets["dev"]["TOKEN"]
	expected := base64.StdEncoding.EncodeToString([]byte("supersecret"))
	if encoded != expected {
		t.Fatalf("encoded mismatch: %s", encoded)
	}

	_, err = transform.Target(cfg, "dev", transform.Base64Decode, []string{"TOKEN"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Targets["dev"]["TOKEN"] != "supersecret" {
		t.Fatalf("roundtrip failed: %s", cfg.Targets["dev"]["TOKEN"])
	}
}

// TestTransformMultipleKeys applies a transform to a subset of keys and
// verifies that untouched keys are unchanged.
func TestTransformMultipleKeys(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"staging": {
				"HOST": "example.com",
				"PORT": "8080",
				"ENV":  "staging",
			},
		},
	}

	_, err := transform.Target(cfg, "staging", transform.Uppercase, []string{"HOST", "ENV"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Targets["staging"]["HOST"] != "EXAMPLE.COM" {
		t.Errorf("HOST not uppercased: %s", cfg.Targets["staging"]["HOST"])
	}
	if cfg.Targets["staging"]["ENV"] != "STAGING" {
		t.Errorf("ENV not uppercased: %s", cfg.Targets["staging"]["ENV"])
	}
	if cfg.Targets["staging"]["PORT"] != "8080" {
		t.Errorf("PORT should be unchanged: %s", cfg.Targets["staging"]["PORT"])
	}
}
