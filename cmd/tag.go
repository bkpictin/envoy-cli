package cmd

import (
	"fmt"
	"strings"

	"envoy-cli/internal/config"
	"envoy-cli/internal/tag"

	"github.com/spf13/cobra"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags on environment variable keys",
	}

	tagAddCmd := &cobra.Command{
		Use:   "add <target> <key> <tag>",
		Short: "Add a tag to a key in a target",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := tag.Add(cfg, args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("Tagged %q in %q with %q\n", args[1], args[0], args[2])
			return config.Save(cfg, config.DefaultPath)
		},
	}

	tagRemoveCmd := &cobra.Command{
		Use:   "remove <target> <key> <tag>",
		Short: "Remove a tag from a key in a target",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := tag.Remove(cfg, args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Printf("Removed tag %q from %q in %q\n", args[2], args[1], args[0])
			return config.Save(cfg, config.DefaultPath)
		},
	}

	tagListCmd := &cobra.Command{
		Use:   "list <target> [--tag <tag> | --key <key>]",
		Short: "List tags or tagged keys in a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			byTag, _ := cmd.Flags().GetString("tag")
			byKey, _ := cmd.Flags().GetString("key")
			switch {
			case byTag != "":
				keys, err := tag.ListByTag(cfg, args[0], byTag)
				if err != nil {
					return err
				}
				fmt.Printf("Keys tagged %q in %q:\n  %s\n", byTag, args[0], strings.Join(keys, ", "))
			case byKey != "":
				tags, err := tag.ListForKey(cfg, args[0], byKey)
				if err != nil {
					return err
				}
				fmt.Printf("Tags on %q in %q:\n  %s\n", byKey, args[0], strings.Join(tags, ", "))
			default:
				return fmt.Errorf("provide --tag or --key flag")
			}
			return nil
		},
	}
	tagListCmd.Flags().String("tag", "", "Filter keys by this tag")
	tagListCmd.Flags().String("key", "", "List tags for this key")

	tagCmd.AddCommand(tagAddCmd, tagRemoveCmd, tagListCmd)
	rootCmd.AddCommand(tagCmd)
}
