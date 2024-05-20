package v1beta1

import (
	"context"
	"errors"

	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type ServiceDataService interface {
	CreateKey(ctx context.Context, key servicedata.Key) (servicedata.Key, error)
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
		case errors.Is(err, user.ErrInvalidEmail):
			return nil, grpcUnauthenticated
		case errors.Is(err, project.ErrNotExist) || errors.Is(err, servicedata.ErrInvalidDetail) ||
			errors.Is(err, relation.ErrInvalidDetail):
			return nil, grpcBadBodyError
		case errors.Is(err, servicedata.ErrConflict):
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

func transformServiceDataKeyToPB(from servicedata.Key) (shieldv1beta1.ServiceDataKey, error) {
	return shieldv1beta1.ServiceDataKey{
		Urn: from.URN,
		Id:  from.ID,
	}, nil
}
