package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	"github.com/goto/shield/pkg/file"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	cli "github.com/spf13/cobra"
)

func RuleCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:     "rule",
		Aliases: []string{"rules"},
		Short:   "Manage rules",
		Long: heredoc.Doc(`
			Work with rules.
		`),
		Example: heredoc.Doc(`
			$ shield rule config upload
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(ruleConfigCommand(cliConfig))

	bindFlagsFromClientConfig(cmd)

	return cmd
}

func ruleConfigCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:   "config",
		Short: "Manage rules config",
		Long: heredoc.Doc(`
			Work with rules config.
		`),
		Example: heredoc.Doc(`
			$ shield rule config upload
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(upsertRuleConfigCommand(cliConfig))

	return cmd
}

func upsertRuleConfigCommand(cliConfig *Config) *cli.Command {
	var name, filePath, header string

	cmd := &cli.Command{
		Use:   "upload",
		Short: "Upload a rule config",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield rule config upload --name --file=<rule-config-body> --header=<key>:<value>
		`),
		Annotations: map[string]string{
			"rule:core": "true",
		},
		RunE: func(cmd *cli.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			var reqBody string
			reqBody, err := file.ReadString(filePath)
			if err != nil {
				return err
			}

			client, cancel, err := createClient(cmd.Context(), cliConfig.Host)
			if err != nil {
				return err
			}
			defer cancel()

			ctx := cmd.Context()
			if header != "" {
				ctx = setCtxHeader(ctx, header)
			}

			res, err := client.UpsertRulesConfig(ctx, &shieldv1beta1.UpsertRulesConfigRequest{
				Name:   name,
				Config: reqBody,
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			fmt.Printf("successfully upserted rule config %s with id %d\nconfig:\n%s", res.GetName(), res.GetId(), res.Config)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Rule config name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the rule body file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&header, "header", "H", "", "Header <key>:<value>")

	return cmd
}
