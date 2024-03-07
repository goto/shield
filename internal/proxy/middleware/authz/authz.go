package authz

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/goto/salt/log"
	"github.com/mitchellh/mapstructure"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/proxy/middleware"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/body_extractor"
	"github.com/goto/shield/pkg/expression"
	"github.com/goto/shield/pkg/uuid"
)

type ResourceService interface {
	CheckAuthz(ctx context.Context, resource resource.Resource, act action.Action) (bool, error)
}

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type GroupService interface {
	GetBySlug(ctx context.Context, slug string) (group.Group, error)
}

type Authz struct {
	log             log.Logger
	userIDHeaderKey string
	next            http.Handler
	resourceService ResourceService
	userService     UserService
	groupService    GroupService
}

type Config struct {
	Actions     []string                        `yaml:"actions" mapstructure:"actions"`
	Permissions []Permission                    `yaml:"permissions" mapstructure:"permissions"`
	Attributes  map[string]middleware.Attribute `yaml:"attributes" mapstructure:"attributes"`
}

type Permission struct {
	Name       string                `yaml:"name" mapstructure:"name"`
	Namespace  string                `yaml:"namespace" mapstructure:"namespace"`
	Attribute  string                `yaml:"attribute" mapstructure:"attribute"`
	Expression expression.Expression `yaml:"expression" mapstructure:"expression"`
}

func New(
	log log.Logger,
	next http.Handler,
	userIDHeaderKey string,
	resourceService ResourceService,
	userService UserService,
	groupService GroupService) *Authz {
	return &Authz{
		log:             log,
		userIDHeaderKey: userIDHeaderKey,
		next:            next,
		resourceService: resourceService,
		userService:     userService,
		groupService:    groupService,
	}
}

func (c Authz) Info() *middleware.MiddlewareInfo {
	return &middleware.MiddlewareInfo{
		Name:        "authz",
		Description: "rule based authorization using casbin",
	}
}

