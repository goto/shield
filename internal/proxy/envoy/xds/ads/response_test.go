package ads_test

import (
	"testing"

	clusterv3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	xds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/goto/shield/internal/proxy/envoy/xds/ads"
	"github.com/goto/shield/internal/proxy/envoy/xds/ads/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var (
	testClusterStream    = &clusterv3.Cluster{}
	testClusterBytes, _  = proto.Marshal(testClusterStream)
	testClusterResources = &anypb.Any{
		TypeUrl: resource.ClusterType,
		Value:   testClusterBytes,
	}

	testListenerStream    = &listenerv3.Listener{}
	testListenerBytes, _  = proto.Marshal(testListenerStream)
	testListenerResources = &anypb.Any{
		TypeUrl: resource.ListenerType,
		Value:   testListenerBytes,
	}

	testRouteStream    = &routev3.RouteConfiguration{}
	testRouteBytes, _  = proto.Marshal(testRouteStream)
	testRouteResources = &anypb.Any{
		TypeUrl: resource.RouteType,
		Value:   testRouteBytes,
	}
)

func TestStreamCDS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cluster []*clusterv3.Cluster
		setup   func(t *testing.T) ads.ResponseStream
		wantErr error
	}{
		{
			name:    "should return error from stream send",
			cluster: []*clusterv3.Cluster{testClusterStream},
			setup: func(t *testing.T) ads.ResponseStream {
				t.Helper()
				stream := mocks.AggregatedDiscoveryService_StreamAggregatedResourcesServer{}
				stream.EXPECT().Send(&xds.DiscoveryResponse{
					VersionInfo: "v1",
					Nonce:       "test",
					Resources:   []*anypb.Any{testClusterResources},
					TypeUrl:     resource.ClusterType,
				}).Return(nil)
				return ads.NewResponseStream(&stream, "v1", "test")
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp := tt.setup(t)

			assert.NotNil(t, resp)
			got := resp.StreamCDS(tt.cluster)

			assert.Equal(t, tt.wantErr, got)
		})
	}
}

func TestStreamLDS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		listener []*listenerv3.Listener
		setup    func(t *testing.T) ads.ResponseStream
		wantErr  error
	}{
		{
			name:     "should return error from stream send",
			listener: []*listenerv3.Listener{testListenerStream},
			setup: func(t *testing.T) ads.ResponseStream {
				t.Helper()
				stream := mocks.AggregatedDiscoveryService_StreamAggregatedResourcesServer{}
				stream.EXPECT().Send(&xds.DiscoveryResponse{
					VersionInfo: "v1",
					Nonce:       "test",
					Resources:   []*anypb.Any{testListenerResources},
					TypeUrl:     resource.ListenerType,
				}).Return(nil)
				return ads.NewResponseStream(&stream, "v1", "test")
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp := tt.setup(t)

			assert.NotNil(t, resp)
			got := resp.StreamLDS(tt.listener)

			assert.Equal(t, tt.wantErr, got)
		})
	}
}

func TestStreamRDS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		route   []*routev3.RouteConfiguration
		setup   func(t *testing.T) ads.ResponseStream
		wantErr error
	}{
		{
			name:  "should return error from stream send",
			route: []*routev3.RouteConfiguration{testRouteStream},
			setup: func(t *testing.T) ads.ResponseStream {
				t.Helper()
				stream := mocks.AggregatedDiscoveryService_StreamAggregatedResourcesServer{}
				stream.EXPECT().Send(&xds.DiscoveryResponse{
					VersionInfo: "v1",
					Nonce:       "test",
					Resources:   []*anypb.Any{testRouteResources},
					TypeUrl:     resource.RouteType,
				}).Return(nil)
				return ads.NewResponseStream(&stream, "v1", "test")
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			resp := tt.setup(t)

			assert.NotNil(t, resp)
			got := resp.StreamRDS(tt.route)

			assert.Equal(t, tt.wantErr, got)
		})
	}
}
