package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/config"
	"github.com/user/envoy-cli/internal/snapshot"
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage environment variable snapshots",
	}

	snapshotCmd.AddCommand(&cobra.Command{
		Use:   "create <target>",
		Short: "Capture current env vars for a target",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath := config.DefaultPath()
			cfg, err := config.Load(cfgPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			snap, err := snapshot.Create(cfg, args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			config.Save(cfgPath, cfg)
			fmt.Printf("snapshot created for %q at %s\n", snap.Target, snap.Timestamp.Format("2006-01-02 15:04:05"))
		},
	})

	snapshotCmd.AddCommand(&cobra.Command{
		Use:   "list <target>",
		Short: "List snapshots for a target",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load(config.DefaultPath())
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			snaps, err := snapshot.List(cfg, args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			if len(snaps) == 0 {
				fmt.Println("no snapshots found")
				return
			}
			for i, s := range snaps {
				fmt.Printf("[%d] %s (%d vars)\n", i, s.Timestamp.Format("2006-01-02 15:04:05"), len(s.Vars))
			}
		},
	})

	snapshotCmd.AddCommand(&cobra.Command{
		Use:   "restore <target> <index>",
		Short: "Restore a target's env vars from a snapshot",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath := config.DefaultPath()
			cfg, err := config.Load(cfgPath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			var idx int
			fmt.Sscanf(args[1], "%d", &idx)
			if err := snapshot.Restore(cfg, args[0], idx); err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			config.Save(cfgPath, cfg)
			fmt.Printf("restored snapshot %d for %q\n", idx, args[0])
		},
	})

	rootCmd.AddCommand(snapshotCmd)
}
