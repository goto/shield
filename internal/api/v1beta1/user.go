package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/exp/maps"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	"github.com/goto/shield/pkg/metadata"

	errorsPkg "github.com/goto/shield/pkg/errors"
	"github.com/goto/shield/pkg/uuid"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
)

var grpcUserNotFoundError = status.Errorf(codes.NotFound, "user doesn't exist")

type UserService interface {
	Get(ctx context.Context, idOrEmail string) (user.User, error)
	GetByIDs(ctx context.Context, userIDs []string) ([]user.User, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
	Create(ctx context.Context, user user.User) (user.User, error)
	List(ctx context.Context, flt user.Filter) (user.PagedUsers, error)
	UpdateByID(ctx context.Context, toUpdate user.User) (user.User, error)
	UpdateByEmail(ctx context.Context, toUpdate user.User) (user.User, error)
	FetchCurrentUser(ctx context.Context) (user.User, error)
	CreateMetadataKey(ctx context.Context, key user.UserMetadataKey) (user.UserMetadataKey, error)
}

func (h Handler) ListUsers(ctx context.Context, request *shieldv1beta1.ListUsersRequest) (*shieldv1beta1.ListUsersResponse, error) {
	logger := grpczap.Extract(ctx)
	var users []*shieldv1beta1.User

	currentUser, err := h.userService.FetchCurrentUser(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcUnauthenticated
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

	userResp, err := h.userService.List(ctx, user.Filter{
		Limit:                     request.GetPageSize(),
		Page:                      request.GetPageNum(),
		Keyword:                   request.GetKeyword(),
		ProjectID:                 prj.ID,
		ServiceDataKeyResourceIds: servicedataKeyResourceIds,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	userList := userResp.Users
	for _, user := range userList {
		userPB, err := transformUserToPB(user)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}

		users = append(users, &userPB)
	}

	return &shieldv1beta1.ListUsersResponse{
		Count: userResp.Count,
		Users: users,
	}, nil
}

func (h Handler) CreateUser(ctx context.Context, request *shieldv1beta1.CreateUserRequest) (*shieldv1beta1.CreateUserResponse, error) {
	logger := grpczap.Extract(ctx)

	currentUser, err := h.userService.FetchCurrentUser(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcUnauthenticated
	}

	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	email := strings.TrimSpace(request.GetBody().GetEmail())
	if email == "" {
		return nil, grpcBadBodyError
	}
	if !isValidEmail(email) {
		return nil, user.ErrInvalidEmail
	}

	metaDataMap, err := metadata.Build(request.GetBody().GetMetadata().AsMap())
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcBadBodyError
	}

	for k := range metaDataMap {
		urn := servicedata.CreateURN(h.serviceDataConfig.DefaultServiceDataProject, k)
		key, err := h.serviceDataService.GetKeyByURN(ctx, urn)
		if err != nil {
			return nil, err
		}

		permission, err := h.relationService.CheckPermission(ctx, currentUser, namespace.Namespace{ID: schema.ServiceDataKeyNamespace}, key.ResourceID, action.Action{ID: schema.EditPermission})
		if err != nil {
			return nil, err
		}
		if !permission {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("you are not authorized to update %s key", k))
		}
	}

	// TODO might need to check the valid email form
	newUser, err := h.userService.Create(ctx, user.User{
		Name:     request.GetBody().GetName(),
		Email:    email,
		Metadata: nil,
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, user.ErrConflict):
			return nil, grpcConflictError
		case errors.Is(errors.Unwrap(err), user.ErrKeyDoesNotExists):
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
			EntityID:    newUser.ID,
			NamespaceID: userNamespaceID,
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

	// TODO: use serviceDataMap
	newUser.Metadata = metaDataMap

	userPB, err := transformUserToPB(newUser)
	if err != nil {
		logger.Error(err.Error())
		return nil, ErrInternalServer
	}

	return &shieldv1beta1.CreateUserResponse{User: &userPB}, nil
}

func (h Handler) CreateMetadataKey(ctx context.Context, request *shieldv1beta1.CreateMetadataKeyRequest) (*shieldv1beta1.CreateMetadataKeyResponse, error) {
	logger := grpczap.Extract(ctx)

	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	newKey, err := h.userService.CreateMetadataKey(ctx, user.UserMetadataKey{
		Key:         request.GetBody().GetKey(),
		Description: request.GetBody().GetDescription(),
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, user.ErrConflict):
			return nil, grpcConflictError
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	return &shieldv1beta1.CreateMetadataKeyResponse{Metadatakey: &shieldv1beta1.MetadataKey{
		Key:         newKey.Key,
		Description: newKey.Description,
	}}, nil
}

func (h Handler) GetUser(ctx context.Context, request *shieldv1beta1.GetUserRequest) (*shieldv1beta1.GetUserResponse, error) {
	logger := grpczap.Extract(ctx)

	_, err := h.userService.FetchCurrentUser(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcUnauthenticated
	}

	fetchedUser, err := h.userService.Get(ctx, request.GetId())
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, user.ErrNotExist), errors.Is(err, user.ErrInvalidUUID), errors.Is(err, user.ErrInvalidID):
			return nil, grpcUserNotFoundError
		default:
			return nil, grpcInternalServerError
		}
	}

	filter := servicedata.Filter{
		ID:        fetchedUser.ID,
		Namespace: userNamespaceID,
		Entities: maps.Values(map[string]string{
			"user": userNamespaceID,
		}),
	}

	userSD, err := h.serviceDataService.Get(ctx, filter)
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
		for _, sd := range userSD {
			metadata[sd.Key.Name] = sd.Value
		}
		fetchedUser.Metadata = metadata
	}

	userPB, err := transformUserToPB(fetchedUser)
	if err != nil {
		logger.Error(err.Error())
		return nil, status.Errorf(codes.Internal, ErrInternalServer.Error())
	}

	return &shieldv1beta1.GetUserResponse{
		User: &userPB,
	}, nil
}

