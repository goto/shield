package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/authzed/authzed-go/proto/authzed/api/v0"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	"github.com/goto/shield/config"
	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/policy"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/role"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/adapter"
	"github.com/goto/shield/internal/api"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/internal/server"
	"github.com/goto/shield/internal/store/blob"
	"github.com/goto/shield/internal/store/inmemory"
	"github.com/goto/shield/internal/store/postgres"
	"github.com/goto/shield/internal/store/spicedb"
	"github.com/goto/shield/pkg/db"

	"github.com/goto/salt/log"
	"github.com/goto/salt/telemetry"
	"github.com/pkg/profile"
	"google.golang.org/grpc/codes"
)

var ruleCacheRefreshDelay = time.Minute * 2

func StartServer(logger *log.Zap, cfg *config.Shield) error {
	if profiling := os.Getenv("SHIELD_PROFILE"); profiling == "true" || profiling == "1" {
		defer profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook).Stop()
	}

	// @TODO: need to inject custom logger wrapper over zap into ctx to use it internally
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancelFunc()

	cleanUpTelemetry, err := telemetry.Init(ctx, cfg.Telemetry, logger)
	if err != nil {
		return err
	}

	defer cleanUpTelemetry()

	dbClient, err := setupDB(cfg.DB, logger)
	if err != nil {
		return err
	}
	defer func() {
		logger.Info("cleaning up db")
		dbClient.Close()
	}()

	// load resource config
	if cfg.App.ResourcesConfigPath == "" {
		return errors.New("resource config path cannot be left empty")
	}

	resourceBlobFS, err := blob.NewStore(ctx, cfg.App.ResourcesConfigPath, cfg.App.ResourcesConfigPathSecret)
	if err != nil {
		return err
	}
	resourceBlobRepository := blob.NewResourcesRepository(logger, resourceBlobFS)
	if err := resourceBlobRepository.InitCache(ctx, ruleCacheRefreshDelay); err != nil {
		return err
	}
	defer func() {
		logger.Info("cleaning up resource blob")
		defer resourceBlobRepository.Close()
	}()

	spiceDBClient, err := spicedb.New(cfg.SpiceDB, logger)
	if err != nil {
		return err
	}

	nrApp, err := setupNewRelic(cfg.NewRelic, logger)
	if err != nil {
		return err
	}

	schemaMigrationConfig := schema.NewSchemaMigrationConfig(cfg.App.DefaultSystemEmail, cfg.App.ServiceData.BootstrapEnabled)

	appConfig := activity.AppConfig{Version: config.Version}
	var activityRepository activity.Repository
	switch cfg.Log.Activity.Sink {
	case activity.SinkTypeDB:
		activityRepository = postgres.NewActivityRepository(dbClient)
	case activity.SinkTypeStdout:
		stdoutLogger, err := zap.NewStdLogAt(logger.GetInternalZapLogger().Desugar(), logger.GetInternalZapLogger().Level())
		if err != nil {
			return err
		}
		activityRepository = activity.NewStdoutRepository(stdoutLogger.Writer())
	default:
		activityRepository = activity.NewStdoutRepository(io.Discard)
	}
	activityService := activity.NewService(appConfig, activityRepository)

	userRepository := postgres.NewUserRepository(dbClient)
	userService := user.NewService(logger, userRepository, activityService)

	actionRepository := postgres.NewActionRepository(dbClient)
	actionService := action.NewService(logger, actionRepository, userService, activityService)

	roleRepository := postgres.NewRoleRepository(dbClient)
	roleService := role.NewService(logger, roleRepository, userService, activityService)

	policyPGRepository := postgres.NewPolicyRepository(dbClient)
	policySpiceRepository := spicedb.NewPolicyRepository(spiceDBClient)
	policyService := policy.NewService(logger, policyPGRepository, userService, activityService)

	namespaceRepository := postgres.NewNamespaceRepository(dbClient)
	namespaceService := namespace.NewService(logger, namespaceRepository, userService, activityService)

	s := schema.NewSchemaMigrationService(
		blob.NewSchemaConfigRepository(resourceBlobFS),
		namespaceService,
		roleService,
		actionService,
		policyService,
		policySpiceRepository,
		userRepository,
		schemaMigrationConfig,
	)

	err = s.RunMigrations(ctx)
	if err != nil {
		return err
	}

	deps, err := BuildAPIDependencies(ctx, logger, activityRepository, resourceBlobRepository, dbClient, spiceDBClient, cfg)
	if err != nil {
		return err
	}

	// serving proxies
	cbs, cps, err := serveProxies(ctx, logger, cfg.App.IdentityProxyHeader, cfg.App.UserIDHeader, cfg.Proxy, deps.ResourceService, deps.RelationService, deps.UserService, deps.GroupService, deps.ProjectService, deps.RelationAdapter)
	if err != nil {
		return err
	}
	defer func() {
		// clean up stage
		logger.Info("cleaning up rules proxy blob")
		for _, f := range cbs {
			if err := f(); err != nil {
				logger.Warn("error occurred during shutdown rules proxy blob storages", "err", err)
			}
		}

		logger.Info("cleaning up proxies")
		for _, f := range cps {
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*20)
			if err := f(shutdownCtx); err != nil {
				shutdownCancel()
				logger.Warn("error occurred during shutdown proxies", "err", err)
				continue
			}
			shutdownCancel()
		}
	}()

	// serving server
	return server.Serve(ctx, logger, cfg.App, nrApp, deps)
}

