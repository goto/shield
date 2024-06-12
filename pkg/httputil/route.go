package httputil

import (
	"net/http"

	"github.com/gorilla/mux"
)

func MuxRouteName(r *http.Request) string {
	route := mux.CurrentRoute(r)
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
