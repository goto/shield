package v1beta1

import (
	"context"
	"errors"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/exp/maps"
)

var (
	userNamespaceID = namespace.DefinitionUser.ID
	groupNamepaceID = namespace.DefinitionTeam.ID
)

type ServiceDataService interface {
	CreateKey(ctx context.Context, key servicedata.Key) (servicedata.Key, error)
	Upsert(ctx context.Context, servicedata servicedata.ServiceData) (servicedata.ServiceData, error)
}

func (h Handler) CreateServiceDataKey(ctx context.Context, request *shieldv1beta1.CreateServiceDataKeyRequest) (*shieldv1beta1.CreateServiceDataKeyResponse, error) {
	logger := grpczap.Extract(ctx)

	requestBody := request.GetBody()
	if requestBody == nil {
		return nil, grpcBadBodyError
	}

	keyResp, err := h.serviceDataService.CreateKey(ctx, servicedata.Key{
		ProjectID:   requestBody.GetProject(),
		Key:         requestBody.GetKey(),
		Description: requestBody.GetDescription(),
	})
	if err != nil {
		logger.Error(err.Error())

		switch {
		case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		case errors.Is(err, project.ErrNotExist), errors.Is(err, servicedata.ErrInvalidDetail),
			errors.Is(err, relation.ErrInvalidDetail):
			return nil, grpcBadBodyError
		case errors.Is(err, servicedata.ErrConflict), errors.Is(err, resource.ErrConflict):
			return nil, grpcConflictError
		default:
			return nil, grpcInternalServerError
		}
	}

	serviceDataKey, err := transformServiceDataKeyToPB(keyResp)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.CreateServiceDataKeyResponse{
		ServiceDataKey: &serviceDataKey,
	}, nil
}

func (h Handler) UpsertUserServiceData(ctx context.Context, request *shieldv1beta1.UpsertUserServiceDataRequest) (*shieldv1beta1.UpsertUserServiceDataResponse, error) {
	logger := grpczap.Extract(ctx)

	requestBody := request.GetBody()
	if requestBody == nil {
		return nil, grpcBadBodyError
	}

	if request.GetId() == "" {
		return nil, grpcBadBodyError
	}

	if len(requestBody.Data) != 1 {
		return nil, grpcBadBodyError
	}

	key := maps.Keys(requestBody.Data)[0]
	value := requestBody.Data[key]

	// get user by id or email
	userEntity, err := h.userService.Get(ctx, request.GetId())
	if err != nil {
		logger.Error(err.Error())

		switch {
		case errors.Is(err, user.ErrNotExist), errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidID):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	serviceDataResp, err := h.serviceDataService.Upsert(ctx, servicedata.ServiceData{
		EntityID:    userEntity.ID,
		NamespaceID: userNamespaceID,
		Key: servicedata.Key{
			Key:       key,
			ProjectID: requestBody.Project,
		},
		Value: value,
	})
	if err != nil {
		logger.Error(err.Error())

		switch {
		case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		case errors.Is(err, project.ErrNotExist), errors.Is(err, servicedata.ErrInvalidDetail),
			errors.Is(err, relation.ErrInvalidDetail):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	return &shieldv1beta1.UpsertUserServiceDataResponse{
		Urn: serviceDataResp.Key.URN,
	}, nil
}

func (h Handler) UpsertGroupServiceData(ctx context.Context, request *shieldv1beta1.UpsertGroupServiceDataRequest) (*shieldv1beta1.UpsertGroupServiceDataResponse, error) {
	logger := grpczap.Extract(ctx)

	requestBody := request.GetBody()
	if requestBody == nil {
		return nil, grpcBadBodyError
	}

	if request.GetId() == "" {
		return nil, grpcBadBodyError
	}

	if len(requestBody.Data) != 1 {
		return nil, grpcBadBodyError
	}

	key := maps.Keys(requestBody.Data)[0]
	value := requestBody.Data[key]

	// get group by id or slug
	groupEntity, err := h.groupService.Get(ctx, request.GetId())
	if err != nil {
		logger.Error(err.Error())

		switch {
		case errors.Is(err, group.ErrNotExist), errors.Is(err, group.ErrInvalidDetail),
			errors.Is(err, group.ErrInvalidID):
			return nil, grpcBadBodyError
		default:
			return nil, grpcInternalServerError
		}
	}

	serviceDataResp, err := h.serviceDataService.Upsert(ctx, servicedata.ServiceData{
		EntityID:    groupEntity.ID,
		NamespaceID: groupNamepaceID,
		Key: servicedata.Key{
			Key:       key,
			ProjectID: requestBody.Project,
		},
		Value: value,
	})
	if err != nil {
		logger.Error(err.Error())

		switch {
		case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		case errors.Is(err, project.ErrNotExist), errors.Is(err, servicedata.ErrInvalidDetail),
			errors.Is(err, relation.ErrInvalidDetail):
			return nil, grpcBadBodyError
		case errors.Is(err, servicedata.ErrConflict), errors.Is(err, resource.ErrConflict):
			return nil, grpcConflictError
		default:
			return nil, grpcInternalServerError
		}
	}

	return &shieldv1beta1.UpsertGroupServiceDataResponse{
		Urn: serviceDataResp.Key.URN,
	}, nil
}

func transformServiceDataKeyToPB(from servicedata.Key) (shieldv1beta1.ServiceDataKey, error) {
	return shieldv1beta1.ServiceDataKey{
		Urn: from.URN,
		Id:  from.ID,
	}, nil
}
