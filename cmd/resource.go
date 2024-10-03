package cmd

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/printer"
	"github.com/goto/shield/pkg/file"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	cli "github.com/spf13/cobra"
)

func ResourceCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:     "resource",
		Aliases: []string{"resources"},
		Short:   "Manage resources",
		Long: heredoc.Doc(`
			Work with resources.
		`),
		Example: heredoc.Doc(`
			$ shield resource config upload
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(resourceConfigCommand(cliConfig))

	bindFlagsFromClientConfig(cmd)

	return cmd
}

func resourceConfigCommand(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:   "config",
		Short: "Manage resources config",
		Long: heredoc.Doc(`
			Work with resources config.
		`),
		Example: heredoc.Doc(`
			$ shield resource config upload
		`),
		Annotations: map[string]string{
			"group":  "core",
			"client": "true",
		},
	}

	cmd.AddCommand(upsertResourcesConfigCommand(cliConfig))

	return cmd
}

func upsertResourcesConfigCommand(cliConfig *Config) *cli.Command {
	var name, filePath, header string

	cmd := &cli.Command{
		Use:   "upload",
		Short: "Upload a resource config",
		Args:  cli.NoArgs,
		Example: heredoc.Doc(`
			$ shield resource config upload --name --file=<resource-config-body> --header=<key>:<value>
		`),
		Annotations: map[string]string{
			"resource:core": "true",
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

			res, err := client.UpsertResourcesConfig(ctx, &shieldv1beta1.UpsertResourcesConfigRequest{
				Name:   name,
				Config: reqBody,
			})
			if err != nil {
				return err
			}

			spinner.Stop()
			fmt.Printf("successfully upserted resource config %s with id %d\nconfig:\n%s", res.GetName(), res.GetId(), res.GetConfig())
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Resource config name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the resource body file")
	cmd.MarkFlagRequired("file")
	cmd.Flags().StringVarP(&header, "header", "H", "", "Header <key>:<value>")

	return cmd
}
