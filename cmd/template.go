package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envoy-cli/internal/config"
	"envoy-cli/internal/template"
)

func init() {
	var target string
	var strict bool
	var listVars bool

	cmd := &cobra.Command{
		Use:   "template <target> <template-string>",
		Short: "Render a template string using variables from a target",
		Example: `  envoy template production "https://${APP_HOST}:${APP_PORT}/api"
  envoy template staging "host=$APP_HOST" --strict
  envoy template production "${DB_URL}" --list-vars`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			targetName := args[0]
			tmpl := args[1]

			if listVars {
				vars := template.ListVars(tmpl)
				if len(vars) == 0 {
					fmt.Println("No variables found in template.")
					return nil
				}
				fmt.Println("Variables referenced:")
				for _, v := range vars {
					fmt.Printf("  %s\n", v)
				}
				return nil
			}

			cfg, err := config.Load(config.DefaultPath)
			if err != nil {
				return err
			}

			if target != "" {
				targetName = target
			}

			res, err := template.Render(cfg, targetName, tmpl, strict)
			if err != nil {
				return err
			}

			fmt.Println(res.Output)

			if len(res.Missing) > 0 {
				fmt.Fprintf(os.Stderr, "warning: unresolved variables: %s\n",
					strings.Join(res.Missing, ", "))
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "Override target name (alias for first positional arg)")
	cmd.Flags().BoolVar(&strict, "strict", false, "Fail if any variable is missing")
	cmd.Flags().BoolVar(&listVars, "list-vars", false, "List all variable references in the template without rendering")

	rootCmd.AddCommand(cmd)
}
