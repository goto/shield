package otelpostprocessor_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/shield/internal/proxy/middleware"
	"github.com/goto/shield/internal/proxy/middleware/otelpostprocessor"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TestWare_ServeHTTP(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		labels, exist := otelhttp.LabelerFromContext(r.Context())
		assert.NotNil(t, labels)
		assert.True(t, exist)
	})

	op := otelpostprocessor.New(next)
	o := otelhttp.NewHandler(op, "")

	assert.Equal(t, &middleware.MiddlewareInfo{
		Name:        "_otel_postprocessor",
		Description: "post process otel metrics and traces",
	}, op.Info())

	req := httptest.NewRequest("GET", "http://testing", nil)
	o.ServeHTTP(httptest.NewRecorder(), req)
}
