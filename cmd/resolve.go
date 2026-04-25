package cmd

import (
	"fmt"
	"sort"

	"envoy-cli/internal/config"
	"envoy-cli/internal/resolve"
	"github.com/spf13/cobra"
)

func init() {
	var target string

	resolveCmd := &cobra.Command{
		Use:   "resolve",
		Short: "Expand ${VAR} references within a target's env values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			results, err := resolve.Target(cfg, target)
			if err != nil {
				return err
			}
			sort.Slice(results, func(i, j int) bool {
				return results[i].Key < results[j].Key
			})
			for _, r := range results {
				if r.Original == r.Resolved {
					fmt.Printf("%s=%s\n", r.Key, r.Resolved)
				} else {
					fmt.Printf("%s=%s  (was: %s)\n", r.Key, r.Resolved, r.Original)
				}
				for _, w := range r.Warnings {
					fmt.Printf("  warning: %s\n", w)
				}
			}
			return nil
		},
	}

	resolveCmd.Flags().StringVarP(&target, "target", "t", "", "Target to resolve (required)")
	_ = resolveCmd.MarkFlagRequired("target")
	rootCmd.AddCommand(resolveCmd)
}
