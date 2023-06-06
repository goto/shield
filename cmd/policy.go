package cmd

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	"github.com/goto/shield/pkg/file"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	cli "github.com/spf13/cobra"
)

func PolicyCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:     "policy",
		Aliases: []string{"policies"},
		Short:   "Manage policies",
		Long: heredoc.Doc(`
			Work with policies.
		`),
		Example: heredoc.Doc(`
			$ shield policy create
			$ shield policy edit
			$ shield policy view
			$ shield policy list
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(createPolicyCommand(cliConfig))
	cmd.AddCommand(listPolicyCommand(cliConfig))

	bindFlagsFromClientConfig(cmd)

	return cmd
}

func createPolicyCommand(cliConfig *Config) *cli.Command {
	var filePath, header string

	cmd := &cli.Command{
		Use:   "create",
		Short: "Create a policy",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield policy create --file=<policy-body> --header=<key>:<value>
		`),
		Annotations: map[string]string{
			"policy:core": "true",
		},
		RunE: func(cmd *cli.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			var reqBody shieldv1beta1.PolicyRequestBody
			if err := file.Parse(filePath, &reqBody); err != nil {
				return err
			}

			err := reqBody.ValidateAll()
			if err != nil {
				return err
			}

			client, cancel, err := createClient(cmd.Context(), cliConfig.Host)
			if err != nil {
				return err
			}
			defer cancel()

			ctx := setCtxHeader(cmd.Context(), header)
			_, err = client.CreatePolicy(ctx, &shieldv1beta1.CreatePolicyRequest{
				Body: &reqBody,
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			fmt.Println("successfully created policy")
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the policy body file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&header, "header", "H", "", "Header <key>:<value>")
	cmd.MarkFlagRequired("header")

	return cmd
}

func listPolicyCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:   "list",
		Short: "List all policies",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield policy list
		`),
		Annotations: map[string]string{
			"policy:core": "true",
		},
		RunE: func(cmd *cli.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd.Context(), cliConfig.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.ListPolicies(cmd.Context(), &shieldv1beta1.ListPoliciesRequest{})
			if err != nil {
				return err
			}

			report := [][]string{}
			policies := res.GetPolicies()

			spinner.Stop()

			if len(policies) == 0 {
				fmt.Printf("No policies found.\n")
				return nil
			}

			fmt.Printf(" \nShowing %d policies\n \n", len(policies))

			report = append(report, []string{"ID", "ACTION", "NAMESPACE"})
			for _, p := range policies {
				report = append(report, []string{
					p.GetId(),
					p.GetActionId(),
					p.GetNamespaceId(),
				})
			}
			printer.Table(os.Stdout, report)

			return nil
		},
	}

	return cmd
}
