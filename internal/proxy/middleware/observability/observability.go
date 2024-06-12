package observability

import (
	"net/http"
	"strings"

	"github.com/goto/shield/internal/proxy/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/goto/salt/log"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

const (
	headerRequestID = "X-Request-Id"
)

type Ware struct {
	log         *log.Zap
	otelHandler http.Handler
}

func New(log *log.Zap, next http.Handler) *Ware {
	return &Ware{
		log: log,
		otelHandler: otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			attr := semconv.HTTPRouteKey.String(r.URL.Path)

			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attr)

			labeler, _ := otelhttp.LabelerFromContext(r.Context())
			labeler.Add(attr)

			next.ServeHTTP(w, r)
		}), ""),
	}
}

func (m Ware) Info() *middleware.MiddlewareInfo {
	return &middleware.MiddlewareInfo{
		Name:        "_observability",
		Description: "to handle observability",
	}
}

func (m *Ware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqID := setRequestID(req)
	ctx := m.log.NewContext(req.Context())
	ctx = log.ZapContextWithFields(ctx,
		zap.String("host", req.Host),
		zap.String("path", req.URL.String()),
		zap.String("method", req.Method),
		zap.String("remote_address", req.RemoteAddr),
		zap.String("scheme", req.Proto),
		zap.String("request_id", reqID),
	)
	req = req.WithContext(ctx)

	m.otelHandler.ServeHTTP(rw, req)
}

func setRequestID(req *http.Request) string {
	reqID := strings.TrimSpace(req.Header.Get(headerRequestID))
	if reqID == "" {
		reqID = xid.New().String()
		req.Header.Set(headerRequestID, reqID)
	}

	return reqID
}
