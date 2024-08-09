// Copyright 2020 Envoyproxy Authors
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package xds

import (
	"context"

	"google.golang.org/grpc"

	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/go-control-plane/pkg/test/v3"
)

type Server struct {
	xdsserver server.Server
}

func NewServer(ctx context.Context) (*Server, error) {
	logger := Logger{}
	nodeID := "node-id"
	// Create a cache
	cache := cache.NewSnapshotCache(false, cache.IDHash{}, logger)

	// Create the snapshot that we'll serve to Envoy
	snapshot := GenerateSnapshot()
	if err := snapshot.Consistent(); err != nil {
		logger.Errorf("snapshot inconsistency: %+v\n%+v", snapshot, err)
		return nil, err
	}
	logger.Debugf("will serve snapshot %+v", snapshot)

	// Add the snapshot to the cache
	if err := cache.SetSnapshot(context.Background(), nodeID, snapshot); err != nil {
		logger.Errorf("snapshot error %q for %+v", err, snapshot)
		return nil, err
	}

	cb := &test.Callbacks{Debug: logger.Debug}
	srv := server.NewServer(ctx, cache, cb)
	return &Server{srv}, nil
}

func (s *Server) Register(grpcServer *grpc.Server) {
	// register services
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, s.xdsserver)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, s.xdsserver)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, s.xdsserver)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, s.xdsserver)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, s.xdsserver)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, s.xdsserver)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, s.xdsserver)
}
