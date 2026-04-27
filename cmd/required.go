package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/required"
)

func init() {
	var cfgPath string
	var failOnMissing bool

	requiredCmd := &cobra.Command{
		Use:   "required",
		Short: "Check that every target defines all shared required keys",
		Long: `Computes the intersection of keys across all targets (the "required" set)
and reports any target that is missing one of those keys.

Exit code 1 is returned when --fail is set and missing keys are found.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			results, err := required.Check(cfg)
			if err != nil {
				return err
			}

			fmt.Print(required.Format(results))

			if failOnMissing {
				for _, r := range results {
					if len(r.Missing) > 0 {
						os.Exit(1)
					}
				}
			}
			return nil
		},
	}

	requiredCmd.Flags().StringVarP(&cfgPath, "config", "c", config.DefaultPath, "path to envoy config file")
	requiredCmd.Flags().BoolVar(&failOnMissing, "fail", false, "exit with code 1 if any required keys are missing")

	rootCmd.AddCommand(requiredCmd)
}
