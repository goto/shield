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

func ActionCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:     "action",
		Aliases: []string{"actions"},
		Short:   "Manage actions",
		Long: heredoc.Doc(`
			Work with actions.
		`),
		Example: heredoc.Doc(`
			$ shield action create
			$ shield action edit
			$ shield action view
			$ shield action list
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(createActionCommand(cliConfig))
	cmd.AddCommand(listActionCommand(cliConfig))

	bindFlagsFromClientConfig(cmd)

	return cmd
}

func createActionCommand(cliConfig *Config) *cli.Command {
	var filePath, header string

	cmd := &cli.Command{
		Use:   "create",
		Short: "Create an action",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield action create --file=<action-body> --header=<key>:<value>
		`),
		Annotations: map[string]string{
			"action:core": "true",
		},
		RunE: func(cmd *cli.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			var reqBody shieldv1beta1.ActionRequestBody
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
			res, err := client.CreateAction(ctx, &shieldv1beta1.CreateActionRequest{
				Body: &reqBody,
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			fmt.Printf("successfully created action %s with id %s\n", res.GetAction().GetName(), res.GetAction().GetId())
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the action body file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&header, "header", "H", "", "Header <key>:<value>")
	cmd.MarkFlagRequired("header")

	return cmd
}

func listActionCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:   "list",
		Short: "List all actions",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield action list
		`),
		Annotations: map[string]string{
			"action:core": "true",
		},
		RunE: func(cmd *cli.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd.Context(), cliConfig.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.ListActions(cmd.Context(), &shieldv1beta1.ListActionsRequest{})
			if err != nil {
				return err
			}

			report := [][]string{}
			actions := res.GetActions()

			spinner.Stop()

			if len(actions) == 0 {
				fmt.Printf("No actions found.\n")
				return nil
			}

			fmt.Printf(" \nShowing %d action(s)\n \n", len(actions))

			report = append(report, []string{"ID", "NAME", "NAMESPACE"})
			for _, a := range actions {
				report = append(report, []string{
					a.GetId(),
					a.GetName(),
					a.GetNamespace().GetId(),
				})
			}
			printer.Table(os.Stdout, report)

			return nil
		},
	}

	return cmd
}
