package v1beta1

import (
	"context"
	"errors"
	"fmt"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"

	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h Handler) CheckResourceUserPermission(ctx context.Context, req *shieldv1beta1.CheckResourceUserPermissionRequest) (*shieldv1beta1.CheckResourceUserPermissionResponse, error) {
	userCtx := user.SetContextWithEmail(ctx, req.GetId()) // id is e-mail here
	resp, err := h.CheckResourcePermission(userCtx, &shieldv1beta1.CheckResourcePermissionRequest{
		ResourcePermissions: req.GetResourcePermissions(),
	})
	if err != nil {
		return nil, err
	}

	var permissionResponses []*shieldv1beta1.CheckResourceUserPermissionResponse_ResourcePermissionResponse
	for _, permissionResp := range resp.GetResourcePermissions() {
		permissionResponses = append(permissionResponses, &shieldv1beta1.CheckResourceUserPermissionResponse_ResourcePermissionResponse{
			ObjectId:        permissionResp.GetObjectId(),
			ObjectNamespace: permissionResp.GetObjectNamespace(),
			Permission:      permissionResp.GetPermission(),
			Allowed:         permissionResp.GetAllowed(),
		})
	}

	return &shieldv1beta1.CheckResourceUserPermissionResponse{
		ResourcePermissions: permissionResponses,
	}, nil
}

func (h Handler) CheckResourcePermission(ctx context.Context, req *shieldv1beta1.CheckResourcePermissionRequest) (*shieldv1beta1.CheckResourcePermissionResponse, error) {
	logger := grpczap.Extract(ctx)
	//if err := req.ValidateAll(); err != nil {
	//	formattedErr := getValidationErrorMessage(err)
	//	logger.Error(formattedErr.Error())
	//	return nil, status.Errorf(codes.NotFound, formattedErr.Error())
	//}
	//  To have backward compatibility
	if req.ObjectId != "" && len(req.ResourcePermissions) == 0 {
		return h.checkSingleResourcePermission(ctx, req)
	}

	if len(req.ResourcePermissions) > h.checkAPILimit || len(req.ResourcePermissions) == 0 {
		formattedErr := fmt.Errorf("%s: %s", ErrRequestBodyValidation, "resource_permissions")
		logger.Error(formattedErr.Error())
		return nil, formattedErr
	}

	var resources []resource.Resource
	var actions []action.Action
	for _, permission := range req.ResourcePermissions {
		resources = append(resources, resource.Resource{
			Name:        permission.GetObjectId(),
			NamespaceID: permission.GetObjectNamespace(),
		})

		actions = append(actions, action.Action{ID: permission.GetPermission()})
	}

	results, err := h.resourceService.BulkCheckAuthz(ctx, resources, actions)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail):
			logger.Error(err.Error())
			return nil, grpcUnauthenticated
		default:
			formattedErr := fmt.Errorf("%s: %w", ErrInternalServer, err)
			logger.Error(formattedErr.Error())
			return nil, status.Errorf(codes.Internal, ErrInternalServer.Error())
		}
	}

	var responseResults []*shieldv1beta1.CheckResourcePermissionResponse_ResourcePermissionResponse
	for _, res := range results {
		responseResults = append(responseResults, &shieldv1beta1.CheckResourcePermissionResponse_ResourcePermissionResponse{
			ObjectId:        res.ObjectID,
			ObjectNamespace: res.ObjectNamespace,
			Permission:      res.Permission,
			Allowed:         res.Allowed,
		})
	}

	return &shieldv1beta1.CheckResourcePermissionResponse{ResourcePermissions: responseResults}, nil
}

// Deprecated: checkSingleResourcePermission is used to check the single resource permission, this method is deprecated it needs to be cleaned up once the clients are onboarded
// with the new check API
func (h Handler) checkSingleResourcePermission(ctx context.Context, req *shieldv1beta1.CheckResourcePermissionRequest) (*shieldv1beta1.CheckResourcePermissionResponse, error) {
	logger := grpczap.Extract(ctx)

	result, err := h.resourceService.CheckAuthz(ctx, resource.Resource{
		Name:        req.GetObjectId(),
		NamespaceID: req.GetObjectNamespace(),
	}, action.Action{ID: req.GetPermission()})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidEmail):
			return nil, grpcUnauthenticated
		default:
			formattedErr := fmt.Errorf("%s: %w", ErrInternalServer, err)
			logger.Error(formattedErr.Error())
			return nil, status.Errorf(codes.Internal, ErrInternalServer.Error())
		}
	}

	if !result {
		return &shieldv1beta1.CheckResourcePermissionResponse{Status: false}, nil
	}

	return &shieldv1beta1.CheckResourcePermissionResponse{Status: true}, nil
}
