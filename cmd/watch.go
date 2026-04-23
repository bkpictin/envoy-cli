package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/envoy-cli/envoy/internal/config"
	"github.com/envoy-cli/envoy/internal/watch"
	"github.com/spf13/cobra"
)

func init() {
	var interval int

	watchCmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch the config file and print a notice on each change",
		Long: `Polls the envoy config file and prints a summary whenever it changes.

Press Ctrl-C to stop watching.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			path := config.DefaultPath()

			fmt.Fprintf(cmd.OutOrStdout(), "Watching %s (interval %ds) — press Ctrl-C to stop\n", path, interval)

			done := make(chan struct{})

			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sig
				close(done)
			}()

			opts := watch.Options{Interval: time.Duration(interval) * time.Millisecond}
			ch := watch.Watch(path, opts, done)

			for ev := range ch {
				if ev.Err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "error: %v\n", ev.Err)
					continue
				}
				targets := ev.Cfg.Targets
				keys := 0
				for _, vars := range targets {
					keys += len(vars)
				}
				fmt.Fprintf(cmd.OutOrStdout(),
					"[%s] config reloaded — %d target(s), %d key(s) total\n",
					time.Now().Format("15:04:05"), len(targets), keys,
				)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Stopped.")
			return nil
		},
	}

	watchCmd.Flags().IntVarP(&interval, "interval", "i", 500, "Poll interval in milliseconds")
	rootCmd.AddCommand(watchCmd)
}
