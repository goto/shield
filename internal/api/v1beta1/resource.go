package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/policy"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/role"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ResourceService interface {
	Get(ctx context.Context, id string) (resource.Resource, error)
	List(ctx context.Context, flt resource.Filter) (resource.PagedResources, error)
	Upsert(ctx context.Context, resource resource.Resource) (resource.Resource, error)
	Update(ctx context.Context, id string, resource resource.Resource) (resource.Resource, error)
	CheckAuthz(ctx context.Context, resource resource.Resource, action action.Action) (bool, error)
	BulkCheckAuthz(ctx context.Context, resources []resource.Resource, actions []action.Action) ([]relation.Permission, error)
	ListUserResourcesByType(ctx context.Context, userID string, resourceType string, permissions []string) (resource.ResourcePermissions, error)
	ListAllUserResources(ctx context.Context, userID string, resourceTypes []string, permissions []string) (map[string]resource.ResourcePermissions, error)
	UpsertResourcesConfig(ctx context.Context, name string, config string) (resource.ResourceConfig, error)
}

var grpcResourceNotFoundErr = status.Errorf(codes.NotFound, "resource doesn't exist")

func (h Handler) ListResources(ctx context.Context, request *shieldv1beta1.ListResourcesRequest) (*shieldv1beta1.ListResourcesResponse, error) {
	logger := grpczap.Extract(ctx)
	var resources []*shieldv1beta1.Resource

	filters := resource.Filter{
		NamespaceID:    request.GetNamespaceId(),
		OrganizationID: request.GetOrganizationId(),
		ProjectID:      request.GetProjectId(),
		GroupID:        request.GetGroupId(),
		Limit:          request.GetPageSize(),
		Page:           request.GetPageNum(),
	}

	resourcesResp, err := h.resourceService.List(ctx, filters)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	for _, r := range resourcesResp.Resources {
		resourcePB, err := transformResourceToPB(r)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}
		resources = append(resources, &resourcePB)
	}

	return &shieldv1beta1.ListResourcesResponse{
		Count:     resourcesResp.Count,
		Resources: resources,
	}, nil
}

func (h Handler) CreateResource(ctx context.Context, request *shieldv1beta1.CreateResourceRequest) (*shieldv1beta1.CreateResourceResponse, error) {
	logger := grpczap.Extract(ctx)
	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	projId := request.GetBody().GetProjectId()
	project, err := h.projectService.Get(ctx, projId)
	if err != nil {
		logger.Error(err.Error())

		switch {
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	newResource, err := h.resourceService.Upsert(ctx, resource.Resource{
		OrganizationID: project.Organization.ID,
		ProjectID:      request.GetBody().GetProjectId(),
		NamespaceID:    request.GetBody().GetNamespaceId(),
		Name:           request.GetBody().GetName(),
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		case errors.Is(err, resource.ErrInvalidUUID),
			errors.Is(err, resource.ErrInvalidDetail):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	relations := request.GetBody().GetRelations()
	for _, r := range relations {
		subject := strings.Split(r.Subject, ":")
		if len(subject) != 2 {
			logger.Error(fmt.Sprintf("inadequate subject format: %s", r.Subject))
			continue
		}

		_, err := h.createRelation(ctx, relation.RelationV2{
			Object: relation.Object{
				ID:          newResource.Idxa,
				NamespaceID: newResource.NamespaceID,
			},
			Subject: relation.Subject{
				RoleID:    r.RoleName,
				ID:        subject[1],
				Namespace: subject[0],
			},
		})
		if err != nil {
			logger.Error(fmt.Sprintf("error creating relation: %s for %s %s", r.RoleName, subject[1], subject[0]))
		} else {
			logger.Info(fmt.Sprintf("created relation: %s for %s %s", r.RoleName, subject[1], subject[0]))
		}
	}

	resourcePB, err := transformResourceToPB(newResource)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.CreateResourceResponse{
		Resource: &resourcePB,
	}, nil
}

func (h Handler) GetResource(ctx context.Context, request *shieldv1beta1.GetResourceRequest) (*shieldv1beta1.GetResourceResponse, error) {
	logger := grpczap.Extract(ctx)

	fetchedResource, err := h.resourceService.Get(ctx, request.GetId())
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, resource.ErrNotExist),
			errors.Is(err, resource.ErrInvalidUUID),
			errors.Is(err, resource.ErrInvalidID):
			return nil, grpcResourceNotFoundErr
		default:
			return nil, grpcInternalServerError
		}
	}

	resourcePB, err := transformResourceToPB(fetchedResource)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.GetResourceResponse{
		Resource: &resourcePB,
	}, nil
}

func (h Handler) UpdateResource(ctx context.Context, request *shieldv1beta1.UpdateResourceRequest) (*shieldv1beta1.UpdateResourceResponse, error) {
	logger := grpczap.Extract(ctx)

	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	projId := request.GetBody().GetProjectId()
	project, err := h.projectService.Get(ctx, projId)
	if err != nil {
		logger.Error(err.Error())

		switch {
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	updatedResource, err := h.resourceService.Update(ctx, request.GetId(), resource.Resource{
		OrganizationID: project.Organization.ID,
		ProjectID:      request.GetBody().GetProjectId(),
		NamespaceID:    request.GetBody().GetNamespaceId(),
		Name:           request.GetBody().GetName(),
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, resource.ErrNotExist),
			errors.Is(err, resource.ErrInvalidUUID),
			errors.Is(err, resource.ErrInvalidID):
			return nil, grpcResourceNotFoundErr
		case errors.Is(err, resource.ErrInvalidDetail),
			errors.Is(err, resource.ErrInvalidURN):
			return nil, grpcBadBodyError
		case errors.Is(err, resource.ErrConflict):
			return nil, grpcConflictError
		default:
			return nil, grpcInternalServerError
		}
	}

	resourcePB, err := transformResourceToPB(updatedResource)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.UpdateResourceResponse{
		Resource: &resourcePB,
	}, nil
}

func (h Handler) ListAllUserResources(ctx context.Context, request *shieldv1beta1.ListAllUserResourcesRequest) (*shieldv1beta1.ListAllUserResourcesResponse, error) {
	logger := grpczap.Extract(ctx)
	resources, err := h.resourceService.ListAllUserResources(ctx, request.UserId, request.Types, request.Permissions)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, resource.ErrNotExist):
			return nil, grpcResourceNotFoundErr
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrNotExist),
			errors.Is(err, resource.ErrInvalidDetail),
			errors.Is(err, resource.ErrInvalidURN),
			errors.Is(err, policy.ErrInvalidDetail),
			errors.Is(err, relation.ErrInvalidDetail):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	result := make(map[string]*structpb.Value)
	for key, value := range resources {
		resourcePB, err := mapToStructpb(value)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}
		result[key] = &structpb.Value{
			Kind: &structpb.Value_StructValue{StructValue: resourcePB},
		}
	}

	resultPB := &structpb.Struct{
		Fields: result,
	}

	return &shieldv1beta1.ListAllUserResourcesResponse{
		Resources: resultPB,
	}, nil
}

