package ads_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	accesslog "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	http_connection_manager "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	uri_template "github.com/envoyproxy/go-control-plane/envoy/extensions/path/match/uri_template/v3"
	matcherv3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/proxy"
	"github.com/goto/shield/internal/proxy/envoy/xds/ads"
	"github.com/goto/shield/internal/proxy/envoy/xds/ads/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	testConfig = proxy.Config{
		Name: "test-proxy",
		Port: 5556,
		Host: "0.0.0.0",
	}

	testRule = rule.Rule{
		Frontend: rule.Frontend{
			URL:    "/shield/test",
			Method: "GET",
		},
		Backend: rule.Backend{
			URL:       "http://localhost:8080",
			Namespace: "shield",
			Prefix:    "/shield",
		},
		Middlewares: rule.MiddlewareSpecs{},
		Hooks:       rule.HookSpecs{},
	}

	testDiscoveryResource = ads.DiscoveryResource{
		Clusters:  []*cluster.Cluster{testCluster},
		Listeners: []*listener.Listener{testListener},
		Routes:    []*route.RouteConfiguration{testRouteConfiguration},
	}

	testCluster = &cluster.Cluster{
		ClusterDiscoveryType: &cluster.Cluster_Type{
			Type: cluster.Cluster_LOGICAL_DNS,
		},
		DnsLookupFamily: cluster.Cluster_V4_PREFERRED,
		Name:            "shield",
		ConnectTimeout:  durationpb.New(1 * time.Second),
		LoadAssignment:  &testCLA,
	}

	testLbEndppoint = &endpoint.LbEndpoint{
		HostIdentifier: &endpoint.LbEndpoint_Endpoint{
			Endpoint: &endpoint.Endpoint{
				Hostname: "localhost",
				Address: &core.Address{
					Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Protocol: core.SocketAddress_TCP,
							Address:  "localhost",
							PortSpecifier: &core.SocketAddress_PortValue{
								PortValue: 8080,
							},
						},
					},
				},
			},
		},
	}

	testLbEndppoints = &endpoint.LocalityLbEndpoints{
		LbEndpoints: []*endpoint.LbEndpoint{testLbEndppoint},
	}

	testCLA = endpoint.ClusterLoadAssignment{
		ClusterName: "shield",
		Endpoints:   []*endpoint.LocalityLbEndpoints{testLbEndppoints},
	}

	testListener = &listener.Listener{
		Name: "test-proxy",
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: 5556,
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{testFilterChain},
	}

	testAds = core.ConfigSource{
		ConfigSourceSpecifier: &core.ConfigSource_Ads{
			Ads: &core.AggregatedConfigSource{},
		},
	}

	testRouterFilter = &http_connection_manager.HttpFilter{
		Name: wellknown.Router,
		ConfigType: &http_connection_manager.HttpFilter_TypedConfig{
			TypedConfig: &anypb.Any{
				TypeUrl: ads.ROUTER_TYPE_URL,
			},
		},
	}

	testAL = accesslog.AccessLog{
		Name: "envoy.access_loggers.stdout",
		ConfigType: &accesslog.AccessLog_TypedConfig{
			TypedConfig: &anypb.Any{
				TypeUrl: ads.STDOUT_LOGGER_TYPE_URL,
			},
		},
	}

	testHttpConnManager = http_connection_manager.HttpConnectionManager{
		CodecType:  http_connection_manager.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		AccessLog:  []*accesslog.AccessLog{&testAL},
		RouteSpecifier: &http_connection_manager.HttpConnectionManager_Rds{
			Rds: &http_connection_manager.Rds{
				ConfigSource:    &testAds,
				RouteConfigName: "test-proxy",
			},
		},
		HttpFilters: []*http_connection_manager.HttpFilter{
			testRouterFilter,
		},
	}

	testHttpConnManagerBytes, _ = proto.Marshal(&testHttpConnManager)

	testFilterChain = &listener.FilterChain{
		Filters: []*listener.Filter{
			{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: &anypb.Any{
						TypeUrl: ads.HTTP_CONNECTION_MANAGER_TYPE_URL,
						Value:   testHttpConnManagerBytes,
					},
				},
			},
		},
	}

	testHeaderMatcher = &route.HeaderMatcher{
		Name: ":method",
		HeaderMatchSpecifier: &route.HeaderMatcher_StringMatch{
			StringMatch: &matcherv3.StringMatcher{
				MatchPattern: &matcherv3.StringMatcher_Exact{
					Exact: "GET",
				},
			},
		},
	}

	testPathTemplate = uri_template.UriTemplateMatchConfig{
		PathTemplate: "/shield/test",
	}

	testPathTemplateBytes, _ = proto.Marshal(&testPathTemplate)

	testRoute = &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_PathMatchPolicy{
				PathMatchPolicy: &core.TypedExtensionConfig{
					Name: "envoy.extensions.path.match.uri_template.v3.UriTemplateMatchConfig",
					TypedConfig: &anypb.Any{
						TypeUrl: ads.URI_TEMPLATE_TYPE_URL,
						Value:   testPathTemplateBytes,
					},
				},
			},
			Headers: []*route.HeaderMatcher{
				testHeaderMatcher,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: "shield",
				},
				HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
					HostRewriteLiteral: "localhost",
				},
				RegexRewrite: &matcherv3.RegexMatchAndSubstitute{
					Pattern: &matcherv3.RegexMatcher{
						Regex: fmt.Sprintf("^(%s)(/.+$)", "/shield"),
					},
					Substitution: "\\2",
				},
			},
		},
	}

	testVH = &route.VirtualHost{
		Name:    "test-proxy",
		Domains: []string{"*"},
		Routes:  []*route.Route{testRoute},
	}

	testRouteConfiguration = &route.RouteConfiguration{
		Name:         "test-proxy",
		VirtualHosts: []*route.VirtualHost{testVH},
	}
)

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(t *testing.T) ads.Service
		want    *ads.DiscoveryResource
		wantErr error
	}{
		{
			name: "should return discovery resource",
			setup: func(t *testing.T) ads.Service {
				t.Helper()
				repository := &mocks.Repository{}
				repository.EXPECT().Fetch(mock.Anything).Return([]rule.Ruleset{
					{
						Rules: []rule.Rule{testRule},
					},
				}, nil)
				return ads.NewService(testConfig, repository)
			},
			want:    &testDiscoveryResource,
			wantErr: nil,
		},
		{
			name: "should return discovery resource",
			setup: func(t *testing.T) ads.Service {
				t.Helper()
				repository := &mocks.Repository{}
				repository.EXPECT().Fetch(mock.Anything).Return([]rule.Ruleset{}, rule.ErrMarshal)
				return ads.NewService(testConfig, repository)
			},
			want:    &ads.DiscoveryResource{},
			wantErr: rule.ErrMarshal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)
			ctx := context.Background()
			got, err := svc.Get(ctx)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
