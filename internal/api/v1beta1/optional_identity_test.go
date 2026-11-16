package v1beta1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/role"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
)

func TestGatewayRoutesWithoutIdentityHeader(t *testing.T) {
	for _, backendErr := range []error{nil, context.DeadlineExceeded} {
		for _, path := range []string{
			"/v1beta1/users/target/groups",
			"/v1beta1/groups/target/relations",
			"/v1beta1/organizations",
			"/v1beta1/roles",
			"/v1beta1/namespaces",
		} {
			name := "success"
			want := http.StatusOK
			if backendErr != nil {
				name = "database timeout"
				want = http.StatusInternalServerError
			}
			t.Run(path+"/"+name, func(t *testing.T) {
				// Any unexpected caller lookup fails the test; these routes use
				// path parameters or list data without a caller identity.
				us := mocks.NewUserService(t)
				h := Handler{userService: us}
				switch path {
				case "/v1beta1/users/target/groups":
					gs := mocks.NewGroupService(t)
					gs.EXPECT().ListUserGroups(mock.Anything, "target", "").Return([]group.Group{}, backendErr).Once()
					h.groupService = gs
				case "/v1beta1/groups/target/relations":
					gs := mocks.NewGroupService(t)
					gs.EXPECT().ListGroupRelations(mock.Anything, "target", "", "").Return([]user.User{}, []group.Group{}, map[string][]string{}, map[string][]string{}, backendErr).Once()
					h.groupService = gs
				case "/v1beta1/organizations":
					os := mocks.NewOrganizationService(t)
					os.EXPECT().List(mock.Anything).Return([]organization.Organization{}, backendErr).Once()
					h.orgService = os
				case "/v1beta1/roles":
					rs := mocks.NewRoleService(t)
					rs.EXPECT().List(mock.Anything).Return([]role.Role{}, backendErr).Once()
					h.roleService = rs
				case "/v1beta1/namespaces":
					ns := mocks.NewNamespaceService(t)
					ns.EXPECT().List(mock.Anything).Return([]namespace.Namespace{}, backendErr).Once()
					h.namespaceService = ns
				}
				mux := newErrorTestGateway()
				require.NoError(t, shieldv1beta1.RegisterShieldServiceHandlerServer(context.Background(), mux, h))
				req := httptest.NewRequest(http.MethodGet, path, nil)
				require.Empty(t, req.Header.Get("X-Auth-Email"))
				rw := httptest.NewRecorder()
				mux.ServeHTTP(rw, req)
				require.Equal(t, want, rw.Code, rw.Body.String())
				require.NotContains(t, rw.Body.String(), "deadline")
				us.AssertNotCalled(t, "FetchCurrentUser", mock.Anything)
			})
		}
	}
}

func TestGetUserStillRequiresIdentity(t *testing.T) {
	us := mocks.NewUserService(t)
	us.EXPECT().FetchCurrentUser(mock.Anything).RunAndReturn(func(ctx context.Context) (user.User, error) {
		// The actual service rejects the absent identity before touching its repository.
		return (user.Service{}).FetchCurrentUser(ctx)
	}).Once()
	mux := newErrorTestGateway()
	require.NoError(t, shieldv1beta1.RegisterShieldServiceHandlerServer(context.Background(), mux, Handler{userService: us}))
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, httptest.NewRequest(http.MethodGet, "/v1beta1/users/target@example.invalid", nil))
	require.Equal(t, http.StatusUnauthorized, rw.Code)
}

func TestUserPermissionCheckUsesPathIdentityWithoutHeader(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"success", nil, http.StatusOK},
		{"unknown target identity", user.ErrInvalidEmail, http.StatusUnauthorized},
		{"backend timeout", context.DeadlineExceeded, http.StatusInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rs := mocks.NewResourceService(t)
			rs.EXPECT().BulkCheckAuthz(mock.MatchedBy(func(ctx context.Context) bool {
				email, ok := user.GetEmailFromContext(ctx)
				return ok && email == "target@example.invalid"
			}), mock.Anything, mock.Anything).Return([]relation.Permission{{Allowed: true}}, tc.err).Once()
			mux := newErrorTestGateway()
			require.NoError(t, shieldv1beta1.RegisterShieldServiceHandlerServer(context.Background(), mux, Handler{resourceService: rs, checkAPILimit: 10}))
			body := `{"resource_permissions":[{"object_id":"target","object_namespace":"test","permission":"view"}]}`
			req := httptest.NewRequest(http.MethodPost, "/v1beta1/users/target@example.invalid/check", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			require.Empty(t, req.Header.Get("X-Auth-Email"))
			rw := httptest.NewRecorder()
			mux.ServeHTTP(rw, req)
			require.Equal(t, tc.status, rw.Code, rw.Body.String())
		})
	}
}
