package v1beta1

import (
	"context"
	"fmt"
	"strings"

	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/errors"
	errorsPkg "github.com/goto/shield/pkg/errors"
	"github.com/goto/shield/pkg/metadata"
	"github.com/goto/shield/pkg/str"
	"github.com/goto/shield/pkg/uuid"
	"golang.org/x/exp/maps"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"

	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GroupService interface {
	Create(ctx context.Context, grp group.Group) (group.Group, error)
	Get(ctx context.Context, id string) (group.Group, error)
	List(ctx context.Context, flt group.Filter) ([]group.Group, error)
	Update(ctx context.Context, grp group.Group) (group.Group, error)
	ListUserGroups(ctx context.Context, userId string, roleId string) ([]group.Group, error)
	ListGroupRelations(ctx context.Context, objectId, subjectType, role string) ([]user.User, []group.Group, map[string][]string, map[string][]string, error)
}

var grpcGroupNotFoundErr = status.Errorf(codes.NotFound, "group doesn't exist")
var grpcInvalidOrgIDErr = status.Errorf(codes.InvalidArgument, "ordIs is not valid uuid")

func (h Handler) ListGroups(ctx context.Context, request *shieldv1beta1.ListGroupsRequest) (*shieldv1beta1.ListGroupsResponse, error) {
	logger := grpczap.Extract(ctx)

	if request.GetOrgId() != "" {
		if !uuid.IsValid(request.GetOrgId()) {
			return nil, grpcInvalidOrgIDErr
		}

		_, err := h.orgService.Get(ctx, request.GetOrgId())
		if err != nil {
			return &shieldv1beta1.ListGroupsResponse{Groups: nil}, nil
		}
	}

	var groups []*shieldv1beta1.Group

	currentUser, err := h.userService.FetchCurrentUser(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	servicedataKeyResourceIds, err := h.relationService.LookupResources(ctx, schema.ServiceDataKeyNamespace, schema.ViewPermission, schema.UserPrincipal, currentUser.ID)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	prj, err := h.projectService.Get(ctx, h.serviceDataConfig.DefaultServiceDataProject)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	groupList, err := h.groupService.List(ctx, group.Filter{
		OrganizationID:            request.GetOrgId(),
		ProjectID:                 prj.ID,
		ServicedataKeyResourceIDs: servicedataKeyResourceIds,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	for _, v := range groupList {
		groupPB, err := transformGroupToPB(v)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}

		groups = append(groups, &groupPB)
	}

	return &shieldv1beta1.ListGroupsResponse{Groups: groups}, nil
}

func (h Handler) CreateGroup(ctx context.Context, request *shieldv1beta1.CreateGroupRequest) (*shieldv1beta1.CreateGroupResponse, error) {
	logger := grpczap.Extract(ctx)

	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	_, err := h.userService.FetchCurrentUser(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	//TODO: change this
	metaDataMap, err := metadata.Build(request.GetBody().GetMetadata().AsMap())
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcBadBodyError
	}

	grp := group.Group{
		Name:           request.GetBody().GetName(),
		Slug:           request.GetBody().GetSlug(),
		OrganizationID: request.GetBody().GetOrgId(),
		Metadata:       nil,
	}

	if strings.TrimSpace(grp.Slug) == "" {
		grp.Slug = str.GenerateSlug(grp.Name)
	}

	newGroup, err := h.groupService.Create(ctx, grp)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, group.ErrConflict):
			return nil, grpcConflictError
		case errors.Is(err, group.ErrInvalidDetail), errors.Is(err, organization.ErrNotExist), errors.Is(err, organization.ErrInvalidUUID):
			return nil, grpcBadBodyError
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	serviceDataMap := map[string]any{}
	for k, v := range metaDataMap {
		serviceDataResp, err := h.serviceDataService.Upsert(ctx, servicedata.ServiceData{
			EntityID:    newGroup.ID,
			NamespaceID: groupNamespaceID,
			Key: servicedata.Key{
				Name:      k,
				ProjectID: h.serviceDataConfig.DefaultServiceDataProject,
			},
			Value: v,
		})
		if err != nil {
			logger.Error(err.Error())

			switch {
			case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrMissingEmail):
				return nil, grpcUnauthenticated
			case errors.Is(err, project.ErrNotExist), errors.Is(err, servicedata.ErrInvalidDetail),
				errors.Is(err, relation.ErrInvalidDetail), errors.Is(err, servicedata.ErrNotExist):
				return nil, grpcBadBodyError
			case errors.Is(err, errorsPkg.ErrForbidden):
				return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("you are not authorized to update %s key", k))
			default:
				return nil, grpcInternalServerError
			}
		}
		serviceDataMap[serviceDataResp.Key.Name] = serviceDataResp.Value
	}

	newGroup.Metadata = metaDataMap

	groupPB, err := transformGroupToPB(newGroup)
	if err != nil {
		logger.Error(err.Error())
		return nil, ErrInternalServer
	}

	return &shieldv1beta1.CreateGroupResponse{Group: &groupPB}, nil
}

func (h Handler) GetGroup(ctx context.Context, request *shieldv1beta1.GetGroupRequest) (*shieldv1beta1.GetGroupResponse, error) {
	logger := grpczap.Extract(ctx)

	fetchedGroup, err := h.groupService.Get(ctx, request.GetId())
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, group.ErrNotExist), errors.Is(err, group.ErrInvalidID), errors.Is(err, group.ErrInvalidUUID):
			return nil, grpcGroupNotFoundErr
		default:
			return nil, grpcInternalServerError
		}
	}

	filter := servicedata.Filter{
		ID:        fetchedGroup.ID,
		Namespace: groupNamespaceID,
		Entities: maps.Values(map[string]string{
			"group": groupNamespaceID,
		}),
	}

	groupSD, err := h.serviceDataService.Get(ctx, filter)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrMissingEmail):
			break
		default:
			return nil, grpcInternalServerError
		}
	} else {
		metadata := map[string]any{}
		for _, sd := range groupSD {
			metadata[sd.Key.Name] = sd.Value
		}
		fetchedGroup.Metadata = metadata
	}

	groupPB, err := transformGroupToPB(fetchedGroup)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.GetGroupResponse{Group: &groupPB}, nil
}