func BuildAPIDependencies(
	ctx context.Context,
	logger log.Logger,
	activityRepository activity.Repository,
	resourceBlobRepository *blob.ResourcesRepository,
	dbc *db.Client,
	sdb *spicedb.SpiceDB,
	cfg *config.Shield,
) (api.Deps, error) {
	cache, err := inmemory.NewCache(cfg.App.CacheConfig)
	if err != nil {
		return api.Deps{}, err
	}

	if err = cache.MonitorCache(otel.Meter("github.com/goto/shield/internal/store/inmemory")); err != nil {
		return api.Deps{}, err
	}

	appConfig := activity.AppConfig{Version: config.Version}
	activityService := activity.NewService(appConfig, activityRepository)

	userRepository := postgres.NewUserRepository(dbc)
	userService := user.NewService(logger, userRepository, activityService)

	actionRepository := postgres.NewActionRepository(dbc)
	actionService := action.NewService(logger, actionRepository, userService, activityService)

	namespaceRepository := postgres.NewNamespaceRepository(dbc)
	namespaceService := namespace.NewService(logger, namespaceRepository, userService, activityService)

	roleRepository := postgres.NewRoleRepository(dbc)
	roleService := role.NewService(logger, roleRepository, userService, activityService)

	relationPGRepository := postgres.NewRelationRepository(dbc)
	relationSpiceRepository := spicedb.NewRelationRepository(sdb)
	relationService := relation.NewService(logger, relationPGRepository, relationSpiceRepository, userService, activityService)

	groupRepository := postgres.NewGroupRepository(dbc)
	cachedGroupRepository := inmemory.NewCachedGroupRepository(cache, groupRepository)
	groupService := group.NewService(logger, groupRepository, cachedGroupRepository, relationService, userService, activityService)

	organizationRepository := postgres.NewOrganizationRepository(dbc)
	organizationService := organization.NewService(logger, organizationRepository, relationService, userService, activityService)

	projectRepository := postgres.NewProjectRepository(dbc)
	projectService := project.NewService(logger, projectRepository, relationService, userService, activityService)

	policyPGRepository := postgres.NewPolicyRepository(dbc)
	policyService := policy.NewService(logger, policyPGRepository, userService, activityService)

	resourcePGRepository := postgres.NewResourceRepository(dbc)
	resourceService := resource.NewService(
		logger, resourcePGRepository, resourceBlobRepository, relationService, userService, projectService, organizationService, groupService, activityService)

	serviceDataRepository := postgres.NewServiceDataRepository(dbc)
	serviceDataService := servicedata.NewService(logger, serviceDataRepository, resourceService, relationService, projectService, userService, activityService)

	relationAdapter := adapter.NewRelation(groupService, userService, relationService, roleService)

	dependencies := api.Deps{
		OrgService:         organizationService,
		UserService:        userService,
		ProjectService:     projectService,
		GroupService:       groupService,
		RelationService:    relationService,
		ResourceService:    resourceService,
		RoleService:        roleService,
		PolicyService:      policyService,
		ActionService:      actionService,
		NamespaceService:   namespaceService,
		RelationAdapter:    relationAdapter,
		ActivityService:    activityService,
		ServiceDataService: serviceDataService,
	}
	return dependencies, nil
}

func setupNewRelic(cfg config.NewRelic, logger log.Logger) (*newrelic.Application, error) {
	nrApp, err := newrelic.NewApplication(func(nrCfg *newrelic.Config) {
		nrCfg.Enabled = cfg.Enabled
		nrCfg.AppName = cfg.AppName
		nrCfg.ErrorCollector.IgnoreStatusCodes = []int{
			http.StatusNotFound,
			http.StatusUnauthorized,
			int(codes.Unauthenticated),
			int(codes.PermissionDenied),
			int(codes.InvalidArgument),
			int(codes.AlreadyExists),
		}
		nrCfg.License = cfg.License
	})
	if err != nil {
		return nil, errors.New("failed to load Newrelic Application")
	}
	return nrApp, nil
}

func setupDB(cfg db.Config, logger log.Logger) (dbc *db.Client, err error) {
	// prefer use pgx instead of lib/pq for postgres to catch pg error
	if cfg.Driver == "postgres" {
		cfg.Driver = "pgx"
	}
	dbc, err = db.New(cfg)
	if err != nil {
		err = fmt.Errorf("failed to setup db: %w", err)
		return
	}

	return
}
