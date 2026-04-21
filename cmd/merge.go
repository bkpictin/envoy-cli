package cmd

import (
	"fmt"
	"os"

	"envoy-cli/internal/audit"
	"envoy-cli/internal/config"
	"envoy-cli/internal/merge"

	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	mergeCmd := &cobra.Command{
		Use:   "merge <dest> <src1> [src2...]",
		Short: "Merge environment variables from one or more sources into a destination target",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := config.DefaultPath()
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			dest := args[0]
			sources := args[1:]

			strategy := merge.StrategySkip
			if overwrite {
				strategy = merge.StrategyOverwrite
			}

			res, err := merge.Targets(cfg, dest, sources, strategy)
			if err != nil {
				return err
			}

			if err := config.Save(cfgPath, cfg); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			audit.Log(cfg, "merge", dest, fmt.Sprintf("sources=%v strategy=%s merged=%d skipped=%d overwrote=%d",
				sources, strategy, res.Merged, res.Skipped, res.Overwrote))
			_ = config.Save(cfgPath, cfg)

			fmt.Fprintf(os.Stdout, "Merged into %q: %d added, %d skipped, %d overwritten\n",
				dest, res.Merged, res.Skipped, res.Overwrote)
			return nil
		},
	}

	mergeCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	rootCmd.AddCommand(mergeCmd)
}