func (c *Authz) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	usr, err := c.userService.FetchCurrentUser(req.Context())
	if err != nil {
		c.log.Error("middleware: failed to get user details", "err", err.Error())
		c.notAllowed(rw, nil)
		return
	}

	req.Header.Set(c.userIDHeaderKey, usr.ID)

	rule, ok := middleware.ExtractRule(req)
	if !ok {
		c.next.ServeHTTP(rw, req)
		return
	}

	wareSpec, ok := middleware.ExtractMiddleware(req, c.Info().Name)
	if !ok {
		c.next.ServeHTTP(rw, req)
		return
	}

	if rule.Backend.Namespace == "" {
		c.log.Error("namespace is not defined for this rule")
		c.notAllowed(rw, nil)
		return
	}

	// TODO: should cache it
	config := Config{}
	if err := mapstructure.Decode(wareSpec.Config, &config); err != nil {
		c.log.Error("middleware: failed to decode authz config", "config", wareSpec.Config)
		c.notAllowed(rw, nil)
		return
	}

	if valid, err := config.validate(); !valid {
		c.log.Error("middleware", c.Info().Name, "rule", rule.Frontend.URLRx, "backend", rule.Backend.Namespace, "err", err)
		c.notAllowed(rw, nil)
		return
	}

	permissionAttributes := map[string]interface{}{}

	permissionAttributes["namespace"] = rule.Backend.Namespace

	permissionAttributes["user"] = req.Header.Get(c.userIDHeaderKey)

	for res, attr := range config.Attributes {
		_ = res

		switch attr.Type {
		case middleware.AttributeTypeGRPCPayload:
			// check if grpc request
			if !strings.HasPrefix(req.Header.Get("Content-Type"), "application/grpc") {
				c.log.Error("middleware: not a grpc request", "attr", attr)
				c.notAllowed(rw, nil)
				return
			}

			// TODO: we can optimise this by parsing all field at once
			payloadField, err := body_extractor.GRPCPayloadHandler{}.Extract(&req.Body, attr.Index)
			if err != nil {
				c.log.Error("middleware: failed to parse grpc payload", "err", err)
				return
			}

			permissionAttributes[res] = payloadField
			c.log.Info("middleware: extracted", "field", payloadField, "attr", attr)

		case middleware.AttributeTypeJSONPayload:
			if attr.Key == "" {
				c.log.Error("middleware: payload key field empty")
				c.notAllowed(rw, nil)
				return
			}
			payloadField, err := body_extractor.JSONPayloadHandler{}.Extract(&req.Body, attr.Key)
			if err != nil {
				c.log.Error("middleware: failed to parse grpc payload", "err", err)
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = payloadField
			c.log.Info("middleware: extracted", "field", payloadField, "attr", attr)

		case middleware.AttributeTypeHeader:
			if attr.Key == "" {
				c.log.Error("middleware: header key field empty")
				c.notAllowed(rw, nil)
				return
			}
			headerAttr := req.Header.Get(attr.Key)
			if headerAttr == "" {
				c.log.Error(fmt.Sprintf("middleware: header %s is empty", attr.Key))
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = headerAttr
			c.log.Info("middleware: extracted", "field", headerAttr, "attr", attr)

		case middleware.AttributeTypeQuery:
			if attr.Key == "" {
				c.log.Error("middleware: query key field empty")
				c.notAllowed(rw, nil)
				return
			}
			queryAttr := req.URL.Query().Get(attr.Key)
			if queryAttr == "" {
				c.log.Error(fmt.Sprintf("middleware: query %s is empty", attr.Key))
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = queryAttr
			c.log.Info("middleware: extracted", "field", queryAttr, "attr", attr)

		case middleware.AttributeTypeConstant:
			if attr.Value == "" {
				c.log.Error("middleware: constant value empty")
				c.notAllowed(rw, nil)
				return
			}

			permissionAttributes[res] = attr.Value
			c.log.Info("middleware: extracted", "constant_key", res, "attr", permissionAttributes[res])

		default:
			c.log.Error("middleware: unknown attribute type", "attr", attr)
			c.notAllowed(rw, nil)
			return
		}
	}

	paramMap, mapExists := middleware.ExtractPathParams(req)
	if !mapExists {
		c.log.Error("middleware: path param map doesn't exist")
		c.notAllowed(rw, nil)
		return
	}

	for key, value := range paramMap {
		permissionAttributes[key] = value
	}

	isAuthorized := true
	for _, permission := range config.Permissions {
		c.log.Info("checking permission", "permission", permission.Name)
		if !permission.Expression.IsEmpty() {
			permission.Expression = enrichExpression(permission.Expression, permissionAttributes)
			c.log.Info("evaluating expression", "expr", permission.Expression)
			output, err := permission.Expression.Evaluate()
			if err != nil {
				c.log.Error("error evaluating expression", "err", err)
				continue
			}
			c.log.Info("successfully evaluated expression", "result", output)

			if output == false {
				continue
			}
		}

		res, err := c.preparePermissionResource(req.Context(), permission, permissionAttributes)
		if err != nil {
			c.log.Error("error while preparing permission resource", "err", err)
			c.notAllowed(rw, err)
			return
		}
		isAuthorized, err = c.resourceService.CheckAuthz(req.Context(), res, action.Action{
			ID: permission.Name,
		})
		if err != nil {
			c.log.Error("error while performing authz permission check", "err", err)
			c.notAllowed(rw, err)
			return
		}
		if isAuthorized {
			break
		}
	}

	c.log.Info("authz check successful", "user", permissionAttributes["user"], "resource", permissionAttributes["resource"], "result", isAuthorized)
	if !isAuthorized {
		c.log.Info("user not allowed to make request", "user", permissionAttributes["user"], "resource", permissionAttributes["resource"], "result", isAuthorized)
		c.notAllowed(rw, nil)
		return
	}

	c.next.ServeHTTP(rw, req)
}

func (c Authz) preparePermissionResource(ctx context.Context, perm Permission, attrs map[string]interface{}) (resource.Resource, error) {
	resourceName := attrs[perm.Attribute].(string)
	res := resource.Resource{
		Name:        resourceName,
		NamespaceID: perm.Namespace,
	}

	if perm.Namespace == schema.GroupNamespace {
		// resolve group id from slug
		if !uuid.IsValid(resourceName) {
			grp, err := c.groupService.GetBySlug(ctx, resourceName)
			if err != nil {
				return resource.Resource{}, err
			}
			res.Name = grp.ID
		}
	}
	return res, nil
}

func (w Authz) notAllowed(rw http.ResponseWriter, err error) {
	if err != nil {
		switch {
		case errors.Is(err, resource.ErrNotExist):
			rw.WriteHeader(http.StatusNotFound)
			return
		}
	}
	rw.WriteHeader(http.StatusUnauthorized)
}

func (cg Config) validate() (bool, error) {
	if len(cg.Permissions) == 0 {
		return false, errors.New("no permissions configured")
	}

	return true, nil
}

func enrichExpression(exp expression.Expression, attributes map[string]interface{}) expression.Expression {
	if val, ok := attributes[exp.Attribute.(string)]; ok {
		exp.Attribute = val
	}
	return exp
}
