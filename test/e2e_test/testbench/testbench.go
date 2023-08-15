package testbench

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/goto/salt/log"
	"github.com/goto/shield/cmd"
	"github.com/goto/shield/config"
	"github.com/goto/shield/internal/proxy"
	"github.com/goto/shield/internal/server"
	"github.com/goto/shield/internal/store/postgres/migrations"
	"github.com/goto/shield/internal/store/spicedb"
	"github.com/goto/shield/pkg/db"
	"github.com/goto/shield/pkg/logger"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
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

func initTestBench(ctx context.Context, appConfig *config.Shield, withMockBackendServer bool) (*TestBench, *config.Shield, error) {
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

	if withMockBackendServer {
		go func() {
			mockServerPort := 4000
			startMockServer(ctx, logger, mockServerPort)
		}()
	}
	go func() {
		if err := cmd.StartServer(logger, appConfig); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	return te, appConfig, nil
}

func (te *TestBench) CleanUp() error {
	return nil
}

func SetupTests(t *testing.T, withMockBackendServer bool) (shieldv1beta1.ShieldServiceClient, *config.Shield, func(), func()) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	parent := filepath.Dir(wd)
	testDataPath := parent + "/testbench/testdata"

	proxyPort, err := GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	apiPort, err := GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	apiGRPCPort, err := GetFreePort()
	if err != nil {
		t.Fatal(err)
	}

	appConfig := &config.Shield{
		Log: logger.Config{
			Level: "fatal",
		},
		App: server.Config{
			Port: apiPort,
			GRPC: server.GRPCConfig{
				Port: apiGRPCPort,
			},
			IdentityProxyHeader: IdentityHeader,
			UserIDHeader:        userIDHeaderKey,
			ResourcesConfigPath: fmt.Sprintf("file://%s/%s", testDataPath, "configs/resources"),
			RulesPath:           fmt.Sprintf("file://%s/%s", testDataPath, "configs/rules"),
		},
		Proxy: proxy.ServicesConfig{
			Services: []proxy.Config{
				{
					Name:      "base",
					Port:      proxyPort,
					RulesPath: fmt.Sprintf("file://%s/%s", testDataPath, "configs/rules"),
				},
			},
		},
	}

	_, _, err = initTestBench(context.Background(), appConfig, withMockBackendServer)
	if err != nil {
		t.Fatal(err.Error())
	}

	ctx, cancelContextFunc := context.WithTimeout(context.Background(), time.Minute*5)

	shieldHost := fmt.Sprintf("localhost:%d", appConfig.App.GRPC.Port)
	client, cancelClient, err := CreateClient(ctx, shieldHost)
	if err != nil {
		t.Fatal(err)
	}

	if err := BootstrapMetadataKey(ctx, client, OrgAdminEmail, testDataPath); err != nil {
		t.Fatal(err)
	}

	if err := BootstrapUser(ctx, client, OrgAdminEmail, testDataPath); err != nil {
		t.Fatal(err)
	}
	if err := BootstrapOrganization(ctx, client, OrgAdminEmail, testDataPath); err != nil {
		t.Fatal(err)
	}
	if err := BootstrapProject(ctx, client, OrgAdminEmail, testDataPath); err != nil {
		t.Fatal(err)
	}
	if err := BootstrapGroup(ctx, client, OrgAdminEmail, testDataPath); err != nil {
		t.Fatal(err)
	}

	return client, appConfig, cancelClient, cancelContextFunc
}
func migrateShield(appConfig *config.Shield) error {
	return db.RunMigrations(db.Config{
		Driver: appConfig.DB.Driver,
		URL:    appConfig.DB.URL,
	}, migrations.MigrationFs, migrations.ResourcePath)
}
