package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/config"
	"github.com/user/envoy-cli/internal/truncate"
)

func init() {
	var target string
	var maxLen int
	var dryRun bool
	var all bool

	truncateCmd := &cobra.Command{
		Use:   "truncate",
		Short: "Truncate environment variable values to a maximum length",
		Example: `  envoy truncate --target production --max 64
  envoy truncate --all --max 128 --dry-run`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := config.DefaultPath()
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			if maxLen <= 0 {
				return fmt.Errorf("--max must be greater than zero")
			}

			var results []truncate.Result

			switch {
			case all:
				results, err = truncate.All(cfg, maxLen, dryRun)
			case target != "":
				results, err = truncate.Target(cfg, target, maxLen, dryRun)
			default:
				return fmt.Errorf("specify --target <name> or --all")
			}

			if err != nil {
				return err
			}

			fmt.Print(truncate.Format(results))

			if !dryRun {
				if err := config.Save(cfg, cfgPath); err != nil {
					return fmt.Errorf("save config: %w", err)
				}
			} else {
				fmt.Fprintln(os.Stderr, "(dry-run: no changes written)")
			}
			return nil
		},
	}

	truncateCmd.Flags().StringVarP(&target, "target", "t", "", "Target to truncate")
	truncateCmd.Flags().BoolVar(&all, "all", false, "Truncate values across all targets")
	truncateCmd.Flags().IntVar(&maxLen, "max", 255, "Maximum allowed value length")
	truncateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")

	rootCmd.AddCommand(truncateCmd)
}
