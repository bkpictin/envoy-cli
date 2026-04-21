package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/schema"
)

func init() {
	var target string
	var required []string
	var optional []string
	var warnOnly bool

	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Validate a target's env vars against a declared schema",
		Long: `Check that a target contains all required keys and flag any
undeclared keys as warnings. Exits non-zero if errors are found
unless --warn-only is set.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			s := make(schema.Schema)
			for _, k := range required {
				s[k] = schema.Rule{Required: true}
			}
			for _, k := range optional {
				s[k] = schema.Rule{Required: false}
			}

			if len(s) == 0 {
				return fmt.Errorf("provide at least one --required or --optional key")
			}

			results, err := schema.Validate(cfg, target, s)
			if err != nil {
				return err
			}

			fmt.Println(schema.Format(results))

			if !warnOnly {
				for _, r := range results {
					if r.Level == "error" {
						os.Exit(1)
					}
				}
			}
			return nil
		},
	}

	schemaCmd.Flags().StringVarP(&target, "target", "t", "", "Target to validate (required)")
	schemaCmd.Flags().StringArrayVar(&required, "required", nil, "Required key names (repeatable)")
	schemaCmd.Flags().StringArrayVar(&optional, "optional", nil, "Optional/declared key names (repeatable)")
	schemaCmd.Flags().BoolVar(&warnOnly, "warn-only", false, "Exit 0 even when errors are found")
	_ = schemaCmd.MarkFlagRequired("target")

	rootCmd.AddCommand(schemaCmd)
}
