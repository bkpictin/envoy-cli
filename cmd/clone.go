package cmd

import (
	"fmt"
	"strings"

	"envoy-cli/internal/audit"
	"envoy-cli/internal/clone"
	"envoy-cli/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool
	var prefix string

	cloneCmd := &cobra.Command{
		Use:   "clone <source> <dest>",
		Short: "Clone a target into a new target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dest := args[0], args[1]

			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			if prefix != "" {
				err = clone.WithFilter(cfg, src, dest, overwrite, func(key string) bool {
					return strings.HasPrefix(key, prefix)
				})
			} else {
				err = clone.Target(cfg, src, dest, overwrite)
			}
			if err != nil {
				return err
			}

			audit.Log(cfg, "clone", fmt.Sprintf("%s -> %s", src, dest))

			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}

			fmt.Printf("Cloned target %q to %q\n", src, dest)
			return nil
		},
	}

	cloneCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite destination if it already exists")
	cloneCmd.Flags().StringVar(&prefix, "prefix", "", "Only clone keys that start with this prefix")

	rootCmd.AddCommand(cloneCmd)
}
