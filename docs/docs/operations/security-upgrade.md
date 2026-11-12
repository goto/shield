# Compatibility review of security changes

The active patch retains SQL parameter binding, diagnostic redaction and guards for
malformed inputs that previously panicked. It also fixes fail-open authorization
defaults and updates documentation dependencies.

Changes that could disrupt existing callers have been deferred and removed from the
active patch: new authorization policies, tenant/target remapping, cross-user
restrictions, request validation, request-body limits, timeouts, cancellation changes, and
disabling reflection. The original server configuration,
interceptor chain and transports are retained. Concierge changes only the obsolete
integration rule described below. No administrator-group
configuration or new caller grants are required by the retained runtime changes.

## Retained behavior and checks

Activity filters bind JSON keys and values as parameters. Arbitrary JSON keys and
existing pagination behavior remain supported; a PostgreSQL regression test exercises
quoted keys/values and page sizes above 500. Malformed subjects and missing non-string
project attributes return an error rather than panicking. Supported input parsing is
preserved, including the existing interpretation of additional subject delimiters.

Diagnostic redaction changes log content, not request admission or forwarding.
Log-based dashboards or parsers that consume removed sensitive fields need review.

## Deferred work

The deferred authorization-policy gaps remain open. Trusted identity headers and
internal cleartext transport remain the existing deployment contract; this review
does not verify the live ingress boundary.

Resume stricter controls only after checking production caller behavior, resource
identity mappings, required grants and tooling dependencies. Prepare code and rules
together at that time; do not apply an earlier archived patch without re-review.

## SHIELD-13 and SHIELD-14 completion

At the user's request, Basic Auth now requires bcrypt with cost at least 10.
MD5, SHA-1, plaintext, malformed hashes and malformed Basic Auth configuration are
rejected before YAML/JSON rule activation and during request handling. Failed reloads
retain the last accepted cache. Previously added panic guards remain in place.

All 28 checked-in Concierge proxy-rule files pass validation and none configures
Basic Auth. Live-only rules were not inspected. Before deploying, validate the actual
loaded rules; migrate any legacy hashes in the operator's terminal with htpasswd
using bcrypt (options -B -C 10). Do not send passwords or hashes through chat.
This narrowly scoped change does not re-enable the other deferred access policies.

Documentation uses explicit Docusaurus plugins, preserving pages, theme and static
SVG assets while removing unused SVG-to-React tooling and its deepmerge dependency.
The Yarn audit gate now fails on any advisory severity, including low. The updated
lockfile has zero reported advisories; no advisory is ignored or suppressed.

## Additional nonbreaking closure review

The fail-open default is corrected without changing routes or grants: at least one
applicable permission must succeed, and expression errors deny. All 182 checked-in
authz routes have an unconditional check and use only supported equality expressions,
so their existing successful/denied paths are preserved. Check live-only rules for
configuration drift before deploying; tests cover fallback success, denial and errors.

Concierge retains the integration shield_grpc filename as an explicit `rules: []`
retirement record. It replaces the obsolete odpf RPC routed to localhost:80, without
changing the current gotocompany RPC or supported HTTP route. The empty record is
intentional: CLI publishers skip deleted files, and GCS uploads do not delete old
objects. Apply the named empty configuration through the existing publication process
when rollout is authorized. Do not infer deployed removal from a local source edit.

Service-data writes already enforce the authenticated actor's edit permission on the
project-scoped service-data key. Tests verify missing identity, missing grant and
lookup errors never write, while authorized delegated cross-user writes keep working.
The HTTP/gRPC handler maps denials to PermissionDenied. This is a verified backend
control, not a new proxy gate or a claim that deployed key grants have been inspected.
