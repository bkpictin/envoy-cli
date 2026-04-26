package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/dedupe"
)

func init() {
	var cfgPath string
	var jsonOut bool

	dedupeCmd := &cobra.Command{
		Use:   "dedupe",
		Short: "Find duplicate values shared across targets",
		Long: `Scans all targets and reports environment variable keys whose values
are identical in two or more targets. Useful for identifying shared secrets
or configuration that could be consolidated.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			result := dedupe.FindCrossTarget(cfg)

			if jsonOut {
				for _, m := range result.Matches {
					fmt.Fprintf(os.Stdout, `{"key":%q,"value":%q,"targets":%q}`+"\n",
						m.Key, m.Value, m.Targets)
				}
				return nil
			}

			fmt.Print(dedupe.Format(result))
			return nil
		},
	}

	dedupeCmd.Flags().StringVarP(&cfgPath, "config", "c", config.DefaultPath, "path to envoy config file")
	dedupeCmd.Flags().BoolVar(&jsonOut, "json", false, "output results as JSON lines")

	rootCmd.AddCommand(dedupeCmd)
}
