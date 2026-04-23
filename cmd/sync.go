package cmd

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/sync"
	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool
	var keys []string

	syncCmd := &cobra.Command{
		Use:   "sync <source> <destination>",
		Short: "Sync env vars from one target into another",
		Long: `Copies keys from the source target into the destination target.
Keys missing in the destination are always added.
Existing keys are only overwritten when --overwrite is set.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			res, err := sync.Targets(cfg, args[0], args[1], sync.Options{
				Overwrite: overwrite,
				Keys:      keys,
			})
			if err != nil {
				return err
			}

			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}

			if len(res.Added) > 0 {
				fmt.Printf("Added   : %s\n", strings.Join(res.Added, ", "))
			}
			if len(res.Updated) > 0 {
				fmt.Printf("Updated : %s\n", strings.Join(res.Updated, ", "))
			}
			if len(res.Skipped) > 0 {
				fmt.Printf("Skipped : %s\n", strings.Join(res.Skipped, ", "))
			}
			fmt.Printf("Sync complete: %d added, %d updated, %d skipped\n",
				len(res.Added), len(res.Updated), len(res.Skipped))
			return nil
		},
	}

	syncCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite existing keys in destination")
	syncCmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "comma-separated list of keys to sync (default: all)")

	rootCmd.AddCommand(syncCmd)
}
