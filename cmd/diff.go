package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/diff"
)

func init() {
	diffCmd := &cobra.Command{
		Use:   "diff <targetA> <targetB>",
		Short: "Show differences in environment variables between two targets",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, _ := cmd.Flags().GetString("config")
			if cfgPath == "" {
				cfgPath = config.DefaultPath()
			}
			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			result, err := diff.Targets(cfg, args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Fprint(os.Stdout, diff.Format(result, args[0], args[1]))
			return nil
		},
	}

	rootCmd.AddCommand(diffCmd)
}
