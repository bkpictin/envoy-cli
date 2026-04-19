package export

import (
	"fmt"
	"sort"
	"strings"

	"envoy-cli/internal/config"
)

// Format represents the output format for exported variables.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatShell  Format = "shell"
	FormatJSON   Format = "json"
)

// ToFile renders environment variables for a given target in the specified format.
func ToFile(cfg *config.Config, target string, format Format) (string, error) {
	envs, ok := cfg.Targets[target]
	if !ok {
		return "", fmt.Errorf("target %q not found", target)
	}

	keys := make([]string, 0, len(envs))
	for k := range envs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	switch format {
	case FormatDotenv:
		return renderDotenv(envs, keys), nil
	case FormatShell:
		return renderShell(envs, keys), nil
	case FormatJSON:
		return renderJSON(envs, keys), nil
	default:
		return "", fmt.Errorf("unsupported format %q", format)
	}
}

func renderDotenv(envs map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, envs[k])
	}
	return sb.String()
}

func renderShell(envs map[string]string, keys []string) string {
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "export %s=%q\n", k, envs[k])
	}
	return sb.String()
}

func renderJSON(envs map[string]string, keys []string) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		comma := ","
		if i == len(keys)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  %q: %q%s\n", k, envs[k], comma)
	}
	sb.WriteString("}\n")
	return sb.String()
}
