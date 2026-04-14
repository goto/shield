package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/goto/shield/core/rule"
	"github.com/goto/shield/pkg/httputil"
)

func EnrichRule(req *http.Request, r *rule.Rule) {
	*req = *req.WithContext(rule.WithContext(req.Context(), r))
}

func EnrichRequestBody(r *http.Request) error {
	// Only buffer JSON and gRPC request bodies. Multipart/form-data (file
	// uploads) and other content types are left untouched so that large
	// payloads are streamed directly to the backend without being copied
	// into memory — preventing OOM on files of 200 MB+.
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") && !strings.HasPrefix(ct, "application/grpc") {
		return nil
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer (r.Body).Close()

	// repopulate body
	(*r).Body = io.NopCloser(bytes.NewBuffer(reqBody))
	*r = *r.WithContext(httputil.SetContextWithRequestBody(r.Context(), reqBody))
	return nil
}

func ExtractRequestBody(r *http.Request) (io.ReadCloser, bool) {
	body, ok := httputil.GetRequestBodyFromContext(r.Context())
	if !ok {
		return nil, false
	}
	return io.NopCloser(bytes.NewBuffer(body)), true
}

func ExtractRule(r *http.Request) (*rule.Rule, bool) {
	return rule.GetFromContext(r.Context())
}

func ExtractMiddleware(r *http.Request, name string) (rule.MiddlewareSpec, bool) {
	rl, ok := ExtractRule(r)
	if !ok {
		return rule.MiddlewareSpec{}, false
	}
	return rl.Middlewares.Get(name)
}

func EnrichRequestWithMuxRoute(r *http.Request, route *mux.Route) {
	*r = *r.WithContext(httputil.SetContextWithMuxRoute(r.Context(), route))
}

func EnrichPathParams(r *http.Request, params map[string]string) {
	*r = *r.WithContext(httputil.SetContextWithPathParams(r.Context(), params))
}

func ExtractPathParams(r *http.Request) (map[string]string, bool) {
	return httputil.GetPathParamsFromContext(r.Context())
}
