package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	"github.com/goto/shield/pkg/file"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	cli "github.com/spf13/cobra"
)

func RoleCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:     "role",
		Aliases: []string{"roles"},
		Short:   "Manage roles",
		Long: heredoc.Doc(`
			Work with roles.
		`),
		Example: heredoc.Doc(`
			$ shield role create
			$ shield role edit
			$ shield role view
			$ shield role list
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(createRoleCommand(cliConfig))
	cmd.AddCommand(listRoleCommand(cliConfig))

	bindFlagsFromClientConfig(cmd)

	return cmd
}

func createRoleCommand(cliConfig *Config) *cli.Command {
	var filePath, header string

	cmd := &cli.Command{
		Use:   "create",
		Short: "Create a role",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield role create --file=<role-body> --header=<key>:<value>
		`),
		Annotations: map[string]string{
			"role:core": "true",
		},
		RunE: func(cmd *cli.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			var reqBody shieldv1beta1.RoleRequestBody
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

			res, err := client.CreateRole(ctx, &shieldv1beta1.CreateRoleRequest{
				Body: &reqBody,
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			fmt.Printf("successfully created role %s with id %s\n", res.GetRole().GetName(), res.GetRole().GetId())
			return nil
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the role body file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&header, "header", "H", "", "Header <key>:<value>")
	cmd.MarkFlagRequired("header")

	return cmd
}

func listRoleCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:   "list",
		Short: "List all roles",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield role list
		`),
		Annotations: map[string]string{
			"role:core": "true",
		},
		RunE: func(cmd *cli.Command, args []string) error {
			spinner := printer.Spin("")
			defer spinner.Stop()

			client, cancel, err := createClient(cmd.Context(), cliConfig.Host)
			if err != nil {
				return err
			}
			defer cancel()

			res, err := client.ListRoles(cmd.Context(), &shieldv1beta1.ListRolesRequest{})
			if err != nil {
				return err
			}

			report := [][]string{}
			roles := res.GetRoles()

			spinner.Stop()

			if len(roles) == 0 {
				fmt.Printf("No roles found.\n")
				return nil
			}

			fmt.Printf(" \nShowing %d roles\n \n", len(roles))

			report = append(report, []string{"ID", "NAME", "TYPE(S)", "NAMESPACE"})
			for _, r := range roles {
				report = append(report, []string{
					r.GetId(),
					r.GetName(),
					strings.Join(r.GetTypes(), ", "),
					r.GetNamespace().GetId(),
				})
			}
			printer.Table(os.Stdout, report)

			return nil
		},
	}

	return cmd
}