func (h Handler) UpdateGroup(ctx context.Context, request *shieldv1beta1.UpdateGroupRequest) (*shieldv1beta1.UpdateGroupResponse, error) {
	logger := grpczap.Extract(ctx)

	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	//TODO: change this implementation
	metaDataMap, err := metadata.Build(request.GetBody().GetMetadata().AsMap())
	if err != nil {
		return nil, grpcBadBodyError
	}

	var updatedGroup group.Group
	if uuid.IsValid(request.GetId()) {
		updatedGroup, err = h.groupService.Update(ctx, group.Group{
			ID:             request.GetId(),
			Name:           request.GetBody().GetName(),
			Slug:           request.GetBody().GetSlug(),
			OrganizationID: request.GetBody().GetOrgId(),
			Metadata:       nil,
		})
	} else {
		updatedGroup, err = h.groupService.Update(ctx, group.Group{
			Name:           request.GetBody().GetName(),
			Slug:           request.GetId(),
			OrganizationID: request.GetBody().GetOrgId(),
			Metadata:       nil,
		})
	}
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, group.ErrNotExist),
			errors.Is(err, group.ErrInvalidUUID),
			errors.Is(err, group.ErrInvalidID):
			return nil, grpcGroupNotFoundErr
		case errors.Is(err, group.ErrConflict):
			return nil, grpcConflictError
		case errors.Is(err, group.ErrInvalidDetail),
			errors.Is(err, organization.ErrInvalidUUID),
			errors.Is(err, organization.ErrNotExist):
			return nil, grpcBadBodyError
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	serviceDataMap := map[string]any{}
	for k, v := range metaDataMap {
		serviceDataResp, err := h.serviceDataService.Upsert(ctx, servicedata.ServiceData{
			EntityID:    updatedGroup.ID,
			NamespaceID: groupNamespaceID,
			Key: servicedata.Key{
				Name:      k,
				ProjectID: h.serviceDataConfig.DefaultServiceDataProject,
			},
			Value: v,
		})
		if err != nil {
			logger.Error(err.Error())

			switch {
			case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrMissingEmail):
				return nil, grpcUnauthenticated
			case errors.Is(err, project.ErrNotExist), errors.Is(err, servicedata.ErrInvalidDetail),
				errors.Is(err, relation.ErrInvalidDetail), errors.Is(err, servicedata.ErrNotExist):
				return nil, grpcBadBodyError
			case errors.Is(err, errorsPkg.ErrForbidden):
				return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("you are not authorized to update %s key", k))
			default:
				return nil, grpcInternalServerError
			}
		}
		serviceDataMap[serviceDataResp.Key.Name] = serviceDataResp.Value
	}

	//Note: this would return only the keys that are updated in the current request
	updatedGroup.Metadata = metaDataMap

	groupPB, err := transformGroupToPB(updatedGroup)
	if err != nil {
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.UpdateGroupResponse{Group: &groupPB}, nil
}

func transformGroupToPB(grp group.Group) (shieldv1beta1.Group, error) {
	metaData, err := grp.Metadata.ToStructPB()
	if err != nil {
		return shieldv1beta1.Group{}, err
	}

	return shieldv1beta1.Group{
		Id:        grp.ID,
		Name:      grp.Name,
		Slug:      grp.Slug,
		OrgId:     grp.OrganizationID,
		Metadata:  metaData,
		CreatedAt: timestamppb.New(grp.CreatedAt),
		UpdatedAt: timestamppb.New(grp.UpdatedAt),
	}, nil
}

func (h Handler) ListGroupRelations(ctx context.Context, request *shieldv1beta1.ListGroupRelationsRequest) (*shieldv1beta1.ListGroupRelationsResponse, error) {
	logger := grpczap.Extract(ctx)
	groupRelations := []*shieldv1beta1.GroupRelation{}

	users, groups, userIDRoleMap, groupIDRoleMap, err := h.groupService.ListGroupRelations(ctx, request.Id, request.SubjectType, request.Role)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	for _, user := range users {
		userPb, err := transformUserToPB(user)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}

		for _, r := range userIDRoleMap[userPb.Id] {
			role := strings.Split(r, ":")

			grprel := &shieldv1beta1.GroupRelation{
				SubjectType: schema.UserPrincipal,
				Role:        role[1],
				Subject: &shieldv1beta1.GroupRelation_User{
					User: &userPb,
				},
			}
			groupRelations = append(groupRelations, grprel)
		}
	}

	for _, group := range groups {
		groupPb, err := transformGroupToPB(group)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}

		for _, r := range groupIDRoleMap[groupPb.Id] {
			role := strings.Split(r, ":")

			grprel := &shieldv1beta1.GroupRelation{
				SubjectType: schema.GroupPrincipal,
				Role:        role[1],
				Subject: &shieldv1beta1.GroupRelation_Group{
					Group: &groupPb,
				},
			}
			groupRelations = append(groupRelations, grprel)
		}
	}

	return &shieldv1beta1.ListGroupRelationsResponse{
		Relations: groupRelations,
	}, nil
}
