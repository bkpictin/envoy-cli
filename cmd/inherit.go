package cmd

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
	"envoy-cli/internal/inherit"

	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	inheritCmd := &cobra.Command{
		Use:   "inherit <base-target> <child-target> [<child-target>...]",
		Short: "Inherit environment variables from a base target into child targets",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			baseTarget := args[0]
			children := args[1:]

			cfgPath := config.DefaultPath()
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			results, err := inherit.Apply(cfg, baseTarget, children, overwrite)
			if err != nil {
				return err
			}

			if err := config.Save(cfgPath, cfg); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			for _, r := range results {
				parts := []string{}
				if r.Added > 0 {
					parts = append(parts, fmt.Sprintf("%d added", r.Added))
				}
				if r.Skipped > 0 {
					parts = append(parts, fmt.Sprintf("%d skipped", r.Skipped))
				}
				summary := "no changes"
				if len(parts) > 0 {
					summary = strings.Join(parts, ", ")
				}
				fmt.Printf("%-20s %s\n", r.Target, summary)
			}
			return nil
		},
	}

	inheritCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing keys in child targets")
	rootCmd.AddCommand(inheritCmd)
}
