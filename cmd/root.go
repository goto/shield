package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/goto/salt/cmdx"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

func New(cliConfig *Config) *cli.Command {
	cmd := &cli.Command{
		Use:   "shield <command> <subcommand> [flags]",
		Short: "A cloud native role-based authorization aware reverse-proxy service",
		Long: heredoc.Doc(`
			A cloud native role-based authorization aware reverse-proxy service.`),
		SilenceUsage:  true,
		SilenceErrors: true,
		Annotations: map[string]string{
			"group": "core",
			"help:learn": heredoc.Doc(`
				Use 'shield <command> <subcommand> --help' for info about a command.
				Read the manual at https://goto.github.io/shield/
			`),
			"help:feedback": heredoc.Doc(`
				Open an issue here https://github.com/goto/shield/issues
			`),
			"help:environment": heredoc.Doc(`
				See 'shield help environment' for the list of supported environment variables.
			`),
		},
	}

	cmd.PersistentPreRunE = func(subCmd *cobra.Command, args []string) error {
		if isClientCLI(subCmd) {
			if err := overrideClientConfigHost(subCmd, cliConfig); err != nil {
				return err
			}
		}
		return nil
	}

	cmd.AddCommand(ServerCommand())
	cmd.AddCommand(NamespaceCommand(cliConfig))
	cmd.AddCommand(UserCommand(cliConfig))
	cmd.AddCommand(OrganizationCommand(cliConfig))
	cmd.AddCommand(GroupCommand(cliConfig))
	cmd.AddCommand(ProjectCommand(cliConfig))
	cmd.AddCommand(RoleCommand(cliConfig))
	cmd.AddCommand(ActionCommand(cliConfig))
	cmd.AddCommand(PolicyCommand(cliConfig))
	cmd.AddCommand(configCommand())
	cmd.AddCommand(RuleCommand(cliConfig))
	cmd.AddCommand(ResourceCommand(cliConfig))

	// Help topics
	cmdx.SetHelp(cmd)
	cmd.AddCommand(cmdx.SetCompletionCmd("shield"))
	cmd.AddCommand(cmdx.SetHelpTopicCmd("environment", envHelp))
	cmd.AddCommand(cmdx.SetHelpTopicCmd("auth", authHelp))
	cmd.AddCommand(cmdx.SetRefCmd(cmd))
	return cmd
}
