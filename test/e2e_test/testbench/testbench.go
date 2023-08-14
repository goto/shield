package testbench

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/goto/salt/log"
	"github.com/goto/shield/cmd"
	"github.com/goto/shield/config"
	"github.com/goto/shield/internal/adapter"
	"github.com/goto/shield/internal/api"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/internal/store/spicedb"
	"github.com/goto/shield/pkg/db"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"context"
	"errors"
	"net/http"

	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1"
	"github.com/goto/shield/internal/proxy"
	"github.com/goto/shield/internal/proxy/hook"
	authz_hook "github.com/goto/shield/internal/proxy/hook/authz"
	"github.com/goto/shield/internal/proxy/middleware/attributes"
	"github.com/goto/shield/internal/proxy/middleware/authz"
	"github.com/goto/shield/internal/proxy/middleware/basic_auth"
	"github.com/goto/shield/internal/proxy/middleware/observability"
	"github.com/goto/shield/internal/proxy/middleware/prefix"
	"github.com/goto/shield/internal/proxy/middleware/rulematch"
	"github.com/goto/shield/internal/store/blob"
)

const (
	preSharedKey         = "shield"
	waitContainerTimeout = 60 * time.Second
)

var (
	RuleCacheRefreshDelay = time.Minute * 2
)

type TestBench struct {
	PGConfig          db.Config
	SpiceDBConfig     spicedb.Config
	bridgeNetworkName string
	pool              *dockertest.Pool
	network           *docker.Network
	resources         []*dockertest.Resource
}

func Init(appConfig *config.Shield) (*TestBench, *config.Shield, error) {
	var (
		err    error
		logger = log.NewZap()
	)

	te := &TestBench{
		bridgeNetworkName: fmt.Sprintf("bridge-%s", uuid.New().String()),
		resources:         []*dockertest.Resource{},
	}

	te.pool, err = dockertest.NewPool("")
	if err != nil {
		return nil, nil, err
	}

	// Create a bridge network for testing.
	te.network, err = te.pool.Client.CreateNetwork(docker.CreateNetworkOptions{
		Name: te.bridgeNetworkName,
	})
	if err != nil {
		return nil, nil, err
	}

	// pg 1
	logger.Info("creating main postgres...")
	_, connMainPGExternal, res, err := initPG(logger, te.network, te.pool, "test_db")
	if err != nil {
		return nil, nil, err
	}
	te.resources = append(te.resources, res)
	logger.Info("main postgres is created")

	// pg 2
	logger.Info("creating spicedb postgres...")
	connSpicePGInternal, _, res, err := initPG(logger, te.network, te.pool, "spicedb")
	if err != nil {
		return nil, nil, err
	}
	te.resources = append(te.resources, res)
	logger.Info("spicedb postgres is created")

	logger.Info("migrating spicedb...")
	if err = migrateSpiceDB(logger, te.network, te.pool, connSpicePGInternal); err != nil {
		return nil, nil, err
	}
	logger.Info("spicedb is migrated")

	logger.Info("starting up spicedb...")
	spiceDBPort, res, err := startSpiceDB(logger, te.network, te.pool, connSpicePGInternal, preSharedKey)
	if err != nil {
		return nil, nil, err
	}
	te.resources = append(te.resources, res)
	logger.Info("spicedb is up")

	te.PGConfig = db.Config{
		Driver:              "postgres",
		URL:                 connMainPGExternal,
		MaxIdleConns:        10,
		MaxOpenConns:        10,
		ConnMaxLifeTime:     time.Millisecond * 100,
		MaxQueryTimeoutInMS: time.Millisecond * 100,
	}

	te.SpiceDBConfig = spicedb.Config{
		Host:         "localhost",
		Port:         spiceDBPort,
		PreSharedKey: preSharedKey,
	}

	appConfig.DB = te.PGConfig
	appConfig.SpiceDB = te.SpiceDBConfig

	logger.Info("migrating shield...")
	if err = migrateShield(appConfig); err != nil {
		return nil, nil, err
	}
	logger.Info("shield is migrated")

	logger.Info("starting up shield...")
	startShield(appConfig)
	logger.Info("shield is up")

	return te, appConfig, nil
}

func (te *TestBench) CleanUp() error {
	return nil
}

