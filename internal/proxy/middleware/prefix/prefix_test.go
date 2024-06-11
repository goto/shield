package prefix_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/proxy/middleware"
	"github.com/goto/shield/internal/proxy/middleware/prefix"
	"github.com/stretchr/testify/assert"
)

func TestWare_Info(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	logger := log.NewZap()
	p := prefix.New(logger, next)

	assert.Equal(t, &middleware.MiddlewareInfo{
		Name:        "_prefix",
		Description: "manipulating prefix middleware",
	}, p.Info())
}

func TestWare_ServeHTTP(t *testing.T) {
	tests := []struct {
		name     string
		next     http.Handler
		req      *http.Request
		withRule bool
	}{
		{
			name: "should trim the prefix if match",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/testing", r.URL.Path)
			}),
			req:      httptest.NewRequest("GET", "http://localhost/abc/testing", nil),
			withRule: true,
		},
		{
			name: "should trim the prefix rawpath if match",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/file one&two", r.URL.Path)
			}),
			req:      httptest.NewRequest("GET", "http://localhost/abc/file%20one%26two", nil),
			withRule: true,
		},
		{
			name: "should not modify path if no rule found",
			next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/abc/testing", r.URL.Path)
			}),
			req: httptest.NewRequest("GET", "http://localhost/abc/testing", nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := log.NewZap()
			w := prefix.New(logger, tt.next)
			var req *http.Request
			if tt.withRule {
				req = tt.req.WithContext(rule.WithContext(tt.req.Context(), &rule.Rule{
					Backend: rule.Backend{
						Prefix: "/abc",
					},
				}))
			} else {
				req = tt.req
			}
			w.ServeHTTP(httptest.NewRecorder(), req)
		})
	}
}
