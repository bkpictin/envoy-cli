package cmd

import (
	"fmt"
	"sort"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/prefix"
	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	prefixCmd := &cobra.Command{
		Use:   "prefix",
		Short: "Add, remove, or list key prefixes in a target",
	}

	addCmd := &cobra.Command{
		Use:   "add <target> <prefix>",
		Short: "Add a prefix to all keys in a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			res, err := prefix.Add(cfg, args[0], args[1], overwrite)
			if err != nil {
				return err
			}
			fmt.Printf("target %q: %d key(s) prefixed, %d skipped\n", res.Target, res.Changed, res.Skipped)
			return config.Save(config.DefaultPath, cfg)
		},
	}
	addCmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing prefixed keys")

	removeCmd := &cobra.Command{
		Use:   "remove <target> <prefix>",
		Short: "Remove a prefix from all matching keys in a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			res, err := prefix.Remove(cfg, args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("target %q: %d key(s) un-prefixed, %d unchanged\n", res.Target, res.Changed, res.Skipped)
			return config.Save(config.DefaultPath, cfg)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <target> <prefix>",
		Short: "List keys that start with the given prefix",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			keys, err := prefix.List(cfg, args[0], args[1])
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Println("no keys found")
				return nil
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		},
	}

	prefixCmd.AddCommand(addCmd, removeCmd, listCmd)
	rootCmd.AddCommand(prefixCmd)
}
