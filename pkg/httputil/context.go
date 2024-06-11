package httputil

import (
	"context"

	"github.com/gorilla/mux"
)

type (
	contextRequestBodyKey struct{}
	contextPathParamsKey  struct{}
	routeKey              struct{}
	varsKey               struct{}
)

func SetContextWithRequestBody(ctx context.Context, body []byte) context.Context {
	return context.WithValue(ctx, contextRequestBodyKey{}, body)
}

func GetRequestBodyFromContext(ctx context.Context) ([]byte, bool) {
	body, ok := ctx.Value(contextRequestBodyKey{}).([]byte)
	return body, ok
}

func SetContextWithPathParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, contextPathParamsKey{}, params)
}

func GetPathParamsFromContext(ctx context.Context) (map[string]string, bool) {
	params, ok := ctx.Value(contextPathParamsKey{}).(map[string]string)
	return params, ok
}

func SetContextWithMuxRouteAndVars(ctx context.Context, route *mux.Route, vars map[string]string) context.Context {
	ctx = context.WithValue(ctx, routeKey{}, route)
	if len(vars) > 0 {
		ctx = context.WithValue(ctx, varsKey{}, vars)
	}
	return ctx
}
