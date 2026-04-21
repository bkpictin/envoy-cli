// Package importenv provides functionality for importing environment variables
// from external file formats (dotenv, shell exports, JSON) into a target.
package importenv

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"envoy-cli/internal/config"
)

// Format represents the input file format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatJSON   Format = "json"
)

// FromFile reads environment variables from a file and imports them into the
// specified target. If overwrite is false, existing keys are skipped.
func FromFile(cfg *config.Config, target, path string, format Format, overwrite bool) (int, error) {
	if _, ok := cfg.Targets[target]; !ok {
		return 0, fmt.Errorf("target %q does not exist", target)
	}

	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var pairs map[string]string
	switch format {
	case FormatDotenv:
		pairs, err = parseDotenv(bufio.NewScanner(f))
	case FormatShell:
		pairs, err = parseShell(bufio.NewScanner(f))
	case FormatJSON:
		pairs, err = parseJSON(f)
	default:
		return 0, fmt.Errorf("unsupported format: %q", format)
	}
	if err != nil {
		return 0, err
	}

	count := 0
	for k, v := range pairs {
		if _, exists := cfg.Targets[target][k]; exists && !overwrite {
			continue
		}
		cfg.Targets[target][k] = v
		count++
	}
	return count, nil
}

func parseDotenv(s *bufio.Scanner) (map[string]string, error) {
	out := map[string]string{}
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		out[strings.TrimSpace(k)] = strings.Trim(strings.TrimSpace(v), `"`)
	}
	return out, s.Err()
}

func parseShell(s *bufio.Scanner) (map[string]string, error) {
	out := map[string]string{}
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		line = strings.TrimPrefix(line, "export ")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		out[strings.TrimSpace(k)] = strings.Trim(strings.TrimSpace(v), `"`)
	}
	return out, s.Err()
}

func parseJSON(f *os.File) (map[string]string, error) {
	out := map[string]string{}
	if err := json.NewDecoder(f).Decode(&out); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return out, nil
}
