package cmd

import (
	"fmt"
	"os"
	"sort"

	"envoy-cli/internal/compare"
	"envoy-cli/internal/config"

	"github.com/spf13/cobra"
)

func init() {
	var snapshotName string

	cmd := &cobra.Command{
		Use:   "compare <targetA> <targetB>",
		Short: "Compare env vars between two targets or a snapshot and a target",
		Example: `  envoy compare staging production
  envoy compare --snapshot snap-1 staging`,
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			var result compare.Result
			var labelA, labelB string

			if snapshotName != "" {
				if len(args) != 1 {
					return fmt.Errorf("provide exactly one target when using --snapshot")
				}
				labelA = "snapshot:" + snapshotName
				labelB = args[0]
				result, err = compare.SnapshotVsTarget(cfg, snapshotName, args[0])
			} else {
				if len(args) != 2 {
					return fmt.Errorf("provide two targets to compare")
				}
				labelA = args[0]
				labelB = args[1]
				result, err = compare.Targets(cfg, args[0], args[1])
			}
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "Comparing %s → %s  (%s)\n\n", labelA, labelB, compare.Summary(result))

			printSection := func(title string, keys map[string]string) {
				if len(keys) == 0 {
					return
				}
				fmt.Fprintf(os.Stdout, "%s:\n", title)
				sorted := make([]string, 0, len(keys))
				for k := range keys {
					sorted = append(sorted, k)
				}
				sort.Strings(sorted)
				for _, k := range sorted {
					fmt.Fprintf(os.Stdout, "  %s=%s\n", k, keys[k])
				}
				fmt.Fprintln(os.Stdout)
			}

			printSection(fmt.Sprintf("Only in %s", labelA), result.OnlyInA)
			printSection(fmt.Sprintf("Only in %s", labelB), result.OnlyInB)

			if len(result.Different) > 0 {
				fmt.Fprintln(os.Stdout, "Changed:")
				keys := make([]string, 0, len(result.Different))
				for k := range result.Different {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					p := result.Different[k]
					fmt.Fprintf(os.Stdout, "  %s: %s → %s\n", k, p.A, p.B)
				}
				fmt.Fprintln(os.Stdout)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&snapshotName, "snapshot", "", "compare a snapshot against a target")
	rootCmd.AddCommand(cmd)
}
