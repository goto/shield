package ads

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Service struct {
	config     proxy.Config
	repository Repository
}

func NewService(config proxy.Config, repository Repository) Service {
	return Service{
		config:     config,
		repository: repository,
	}
}

func (s Service) Get(ctx context.Context) (*DiscoveryResource, error) {
	ruleset, err := s.repository.GetAll(ctx)
	if err != nil {
		return &DiscoveryResource{}, err
	}

	var clusters []*cluster.Cluster
	var listeners []*listener.Listener
	var routes []*route.RouteConfiguration
	backendmap := make(map[string]bool)
	for _, rule := range ruleset {
		for _, r := range rule.Rules {
			if _, ok := backendmap[r.Backend.Namespace]; ok {
				continue
			}
			backendmap[r.Backend.Namespace] = true
			clusters = append(clusters, s.getCluster(r))
		}
	}

	routes = append(routes, s.getRoute(ruleset))

	ls, err := s.getListener()
	if err != nil {
		return &DiscoveryResource{}, err
	}

	listeners = append(listeners, ls)

	return &DiscoveryResource{
		Clusters:  clusters,
		Listeners: listeners,
		Routes:    routes,
	}, nil
}

func (s Service) getCluster(rule rule.Rule) *cluster.Cluster {
	return &cluster.Cluster{
		ClusterDiscoveryType: &cluster.Cluster_Type{
			Type: cluster.Cluster_LOGICAL_DNS,
		},
		DnsLookupFamily: cluster.Cluster_V4_ONLY,
		Name:            rule.Backend.Namespace,
		ConnectTimeout:  durationpb.New(1 * time.Second),
		LoadAssignment:  s.getEndpoint(rule),
	}
}

func (s Service) getEndpoint(rule rule.Rule) *endpoint.ClusterLoadAssignment {
	host, port, err := resolveHostPort(rule.Backend.URL)
	if err != nil {
		return nil
	}

	lbEndpoint := &endpoint.LbEndpoint{
		HostIdentifier: &endpoint.LbEndpoint_Endpoint{
			Endpoint: &endpoint.Endpoint{
				Hostname: host,
				Address: &core.Address{
					Address: &core.Address_SocketAddress{
						SocketAddress: &core.SocketAddress{
							Protocol: core.SocketAddress_TCP,
							Address:  host,
							PortSpecifier: &core.SocketAddress_PortValue{
								PortValue: port,
							},
						},
					},
				},
			},
		},
	}

	lbEndpoints := &endpoint.LocalityLbEndpoints{
		LbEndpoints: []*endpoint.LbEndpoint{lbEndpoint},
	}

	return &endpoint.ClusterLoadAssignment{
		ClusterName: rule.Backend.Namespace,
		Endpoints:   []*endpoint.LocalityLbEndpoints{lbEndpoints},
	}
}

