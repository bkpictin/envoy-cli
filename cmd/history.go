package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/envoy-cli/internal/config"
	"github.com/your-org/envoy-cli/internal/history"
)

func init() {
	var target string
	var key string

	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Show audit history for a target or specific key",
		Example: `  envoy history --target production
  envoy history --target production --key API_URL`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if target == "" {
				return fmt.Errorf("--target is required")
			}

			cfgPath, _ := cmd.Flags().GetString("config")
			if cfgPath == "" {
				cfgPath = config.DefaultPath()
			}

			cfg, err := config.Load(cfgPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			var entries []config.AuditEntry
			if key != "" {
				entries = history.ForKey(cfg, target, key)
				if len(entries) == 0 {
					fmt.Fprintf(os.Stderr, "no history found for key %q in target %q\n", key, target)
					return nil
				}
			} else {
				entries = history.ForTarget(cfg, target)
				if len(entries) == 0 {
					fmt.Fprintf(os.Stderr, "no history found for target %q\n", target)
					return nil
				}
			}

			fmt.Print(history.Format(entries))
			return nil
		},
	}

	historyCmd.Flags().StringVarP(&target, "target", "t", "", "target name (required)")
	historyCmd.Flags().StringVarP(&key, "key", "k", "", "filter history to a specific key")

	rootCmd.AddCommand(historyCmd)
}
