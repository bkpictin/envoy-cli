package cmd

import (
	"fmt"

	"github.com/envoy-cli/envoy-cli/internal/config"
	"github.com/envoy-cli/envoy-cli/internal/env"
	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables for deployment targets",
}

var envSetCmd = &cobra.Command{
	Use:   "set <target> <key> <value>",
	Short: "Set an environment variable for a target",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(config.DefaultPath())
		if err != nil {
			return err
		}
		if err := env.Set(cfg, args[0], args[1], args[2]); err != nil {
			return err
		}
		return config.Save(cfg, config.DefaultPath())
	},
}

var envGetCmd = &cobra.Command{
	Use:   "get <target> <key>",
	Short: "Get an environment variable from a target",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(config.DefaultPath())
		if err != nil {
			return err
		}
		val, err := env.Get(cfg, args[0], args[1])
		if err != nil {
			return err
		}
		fmt.Println(val)
		return nil
	},
}

var envDeleteCmd = &cobra.Command{
	Use:   "delete <target> <key>",
	Short: "Delete an environment variable from a target",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(config.DefaultPath())
		if err != nil {
			return err
		}
		if err := env.Delete(cfg, args[0], args[1]); err != nil {
			return err
		}
		return config.Save(cfg, config.DefaultPath())
	},
}

func init() {
	envCmd.AddCommand(envSetCmd, envGetCmd, envDeleteCmd)
	rootCmd.AddCommand(envCmd)
}