func (h Handler) GetCurrentUser(ctx context.Context, request *shieldv1beta1.GetCurrentUserRequest) (*shieldv1beta1.GetCurrentUserResponse, error) {
	logger := grpczap.Extract(ctx)

	email, ok := user.GetEmailFromContext(ctx)
	if !ok {
		return nil, grpcUnauthenticated
	}

	email = strings.TrimSpace(email)
	if email == "" {
		logger.Error(ErrEmptyEmailID.Error())
		return nil, grpcUnauthenticated
	}

	fetchedUser, err := h.userService.GetByEmail(ctx, email)
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, user.ErrNotExist), errors.Is(err, user.ErrInvalidID), errors.Is(err, user.ErrInvalidEmail):
			return nil, grpcUserNotFoundError
		default:
			return nil, grpcInternalServerError
		}
	}

	userPB, err := transformUserToPB(fetchedUser)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.GetCurrentUserResponse{
		User: &userPB,
	}, nil
}

func (h Handler) UpdateUser(ctx context.Context, request *shieldv1beta1.UpdateUserRequest) (*shieldv1beta1.UpdateUserResponse, error) {
	logger := grpczap.Extract(ctx)
	var updatedUser user.User

	currentUser, err := h.userService.FetchCurrentUser(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcUnauthenticated
	}

	if strings.TrimSpace(request.GetId()) == "" {
		return nil, grpcUserNotFoundError
	}

	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	email := strings.TrimSpace(request.GetBody().GetEmail())
	if email == "" {
		return nil, grpcBadBodyError
	}

	if !isValidEmail(email) {
		return nil, user.ErrInvalidEmail
	}

	metaDataMap, err := metadata.Build(request.GetBody().GetMetadata().AsMap())
	if err != nil {
		return nil, grpcBadBodyError
	}

	for k := range metaDataMap {
		urn := servicedata.CreateURN(h.serviceDataConfig.DefaultServiceDataProject, k)
		key, err := h.serviceDataService.GetKeyByURN(ctx, urn)
		if err != nil {
			return nil, err
		}

		permission, err := h.relationService.CheckPermission(ctx, currentUser, namespace.Namespace{ID: schema.ServiceDataKeyNamespace}, key.ResourceID, action.Action{ID: schema.EditPermission})
		if err != nil {
			return nil, err
		}
		if !permission {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("you are not authorized to update %s key", k))
		}
	}

	id := request.GetId()
	if uuid.IsValid(id) {
		updatedUser, err = h.userService.UpdateByID(ctx, user.User{
			ID:       id,
			Name:     request.GetBody().GetName(),
			Email:    email,
			Metadata: nil,
		})
		if err != nil {
			logger.Error(err.Error())
			switch {
			case errors.Is(err, user.ErrNotExist), errors.Is(err, user.ErrInvalidID), errors.Is(err, user.ErrInvalidUUID):
				return nil, grpcUserNotFoundError
			case errors.Is(err, user.ErrInvalidEmail):
				return nil, grpcBadBodyError
			case errors.Is(err, user.ErrConflict):
				return nil, grpcConflictError
			case errors.Is(err, user.ErrInvalidEmail),
				errors.Is(err, user.ErrMissingEmail):
				return nil, grpcUnauthenticated
			default:
				return nil, grpcInternalServerError
			}
		}
	} else {
		_, err := h.userService.GetByEmail(ctx, id)
		if err != nil {
			if err == user.ErrNotExist {
				createUserResponse, err := h.CreateUser(ctx, &shieldv1beta1.CreateUserRequest{Body: request.GetBody()})
				if err != nil {
					return nil, grpcInternalServerError
				}
				return &shieldv1beta1.UpdateUserResponse{User: createUserResponse.User}, nil
			} else {
				return nil, grpcInternalServerError
			}
		}

		updatedUser, err = h.userService.UpdateByEmail(ctx, user.User{
			Name:     request.GetBody().GetName(),
			Email:    email,
			Metadata: nil,
		})
		if err != nil {
			logger.Error(err.Error())
			switch {
			case errors.Is(err, user.ErrNotExist), errors.Is(err, user.ErrInvalidID), errors.Is(err, user.ErrInvalidUUID):
				return nil, grpcUserNotFoundError
			case errors.Is(err, user.ErrInvalidEmail):
				return nil, grpcBadBodyError
			case errors.Is(err, user.ErrConflict):
				return nil, grpcConflictError
			case errors.Is(err, user.ErrInvalidEmail),
				errors.Is(err, user.ErrMissingEmail):
				return nil, grpcUnauthenticated
			default:
				return nil, grpcInternalServerError
			}
		}
	}

	serviceDataMap := map[string]any{}
	for k, v := range metaDataMap {
		serviceDataResp, err := h.serviceDataService.Upsert(ctx, servicedata.ServiceData{
			EntityID:    updatedUser.ID,
			NamespaceID: userNamespaceID,
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

	// Note: this would return only the keys that are updated in the current request
	updatedUser.Metadata = metaDataMap

	userPB, err := transformUserToPB(updatedUser)
	if err != nil {
		logger.Error(err.Error())
		return nil, ErrInternalServer
	}

	return &shieldv1beta1.UpdateUserResponse{User: &userPB}, nil
}

func (h Handler) UpdateCurrentUser(ctx context.Context, request *shieldv1beta1.UpdateCurrentUserRequest) (*shieldv1beta1.UpdateCurrentUserResponse, error) {
	logger := grpczap.Extract(ctx)

	email, ok := user.GetEmailFromContext(ctx)
	if !ok {
		return nil, grpcUnauthenticated
	}

	email = strings.TrimSpace(email)
	if email == "" {
		logger.Error(ErrEmptyEmailID.Error())
		return nil, grpcUnauthenticated
	}

	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	metaDataMap, err := metadata.Build(request.GetBody().GetMetadata().AsMap())
	if err != nil {
		return nil, grpcBadBodyError
	}

	// if email in request body is different from the email in the header
	if request.GetBody().GetEmail() != email {
		return nil, grpcBadBodyError
	}

	updatedUser, err := h.userService.UpdateByEmail(ctx, user.User{
		Name:     request.GetBody().GetName(),
		Email:    email,
		Metadata: metaDataMap,
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, user.ErrNotExist):
			return nil, grpcUserNotFoundError
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	userPB, err := transformUserToPB(updatedUser)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.UpdateCurrentUserResponse{User: &userPB}, nil
}

func (h Handler) ListUserGroups(ctx context.Context, request *shieldv1beta1.ListUserGroupsRequest) (*shieldv1beta1.ListUserGroupsResponse, error) {
	logger := grpczap.Extract(ctx)
	var groups []*shieldv1beta1.Group

	groupsList, err := h.groupService.ListUserGroups(ctx, request.GetId(), request.GetRole())
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	for _, group := range groupsList {
		groupPB, err := transformGroupToPB(group)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}

		groups = append(groups, &groupPB)
	}

	return &shieldv1beta1.ListUserGroupsResponse{
		Groups: groups,
	}, nil
}

func transformUserToPB(usr user.User) (shieldv1beta1.User, error) {
	metaData, err := usr.Metadata.ToStructPB()
	if err != nil {
		return shieldv1beta1.User{}, err
	}

	return shieldv1beta1.User{
		Id:        usr.ID,
		Name:      usr.Name,
		Email:     usr.Email,
		Metadata:  metaData,
		CreatedAt: timestamppb.New(usr.CreatedAt),
		UpdatedAt: timestamppb.New(usr.UpdatedAt),
	}, nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
