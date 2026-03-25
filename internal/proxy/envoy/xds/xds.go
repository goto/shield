package xds

import (
	"context"
	"fmt"
	"net"

	xds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/goto/salt/log"
	"github.com/goto/shield/internal/proxy"
	"github.com/goto/shield/internal/proxy/envoy/xds/ads"
	"google.golang.org/grpc"
)

func Serve(ctx context.Context, logger log.Logger, cfg proxy.ServicesConfig, repositories map[string]ads.Repository) error {
	xdsURL := fmt.Sprintf("%s:%d", cfg.EnvoyAgent.XDS.Host, cfg.EnvoyAgent.XDS.Port)
	logger.Info("[envoy agent] starting envoy xds", "url", xdsURL)

	server := grpc.NewServer()

	services := make(map[string]ads.Service)
	for _, c := range cfg.Services {
		if repo, ok := repositories[c.Name]; ok {
			services[c.Name] = ads.NewService(c, repo)
		}
	}
	xds.RegisterAggregatedDiscoveryServiceServer(server, ads.New(logger, services, cfg.EnvoyAgent.XDS.RefreshInterval))

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.EnvoyAgent.XDS.Host, cfg.EnvoyAgent.XDS.Port))
	if err != nil {
		logger.Error("[envoy agent] envoy xds failed to listen: %v\n", err)
		return err
	}

	return server.Serve(lis)
}
