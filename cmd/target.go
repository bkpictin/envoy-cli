package cmd

import (
	"fmt"
	"sort"

	"github.com/envoy-cli/envoy-cli/internal/config"
	"github.com/envoy-cli/envoy-cli/internal/target"
	"github.com/spf13/cobra"
)

var targetCmd = &cobra.Command{
	Use:   "target",
	Short: "Manage deployment targets",
}

var targetListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(config.DefaultPath())
		if err != nil {
			return err
		}
		names := target.List(cfg)
		sort.Strings(names)
		for _, n := range names {
			fmt.Println(n)
		}
		return nil
	},
}

var targetAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new target",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(config.DefaultPath())
		if err != nil {
			return err
		}
		if err := target.Add(cfg, args[0]); err != nil {
			return err
		}
		return config.Save(cfg, config.DefaultPath())
	},
}

var targetRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a target",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(config.DefaultPath())
		if err != nil {
			return err
		}
		if err := target.Remove(cfg, args[0]); err != nil {
			return err
		}
		return config.Save(cfg, config.DefaultPath())
	},
}

var targetRenameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "Rename a target",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(config.DefaultPath())
		if err != nil {
			return err
		}
		if err := target.Rename(cfg, args[0], args[1]); err != nil {
			return err
		}
		return config.Save(cfg, config.DefaultPath())
	},
}

func init() {
	targetCmd.AddCommand(targetListCmd, targetAddCmd, targetRemoveCmd, targetRenameCmd)
	rootCmd.AddCommand(targetCmd)
}
