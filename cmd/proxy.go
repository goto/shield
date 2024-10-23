package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/goto/shield/config"
	"github.com/goto/shield/internal/proxy/envoy/xds"
	"github.com/goto/shield/internal/store/postgres"
	shieldlogger "github.com/goto/shield/pkg/logger"
	"github.com/spf13/cobra"
	cli "github.com/spf13/cobra"
)

func ProxyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy <command>",
		Short: "Proxy management",
		Long:  "Server management commands.",
		Example: heredoc.Doc(`
			$ shield proxy envoy-xds start -c ./config.yaml
		`),
	}

	cmd.AddCommand(proxyEnvoyXDSCommand())

	return cmd
}

func proxyEnvoyXDSCommand() *cobra.Command {
	c := &cli.Command{
		Use:   "envoy-xds",
		Short: "Envoy Agent xDS management",
		Long:  "Envoy Agent xDS management commands.",
		Example: heredoc.Doc(`
			$ shield proxy envoy-xds start
		`),
	}

	c.AddCommand(envoyXDSStartCommand())

	return c
}

func envoyXDSStartCommand() *cobra.Command {
	var configFile string

	c := &cli.Command{
		Use:     "start",
		Short:   "Start Envoy Agent xDS server",
		Long:    "Start Envoy Agent xDS server commands.",
		Example: "shield proxy envoy-xds start",
		RunE: func(cmd *cli.Command, args []string) error {
			appConfig, err := config.Load(configFile)
			if err != nil {
				panic(err)
			}

			logger := shieldlogger.InitLogger(shieldlogger.Config{Level: appConfig.Log.Level})

			dbClient, err := setupDB(appConfig.DB, logger)
			if err != nil {
				return err
			}
			defer func() {
				logger.Info("cleaning up db")
				dbClient.Close()
			}()

			ctx := cmd.Context()

			pgRuleRepository := postgres.NewRuleRepository(dbClient)
			if err := pgRuleRepository.InitCache(ctx); err != nil {
				return err
			}

			cbs, repositories, err := buildXDSDependencies(ctx, logger, appConfig.Proxy, pgRuleRepository)
			if err != nil {
				return err
			}
			defer func() {
				logger.Info("cleaning up rules proxy blob")
				for _, f := range cbs {
					if err := f(); err != nil {
						logger.Warn("error occurred during shutdown rules proxy blob storages", "err", err)
					}
				}
			}()

			return xds.Serve(ctx, logger, appConfig.Proxy, repositories)
		},
	}

	c.Flags().StringVarP(&configFile, "config", "c", "", "Config file path")
	return c
}
