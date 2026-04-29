package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/suffix"
	"github.com/spf13/cobra"
)

func init() {
	var dryRun bool

	suffixCmd := &cobra.Command{
		Use:   "suffix",
		Short: "Add or remove key suffixes within a target",
	}

	addCmd := &cobra.Command{
		Use:   "add <target> <suffix>",
		Short: "Append a suffix to all keys in a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			results, err := suffix.Add(cfg, args[0], args[1], dryRun)
			if err != nil {
				return err
			}
			sort.Slice(results, func(i, j int) bool { return results[i].OldKey < results[j].OldKey })
			for _, r := range results {
				if r.Skipped {
					fmt.Fprintf(os.Stderr, "skip  %s → %s (%s)\n", r.OldKey, r.NewKey, r.Reason)
				} else {
					fmt.Printf("rename %s → %s\n", r.OldKey, r.NewKey)
				}
			}
			if dryRun {
				return nil
			}
			return config.Save(config.DefaultPath, cfg)
		},
	}
	addCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying them")

	removeCmd := &cobra.Command{
		Use:   "remove <target> <suffix>",
		Short: "Strip a suffix from all matching keys in a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			results, err := suffix.Remove(cfg, args[0], args[1], dryRun)
			if err != nil {
				return err
			}
			sort.Slice(results, func(i, j int) bool { return results[i].OldKey < results[j].OldKey })
			for _, r := range results {
				if r.Skipped {
					fmt.Fprintf(os.Stderr, "skip  %s (%s)\n", r.OldKey, r.Reason)
				} else {
					fmt.Printf("rename %s → %s\n", r.OldKey, r.NewKey)
				}
			}
			if dryRun {
				return nil
			}
			return config.Save(config.DefaultPath, cfg)
		},
	}
	removeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying them")

	listCmd := &cobra.Command{
		Use:   "list <target> <suffix>",
		Short: "List all keys in a target that end with the given suffix",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			keys, err := suffix.List(cfg, args[0], args[1])
			if err != nil {
				return err
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		},
	}

	suffixCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(suffixCmd)
}
