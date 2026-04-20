package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/envoy-cli/internal/config"
	"github.com/user/envoy-cli/internal/encrypt"
)

func init() {
	var passphrase string

	encryptCmd := &cobra.Command{
		Use:   "encrypt <target>",
		Short: "Encrypt all plain-text values in a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := encrypt.EncryptTarget(cfg, args[0], passphrase); err != nil {
				return err
			}
			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "target %q encrypted\n", args[0])
			return nil
		},
	}
	encryptCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "encryption passphrase (required)")
	_ = encryptCmd.MarkFlagRequired("passphrase")

	decryptCmd := &cobra.Command{
		Use:   "decrypt <target>",
		Short: "Decrypt all encrypted values in a target",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}
			if err := encrypt.DecryptTarget(cfg, args[0], passphrase); err != nil {
				return err
			}
			if err := config.Save(cfg, config.DefaultPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "target %q decrypted\n", args[0])
			return nil
		},
	}
	decryptCmd.Flags().StringVarP(&passphrase, "passphrase", "p", "", "decryption passphrase (required)")
	_ = decryptCmd.MarkFlagRequired("passphrase")

	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}
