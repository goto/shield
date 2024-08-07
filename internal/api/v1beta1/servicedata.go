package v1beta1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/project"
	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/servicedata"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"
	errPkg "github.com/goto/shield/pkg/errors"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	userNamespaceID      = schema.UserPrincipal
	groupNamespaceID     = schema.GroupPrincipal
	projectNamespaceID   = schema.ProjectNamespace
	entitiesNamespaceMap = map[string]string{
		"user":  userNamespaceID,
		"group": groupNamespaceID,
	}
)

type ServiceDataService interface {
	CreateKey(ctx context.Context, key servicedata.Key) (servicedata.Key, error)
	Upsert(ctx context.Context, serviceData servicedata.ServiceData) (servicedata.ServiceData, error)
	Get(ctx context.Context, filter servicedata.Filter) ([]servicedata.ServiceData, error)
	GetKeyByURN(ctx context.Context, urn string) (servicedata.Key, error)
}

func (h Handler) CreateServiceDataKey(ctx context.Context, request *shieldv1beta1.CreateServiceDataKeyRequest) (*shieldv1beta1.CreateServiceDataKeyResponse, error) {
	logger := grpczap.Extract(ctx)

	requestBody := request.GetBody()
	if requestBody == nil {
		return nil, grpcBadBodyError
	}

	keyResp, err := h.serviceDataService.CreateKey(ctx, servicedata.Key{
		ProjectID:   requestBody.GetProject(),
		Name:        requestBody.GetKey(),
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

	if request.GetUserId() == "" {
		return nil, grpcBadBodyError
	}

	data := requestBody.GetData()
	if data == nil {
		return nil, grpcBadBodyError
	}
	sdMap := data.AsMap()

	if len(sdMap) > h.serviceDataConfig.MaxUpsert {
		return nil, grpcBadBodyError
	}

	// get user by id or email
	userEntity, err := h.userService.Get(ctx, request.GetUserId())
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
	serviceDataMap := map[string]any{}
	for k, v := range sdMap {
		serviceDataResp, err := h.serviceDataService.Upsert(ctx, servicedata.ServiceData{
			EntityID:    userEntity.ID,
			NamespaceID: userNamespaceID,
			Key: servicedata.Key{
				Name:      k,
				ProjectID: requestBody.Project,
			},
			Value: v,
		})
		if err != nil {
			logger.Error(err.Error())

			switch {
			case errors.Is(err, user.ErrInvalidEmail), errors.Is(err, user.ErrMissingEmail):
				return nil, grpcUnauthenticated
			case errors.Is(err, errPkg.ErrForbidden):
				return nil, grpcPermissionDenied
			case errors.Is(err, project.ErrNotExist), errors.Is(err, servicedata.ErrInvalidDetail),
				errors.Is(err, relation.ErrInvalidDetail), errors.Is(err, servicedata.ErrNotExist):
				return nil, grpcBadBodyError
			default:
				return nil, grpcInternalServerError
			}
		}
		serviceDataMap[serviceDataResp.Key.Name] = serviceDataResp.Value
	}

	serviceDataMapPB, err := structpb.NewStruct(serviceDataMap)
	if err != nil {
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.UpsertUserServiceDataResponse{
		Data: serviceDataMapPB,
	}, nil
}

func (h Handler) UpsertGroupServiceData(ctx context.Context, request *shieldv1beta1.UpsertGroupServiceDataRequest) (*shieldv1beta1.UpsertGroupServiceDataResponse, error) {
	logger := grpczap.Extract(ctx)

	requestBody := request.GetBody()
	if requestBody == nil {
		return nil, grpcBadBodyError
	}

	if request.GetGroupId() == "" {
		return nil, grpcBadBodyError
	}

	data := requestBody.GetData()
	if data == nil {
		return nil, grpcBadBodyError
	}
	sdMap := data.AsMap()

	if len(sdMap) > h.serviceDataConfig.MaxUpsert {
		return nil, grpcBadBodyError
	}

	// get group by id or slug
	groupEntity, err := h.groupService.Get(ctx, request.GetGroupId())
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
	serviceDataMap := map[string]any{}
	for k, v := range sdMap {
		serviceDataResp, err := h.serviceDataService.Upsert(ctx, servicedata.ServiceData{
			EntityID:    groupEntity.ID,
			NamespaceID: groupNamespaceID,
			Key: servicedata.Key{
				Name:      k,
				ProjectID: requestBody.Project,
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
			default:
				return nil, grpcInternalServerError
			}
		}
		serviceDataMap[serviceDataResp.Key.Name] = serviceDataResp.Value
	}

	serviceDataMapPB, err := structpb.NewStruct(serviceDataMap)
	if err != nil {
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.UpsertGroupServiceDataResponse{
		Data: serviceDataMapPB,
	}, nil
}

func (h Handler) GetUserServiceData(ctx context.Context, request *shieldv1beta1.GetUserServiceDataRequest) (*shieldv1beta1.GetUserServiceDataResponse, error) {
	logger := grpczap.Extract(ctx)

	usr, err := h.userService.Get(ctx, request.GetUserId())
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

	entities := []string{}
	if request.GetEntity() != nil {
		for _, ent := range request.GetEntity() {
			if entNamespace, ok := entitiesNamespaceMap[ent]; ok {
				entities = append(entities, entNamespace)
			}
		}
	}

	if len(entities) == 0 {
		entities = maps.Values(entitiesNamespaceMap)
	}

	filter := servicedata.Filter{
		ID:        usr.ID,
		Namespace: userNamespaceID,
		Entities:  entities,
		Project:   request.GetProject(),
	}

	serviceData, err := h.serviceDataService.Get(ctx, filter)
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

	serviceDataPB, err := transformServiceDataListToPB(serviceData)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.GetUserServiceDataResponse{
		Data: serviceDataPB,
	}, nil
}

func (h Handler) GetGroupServiceData(ctx context.Context, request *shieldv1beta1.GetGroupServiceDataRequest) (*shieldv1beta1.GetGroupServiceDataResponse, error) {
	logger := grpczap.Extract(ctx)

	grp, err := h.groupService.Get(ctx, request.GetGroupId())
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

	filter := servicedata.Filter{
		ID:        grp.ID,
		Namespace: groupNamespaceID,
		Project:   request.GetProject(),
	}

	serviceData, err := h.serviceDataService.Get(ctx, filter)
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

	serviceDataPB, err := transformServiceDataListToPB(serviceData)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.GetGroupServiceDataResponse{
		Data: serviceDataPB,
	}, nil
}

func transformServiceDataKeyToPB(from servicedata.Key) (shieldv1beta1.ServiceDataKey, error) {
	return shieldv1beta1.ServiceDataKey{
		Urn: from.URN,
		Id:  from.ID,
	}, nil
}

func transformServiceDataListToPB(from []servicedata.ServiceData) (*structpb.Struct, error) {
	data := map[string]map[string]map[string]any{}

	for _, sd := range from {
		prjKey := fmt.Sprintf("%s:%s", projectNamespaceID, sd.Key.ProjectID)
		entKey := fmt.Sprintf("%s:%s", sd.NamespaceID, sd.EntityID)
		prj, ok := data[prjKey]
		if ok {
			ent, ok := prj[entKey]
			if ok {
				ent[sd.Key.Name] = sd.Value
			} else {
				prj[entKey] = map[string]any{
					sd.Key.Name: sd.Value,
				}
			}
		} else {
			kv := map[string]any{sd.Key.Name: sd.Value}
			data[prjKey] = map[string]map[string]any{
				entKey: kv,
			}
		}
	}

	var decodedData map[string]interface{}
	encodedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(encodedData, &decodedData)
	if err != nil {
		return nil, err
	}

	serviceData, err := structpb.NewStruct(decodedData)
	if err != nil {
		return nil, err
	}

	return serviceData, nil
}
