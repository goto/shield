package rulematch_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/internal/proxy/middleware"
	"github.com/goto/shield/internal/proxy/middleware/rulematch"
	"github.com/stretchr/testify/assert"
)

type mockRuleMatcher struct {
	rule *rule.Rule
	err  error
}

func (m *mockRuleMatcher) Match(_ *http.Request) (*rule.Rule, error) {
	return m.rule, m.err
}

func TestWare_Info(t *testing.T) {
	logger := log.NewZap()
	w := rulematch.New(logger, nil, nil)
	info := w.Info()
	assert.Equal(t, "_rulematch", info.Name)
}

func requestWithLogger(method, url string, body []byte) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	logger := log.NewZap()
	ctx := logger.NewContext(req.Context())
	return req.WithContext(ctx)
}

func TestWare_ServeHTTP_SkipReadBody(t *testing.T) {
	t.Run("should skip body buffering when SkipReadBody is true", func(t *testing.T) {
		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			// body should NOT be in context (EnrichRequestBody was skipped)
			_, ok := middleware.ExtractRequestBody(r)
			assert.False(t, ok, "request body should not be in context when SkipReadBody=true")
			w.WriteHeader(http.StatusOK)
		})

		matcher := &mockRuleMatcher{
			rule: &rule.Rule{
				Frontend: rule.Frontend{
					URL:    "/upload",
					Method: "POST",
				},
				SkipReadBody: true,
			},
		}

		logger := log.NewZap()
		ware := rulematch.New(logger, next, matcher)

		req := requestWithLogger("POST", "http://localhost/upload", []byte("large file content"))
		rr := httptest.NewRecorder()

		ware.ServeHTTP(rr, req)
		assert.True(t, nextCalled, "next handler should have been called")
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should buffer body when SkipReadBody is false", func(t *testing.T) {
		payload := []byte(`{"key":"value"}`)
		nextCalled := false
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nextCalled = true
			// body SHOULD be in context (EnrichRequestBody was called)
			body, ok := middleware.ExtractRequestBody(r)
			assert.True(t, ok, "request body should be in context when SkipReadBody=false")
			if ok {
				buf := new(bytes.Buffer)
				buf.ReadFrom(body)
				assert.Equal(t, string(payload), buf.String())
			}
			w.WriteHeader(http.StatusOK)
		})

		matcher := &mockRuleMatcher{
			rule: &rule.Rule{
				Frontend: rule.Frontend{
					URL:    "/api/resource",
					Method: "POST",
				},
				SkipReadBody: false,
			},
		}

		logger := log.NewZap()
		ware := rulematch.New(logger, next, matcher)

		req := requestWithLogger("POST", "http://localhost/api/resource", payload)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		ware.ServeHTTP(rr, req)
		assert.True(t, nextCalled, "next handler should have been called")
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("should return 400 when rule matching fails", func(t *testing.T) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Fatal("next handler should not have been called")
		})

		matcher := &mockRuleMatcher{
			rule: nil,
			err:  assert.AnError,
		}

		logger := log.NewZap()
		ware := rulematch.New(logger, next, matcher)

		req := requestWithLogger("GET", "http://localhost/unknown", nil)
		rr := httptest.NewRecorder()

		ware.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
