package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/transform"
)

func init() {
	var dryRun bool
	var keys []string

	transformCmd := &cobra.Command{
		Use:   "transform <target> <kind>",
		Short: "Apply a value transformation to env vars in a target",
		Long: `Apply a transformation to environment variable values in a target.

Supported kinds: uppercase, lowercase, base64encode, base64decode, trimspace

Examples:
  envoy transform dev uppercase
  envoy transform staging trimspace --keys DB_HOST,API_URL
  envoy transform prod base64encode --dry-run`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			kind := transform.Kind(args[1])

			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			results, err := transform.Target(cfg, target, kind, keys, dryRun)
			if err != nil {
				return err
			}

			changed := 0
			for _, r := range results {
				if r.Changed {
					changed++
					label := ""
					if dryRun {
						label = " (dry-run)"
					}
					fmt.Fprintf(os.Stdout, "  %s: %q → %q%s\n", r.Key, r.Before, r.After, label)
				}
			}

			if changed == 0 {
				fmt.Println("no values changed")
				return nil
			}

			if dryRun {
				fmt.Printf("\n%d key(s) would be changed (dry-run)\n", changed)
				return nil
			}

			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}
			fmt.Printf("\n%d key(s) transformed in %q\n", changed, target)
			return nil
		},
	}

	transformCmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without applying them")
	transformCmd.Flags().StringSliceVar(&keys, "keys", nil, "comma-separated list of keys to transform (default: all)")
	_ = strings.Join // suppress unused import in minimal build

	rootCmd.AddCommand(transformCmd)
}
