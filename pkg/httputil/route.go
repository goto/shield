package httputil

import (
	"net/http"
)

func MuxRouteName(r *http.Request) string {
	route := GetMuxRoute(r.Context())
	if nil == route {
		return "NotFoundHandler"
	}
	if n := route.GetName(); n != "" {
		return n
	}
	if n, _ := route.GetPathTemplate(); n != "" {
		return n
	}
	n, _ := route.GetHostTemplate()
	return n
}