func (s Service) getRoute(ruleset []rule.Ruleset) *route.RouteConfiguration {
	vh := &route.VirtualHost{
		Name:    s.config.Name,
		Domains: []string{"*"},
		Routes:  []*route.Route{},
	}

	rc := &route.RouteConfiguration{
		Name:         s.config.Name,
		VirtualHosts: []*route.VirtualHost{vh},
	}

	for _, rule := range ruleset {
		for _, r := range rule.Rules {
			host, _, err := resolveHostPort(r.Backend.URL)
			if err != nil {
				continue
			}
			headerMatcher := &route.HeaderMatcher{
				Name: ":method",
				HeaderMatchSpecifier: &route.HeaderMatcher_StringMatch{
					StringMatch: &matcherv3.StringMatcher{
						MatchPattern: &matcherv3.StringMatcher_Exact{
							Exact: r.Frontend.Method,
						},
					},
				},
			}

			pathTemplate := uri_template.UriTemplateMatchConfig{
				PathTemplate: r.Frontend.URL,
			}

			pathTemplateBytes, err := proto.Marshal(&pathTemplate)
			if err != nil {
				continue
			}

			rt := &route.Route{
				Match: &route.RouteMatch{
					PathSpecifier: &route.RouteMatch_PathMatchPolicy{
						PathMatchPolicy: &core.TypedExtensionConfig{
							Name: "envoy.extensions.path.match.uri_template.v3.UriTemplateMatchConfig",
							TypedConfig: &anypb.Any{
								TypeUrl: URI_TEMPLATE_TYPE_URL,
								Value:   pathTemplateBytes,
							},
						},
					},
					Headers: []*route.HeaderMatcher{
						headerMatcher,
					},
				},
			}
			if r.Backend.Prefix != "" {
				rt.Action = &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: r.Backend.Namespace,
						},
						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: host,
						},
						RegexRewrite: &matcherv3.RegexMatchAndSubstitute{
							Pattern: &matcherv3.RegexMatcher{
								Regex: fmt.Sprintf("^(%s)(/.+$)", r.Backend.Prefix),
							},
							Substitution: "\\2",
						},
					},
				}
			} else {
				rt.Action = &route.Route_Route{
					Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: r.Backend.Namespace,
						},
						HostRewriteSpecifier: &route.RouteAction_HostRewriteLiteral{
							HostRewriteLiteral: host,
						},
					},
				}
			}
			vh.Routes = append(vh.Routes, rt)
		}
	}

	return rc
}

func (s Service) getListener() (*listener.Listener, error) {
	ads := core.ConfigSource{
		ConfigSourceSpecifier: &core.ConfigSource_Ads{
			Ads: &core.AggregatedConfigSource{},
		},
	}

	routerFilter := &http_connection_manager.HttpFilter{
		Name: wellknown.Router,
		ConfigType: &http_connection_manager.HttpFilter_TypedConfig{
			TypedConfig: &anypb.Any{
				TypeUrl: ROUTER_TYPE_URL,
			},
		},
	}

	al := accesslog.AccessLog{
		Name: "envoy.access_loggers.stdout",
		ConfigType: &accesslog.AccessLog_TypedConfig{
			TypedConfig: &anypb.Any{
				TypeUrl: STDOUT_LOGGER_TYPE_URL,
			},
		},
	}

	httpConnManager := http_connection_manager.HttpConnectionManager{
		CodecType:  http_connection_manager.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		AccessLog:  []*accesslog.AccessLog{&al},
		RouteSpecifier: &http_connection_manager.HttpConnectionManager_Rds{
			Rds: &http_connection_manager.Rds{
				ConfigSource:    &ads,
				RouteConfigName: s.config.Name,
			},
		},
		HttpFilters: []*http_connection_manager.HttpFilter{
			routerFilter,
		},
	}

	httpConnManagerBytes, err := proto.Marshal(&httpConnManager)
	if err != nil {
		return &listener.Listener{}, err
	}

	filterChain := &listener.FilterChain{
		Filters: []*listener.Filter{
			{
				Name: wellknown.HTTPConnectionManager,
				ConfigType: &listener.Filter_TypedConfig{
					TypedConfig: &anypb.Any{
						TypeUrl: HTTP_CONNECTION_MANAGER_TYPE_URL,
						Value:   httpConnManagerBytes,
					},
				},
			},
		},
	}

	ls := &listener.Listener{
		Name: s.config.Name,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.SocketAddress_TCP,
					Address:  s.config.Host,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(s.config.Port),
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{filterChain},
	}

	return ls, nil
}

func (s Service) IsUpdated(ctx context.Context, since time.Time) bool {
	return s.repository.IsUpdated(ctx, since)
}

func resolveHostPort(urlString string) (string, uint32, error) {
	parsed, err := url.Parse(urlString)
	if err != nil {
		return "", 0, err
	}

	port := parsed.Port()
	if parsed.Port() == "" {
		switch parsed.Scheme {
		case "https":
			return parsed.Host, 443, nil
		default:
			return parsed.Host, 80, nil
		}
	}

	uintPort, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return "", 0, err
	}

	return parsed.Hostname(), uint32(uintPort), nil
}
