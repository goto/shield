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
	NodeID            string
	LatestVersionSent string
	LatestVersionACK  string
	LatestNonceSent   string
	LatestNonceACK    string
	LastUpdated       time.Time
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
		refreshInterval: refreshInterval,
		stream:          stream,
		messageChan:     make(MessageChan),
		services:        services,
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
					return
				}

				if in.ResponseNonce == "" {
					versionInfo := strconv.FormatInt(time.Now().UnixNano(), 10)
					nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
					message := Message{
						NodeID:      in.Node.Id,
						VersionInfo: versionInfo,
						Nonce:       nonce,
					}
					s.messageChan.Push(message)
					s.client.LastUpdated = time.Now()
					s.client.LatestVersionSent = versionInfo
					s.client.LatestNonceSent = nonce

					if s.client.NodeID == "" {
						s.client.NodeID = in.Node.Id
						s.PushUpdatePeriodically()
					}
				} else {
					if in.ResponseNonce == s.client.LatestNonceSent {
						s.client.LatestVersionACK = in.VersionInfo
						s.client.LatestNonceACK = in.ResponseNonce
						s.logger.Info("received ACK on stream", in)
					} else {
						s.logger.Info("received NACK on stream", in.ErrorDetail)
						nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
						message := Message{
							NodeID:      s.client.NodeID,
							VersionInfo: s.client.LatestVersionSent,
							Nonce:       nonce,
						}
						s.client.LatestNonceSent = nonce
						s.messageChan.Push(message)
						s.client.LastUpdated = time.Now()
					}
				}
			}
		}
	}()

	go func() {
		for e := range s.messageChan {
			if err := s.streamResponses(e); err != nil {
				s.logger.Debug("error while streaming response", err)
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

	// When using ADS we need to order responses.
	// https://www.envoyproxy.io/docs/envoy/latest/api-docs/xds_protocol#eventual-consistency-considerations
	responseStream := NewResponseStream(s.stream, message.VersionInfo, message.Nonce)
	if err := responseStream.StreamCDS(cfg.Clusters); err != nil {
		return err
	}
	if err := responseStream.StreamLDS(cfg.Listeners); err != nil {
		return err
	}
	if err := responseStream.StreamRDS(cfg.Routes); err != nil {
		return err
	}

	return nil
}

func (s Stream) PushUpdatePeriodically() {
	ticker := time.NewTicker(s.refreshInterval)
	defer ticker.Stop()

	service, ok := s.services[s.client.NodeID]
	if !ok {
		s.logger.Debug("service not found for node id", s.client.NodeID)
		return
	}

	for {
		select {
		case <-ticker.C:
			if service.IsUpdated(s.ctx, s.client.LastUpdated) {
				s.logger.Debug("discovery resource update found")
				versionInfo := strconv.FormatInt(time.Now().UnixNano(), 10)
				nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
				message := Message{
					NodeID:      s.client.NodeID,
					VersionInfo: versionInfo,
					Nonce:       nonce,
				}
				s.messageChan.Push(message)
				s.client.LatestVersionSent = versionInfo
				s.client.LatestNonceSent = nonce
				s.client.LastUpdated = time.Now()
			} else {
				s.logger.Debug("no discovery resource update")
			}
		case <-s.ctx.Done():
			return
		}
	}
}
