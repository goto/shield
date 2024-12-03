package cmd

import (
	"github.com/raystack/salt/log"
	"github.com/raystack/shield/config"
	"github.com/raystack/shield/pkg/sql"
	"github.com/raystack/shield/store/postgres/migrations"
	cli "github.com/spf13/cobra"
)

func migrationsCommand(logger log.Logger, appConfig *config.Shield) *cli.Command {
	c := &cli.Command{
		Use:     "migrate",
		Short:   "Run DB Schema Migrations",
		Example: "shield migrate",
		RunE: func(c *cli.Command, args []string) error {
			return sql.RunMigrations(sql.Config{
				Driver: appConfig.DB.Driver,
				URL:    appConfig.DB.URL,
			}, migrations.MigrationFs, migrations.ResourcePath)
		},
	}
	return c
}

func migrationsRollbackCommand(logger log.Logger, appConfig *config.Shield) *cli.Command {
	c := &cli.Command{
		Use:     "migration-rollback",
		Short:   "Run DB Schema Migrations Rollback to last state",
		Example: "shield migration-rollback",
		RunE: func(c *cli.Command, args []string) error {
			return sql.RunRollback(sql.Config{
				Driver: appConfig.DB.Driver,
				URL:    appConfig.DB.URL,
			}, migrations.MigrationFs, migrations.ResourcePath)
		},
	}
	return c
}
