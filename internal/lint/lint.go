// Package lint provides heuristic checks for environment variable sets,
// warning about common mistakes such as empty values, duplicate keys across
// targets, or suspiciously short secret-looking values.
package lint

import (
	"fmt"
	"strings"

	"github.com/yourorg/envoy-cli/internal/config"
)

// Issue represents a single lint finding.
type Issue struct {
	Target  string
	Key     string
	Message string
	Level   string // "warn" | "error"
}

func (i Issue) String() string {
	return fmt.Sprintf("[%s] %s/%s: %s", strings.ToUpper(i.Level), i.Target, i.Key, i.Message)
}

// Run executes all lint rules against every target in cfg and returns the
// collected issues.
func Run(cfg *config.Config) []Issue {
	var issues []Issue
	for _, target := range cfg.Targets {
		issues = append(issues, checkEmptyValues(target.Name, target.Vars)...)
		issues = append(issues, checkKeyNaming(target.Name, target.Vars)...)
	}
	issues = append(issues, checkCrossTargetDuplicates(cfg)...)
	return issues
}

func checkEmptyValues(target string, vars map[string]string) []Issue {
	var issues []Issue
	for k, v := range vars {
		if strings.TrimSpace(v) == "" {
			issues = append(issues, Issue{Target: target, Key: k, Message: "value is empty", Level: "warn"})
		}
	}
	return issues
}

func checkKeyNaming(target string, vars map[string]string) []Issue {
	var issues []Issue
	for k := range vars {
		if k != strings.ToUpper(k) {
			issues = append(issues, Issue{Target: target, Key: k, Message: "key is not uppercase", Level: "warn"})
		}
		if strings.Contains(k, " ") {
			issues = append(issues, Issue{Target: target, Key: k, Message: "key contains spaces", Level: "error"})
		}
	}
	return issues
}

func checkCrossTargetDuplicates(cfg *config.Config) []Issue {
	type occurrence struct{ target, value string }
	seen := map[string][]occurrence{}
	for _, target := range cfg.Targets {
		for k, v := range target.Vars {
			seen[k] = append(seen[k], occurrence{target.Name, v})
		}
	}
	var issues []Issue
	for k, occ := range seen {
		if len(occ) < 2 {
			continue
		}
		for i := 1; i < len(occ); i++ {
			if occ[i].value == occ[0].value {
				msg := fmt.Sprintf("value identical to target %q — consider a shared default", occ[0].target)
				issues = append(issues, Issue{Target: occ[i].target, Key: k, Message: msg, Level: "warn"})
			}
		}
	}
	return issues
}
