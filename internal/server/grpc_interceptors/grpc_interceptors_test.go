package grpc_interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/goto/shield/core/user"
)

func TestIdentityHeaderIsOptionalInInterceptor(t *testing.T) {
	for _, email := range []string{"", "caller@example.invalid"} {
		t.Run(email, func(t *testing.T) {
			md := metadata.MD{}
			if email != "" {
				md.Set("x-auth-email", email)
			}
			ctx := metadata.NewIncomingContext(context.Background(), md)
			called := false
			_, err := EnrichCtxWithIdentity("X-Auth-Email")(ctx, nil, &grpc.UnaryServerInfo{}, func(ctx context.Context, _ interface{}) (interface{}, error) {
				called = true
				got, ok := user.GetEmailFromContext(ctx)
				require.True(t, ok)
				require.Equal(t, email, got)
				return nil, nil
			})
			require.NoError(t, err)
			require.True(t, called)
		})
	}
}
