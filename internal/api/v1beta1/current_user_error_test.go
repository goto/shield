package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/api/v1beta1/mocks"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
)

func TestCallerLookupFailures(t *testing.T) {
	ctx := user.SetContextWithEmail(context.Background(), "caller@example.invalid")
	handlers := map[string]func(Handler) error{
		"GetUser":    func(h Handler) error { _, err := h.GetUser(ctx, &shieldv1beta1.GetUserRequest{}); return err },
		"ListUsers":  func(h Handler) error { _, err := h.ListUsers(ctx, &shieldv1beta1.ListUsersRequest{}); return err },
		"UpdateUser": func(h Handler) error { _, err := h.UpdateUser(ctx, &shieldv1beta1.UpdateUserRequest{}); return err },
		"DeleteUser": func(h Handler) error { _, err := h.DeleteUser(ctx, &shieldv1beta1.DeleteUserRequest{}); return err },
		"GetGroup":   func(h Handler) error { _, err := h.GetGroup(ctx, &shieldv1beta1.GetGroupRequest{}); return err },
		"ListGroups": func(h Handler) error { _, err := h.ListGroups(ctx, &shieldv1beta1.ListGroupsRequest{}); return err },
		"CreateGroup": func(h Handler) error {
			_, err := h.CreateGroup(ctx, &shieldv1beta1.CreateGroupRequest{Body: &shieldv1beta1.GroupRequestBody{}})
			return err
		},
		"UpdateGroup": func(h Handler) error {
			_, err := h.UpdateGroup(ctx, &shieldv1beta1.UpdateGroupRequest{Body: &shieldv1beta1.GroupRequestBody{}})
			return err
		},
		"CreateUser": func(h Handler) error {
			_, err := h.CreateUser(ctx, &shieldv1beta1.CreateUserRequest{Body: &shieldv1beta1.UserRequestBody{Email: "other@example.invalid"}})
			return err
		},
	}
	cases := []struct {
		name string
		err  error
		code codes.Code
	}{
		{"missing identity", user.ErrMissingEmail, codes.Unauthenticated},
		{"invalid identity", fmt.Errorf("caller: %w", user.ErrInvalidEmail), codes.Unauthenticated},
		{"deadline", fmt.Errorf("query: %w", context.DeadlineExceeded), codes.Internal},
		// PostgreSQL currently retains the DB sentinel but formats the driver error as text.
		{"repository timeout", errors.New("error while running query: timeout: context deadline exceeded"), codes.Internal},
		{"connection failure", errors.New("database unavailable"), codes.Internal},
		{"canceled", context.Canceled, codes.Internal},
	}
	for name, call := range handlers {
		for _, tc := range cases {
			t.Run(name+"/"+tc.name, func(t *testing.T) {
				us := mocks.NewUserService(t)
				us.EXPECT().FetchCurrentUser(ctx).Return(user.User{}, tc.err).Once()
				got := call(Handler{userService: us})
				want := tc.code
				if name == "CreateUser" && tc.name == "invalid identity" {
					// Existing self-onboarding rules reject creation for another identity.
					want = codes.PermissionDenied
				}
				require.Equal(t, want, status.Code(got))
				if want == codes.Internal {
					require.Equal(t, grpcInternalServerError, got)
				}
			})
		}
	}
}

func TestGetUserGatewayCallerLookupFailure(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"timeout", errors.New("error while running query: timeout: context deadline exceeded"), http.StatusInternalServerError},
		{"missing identity", user.ErrMissingEmail, http.StatusUnauthorized},
		{"invalid identity", fmt.Errorf("caller: %w", user.ErrInvalidEmail), http.StatusUnauthorized},
	} {
		t.Run(tc.name, func(t *testing.T) {
			us := mocks.NewUserService(t)
			us.EXPECT().FetchCurrentUser(mock.Anything).Return(user.User{}, tc.err).Once()
			mux := newErrorTestGateway()
			require.NoError(t, shieldv1beta1.RegisterShieldServiceHandlerServer(context.Background(), mux, Handler{userService: us}))
			rw := httptest.NewRecorder()
			mux.ServeHTTP(rw, httptest.NewRequest(http.MethodGet, "/v1beta1/users/target@example.invalid", nil))
			require.Equal(t, tc.status, rw.Code)
			require.NotContains(t, rw.Body.String(), "query")
			require.NotContains(t, rw.Body.String(), "deadline")
		})
	}
}

func newErrorTestGateway() *runtime.ServeMux {
	// Serialize a fresh status so protobuf's lazy initialization does not
	// mutate the package-level errors compared by existing handler tests.
	return runtime.NewServeMux(runtime.WithErrorHandler(func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		s := status.Convert(err)
		runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, status.Error(s.Code(), s.Message()))
	}))
}
