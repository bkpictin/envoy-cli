package cmd

import (
	"fmt"
	"strings"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/group"
	"github.com/spf13/cobra"
)

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage logical key groups within a target",
	}

	groupCmd.AddCommand(
		groupCreateCmd(),
		groupDeleteCmd(),
		groupAddKeyCmd(),
		groupRemoveKeyCmd(),
		groupListCmd(),
		groupKeysCmd(),
	)

	rootCmd.AddCommand(groupCmd)
}

func groupCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create <target> <group>",
		Short: "Create a new group in a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := group.Create(cfg, args[0], args[1]); err != nil {
				return err
			}
			return config.Save(cfg, config.DefaultPath)
		},
	}
}

func groupDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <target> <group>",
		Short: "Delete a group from a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := group.Delete(cfg, args[0], args[1]); err != nil {
				return err
			}
			return config.Save(cfg, config.DefaultPath)
		},
	}
}

func groupAddKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add-key <target> <group> <key>",
		Short: "Add a key to a group",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := group.AddKey(cfg, args[0], args[1], args[2]); err != nil {
				return err
			}
			return config.Save(cfg, config.DefaultPath)
		},
	}
}

func groupRemoveKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove-key <target> <group> <key>",
		Short: "Remove a key from a group",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := group.RemoveKey(cfg, args[0], args[1], args[2]); err != nil {
				return err
			}
			return config.Save(cfg, config.DefaultPath)
		},
	}
}

func groupListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <target>",
		Short: "List all groups in a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			names, err := group.ListGroups(cfg, args[0])
			if err != nil {
				return err
			}
			if len(names) == 0 {
				fmt.Println("(no groups)")
				return nil
			}
			fmt.Println(strings.Join(names, "\n"))
			return nil
		},
	}
}

func groupKeysCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "keys <target> <group>",
		Short: "List keys belonging to a group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			keys, err := group.GetKeys(cfg, args[0], args[1])
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Println("(no keys)")
				return nil
			}
			fmt.Println(strings.Join(keys, "\n"))
			return nil
		},
	}
}
