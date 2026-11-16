package authz

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/salt/log"
	"github.com/stretchr/testify/require"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/proxy/middleware"
)

type userServiceFunc func(context.Context) (user.User, error)

func (f userServiceFunc) FetchCurrentUser(ctx context.Context) (user.User, error) { return f(ctx) }

type resourceServiceFunc func(context.Context, resource.Resource, action.Action) (bool, error)

func (f resourceServiceFunc) CheckAuthz(ctx context.Context, res resource.Resource, act action.Action) (bool, error) {
	return f(ctx, res, act)
}

func TestServeHTTPCallerLookupFailure(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"missing identity", user.ErrMissingEmail, http.StatusUnauthorized},
		{"invalid identity", fmt.Errorf("caller: %w", user.ErrInvalidEmail), http.StatusUnauthorized},
		{"deadline", fmt.Errorf("query: %w", context.DeadlineExceeded), http.StatusInternalServerError},
		{"repository timeout", errors.New("error while running query: timeout: context deadline exceeded"), http.StatusInternalServerError},
		{"connection failure", errors.New("database unavailable"), http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			forwarded := false
			h := New(log.NewNoop(), http.HandlerFunc(func(http.ResponseWriter, *http.Request) { forwarded = true }), "X-User-ID", nil,
				userServiceFunc(func(context.Context) (user.User, error) { return user.User{}, tc.err }), nil)
			rw := httptest.NewRecorder()
			h.ServeHTTP(rw, httptest.NewRequest(http.MethodGet, "/", nil))
			require.Equal(t, tc.status, rw.Code)
			require.False(t, forwarded)
			require.Empty(t, rw.Body.String())
		})
	}
}

func TestServeHTTPPermissionCheckResult(t *testing.T) {
	for _, tc := range []struct {
		name    string
		allowed bool
		err     error
		status  int
	}{
		{"allowed", true, nil, http.StatusNoContent},
		{"denied", false, nil, http.StatusUnauthorized},
		{"missing resource", false, fmt.Errorf("resource: %w", resource.ErrNotExist), http.StatusNotFound},
		{"backend timeout", false, context.DeadlineExceeded, http.StatusInternalServerError},
		{"backend failure", false, errors.New("authorization backend unavailable"), http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			forwarded := false
			checked := false
			h := New(log.NewNoop(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				forwarded = true
				require.Equal(t, "caller-id", r.Header.Get("X-User-ID"))
				w.WriteHeader(http.StatusNoContent)
			}), "X-User-ID", resourceServiceFunc(func(context.Context, resource.Resource, action.Action) (bool, error) {
				checked = true
				return tc.allowed, tc.err
			}), userServiceFunc(func(context.Context) (user.User, error) { return user.User{ID: "caller-id"}, nil }), nil)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			middleware.EnrichRule(req, &rule.Rule{
				Backend: rule.Backend{Namespace: "test"},
				Middlewares: rule.MiddlewareSpecs{{Name: "authz", Config: map[string]interface{}{
					"permissions": []map[string]interface{}{{"name": "view", "namespace": "test", "attribute": "resource"}},
				}}},
			})
			middleware.EnrichPathParams(req, map[string]string{"resource": "target"})
			rw := httptest.NewRecorder()
			h.ServeHTTP(rw, req)
			require.True(t, checked)
			require.Equal(t, tc.status, rw.Code)
			require.Equal(t, tc.allowed, forwarded)
		})
	}
}
