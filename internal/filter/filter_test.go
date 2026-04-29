package filter_test

import (
	"testing"

	"github.com/yourusername/envoy-cli/internal/config"
	"github.com/yourusername/envoy-cli/internal/filter"
)

func newCfg() *config.Config {
	cfg := &config.Config{
		Targets: map[string]map[string]string{
			"dev": {
				"APP_HOST":     "localhost",
				"APP_PORT":     "8080",
				"DB_HOST":      "localhost",
				"DB_PORT":      "5432",
				"SECRET_KEY":   "dev-secret",
				"FEATURE_FLAG": "true",
			},
			"prod": {
				"APP_HOST":   "prod.example.com",
				"APP_PORT":   "443",
				"DB_HOST":    "db.prod.example.com",
				"DB_PORT":    "5432",
				"SECRET_KEY": "prod-secret",
			},
		},
		Tags: map[string]map[string][]string{
			"dev": {
				"APP_HOST":     {"app", "network"},
				"APP_PORT":     {"app", "network"},
				"DB_HOST":      {"db", "network"},
				"DB_PORT":      {"db", "network"},
				"SECRET_KEY":   {"secret"},
				"FEATURE_FLAG": {"feature"},
			},
			"prod": {
				"APP_HOST": {"app", "network"},
				"APP_PORT": {"app", "network"},
				"DB_HOST":  {"db", "network"},
				"DB_PORT":  {"db", "network"},
			},
		},
	}
	return cfg
}

func TestByPatternGlob(t *testing.T) {
	cfg := newCfg()
	opts := filter.Options{Pattern: "APP_*"}
	result, err := filter.ByPattern(cfg, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in result")
	}
	if _, ok := result["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in result")
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("did not expect DB_HOST in result")
	}
}

func TestByPatternPrefix(t *testing.T) {
	cfg := newCfg()
	opts := filter.Options{Pattern: "DB_*"}
	result, err := filter.ByPattern(cfg, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 keys, got %d", len(result))
	}
}

func TestByPatternMissingTarget(t *testing.T) {
	cfg := newCfg()
	opts := filter.Options{Pattern: "*"}
	_, err := filter.ByPattern(cfg, "staging", opts)
	if err == nil {
		t.Fatal("expected error for missing target")
	}
}

func TestByPatternWithTags(t *testing.T) {
	cfg := newCfg()
	opts := filter.Options{
		Pattern: "*",
		Tags:    []string{"db"},
	}
	result, err := filter.ByPattern(cfg, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := result["DB_PORT"]; !ok {
		t.Error("expected DB_PORT in result")
	}
	if _, ok := result["APP_HOST"]; ok {
		t.Error("did not expect APP_HOST when filtering by db tag")
	}
}

func TestByPatternWithMultipleTags(t *testing.T) {
	cfg := newCfg()
	opts := filter.Options{
		Pattern: "*",
		Tags:    []string{"app", "network"},
	}
	result, err := filter.ByPattern(cfg, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Only keys that have BOTH app and network tags
	if _, ok := result["APP_HOST"]; !ok {
		t.Error("expected APP_HOST")
	}
	if _, ok := result["DB_HOST"]; ok {
		t.Error("did not expect DB_HOST — it has db+network, not app+network")
	}
}

func TestByPatternExclude(t *testing.T) {
	cfg := newCfg()
	opts := filter.Options{
		Pattern: "*",
		Exclude: "SECRET_*",
	}
	result, err := filter.ByPattern(cfg, "dev", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["SECRET_KEY"]; ok {
		t.Error("did not expect SECRET_KEY in result after exclusion")
	}
	if len(result) == 0 {
		t.Error("expected non-empty result after excluding only SECRET_*")
	}
}

func TestFormat(t *testing.T) {
	matched := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
	}
	out := filter.Format(matched)
	if out == "" {
		t.Error("expected non-empty formatted output")
	}
}

func TestFormatEmpty(t *testing.T) {
	out := filter.Format(map[string]string{})
	if out != "(no keys matched)" {
		t.Errorf("expected '(no keys matched)', got %q", out)
	}
}
