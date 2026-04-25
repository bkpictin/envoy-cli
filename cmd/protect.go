package cmd

import (
	"fmt"
	"os"

	"envoy-cli/internal/config"
	"envoy-cli/internal/protect"
	"github.com/spf13/cobra"
)

func init() {
	protectCmd := &cobra.Command{
		Use:   "protect",
		Short: "Manage protected (read-only) keys",
	}

	protectAddCmd := &cobra.Command{
		Use:   "add <target> <key>",
		Short: "Mark a key as protected",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := protect.Protect(cfg, args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Key %q in target %q is now protected.\n", args[1], args[0])
			return config.Save(cfg, config.DefaultPath)
		},
	}

	protectRemoveCmd := &cobra.Command{
		Use:   "remove <target> <key>",
		Short: "Remove protection from a key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := protect.Unprotect(cfg, args[0], args[1]); err != nil {
				return err
			}
			fmt.Printf("Key %q in target %q is no longer protected.\n", args[1], args[0])
			return config.Save(cfg, config.DefaultPath)
		},
	}

	protectListCmd := &cobra.Command{
		Use:   "list <target>",
		Short: "List protected keys for a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			keys, err := protect.List(cfg, args[0])
			if err != nil {
				return err
			}
			if len(keys) == 0 {
				fmt.Fprintf(os.Stdout, "No protected keys in target %q.\n", args[0])
				return nil
			}
			for _, k := range keys {
				fmt.Println(k)
			}
			return nil
		},
	}

	protectCmd.AddCommand(protectAddCmd, protectRemoveCmd, protectListCmd)
	rootCmd.AddCommand(protectCmd)
}
