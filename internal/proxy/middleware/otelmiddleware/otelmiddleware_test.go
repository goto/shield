package otelmiddleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/shield/internal/proxy/middleware/otelmiddleware"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func TestWare_ServeHTTP(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sp := trace.SpanFromContext(r.Context())
		lb, exist := otelhttp.LabelerFromContext(r.Context())
		assert.NotNil(t, sp)
		assert.True(t, exist)
		assert.NotNil(t, lb)
	})

	m := otelmiddleware.New(next)
	req := httptest.NewRequest("GET", "http://testing", nil)
	m.ServeHTTP(httptest.NewRecorder(), req)
}
