package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/watch"
)

func writeCfg(t *testing.T, path string, cfg *config.Config) {
	t.Helper()
	if err := config.Save(cfg, path); err != nil {
		t.Fatalf("save: %v", err)
	}
}

func TestWatchDetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "envoy.json")

	cfg := &config.Config{Targets: map[string]map[string]string{
		"dev": {"FOO": "bar"},
	}}
	writeCfg(t, path, cfg)

	done := make(chan struct{})
	defer close(done)

	ch := watch.Watch(path, watch.Options{Interval: 50 * time.Millisecond}, done)

	// Give the watcher time to record the initial hash.
	time.Sleep(120 * time.Millisecond)

	// Mutate the file.
	cfg.Targets["dev"]["FOO"] = "changed"
	writeCfg(t, path, cfg)

	select {
	case ev := <-ch:
		if ev.Err != nil {
			t.Fatalf("unexpected error: %v", ev.Err)
		}
		if ev.Cfg == nil {
			t.Fatal("expected non-nil config")
		}
		if ev.Cfg.Targets["dev"]["FOO"] != "changed" {
			t.Errorf("got %q, want %q", ev.Cfg.Targets["dev"]["FOO"], "changed")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for watch event")
	}
}

func TestWatchNoSpuriousEvents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "envoy.json")

	cfg := &config.Config{Targets: map[string]map[string]string{
		"prod": {"KEY": "value"},
	}}
	writeCfg(t, path, cfg)

	done := make(chan struct{})
	defer close(done)

	ch := watch.Watch(path, watch.Options{Interval: 50 * time.Millisecond}, done)

	// No changes — expect silence for 300 ms.
	select {
	case ev := <-ch:
		t.Fatalf("unexpected event: %+v", ev)
	case <-time.After(300 * time.Millisecond):
		// pass
	}
}

func TestWatchMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	done := make(chan struct{})
	defer close(done)

	ch := watch.Watch(path, watch.Options{Interval: 50 * time.Millisecond}, done)

	select {
	case ev := <-ch:
		if ev.Err == nil {
			t.Fatal("expected error for missing file")
		}
		if !os.IsNotExist(ev.Err) {
			t.Fatalf("unexpected error type: %v", ev.Err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for error event")
	}
}
