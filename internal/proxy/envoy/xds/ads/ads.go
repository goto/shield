package ads

import (
	"context"
	"time"

	"github.com/goto/shield/core/rule"
)

const (
	CLUSTER_TYPE_URL             = "type.googleapis.com/envoy.config.cluster.v3.Cluster"
	LISTENER_TYPE_URL            = "type.googleapis.com/envoy.config.listener.v3.Listener"
	ROUTE_CONFIGURATION_TYPE_URL = "type.googleapis.com/envoy.config.route.v3.RouteConfiguration"

	HTTP_CONNECTION_MANAGER_TYPE_URL = "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
	ROUTER_TYPE_URL                  = "type.googleapis.com/envoy.extensions.filters.http.router.v3.Router"
	URI_TEMPLATE_TYPE_URL            = "type.googleapis.com/envoy.extensions.path.match.uri_template.v3.UriTemplateMatchConfig"
	STDOUT_LOGGER_TYPE_URL           = "type.googleapis.com/envoy.extensions.access_loggers.stream.v3.StdoutAccessLog"
)

type Repository interface {
	GetAll(ctx context.Context) ([]rule.Ruleset, error)
	IsUpdated(ctx context.Context, since time.Time) bool
}
