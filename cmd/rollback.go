package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/rollback"
)

func init() {
	var snapshotName string
	var listOnly bool

	rollbackCmd := &cobra.Command{
		Use:   "rollback <target>",
		Short: "Revert a target to a previous snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			if listOnly {
				names, err := rollback.ListAvailable(cfg, target)
				if err != nil {
					return err
				}
				fmt.Printf("Available rollback points for %q:\n", target)
				for i, n := range names {
					fmt.Printf("  [%d] %s\n", i+1, n)
				}
				return nil
			}

			var restoredName string
			if snapshotName != "" {
				if err := rollback.ToSnapshot(cfg, target, snapshotName); err != nil {
					return err
				}
				restoredName = snapshotName
			} else {
				restoredName, err = rollback.ToPrevious(cfg, target)
				if err != nil {
					return err
				}
			}

			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Rolled back %q to snapshot %q\n", target, restoredName)
			return nil
		},
	}

	rollbackCmd.Flags().StringVarP(&snapshotName, "snapshot", "s", "", "Name of snapshot to restore (default: most recent)")
	rollbackCmd.Flags().BoolVarP(&listOnly, "list", "l", false, "List available rollback points without restoring")

	rootCmd.AddCommand(rollbackCmd)
}
