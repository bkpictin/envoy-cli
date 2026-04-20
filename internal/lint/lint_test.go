package lint

import (
	"testing"

	"github.com/yourorg/envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: []config.Target{
			{Name: "dev", Vars: map[string]string{}},
			{Name: "prod", Vars: map[string]string{}},
		},
	}
}

func TestEmptyValue(t *testing.T) {
	cfg := newCfg()
	cfg.Targets[0].Vars["API_KEY"] = ""
	issues := Run(cfg)
	if len(issues) == 0 {
		t.Fatal("expected at least one issue for empty value")
	}
	if issues[0].Level != "warn" {
		t.Errorf("expected warn, got %s", issues[0].Level)
	}
}

func TestLowercaseKey(t *testing.T) {
	cfg := newCfg()
	cfg.Targets[0].Vars["api_key"] = "secret"
	issues := Run(cfg)
	found := false
	for _, iss := range issues {
		if iss.Key == "api_key" && iss.Level == "warn" {
			found = true
		}
	}
	if !found {
		t.Error("expected warn for lowercase key")
	}
}

func TestKeyWithSpaces(t *testing.T) {
	cfg := newCfg()
	cfg.Targets[0].Vars["BAD KEY"] = "value"
	issues := Run(cfg)
	found := false
	for _, iss := range issues {
		if iss.Key == "BAD KEY" && iss.Level == "error" {
			found = true
		}
	}
	if !found {
		t.Error("expected error for key with spaces")
	}
}

func TestCrossTargetDuplicate(t *testing.T) {
	cfg := newCfg()
	cfg.Targets[0].Vars["DB_URL"] = "postgres://localhost/db"
	cfg.Targets[1].Vars["DB_URL"] = "postgres://localhost/db"
	issues := Run(cfg)
	found := false
	for _, iss := range issues {
		if iss.Key == "DB_URL" {
			found = true
		}
	}
	if !found {
		t.Error("expected cross-target duplicate warning")
	}
}

func TestNoIssues(t *testing.T) {
	cfg := newCfg()
	cfg.Targets[0].Vars["APP_ENV"] = "development"
	cfg.Targets[1].Vars["APP_ENV"] = "production"
	issues := Run(cfg)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}