func ServeProxies(
	ctx context.Context,
	logger *log.Zap,
	identityProxyHeaderKey,
	userIDHeaderKey string,
	cfg proxy.ServicesConfig,
	resourceService *resource.Service,
	relationService *relation.Service,
	userService *user.Service,
	projectService *project.Service,
	relationAdapter *adapter.Relation,
) ([]func() error, []func(ctx context.Context) error, error) {
	var cleanUpBlobs []func() error
	var cleanUpProxies []func(ctx context.Context) error

	for _, svcConfig := range cfg.Services {
		hookPipeline := buildHookPipeline(logger, resourceService, relationService, relationAdapter, identityProxyHeaderKey)

		h2cProxy := proxy.NewH2c(
			proxy.NewH2cRoundTripper(logger, hookPipeline),
			proxy.NewDirector(),
		)

		// load rules sets
		if svcConfig.RulesPath == "" {
			return nil, nil, errors.New("ruleset field cannot be left empty")
		}

		ruleBlobFS, err := blob.NewStore(ctx, svcConfig.RulesPath, svcConfig.RulesPathSecret)
		if err != nil {
			return nil, nil, err
		}

		ruleBlobRepository := blob.NewRuleRepository(logger, ruleBlobFS)
		if err := ruleBlobRepository.InitCache(ctx, RuleCacheRefreshDelay); err != nil {
			return nil, nil, err
		}
		cleanUpBlobs = append(cleanUpBlobs, ruleBlobRepository.Close)

		ruleService := rule.NewService(ruleBlobRepository)

		middlewarePipeline := buildMiddlewarePipeline(logger, h2cProxy, identityProxyHeaderKey, userIDHeaderKey, resourceService, userService, ruleService, projectService)

		cps := proxy.Serve(ctx, logger, svcConfig, middlewarePipeline)
		cleanUpProxies = append(cleanUpProxies, cps)
	}

	logger.Info("[shield] proxy is up")
	return cleanUpBlobs, cleanUpProxies, nil
}

func buildHookPipeline(log log.Logger, resourceService v1beta1.ResourceService, relationService v1beta1.RelationService, relationAdapter *adapter.Relation, identityProxyHeaderKey string) hook.Service {
	rootHook := hook.New()
	return authz_hook.New(log, rootHook, rootHook, resourceService, relationService, relationAdapter, identityProxyHeaderKey)
}

// buildPipeline builds middleware sequence
func buildMiddlewarePipeline(
	logger *log.Zap,
	proxy http.Handler,
	identityProxyHeaderKey, userIDHeaderKey string,
	resourceService *resource.Service,
	userService *user.Service,
	ruleService *rule.Service,
	projectService *project.Service,
) http.Handler {
	// Note: execution order is bottom up
	prefixWare := prefix.New(logger, proxy)
	casbinAuthz := authz.New(logger, prefixWare, userIDHeaderKey, resourceService, userService)
	basicAuthn := basic_auth.New(logger, casbinAuthz)
	attributeExtractor := attributes.New(logger, basicAuthn, identityProxyHeaderKey, projectService)
	matchWare := rulematch.New(logger, attributeExtractor, rulematch.NewRouteMatcher(ruleService))
	observability := observability.New(logger, matchWare)
	return observability
}

func BuildAPIDependenciesAndMigrate(
	ctx context.Context,
	logger *log.Zap,
	resourceBlobRepository *blob.ResourcesRepository,
	dbc *db.Client,
	sdb *spicedb.SpiceDB,
	rbfs blob.Bucket,
) (api.Deps, error) {
	policySpiceRepository := spicedb.NewPolicyRepository(sdb)

	dependencies, err := cmd.BuildAPIDependencies(ctx, logger, resourceBlobRepository, dbc, sdb)
	if err != nil {
		return api.Deps{}, err
	}

	s := schema.NewSchemaMigrationService(
		blob.NewSchemaConfigRepository(rbfs),
		dependencies.NamespaceService,
		dependencies.RoleService,
		dependencies.ActionService,
		dependencies.PolicyService,
		policySpiceRepository,
	)

	if err := s.RunMigrations(ctx); err != nil {
		return api.Deps{}, err
	}

	return dependencies, nil
}
