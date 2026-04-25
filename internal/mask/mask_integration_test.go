package mask_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/mask"
)

func TestAllReturnsEveryTarget(t *testing.T) {
	cfg := newCfg()
	results := mask.All(cfg, 4)

	targetsSeen := map[string]bool{}
	for _, r := range results {
		targetsSeen[r.Target] = true
	}
	if !targetsSeen["prod"] || !targetsSeen["staging"] {
		t.Fatalf("expected both 'prod' and 'staging' in results, got %v", targetsSeen)
	}
}

func TestMaskNeverExposesFullValue(t *testing.T) {
	cfg := newCfg()
	results, err := mask.Target(cfg, "prod", 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	originals := cfg.Targets["prod"]
	for _, r := range results {
		original := originals[r.Key]
		if len(original) > 4 && r.Masked == original {
			t.Errorf("key %q: masked value equals original — value not masked", r.Key)
		}
	}
}

func TestZeroVisibleLen(t *testing.T) {
	out := mask.MaskValue("topsecret", 0)
	if strings.ContainsAny(out, "abcdefghijklmnopqrstuvwxyz") {
		t.Errorf("expected fully masked value, got %q", out)
	}
}

func TestVisibleLenLargerThanValue(t *testing.T) {
	out := mask.MaskValue("hi", 10)
	if out != "**" {
		t.Errorf("expected '**' for short value with large visibleLen, got %q", out)
	}
}
