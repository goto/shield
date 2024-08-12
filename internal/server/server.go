package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/goto/salt/log"
	"github.com/goto/salt/mux"
	"github.com/goto/shield/internal/api"
	"github.com/goto/shield/internal/api/v1beta1"
	"github.com/goto/shield/internal/server/grpc_interceptors"
	"github.com/goto/shield/internal/server/health"
	envoyauthz "github.com/goto/shield/proxy/envoy/authz"
	envoyextproc "github.com/goto/shield/proxy/envoy/proc"
	"github.com/goto/shield/proxy/envoy/xds"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_service_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func Serve(
	ctx context.Context,
	logger log.Logger,
	cfg Config,
	nrApp *newrelic.Application,
	deps api.Deps,
) error {
	httpMux := http.NewServeMux()

	grpcDialCtx, grpcDialCancel := context.WithTimeout(ctx, time.Second*5)
	defer grpcDialCancel()

	grpcConn, err := grpc.DialContext(
		grpcDialCtx,
		cfg.grpcAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(cfg.GRPC.MaxRecvMsgSize),
			grpc.MaxCallSendMsgSize(cfg.GRPC.MaxSendMsgSize),
		))
	if err != nil {
		return err
	}

	grpcGateway := runtime.NewServeMux(
		runtime.WithHealthEndpointAt(grpc_health_v1.NewHealthClient(grpcConn), "/ping"),
		runtime.WithIncomingHeaderMatcher(customHeaderMatcherFunc(map[string]bool{cfg.IdentityProxyHeader: true})),
	)

	httpMux.Handle("/admin/", http.StripPrefix("/admin", grpcGateway))

	if err := shieldv1beta1.RegisterShieldServiceHandler(ctx, grpcGateway, grpcConn); err != nil {
		return err
	}

	grpcServiceDataGateway := runtime.NewServeMux(
		runtime.WithHealthEndpointAt(grpc_health_v1.NewHealthClient(grpcConn), "/ping"),
		runtime.WithIncomingHeaderMatcher(customHeaderMatcherFunc(map[string]bool{cfg.IdentityProxyHeader: true})),
	)

	httpMux.Handle(fmt.Sprintf("%s/", cfg.PublicAPIPrefix), http.StripPrefix(cfg.PublicAPIPrefix, grpcServiceDataGateway))

	if err := shieldv1beta1.RegisterServiceDataServiceHandler(ctx, grpcServiceDataGateway, grpcConn); err != nil {
		return err
	}

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		getGRPCMiddleware(cfg, logger, nrApp),
	)
	reflection.Register(grpcServer)

	healthHandler := health.NewHandler()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthHandler)

	// Envoy control plane
	envoy_service_auth_v3.RegisterAuthorizationServer(grpcServer, envoyauthz.New(logger))
	envoy_service_ext_proc_v3.RegisterExternalProcessorServer(grpcServer, envoyextproc.NewServer(logger))

	xdsServer, err := xds.NewServer(ctx)
	if err != nil {
		return err
	}
	xdsServer.Register(grpcServer)

	serviceDataConfig := v1beta1.ServiceDataConfig{MaxUpsert: cfg.ServiceData.MaxNumUpsertData, DefaultServiceDataProject: cfg.ServiceData.DefaultServiceDataProject}
	err = v1beta1.Register(ctx, grpcServer, deps, cfg.CheckAPILimit, serviceDataConfig)
	if err != nil {
		return err
	}

	httpMuxMetrics := http.NewServeMux()

	logger.Info("[shield] api server starting", "http-port", cfg.Port, "grpc-port", cfg.GRPC.Port, "metrics-port", cfg.MetricsPort)

	if err := mux.Serve(
		ctx,
		mux.WithHTTPTarget(fmt.Sprintf(":%d", cfg.Port), &http.Server{
			Handler:        httpMux,
			ReadTimeout:    120 * time.Second,
			WriteTimeout:   120 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}),
		mux.WithHTTPTarget(fmt.Sprintf(":%d", cfg.MetricsPort), &http.Server{
			Handler:        httpMuxMetrics,
			ReadTimeout:    120 * time.Second,
			WriteTimeout:   120 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}),
		mux.WithGRPCTarget(fmt.Sprintf(":%d", cfg.GRPC.Port), grpcServer),
		mux.WithGracePeriod(5*time.Second),
	); !errors.Is(err, context.Canceled) {
		logger.Error("mux serve error", "err", err)
		return nil
	}

	logger.Info("server stopped gracefully")
	return nil
}

// REVISIT: passing config.Shield as reference
func getGRPCMiddleware(cfg Config, logger log.Logger, nrApp *newrelic.Application) grpc.ServerOption {
	recoveryFunc := func(p interface{}) (err error) {
		fmt.Println("-----------------------------")
		return status.Errorf(codes.Internal, "internal server error")
	}

	grpcRecoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	grpcZapLogger := zap.NewExample().Sugar()
	loggerZap, ok := logger.(*log.Zap)
	if ok {
		grpcZapLogger = loggerZap.GetInternalZapLogger()
	}
	return grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_interceptors.EnrichCtxWithIdentity(cfg.IdentityProxyHeader),
			grpc_zap.UnaryServerInterceptor(grpcZapLogger.Desugar()),
			grpc_recovery.UnaryServerInterceptor(grpcRecoveryOpts...),
			grpc_ctxtags.UnaryServerInterceptor(),
			nrgrpc.UnaryServerInterceptor(nrApp),
		))
}

func customHeaderMatcherFunc(headerKeys map[string]bool) func(key string) (string, bool) {
	return func(key string) (string, bool) {
		if _, ok := headerKeys[key]; ok {
			return key, true
		}
		return runtime.DefaultHeaderMatcher(key)
	}
}
