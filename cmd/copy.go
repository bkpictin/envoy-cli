package cmd

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy-cli/internal/config"
	"github.com/envoy-cli/envoy-cli/internal/copy"
	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	copyCmd := &cobra.Command{
		Use:   "copy <src> <dst>",
		Short: "Copy environment variables from one target to another",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(config.DefaultPath())
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			if err := copy.CopyEnvs(cfg, args[0], args[1], overwrite); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			if err := config.Save(cfg, config.DefaultPath()); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			fmt.Printf("Copied envs from %q to %q\n", args[0], args[1])
		},
	}

	copyCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in destination")
	rootCmd.AddCommand(copyCmd)
}