func (h Handler) ListUserResourcesByType(ctx context.Context, request *shieldv1beta1.ListUserResourcesByTypeRequest) (*shieldv1beta1.ListUserResourcesByTypeResponse, error) {
	logger := grpczap.Extract(ctx)

	resources, err := h.resourceService.ListUserResourcesByType(ctx, request.UserId, fmt.Sprintf("%s/%s", request.Namespace, request.Type), request.Permissions)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, resource.ErrNotExist):
			return nil, grpcResourceNotFoundErr
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrNotExist),
			errors.Is(err, resource.ErrInvalidDetail),
			errors.Is(err, resource.ErrInvalidURN),
			errors.Is(err, policy.ErrInvalidDetail),
			errors.Is(err, relation.ErrInvalidDetail):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	resourcesPB, err := mapToStructpb(resources)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.ListUserResourcesByTypeResponse{
		Resources: resourcesPB,
	}, nil
}

func (h Handler) UpsertResourcesConfig(ctx context.Context, request *shieldv1beta1.UpsertResourcesConfigRequest) (*shieldv1beta1.UpsertResourcesConfigResponse, error) {
	logger := grpczap.Extract(ctx)

	rc, err := h.resourceService.UpsertResourcesConfig(ctx, request.Name, request.Config)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, resource.ErrUpsertConfigNotSupported):
			return nil, grpcUnsupportedError
		case errors.Is(err, resource.ErrInvalidDetail), errors.Is(err, role.ErrNotExist),
			errors.Is(err, schema.ErrInvalidDetail):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	rc.Config = request.Config
	return resourceConfigToPB(rc), nil
}

func resourceConfigToPB(from resource.ResourceConfig) *shieldv1beta1.UpsertResourcesConfigResponse {
	return &shieldv1beta1.UpsertResourcesConfigResponse{
		Id:        from.ID,
		Name:      from.Name,
		Config:    from.Config,
		CreatedAt: timestamppb.New(from.CreatedAt),
		UpdatedAt: timestamppb.New(from.UpdatedAt),
	}
}

func transformResourceToPB(from resource.Resource) (shieldv1beta1.Resource, error) {
	// TODO(krtkvrm): will be replaced with IDs
	return shieldv1beta1.Resource{
		Id:   from.Idxa,
		Urn:  from.URN,
		Name: from.Name,
		Project: &shieldv1beta1.Project{
			Id: from.ProjectID,
		},
		Organization: &shieldv1beta1.Organization{
			Id: from.OrganizationID,
		},
		Namespace: &shieldv1beta1.Namespace{
			Id: from.NamespaceID,
		},
		User: &shieldv1beta1.User{
			Id: from.UserID,
		},
		CreatedAt: timestamppb.New(from.CreatedAt),
		UpdatedAt: timestamppb.New(from.UpdatedAt),
	}, nil
}

func mapToStructpb(p resource.ResourcePermissions) (*structpb.Struct, error) {
	fields := make(map[string]*structpb.Value)
	for key, values := range p {
		listValue := &structpb.ListValue{}
		for _, value := range values {
			listValue.Values = append(listValue.Values, &structpb.Value{
				Kind: &structpb.Value_StringValue{StringValue: value},
			})
		}

		fields[key] = &structpb.Value{
			Kind: &structpb.Value_ListValue{ListValue: listValue},
		}
	}

	return &structpb.Struct{Fields: fields}, nil
}
