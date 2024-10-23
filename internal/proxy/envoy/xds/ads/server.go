package ads

import (
	"errors"
	"time"

	xds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/goto/salt/log"
)

type Server struct {
	Logger          log.Logger
	Services        map[string]Service
	RefreshInterval time.Duration
}

func (a *Server) DeltaAggregatedResources(xds.AggregatedDiscoveryService_DeltaAggregatedResourcesServer) error {
	return errors.New("not implemented")
}

func (a *Server) StreamAggregatedResources(stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	err := NewStream(a.Logger, a.RefreshInterval, stream, a.Services).Stream()
	return err
}

func New(logger log.Logger, services map[string]Service, refreshInterval time.Duration) *Server {
	return &Server{
		Logger:          logger,
		Services:        services,
		RefreshInterval: refreshInterval,
	}
}
