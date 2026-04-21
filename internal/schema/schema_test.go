package schema_test

import (
	"strings"
	"testing"

	"envoy-cli/internal/config"
	"envoy-cli/internal/schema"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"production": {
				"DATABASE_URL": "postgres://prod",
				"API_KEY":      "secret",
				"EXTRA_KEY":    "unexpected",
			},
		},
	}
}

func baseSchema() schema.Schema {
	return schema.Schema{
		"DATABASE_URL": {Required: true, Description: "Primary database connection string"},
		"API_KEY":      {Required: true, Description: "External API key"},
		"LOG_LEVEL":    {Required: false, Description: "Optional log verbosity"},
	}
}

func TestValidateMissingRequired(t *testing.T) {
	cfg := newCfg()
	s := baseSchema()
	s["REQUIRED_MISSING"] = schema.Rule{Required: true, Description: "Must be set"}

	results, err := schema.Validate(cfg, "production", s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var found bool
	for _, r := range results {
		if r.Key == "REQUIRED_MISSING" && r.Level == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected error for missing required key")
	}
}

func TestValidateUndeclaredKey(t *testing.T) {
	cfg := newCfg()
	results, err := schema.Validate(cfg, "production", baseSchema())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var found bool
	for _, r := range results {
		if r.Key == "EXTRA_KEY" && r.Level == "warning" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning for undeclared key EXTRA_KEY")
	}
}

func TestValidateMissingTarget(t *testing.T) {
	cfg := newCfg()
	_, err := schema.Validate(cfg, "staging", baseSchema())
	if err == nil {
		t.Error("expected error for missing target")
	}
}

func TestValidateClean(t *testing.T) {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"staging": {"DATABASE_URL": "postgres://stg", "API_KEY": "key123"},
		},
	}
	results, err := schema.Validate(cfg, "staging", baseSchema())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, r := range results {
		if r.Level == "error" {
			t.Errorf("unexpected error: %s", r.Message)
		}
	}
}

func TestFormat(t *testing.T) {
	results := []schema.Result{
		{Target: "prod", Key: "FOO", Message: "required key \"FOO\" is missing", Level: "error"},
	}
	out := schema.Format(results)
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", out)
	}
}

func TestFormatEmpty(t *testing.T) {
	out := schema.Format(nil)
	if !strings.Contains(out, "passed") {
		t.Errorf("expected pass message, got: %s", out)
	}
}
