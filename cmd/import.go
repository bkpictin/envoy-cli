package cmd

import (
	fmt2 "fmt"

	"github.com/spf13/cobra"

	"envoy-cli/internal/audit"
	"envoy-cli/internal/config"
	importenv "envoy-cli/internal/import"
)

func init() {
	var format string
	var overwrite bool

	importCmd := &cobra.Command{
		Use:   "import <target> <file>",
		Short: "Import environment variables from a file into a target",
		Long: `Import environment variables from a dotenv, shell-export, or JSON file
into the specified target. Existing keys are skipped unless --overwrite is set.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			path := args[1]

			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			n, err := importenv.FromFile(cfg, target, path, importenv.Format(format), overwrite)
			if err != nil {
				return err
			}

			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}

			_ = audit.Log(cfg, "import", target, map[string]string{
				"file":      path,
				"format":    format,
				"imported":  fmt2.Sprintf("%d", n),
				"overwrite": fmt2.Sprintf("%v", overwrite),
			})
			_ = config.Save(cfg, config.DefaultPath)

			cmd.Printf("Imported %d key(s) into target %q from %s\n", n, target, path)
			return nil
		},
	}

	importCmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Input format: dotenv, shell, json")
	importCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys")

	rootCmd.AddCommand(importCmd)
}
