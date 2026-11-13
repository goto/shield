package authz

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/proxy/middleware"
)

type defaultTestUser struct{}

func (defaultTestUser) FetchCurrentUser(context.Context) (user.User, error) {
	return user.User{ID: "synthetic-actor"}, nil
}

type defaultTestChecker struct {
	allowed bool
	calls   int
}

func (c *defaultTestChecker) CheckAuthz(context.Context, resource.Resource, action.Action) (bool, error) {
	c.calls++
	return c.allowed, nil
}
func permissionWith(operator, value string) map[string]interface{} {
	p := map[string]interface{}{"name": "view", "namespace": "test/item", "attribute": "resource"}
	if operator != "" {
		p["expression"] = map[string]interface{}{"attribute": "selector", "operator": operator, "value": value}
	}
	return p
}
func TestAuthorizationDefaultsAndExistingFallbackSemantics(t *testing.T) {
	for _, tc := range []struct {
		name           string
		perms          []interface{}
		grant, forward bool
		checks         int
	}{
		{"unconditional-grant", []interface{}{permissionWith("", "")}, true, true, 1},
		{"unconditional-denial", []interface{}{permissionWith("", "")}, false, false, 1},
		{"existing-fallback-grant", []interface{}{permissionWith("==", "other"), permissionWith("", "")}, true, true, 1},
		{"existing-fallback-denial", []interface{}{permissionWith("==", "other"), permissionWith("", "")}, false, false, 1},
		{"matching-expression", []interface{}{permissionWith("==", "current")}, true, true, 1},
		{"no-applicable-permission", []interface{}{permissionWith("==", "other")}, true, false, 0},
		{"expression-error-before-fallback", []interface{}{permissionWith("unsupported", "current"), permissionWith("", "")}, true, false, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checker := &defaultTestChecker{allowed: tc.grant}
			forwarded := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { forwarded = true; w.WriteHeader(http.StatusNoContent) })
			gate := New(log.NewNoop(), next, "X-Shield-User-Id", checker, defaultTestUser{}, nil)
			r := httptest.NewRequest(http.MethodGet, "/item", nil)
			cfg := rule.Rule{Backend: rule.Backend{Namespace: "test"}, Middlewares: rule.MiddlewareSpecs{{Name: "authz", Config: map[string]interface{}{"permissions": tc.perms}}}}
			middleware.EnrichRule(r, &cfg)
			middleware.EnrichPathParams(r, map[string]string{"resource": "item", "selector": "current"})
			w := httptest.NewRecorder()
			gate.ServeHTTP(w, r)
			if forwarded != tc.forward || checker.calls != tc.checks {
				t.Fatalf("forwarded=%v checks=%d", forwarded, checker.calls)
			}
			if !tc.forward && w.Code != http.StatusUnauthorized {
				t.Fatalf("status=%d", w.Code)
			}
		})
	}
}
