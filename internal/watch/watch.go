// Package watch provides file-watching utilities that reload the envoy
// configuration whenever the backing file changes on disk.
package watch

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/envoy-cli/envoy/internal/config"
)

// Event is emitted each time a reload occurs.
type Event struct {
	Cfg *config.Config
	Err error
}

// Options controls Watch behaviour.
type Options struct {
	// Interval between file-stat polls. Defaults to 500 ms.
	Interval time.Duration
}

// Watch polls the file at path on the given interval and sends an Event on
// the returned channel whenever the file changes. The caller must close done
// to stop the watcher.
func Watch(path string, opts Options, done <-chan struct{}) <-chan Event {
	if opts.Interval <= 0 {
		opts.Interval = 500 * time.Millisecond
	}

	ch := make(chan Event, 1)

	go func() {
		defer close(ch)

		last, _ := hashFile(path)

		for {
			select {
			case <-done:
				return
			case <-time.After(opts.Interval):
				current, err := hashFile(path)
				if err != nil {
					ch <- Event{Err: err}
					continue
				}
				if current == last {
					continue
				}
				last = current
				cfg, loadErr := config.Load(path)
				ch <- Event{Cfg: cfg, Err: loadErr}
			}
		}
	}()

	return ch
}

// hashFile returns an MD5 hex digest of the file contents.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
