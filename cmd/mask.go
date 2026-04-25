package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"envoy-cli/internal/config"
	"envoy-cli/internal/mask"

	"github.com/spf13/cobra"
)

func init() {
	var target string
	var visibleLen int

	maskCmd := &cobra.Command{
		Use:   "mask",
		Short: "Display environment variables with masked values",
		Long:  "Print all keys for a target (or all targets) with secret values partially or fully hidden.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			var results []mask.Result
			if target != "" {
				results, err = mask.Target(cfg, target, visibleLen)
				if err != nil {
					return err
				}
			} else {
				results = mask.All(cfg, visibleLen)
			}

			sort.Slice(results, func(i, j int) bool {
				if results[i].Target != results[j].Target {
					return results[i].Target < results[j].Target
				}
				return results[i].Key < results[j].Key
			})

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TARGET\tKEY\tVALUE")
			for _, r := range results {
				fmt.Fprintf(w, "%s\t%s\t%s\n", r.Target, r.Key, r.Masked)
			}
			return w.Flush()
		},
	}

	maskCmd.Flags().StringVarP(&target, "target", "t", "", "Target to mask (default: all targets)")
	maskCmd.Flags().IntVarP(&visibleLen, "visible", "n", 4, "Number of trailing characters to leave visible")
	rootCmd.AddCommand(maskCmd)
}
