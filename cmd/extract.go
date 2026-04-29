package cmd

import (
	"fmt"
	"os"
	"strings"

	"envoy-cli/internal/config"
	"envoy-cli/internal/extract"
	"github.com/spf13/cobra"
)

func init() {
	var (
		keys      []string
		dest      string
		overwrite bool
		strict    bool
	)

	cmd := &cobra.Command{
		Use:   "extract <target> [--key KEY]... [--into DEST]",
		Short: "Extract specific keys from a target",
		Long: `Extract one or more keys from a target and display them or write
them into another target.

If --into is omitted the extracted keys are printed to stdout.
Use --strict to fail if any requested key is absent.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]

			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return fmt.Errorf("load config: %w", err)
			}

			// Flatten comma-separated keys passed as a single flag value.
			var flatKeys []string
			for _, k := range keys {
				for _, part := range strings.Split(k, ",") {
					if p := strings.TrimSpace(part); p != "" {
						flatKeys = append(flatKeys, p)
					}
				}
			}

			result, err := extract.FromTarget(cfg, src, flatKeys, strict)
			if err != nil {
				return err
			}

			if dest == "" {
				fmt.Print(extract.Format(result))
				return nil
			}

			if err := extract.IntoTarget(cfg, dest, result, overwrite); err != nil {
				return err
			}

			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return fmt.Errorf("save config: %w", err)
			}

			fmt.Fprintf(os.Stdout, "extracted %d key(s) from %q into %q\n", len(result.Keys), src, dest)
			return nil
		},
	}

	cmd.Flags().StringArrayVarP(&keys, "key", "k", nil, "key(s) to extract (repeatable, comma-separated)")
	cmd.Flags().StringVar(&dest, "into", "", "destination target to write extracted keys into")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite existing keys in destination")
	cmd.Flags().BoolVar(&strict, "strict", false, "fail if a requested key is missing")

	rootCmd.AddCommand(cmd)
}
