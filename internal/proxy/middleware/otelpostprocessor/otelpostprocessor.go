package otelpostprocessor

import (
	"net/http"

	"github.com/goto/shield/internal/proxy/middleware"
	"github.com/goto/shield/pkg/httputil"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Ware struct {
	next http.Handler
}

func New(next http.Handler) *Ware {
	return &Ware{
		next: next,
	}
}

func (m Ware) Info() *middleware.MiddlewareInfo {
	return &middleware.MiddlewareInfo{
		Name:        "_otel_postprocessor",
		Description: "post process otel metrics and traces",
	}
}

func (m *Ware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	routeName := httputil.MuxRouteName(r)
	labeler, _ := otelhttp.LabelerFromContext(r.Context())
	labeler.Add(semconv.HTTPRouteKey.String(routeName))

	span := trace.SpanFromContext(r.Context())
	span.SetName(m.SpanName(r, routeName))

	ctx := trace.ContextWithSpan(r.Context(), span)
	r = r.WithContext(ctx)

	m.next.ServeHTTP(rw, r)
}

func (w *Ware) SpanName(r *http.Request, routeName string) string {
	return r.Method + " " + routeName
}
