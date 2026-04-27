package transform_test

import (
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/transform"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"APP_NAME": "myapp",
				"DB_HOST": "  localhost  ",
				"SECRET": "aGVsbG8=", // base64("hello")
			},
		},
	}
}

func TestUppercase(t *testing.T) {
	cfg := newCfg()
	cfg.Targets["dev"]["APP_NAME"] = "myapp"
	res, err := transform.Target(cfg, "dev", transform.Uppercase, []string{"APP_NAME"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 1 || res[0].After != "MYAPP" {
		t.Fatalf("expected MYAPP, got %v", res)
	}
	if cfg.Targets["dev"]["APP_NAME"] != "MYAPP" {
		t.Fatal("value not updated in config")
	}
}

func TestLowercase(t *testing.T) {
	cfg := newCfg()
	res, err := transform.Target(cfg, "dev", transform.Lowercase, []string{"APP_NAME"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if res[0].After != "myapp" {
		t.Fatalf("expected myapp, got %s", res[0].After)
	}
}

func TestBase64Decode(t *testing.T) {
	cfg := newCfg()
	res, err := transform.Target(cfg, "dev", transform.Base64Decode, []string{"SECRET"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if res[0].After != "hello" {
		t.Fatalf("expected hello, got %s", res[0].After)
	}
}

func TestTrimSpace(t *testing.T) {
	cfg := newCfg()
	res, err := transform.Target(cfg, "dev", transform.TrimSpace, []string{"DB_HOST"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if res[0].After != "localhost" {
		t.Fatalf("expected localhost, got %q", res[0].After)
	}
}

func TestDryRunDoesNotMutate(t *testing.T) {
	cfg := newCfg()
	_, err := transform.Target(cfg, "dev", transform.Uppercase, []string{"APP_NAME"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Targets["dev"]["APP_NAME"] != "myapp" {
		t.Fatal("dry-run should not mutate config")
	}
}

func TestMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := transform.Target(cfg, "prod", transform.Uppercase, nil, false)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestAllKeysWhenNoneSpecified(t *testing.T) {
	cfg := newCfg()
	res, err := transform.Target(cfg, "dev", transform.TrimSpace, nil, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != len(cfg.Targets["dev"]) {
		t.Fatalf("expected %d results, got %d", len(cfg.Targets["dev"]), len(res))
	}
}
