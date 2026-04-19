package export

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"prod": {
				"APP_ENV": "production",
				"DB_URL":  "postgres://prod",
			},
		},
		Snapshots: map[string]map[string]map[string]string{},
	}
}

func TestDotenv(t *testing.T) {
	out, err := ToFile(newCfg(), "prod", FormatDotenv)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected dotenv line, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_URL=postgres://prod") {
		t.Errorf("expected dotenv line, got:\n%s", out)
	}
}

func TestShell(t *testing.T) {
	out, err := ToFile(newCfg(), "prod", FormatShell)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected shell export, got:\n%s", out)
	}
}

func TestJSON(t *testing.T) {
	out, err := ToFile(newCfg(), "prod", FormatJSON)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(out, "{") {
		t.Errorf("expected JSON object, got:\n%s", out)
	}
	if !strings.Contains(out, `"APP_ENV"`) {
		t.Errorf("expected key in JSON, got:\n%s", out)
	}
}

func TestMissingTarget(t *testing.T) {
	_, err := ToFile(newCfg(), "dev", FormatDotenv)
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestUnsupportedFormat(t *testing.T) {
	_, err := ToFile(newCfg(), "prod", Format("xml"))
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}
