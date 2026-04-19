package cmd

import (
	"fmt"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new envoy config file in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := config.DefaultPath()
		if err := config.Init(path); err != nil {
			return err
		}
		fmt.Printf("Initialized envoy config at %s\n", path)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
