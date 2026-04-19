package cmd

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/validate"
	"github.com/spf13/cobra"
)

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate <target> <key>",
		Short: "Validate a target name and environment variable key",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			target, key := args[0], args[1]
			if err := validate.All(cfg, target, key); err != nil {
				return err
			}
			fmt.Printf("✓ target %q and key %q are valid\n", target, key)
			return nil
		},
	}

	keyCmd := &cobra.Command{
		Use:   "key <key>",
		Short: "Validate only the format of an environment variable key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate.KeyFormat(args[0]); err != nil {
				return err
			}
			fmt.Printf("✓ key %q is valid\n", args[0])
			return nil
		},
	}

	validateCmd.AddCommand(keyCmd)
	rootCmd.AddCommand(validateCmd)
}
