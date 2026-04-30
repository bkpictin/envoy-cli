package cmd

import (
	"fmt"
	"strings"

	"envoy-cli/internal/cascade"
	"envoy-cli/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	cascadeCmd := &cobra.Command{
		Use:   "cascade <base> <target1> [target2 ...]",
		Short: "Propagate env vars from a base target through a chain of targets",
		Long: `Cascade walks a chain of targets in order, copying environment variables
from the base (and each successive target) into the next one.
Keys already present in a target are skipped unless --overwrite is set.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			base := args[0]
			chain := args[1:]

			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			results, err := cascade.Apply(cfg, base, chain, overwrite)
			if err != nil {
				return err
			}

			for _, r := range results {
				fmt.Printf("%-20s applied=%-4d skipped=%d\n", r.Target, r.Applied, r.Skipped)
			}

			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}

			fmt.Printf("\nCascaded %q → %s\n", base, strings.Join(chain, " → "))
			return nil
		},
	}

	cascadeCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in destination targets")
	rootCmd.AddCommand(cascadeCmd)
}
