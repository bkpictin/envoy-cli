package cmd

import (
	"fmt"
	"os"

	"envoy-cli/internal/config"
	"envoy-cli/internal/export"
	"github.com/spf13/cobra"
)

func init() {
	var format string
	var output string

	cmd := &cobra.Command{
		Use:   "export <target>",
		Short: "Export environment variables for a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			result, err := export.ToFile(cfg, args[0], export.Format(format))
			if err != nil {
				return err
			}

			if output != "" {
				if err := os.WriteFile(output, []byte(result), 0644); err != nil {
					return fmt.Errorf("writing file: %w", err)
				}
				fmt.Printf("Exported %q to %s\n", args[0], output)
				return nil
			}

			fmt.Print(result)
			return nil
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Output format: dotenv, shell, json")
	cmd.Flags().StringVarP(&output, "out", "o", "", "Write output to file instead of stdout")

	rootCmd.AddCommand(cmd)
}
