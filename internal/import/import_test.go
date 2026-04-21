package importenv_test

import (
	"os"
	"path/filepath"
	"testing"

	importenv "envoy-cli/internal/import"
	"envoy-cli/internal/config"
)

func newCfg() *config.Config {
	return &config.Config{
		Targets: map[string]map[string]string{
			"dev": {"EXISTING": "old"},
		},
		Snapshots: map[string]map[string]map[string]string{},
	}
}

func writeFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "envoy-import-*")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	return f.Name()
}

func TestImportDotenv(t *testing.T) {
	cfg := newCfg()
	path := writeFile(t, "# comment\nFOO=bar\nBAZ=\"quoted\"\n")
	n, err := importenv.FromFile(cfg, "dev", path, importenv.FormatDotenv, false)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 imported, got %d", n)
	}
	if cfg.Targets["dev"]["FOO"] != "bar" {
		t.Errorf("FOO mismatch")
	}
	if cfg.Targets["dev"]["BAZ"] != "quoted" {
		t.Errorf("BAZ mismatch")
	}
}

func TestImportShell(t *testing.T) {
	cfg := newCfg()
	path := writeFile(t, "export APP=myapp\nexport PORT=8080\n")
	n, err := importenv.FromFile(cfg, "dev", path, importenv.FormatShell, false)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 imported, got %d", n)
	}
	if cfg.Targets["dev"]["APP"] != "myapp" {
		t.Errorf("APP mismatch")
	}
}

func TestImportJSON(t *testing.T) {
	cfg := newCfg()
	path := writeFile(t, `{"DB_HOST":"localhost","DB_PORT":"5432"}`)
	n, err := importenv.FromFile(cfg, "dev", path, importenv.FormatJSON, false)
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatalf("expected 2 imported, got %d", n)
	}
	if cfg.Targets["dev"]["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST mismatch")
	}
}

func TestImportNoOverwrite(t *testing.T) {
	cfg := newCfg()
	path := writeFile(t, "EXISTING=new\n")
	n, _ := importenv.FromFile(cfg, "dev", path, importenv.FormatDotenv, false)
	if n != 0 {
		t.Fatalf("expected 0 imported (no overwrite), got %d", n)
	}
	if cfg.Targets["dev"]["EXISTING"] != "old" {
		t.Error("existing key should not be overwritten")
	}
}

func TestImportOverwrite(t *testing.T) {
	cfg := newCfg()
	path := writeFile(t, "EXISTING=new\n")
	n, _ := importenv.FromFile(cfg, "dev", path, importenv.FormatDotenv, true)
	if n != 1 {
		t.Fatalf("expected 1 imported, got %d", n)
	}
	if cfg.Targets["dev"]["EXISTING"] != "new" {
		t.Error("existing key should be overwritten")
	}
}

func TestImportMissingTarget(t *testing.T) {
	cfg := newCfg()
	path := filepath.Join(t.TempDir(), "x.env")
	_, err := importenv.FromFile(cfg, "prod", path, importenv.FormatDotenv, false)
	if err == nil {
		t.Error("expected error for missing target")
	}
}
