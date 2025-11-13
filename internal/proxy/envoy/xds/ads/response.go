package ads

import (
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	xds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

type ResponseStream struct {
	stream      xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer
	versionInfo string
	nonce       string
}

func (s ResponseStream) StreamCDS(clusters []*cluster.Cluster) error {
	if len(clusters) == 0 {
		return nil
	}

	var resources []*anypb.Any
	for _, cls := range clusters {
		res, err := proto.Marshal(cls)
		if err != nil {
			return err
		}

		resources = append(resources, &anypb.Any{
			TypeUrl: resource.ClusterType,
			Value:   res,
		})
	}

	resp := &xds.DiscoveryResponse{
		VersionInfo: s.versionInfo,
		Nonce:       s.nonce,
		Resources:   resources,
		TypeUrl:     resource.ClusterType,
	}

	return s.stream.Send(resp)
}

func (s ResponseStream) StreamLDS(listeners []*listener.Listener) error {
	if len(listeners) == 0 {
		return nil
	}

	var resources []*anypb.Any
	for _, ls := range listeners {
		res, err := proto.Marshal(ls)
		if err != nil {
			return err
		}

		resources = append(resources, &anypb.Any{
			TypeUrl: resource.ListenerType,
			Value:   res,
		})
	}

	resp := &xds.DiscoveryResponse{
		VersionInfo: s.versionInfo,
		Nonce:       s.nonce,
		Resources:   resources,
		TypeUrl:     resource.ListenerType,
	}
	return s.stream.Send(resp)
}

func (s ResponseStream) StreamRDS(routes []*route.RouteConfiguration) error {
	if len(routes) == 0 {
		return nil
	}

	var resources []*anypb.Any
	for _, r := range routes {
		res, err := proto.Marshal(r)
		if err != nil {
			return err
		}

		resources = append(resources, &anypb.Any{
			TypeUrl: resource.RouteType,
			Value:   res,
		})
	}

	resp := &xds.DiscoveryResponse{
		VersionInfo: s.versionInfo,
		Nonce:       s.nonce,
		Resources:   resources,
		TypeUrl:     resource.RouteType,
	}

	return s.stream.Send(resp)
}

func NewResponseStream(stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer, versionInfo, nonce string) ResponseStream {
	return ResponseStream{
		stream:      stream,
		versionInfo: versionInfo,
		nonce:       nonce,
	}
}
