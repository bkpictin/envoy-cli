package cmd

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
	"envoy-cli/internal/pin"
	"github.com/spf13/cobra"
)

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin",
		Short: "Manage pinned keys that are protected from bulk overwrites",
	}

	pinSetCmd := &cobra.Command{
		Use:   "set <target> <key>",
		Short: "Pin a key in a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := pin.Pin(cfg, args[0], args[1]); err != nil {
				return err
			}
			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}
			fmt.Printf("Pinned %q in target %q\n", args[1], args[0])
			return nil
		},
	}

	pinUnsetCmd := &cobra.Command{
		Use:   "unset <target> <key>",
		Short: "Unpin a key in a target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := pin.Unpin(cfg, args[0], args[1]); err != nil {
				return err
			}
			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}
			fmt.Printf("Unpinned %q in target %q\n", args[1], args[0])
			return nil
		},
	}

	pinListCmd := &cobra.Command{
		Use:   "list <target>",
		Short: "List all pinned keys in a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			keys, err := pin.List(cfg, args[0])
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Printf("No pinned keys in target %q\n", args[0])
				return nil
			}
			fmt.Printf("Pinned keys in %q:\n  %s\n", args[0], strings.Join(keys, "\n  "))
			return nil
		},
	}

	pinCmd.AddCommand(pinSetCmd, pinUnsetCmd, pinListCmd)
	rootCmd.AddCommand(pinCmd)
}
