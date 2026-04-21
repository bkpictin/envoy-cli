package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/internal/config"
	"github.com/envoy-cli/internal/prune"
	"github.com/spf13/cobra"
)

func init() {
	var dryRun bool
	var mode string

	pruneCmd := &cobra.Command{
		Use:   "prune <target>",
		Short: "Remove stale or empty keys from a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetName := args[0]

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			var results []prune.Result

			switch mode {
			case "orphaned":
				results, err = prune.OrphanedKeys(cfg, targetName, dryRun)
			case "empty":
				results, err = prune.EmptyValues(cfg, targetName, dryRun)
			default:
				return fmt.Errorf("unknown mode %q: use 'orphaned' or 'empty'", mode)
			}
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Println("Nothing to prune.")
				return nil
			}

			verb := "Pruned"
			if dryRun {
				verb = "Would prune"
			}
			for _, r := range results {
				fmt.Fprintf(os.Stdout, "%s [%s] %s\n", verb, r.Target, r.Key)
			}

			if !dryRun {
				if err := config.Save(cfg, cfgPath); err != nil {
					return err
				}
			}
			return nil
		},
	}

	pruneCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Report keys that would be pruned without modifying the config")
	pruneCmd.Flags().StringVar(&mode, "mode", "empty", "Pruning mode: 'orphaned' or 'empty'")

	rootCmd.AddCommand(pruneCmd)
}
