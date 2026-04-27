package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-cli/internal/config"
	"envoy-cli/internal/placeholder"
	"github.com/spf13/cobra"
)

func init() {
	var targets []string
	var extraPatterns []string
	var failOnFound bool

	placeholderCmd := &cobra.Command{
		Use:   "placeholder",
		Short: "Detect placeholder or stub values in environment variables",
		Long: `Scans environment variable values for common placeholder patterns
such as TODO, FIXME, CHANGEME, <value>, etc.

Additional patterns can be provided with --pattern.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			// Allow targets as positional args or --target flag
			scope := targets
			if len(args) > 0 {
				scope = args
			}

			// Normalise extra patterns to lowercase
			var normalised []string
			for _, p := range extraPatterns {
				normalised = append(normalised, strings.ToLower(p))
			}

			results, err := placeholder.Find(cfg, scope, normalised)
			if err != nil {
				return err
			}

			fmt.Print(placeholder.Format(results))

			if failOnFound && len(results) > 0 {
				os.Exit(1)
			}
			return nil
		},
	}

	placeholderCmd.Flags().StringSliceVarP(&targets, "target", "t", nil, "Targets to scan (default: all)")
	placeholderCmd.Flags().StringSliceVarP(&extraPatterns, "pattern", "p", nil, "Additional placeholder patterns to detect")
	placeholderCmd.Flags().BoolVar(&failOnFound, "fail", false, "Exit with code 1 if any placeholders are found")

	rootCmd.AddCommand(placeholderCmd)
}
