package encrypt_test

import (
	"testing"

	"github.com/user/envoy-cli/internal/config"
	"github.com/user/envoy-cli/internal/encrypt"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {"DB_PASS": "secret", "API_KEY": "abc123"},
			"staging":    {"DB_PASS": "staging-secret"},
		},
	}
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	plain := "my-super-secret"
	pass := "passphrase123"

	enc, err := encrypt.EncryptValue(pass, plain)
	if err != nil {
		t.Fatalf("EncryptValue: %v", err)
	}
	if !encrypt.IsEncrypted(enc) {
		t.Fatal("expected encrypted prefix")
	}

	dec, err := encrypt.DecryptValue(pass, enc)
	if err != nil {
		t.Fatalf("DecryptValue: %v", err)
	}
	if dec != plain {
		t.Fatalf("expected %q, got %q", plain, dec)
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	enc, _ := encrypt.EncryptValue("correct", "value")
	_, err := encrypt.DecryptValue("wrong", enc)
	if err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestDecryptNonEncryptedValue(t *testing.T) {
	_, err := encrypt.DecryptValue("pass", "plaintext")
	if err == nil {
		t.Fatal("expected error for non-encrypted value")
	}
}

func TestIsEncrypted(t *testing.T) {
	if encrypt.IsEncrypted("hello") {
		t.Fatal("plain text should not be encrypted")
	}
	enc, _ := encrypt.EncryptValue("p", "v")
	if !encrypt.IsEncrypted(enc) {
		t.Fatal("encrypted value should be detected")
	}
}

func TestEncryptTarget(t *testing.T) {
	cfg := newCfg()
	if err := encrypt.EncryptTarget(cfg, "production", "pass"); err != nil {
		t.Fatalf("EncryptTarget: %v", err)
	}
	for k, v := range cfg.Targets["production"] {
		if !encrypt.IsEncrypted(v) {
			t.Errorf("key %s not encrypted", k)
		}
	}
}

func TestEncryptTargetMissing(t *testing.T) {
	cfg := newCfg()
	err := encrypt.EncryptTarget(cfg, "nonexistent", "pass")
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestDecryptTarget(t *testing.T) {
	cfg := newCfg()
	original := map[string]string{}
	for k, v := range cfg.Targets["production"] {
		original[k] = v
	}
	_ = encrypt.EncryptTarget(cfg, "production", "pass")
	if err := encrypt.DecryptTarget(cfg, "production", "pass"); err != nil {
		t.Fatalf("DecryptTarget: %v", err)
	}
	for k, v := range cfg.Targets["production"] {
		if v != original[k] {
			t.Errorf("key %s: expected %q, got %q", k, original[k], v)
		}
	}
}
