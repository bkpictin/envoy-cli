package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/envoy-cli/internal/config"
	"github.com/yourorg/envoy-cli/internal/lint"
)

func init() {
	var cfgPath string
	var strict bool

	lintCmd := &cobra.Command{
		Use:   "lint",
		Short: "Check environment variable sets for common issues",
		Long: `Runs heuristic checks across all targets and reports warnings and errors.

Use --strict to exit with a non-zero status code when any issue is found.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := cfgPath
			if path == "" {
				path = config.DefaultPath()
			}
			cfg, err := config.Load(path)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			issues := lint.Run(cfg)
			if len(issues) == 0 {
				fmt.Println("No issues found.")
				return nil
			}

			hasError := false
			for _, iss := range issues {
				fmt.Println(iss.String())
				if iss.Level == "error" {
					hasError = true
				}
			}

			if strict || hasError {
				os.Exit(1)
			}
			return nil
		},
	}

	lintCmd.Flags().StringVar(&cfgPath, "config", "", "path to envoy config file")
	lintCmd.Flags().BoolVar(&strict, "strict", false, "exit non-zero on warnings as well as errors")
	rootCmd.AddCommand(lintCmd)
}
