package mask_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/mask"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"prod": {
				"API_KEY":    "supersecretvalue",
				"SHORT":      "ab",
				"DB_PASS":    "enc:abc123encrypted",
			},
			"staging": {
				"TOKEN": "stagingtoken1234",
			},
		},
	}
}

func TestMaskValueDefault(t *testing.T) {
	out := mask.MaskValue("supersecretvalue", 4)
	if !strings.HasSuffix(out, "alue") {
		t.Fatalf("expected suffix 'alue', got %q", out)
	}
	if !strings.HasPrefix(out, "************") {
		t.Fatalf("expected asterisk prefix, got %q", out)
	}
}

func TestMaskValueShort(t *testing.T) {
	out := mask.MaskValue("ab", 4)
	if out != "**" {
		t.Fatalf("expected '**', got %q", out)
	}
}

func TestMaskValueEncrypted(t *testing.T) {
	out := mask.MaskValue("enc:abc123encrypted", 4)
	if !strings.Contains(out, "[encrypted]") {
		t.Fatalf("expected '[encrypted]' marker, got %q", out)
	}
}

func TestTarget(t *testing.T) {
	cfg := newCfg()
	results, err := mask.Target(cfg, "prod", 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Target != "prod" {
			t.Errorf("expected target 'prod', got %q", r.Target)
		}
	}
}

func TestTargetMissing(t *testing.T) {
	cfg := newCfg()
	_, err := mask.Target(cfg, "nonexistent", 4)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestAll(t *testing.T) {
	cfg := newCfg()
	results := mask.All(cfg, 4)
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}
