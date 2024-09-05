package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
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
	ListResourceOfUser(ctx context.Context, userID string, resourceType string) ([]resource.ResourcePermission, error)
	ListResourceOfUserGlobal(ctx context.Context, userID string, resourceType []string) (map[string][]resource.ResourcePermission, error)
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

func (h Handler) ListResourceOfUserGlobal(ctx context.Context, request *shieldv1beta1.ListResourceOfUserGlobalRequest) (*shieldv1beta1.ListResourceOfUserGlobalResponse, error) {
	logger := grpczap.Extract(ctx)
	resources, err := h.resourceService.ListResourceOfUserGlobal(ctx, request.UserId, request.Types)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, resource.ErrNotExist),
			errors.Is(err, resource.ErrInvalidUUID),
			errors.Is(err, resource.ErrInvalidID):
			return nil, grpcResourceNotFoundErr
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrNotExist),
			errors.Is(err, resource.ErrInvalidDetail),
			errors.Is(err, resource.ErrInvalidURN):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	result := make(map[string]*structpb.Value)
	for key, value := range resources {
		resourcePB, err := transformListResourcePrincipalToPB(value)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}

		if len(resourcePB.Fields) == 0 {
			continue
		}

		result[key] = &structpb.Value{
			Kind: &structpb.Value_StructValue{StructValue: resourcePB},
		}
	}

	resultPB := &structpb.Struct{
		Fields: result,
	}

	return &shieldv1beta1.ListResourceOfUserGlobalResponse{
		Resources: resultPB,
	}, nil
}

func (h Handler) ListResourceOfUser(ctx context.Context, request *shieldv1beta1.ListResourceOfUserRequest) (*shieldv1beta1.ListResourceOfUserResponse, error) {
	logger := grpczap.Extract(ctx)

	resources, err := h.resourceService.ListResourceOfUser(ctx, request.UserId, fmt.Sprintf("%s/%s", request.Namespace, request.Type))
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, resource.ErrNotExist),
			errors.Is(err, resource.ErrInvalidUUID),
			errors.Is(err, resource.ErrInvalidID):
			return nil, grpcResourceNotFoundErr
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrNotExist),
			errors.Is(err, resource.ErrInvalidDetail),
			errors.Is(err, resource.ErrInvalidURN):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	resourcesPB, err := transformListResourcePrincipalToPB(resources)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.ListResourceOfUserResponse{
		Resources: resourcesPB,
	}, nil
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

func transformListResourcePrincipalToPB(from []resource.ResourcePermission) (*structpb.Struct, error) {
	to := make(map[string][]string)
	for _, f := range from {
		for _, res := range f.ResourceIDs {
			if _, ok := to[res]; !ok {
				to[res] = []string{f.Permission}
			} else {
				to[res] = append(to[res], f.Permission)
			}
		}
	}

	toPB, err := mapToStructpb(to)
	if err != nil {
		return nil, err
	}

	return toPB, nil
}

func mapToStructpb(m map[string][]string) (*structpb.Struct, error) {
	fields := make(map[string]*structpb.Value)
	for key, values := range m {
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
