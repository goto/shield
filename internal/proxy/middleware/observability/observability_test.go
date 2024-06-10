package observability_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/salt/log"
	"github.com/goto/shield/internal/proxy/middleware"
	"github.com/goto/shield/internal/proxy/middleware/observability"
	"github.com/stretchr/testify/assert"
)

func TestWare_ServeHTTP(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zp := log.ZapFromContext(r.Context())
		assert.NotNil(t, zp)
	})

	logger := log.NewZap()
	o := observability.New(logger, next)

	assert.Equal(t, &middleware.MiddlewareInfo{
		Name:        "_observability",
		Description: "to handle observability",
	}, o.Info())

	req := httptest.NewRequest("GET", "http://testing", nil)
	o.ServeHTTP(httptest.NewRecorder(), req)
}
