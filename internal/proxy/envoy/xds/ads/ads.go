package ads

import (
	"context"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/goto/shield/core/rule"
)

const (
	HTTP_CONNECTION_MANAGER_TYPE_URL = resource.APITypePrefix + "envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
	ROUTER_TYPE_URL                  = resource.APITypePrefix + "envoy.extensions.filters.http.router.v3.Router"
	URI_TEMPLATE_TYPE_URL            = resource.APITypePrefix + "envoy.extensions.path.match.uri_template.v3.UriTemplateMatchConfig"
	STDOUT_LOGGER_TYPE_URL           = resource.APITypePrefix + "envoy.extensions.access_loggers.stream.v3.StdoutAccessLog"
)

type Repository interface {
	Fetch(ctx context.Context) ([]rule.Ruleset, error)
	IsUpdated(ctx context.Context, since time.Time) bool
}
