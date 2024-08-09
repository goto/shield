package proc

import (
	"fmt"
	"io"

	"github.com/goto/salt/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	corepb "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	envoy_service_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
)

type server struct {
	logger log.Logger
}

// var _ envoy_service_auth_v3.AuthorizationServer = &server{}

// New creates a new authorization server.
func NewServer(logger log.Logger) envoy_service_ext_proc_v3.ExternalProcessorServer {
	return &server{logger: logger}
}

func (s *server) Process(srv envoy_service_ext_proc_v3.ExternalProcessor_ProcessServer) error {
	ctx := srv.Context()
	// state := 0
	reqHeaders := []*corepb.HeaderValue{}
	resHeaders := []*corepb.HeaderValue{}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive stream request: %v", err)
		}
		resp := &envoy_service_ext_proc_v3.ProcessingResponse{}
		switch value := req.Request.(type) {
		case *envoy_service_ext_proc_v3.ProcessingRequest_RequestHeaders:
			s.logger.Debug(fmt.Sprintf("Handle (REQ_HEAD): downstream -> ext_proc -> upstream, Headers: %v", value.RequestHeaders.Headers.Headers))
			reqHeaders = value.RequestHeaders.Headers.Headers
			resp = &envoy_service_ext_proc_v3.ProcessingResponse{
				Response: &envoy_service_ext_proc_v3.ProcessingResponse_RequestHeaders{},
				ModeOverride: &v3.ProcessingMode{
					RequestBodyMode: v3.ProcessingMode_BUFFERED,
				},
			}
		case *envoy_service_ext_proc_v3.ProcessingRequest_RequestBody:
			s.logger.Debug(fmt.Sprintf("Handle (REQ_BODY): downstream -> ext_proc -> upstream, Body: %s", string(value.RequestBody.Body)))
			s.logger.Debug(fmt.Sprintf("Cached (REQ_HEAD): downstream -> ext_proc -> upstream, Headers: %v", reqHeaders))
			resp = &envoy_service_ext_proc_v3.ProcessingResponse{
				Response: &envoy_service_ext_proc_v3.ProcessingResponse_RequestBody{},
			}
		case *envoy_service_ext_proc_v3.ProcessingRequest_ResponseHeaders:
			s.logger.Debug(fmt.Sprintf("Handle (RES_HEAD): upstream -> ext_proc -> downstream, Headers: %v", value.ResponseHeaders.Headers.Headers))
			resHeaders = value.ResponseHeaders.Headers.Headers
			resp = &envoy_service_ext_proc_v3.ProcessingResponse{
				Response: &envoy_service_ext_proc_v3.ProcessingResponse_ResponseHeaders{},
				ModeOverride: &v3.ProcessingMode{
					ResponseBodyMode: v3.ProcessingMode_BUFFERED,
				},
			}
		case *envoy_service_ext_proc_v3.ProcessingRequest_ResponseBody:
			s.logger.Debug(fmt.Sprintf("Handle (RES_BODY): upstream -> ext_proc -> downstream, Body: %s", string(value.ResponseBody.Body)))
			s.logger.Debug(fmt.Sprintf("Cached (RES_HEAD): upstream -> ext_proc -> downstream, Headers: %v", resHeaders))
			resp = &envoy_service_ext_proc_v3.ProcessingResponse{
				Response: &envoy_service_ext_proc_v3.ProcessingResponse_ResponseBody{},
			}
		default:
			s.logger.Debug(fmt.Sprintf("Unknown Request type %v\n", value))
		}
		if err := srv.Send(resp); err != nil {
			s.logger.Debug(fmt.Sprintf("send error %v", err))
		}
	}
}
