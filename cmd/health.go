package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/health"
)

func init() {
	var cfgPath string
	var errorsOnly bool

	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Check the integrity of the envoy configuration",
		Long: `Runs a series of checks against the configuration file and reports
warnings or errors such as empty targets, blank values, or snapshots
that reference non-existent targets.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfgPath == "" {
				cfgPath = config.DefaultPath()
			}
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			report := health.Check(cfg)

			if errorsOnly {
				errs := health.FilterByLevel(report, health.Error)
				if len(errs) == 0 {
					fmt.Println("no errors found")
					return nil
				}
				for _, e := range errs {
					if e.Target != "" {
						fmt.Fprintf(os.Stderr, "[ERROR] (%s) %s\n", e.Target, e.Message)
					} else {
						fmt.Fprintf(os.Stderr, "[ERROR] %s\n", e.Message)
					}
				}
				os.Exit(1)
			}

			fmt.Println(health.Format(report))
			if !report.OK() {
				os.Exit(1)
			}
			return nil
		},
	}

	healthCmd.Flags().StringVar(&cfgPath, "config", "", "path to envoy config file")
	healthCmd.Flags().BoolVar(&errorsOnly, "errors-only", false, "only report ERROR level issues")

	rootCmd.AddCommand(healthCmd)
}
