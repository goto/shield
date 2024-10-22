package ads

import (
	"context"
	"io"
	"strconv"
	"time"

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	xds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/goto/salt/log"
)

type DiscoveryResource struct {
	Clusters  []*cluster.Cluster
	Listeners []*listener.Listener
	Routes    []*route.RouteConfiguration
}

type Client struct {
	NodeID      string
	LastUpdated time.Time
}

type Stream struct {
	ctx             context.Context
	cancel          func()
	logger          log.Logger
	stream          xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer
	client          Client
	services        map[string]Service
	messageChan     MessageChan
	refreshInterval time.Duration
}

func NewStream(logger log.Logger, refreshInterval time.Duration, stream xds.AggregatedDiscoveryService_StreamAggregatedResourcesServer, services map[string]Service) Stream {
	ctx, cancel := context.WithCancel(context.Background())
	return Stream{
		ctx:             ctx,
		cancel:          cancel,
		logger:          logger,
		stream:          stream,
		services:        services,
		messageChan:     make(MessageChan),
		refreshInterval: refreshInterval,
	}
}

func (s Stream) Stream() error {
	terminate := make(chan bool)

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			default:
				in, err := s.stream.Recv()
				if err == io.EOF {
					return
				}

				if err != nil {
					s.logger.Error(err.Error())
					return
				}

				if in.ResponseNonce == "" {
					s.logger.Info("received request on stream", "typeurl", in.TypeUrl)
					message := Message{
						NodeID:      in.Node.Id,
						VersionInfo: strconv.FormatInt(time.Now().UnixNano(), 10),
						Nonce:       strconv.FormatInt(time.Now().UnixNano(), 10),
						TypeUrl:     in.TypeUrl,
					}
					s.messageChan.Push(message)
					s.client.LastUpdated = time.Now()

					if s.client.NodeID == "" {
						s.client.NodeID = in.Node.Id
						go s.PushUpdatePeriodically()
					}
				} else if in.ErrorDetail == nil {
					s.logger.Info("received ACK on stream", "typeurl", in.TypeUrl, "version_info", in.VersionInfo)
				} else {
					s.logger.Info("received NACK on stream", "typeurl", in.TypeUrl, "version_info", in.VersionInfo, "error", in.ErrorDetail)
				}
			}
		}
	}()

	go func() {
		for e := range s.messageChan {
			if err := s.streamResponses(e); err != nil {
				s.logger.Debug("error while streaming response", "error", err)
			}
		}
	}()

	go func() {
		<-s.stream.Context().Done()
		close(s.messageChan)
		s.cancel()
		terminate <- true
	}()
	<-terminate
	return nil
}

func (s Stream) streamResponses(message Message) error {
	cfg := &DiscoveryResource{}
	var err error
	if repo, ok := s.services[message.NodeID]; ok {
		cfg, err = repo.Get(s.ctx)
		if err != nil {
			return err
		}
	}

	responseStream := NewResponseStream(s.stream, message.VersionInfo, message.Nonce)
	switch message.TypeUrl {
	case CLUSTER_TYPE_URL:
		if err := responseStream.StreamCDS(cfg.Clusters); err != nil {
			return err
		}
	case LISTENER_TYPE_URL:
		if err := responseStream.StreamLDS(cfg.Listeners); err != nil {
			return err
		}
	case ROUTER_TYPE_URL:
		if err := responseStream.StreamRDS(cfg.Routes); err != nil {
			return err
		}
	default:
		if err := responseStream.StreamCDS(cfg.Clusters); err != nil {
			return err
		}
		if err := responseStream.StreamLDS(cfg.Listeners); err != nil {
			return err
		}
		if err := responseStream.StreamRDS(cfg.Routes); err != nil {
			return err
		}
	}

	return nil
}

func (s Stream) PushUpdatePeriodically() {
	service, ok := s.services[s.client.NodeID]
	if !ok {
		s.logger.Debug("service not found", "node_id", s.client.NodeID)
		return
	}

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			time.Sleep(s.refreshInterval)
			if service.IsUpdated(s.ctx, s.client.LastUpdated) {
				s.logger.Debug("discovery resource update found", "node_id", s.client.NodeID)
				message := Message{
					NodeID:      s.client.NodeID,
					VersionInfo: strconv.FormatInt(time.Now().UnixNano(), 10),
					Nonce:       strconv.FormatInt(time.Now().UnixNano(), 10),
				}
				s.messageChan.Push(message)
				s.client.LastUpdated = time.Now()
			} else {
				s.logger.Debug("no discovery resource update", "node_id", s.client.NodeID)
			}
		}
	}
}
