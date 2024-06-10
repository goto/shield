package otelmiddleware

import (
	"net/http"

	"github.com/goto/shield/internal/proxy/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type Ware struct {
	otelHandler http.Handler
}

func New(next http.Handler, opts ...otelhttp.Option) *Ware {
	return &Ware{
		otelHandler: otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := r.URL.Path
			attr := semconv.HTTPRouteKey.String(route)

			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attr)

			labeler, _ := otelhttp.LabelerFromContext(r.Context())
			labeler.Add(attr)

			next.ServeHTTP(w, r)
		}), "", opts...),
	}
}

func (m Ware) Info() *middleware.MiddlewareInfo {
	return &middleware.MiddlewareInfo{
		Name:        "_opentelemetry",
		Description: "handling opentelemetry middleware",
	}
}

func (m *Ware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	m.otelHandler.ServeHTTP(rw, req)
}
