package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/goto/shield/core/action"
	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/internal/schema"

	"github.com/goto/shield/core/relation"
	errpkg "github.com/goto/shield/pkg/errors"
	"github.com/goto/shield/pkg/uuid"
	shieldv1beta1 "github.com/goto/shield/proto/v1beta1"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RelationService interface {
	Get(ctx context.Context, id string) (relation.RelationV2, error)
	Create(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	List(ctx context.Context) ([]relation.RelationV2, error)
	DeleteV2(ctx context.Context, rel relation.RelationV2) error
	GetRelationByFields(ctx context.Context, rel relation.RelationV2) (relation.RelationV2, error)
	LookupResources(ctx context.Context, resourceType, permission, subjectType, subjectID string) ([]string, error)
	CheckPermission(ctx context.Context, usr user.User, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error)
	CheckIsPublic(ctx context.Context, resourceNS namespace.Namespace, resourceIdxa string, action action.Action) (bool, error)
}

var grpcRelationNotFoundErr = status.Errorf(codes.NotFound, "relation doesn't exist")

func (h Handler) ListRelations(ctx context.Context, request *shieldv1beta1.ListRelationsRequest) (*shieldv1beta1.ListRelationsResponse, error) {
	logger := grpczap.Extract(ctx)
	var relations []*shieldv1beta1.Relation

	relationsList, err := h.relationService.List(ctx)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	for _, r := range relationsList {
		relationPB, err := transformRelationV2ToPB(r)
		if err != nil {
			logger.Error(err.Error())
			return nil, grpcInternalServerError
		}

		relations = append(relations, &relationPB)
	}

	return &shieldv1beta1.ListRelationsResponse{
		Relations: relations,
	}, nil
}

func (h Handler) CreateRelation(ctx context.Context, request *shieldv1beta1.CreateRelationRequest) (*shieldv1beta1.CreateRelationResponse, error) {
	logger := grpczap.Extract(ctx)
	if request.GetBody() == nil {
		return nil, grpcBadBodyError
	}

	if !uuid.IsValid(request.GetBody().GetObjectId()) {
		return nil, grpcBadBodyError
	}

	if !namespace.IsSystemNamespaceID(request.GetBody().GetObjectNamespace()) {
		_, err := h.resourceService.Get(ctx, request.GetBody().GetObjectId())
		if err != nil {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
	}

	principal, subjectID := extractSubjectFromPrincipal(request.GetBody().GetSubject())
	if err := h.verifyPrincipal(ctx, principal, subjectID); err != nil {
		return nil, err
	}

	result, err := h.resourceService.CheckAuthz(ctx, resource.Resource{
		Name:        request.GetBody().GetObjectId(),
		NamespaceID: request.GetBody().ObjectNamespace,
	}, action.Action{ID: schema.EditPermission})
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
		return nil, status.Errorf(codes.PermissionDenied, errpkg.ErrForbidden.Error())
	}

	newRelation, err := h.createRelation(ctx, relation.RelationV2{
		Object: relation.Object{
			ID:          request.GetBody().GetObjectId(),
			NamespaceID: request.GetBody().GetObjectNamespace(),
		},
		Subject: relation.Subject{
			ID:        subjectID,
			Namespace: principal,
			RoleID:    request.GetBody().GetRoleName(),
		},
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, relation.ErrInvalidDetail):
			return nil, grpcBadBodyError
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	relationPB, err := transformRelationV2ToPB(newRelation)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.CreateRelationResponse{
		Relation: &relationPB,
	}, nil
}

func (h Handler) createRelation(ctx context.Context, rlt relation.RelationV2) (relation.RelationV2, error) {
	rel, err := h.relationAdapter.TransformRelation(ctx, rlt)
	if err != nil {
		return relation.RelationV2{}, err
	}

	rel, err = h.relationService.Create(ctx, rel)
	if err != nil {
		return relation.RelationV2{}, err
	}

	return rel, nil
}

func (h Handler) GetRelation(ctx context.Context, request *shieldv1beta1.GetRelationRequest) (*shieldv1beta1.GetRelationResponse, error) {
	logger := grpczap.Extract(ctx)

	fetchedRelation, err := h.relationService.Get(ctx, request.GetId())
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, relation.ErrNotExist),
			errors.Is(err, relation.ErrInvalidUUID),
			errors.Is(err, relation.ErrInvalidID):
			return nil, grpcRelationNotFoundErr
		default:
			return nil, grpcInternalServerError
		}
	}

	relationPB, err := transformRelationV2ToPB(fetchedRelation)
	if err != nil {
		logger.Error(err.Error())
		return nil, grpcInternalServerError
	}

	return &shieldv1beta1.GetRelationResponse{
		Relation: &relationPB,
	}, nil
}

func (h Handler) DeleteRelation(ctx context.Context, request *shieldv1beta1.DeleteRelationRequest) (*shieldv1beta1.DeleteRelationResponse, error) {
	logger := grpczap.Extract(ctx)

	reln, err := h.relationService.GetRelationByFields(ctx, relation.RelationV2{
		Object: relation.Object{
			ID: request.GetObjectId(),
		},
		Subject: relation.Subject{
			ID:     request.GetSubjectId(),
			RoleID: request.GetRole(),
		},
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, relation.ErrNotExist):
			return nil, grpcRelationNotFoundErr
		default:
			return nil, grpcInternalServerError
		}
	}

	result, err := h.resourceService.CheckAuthz(ctx, resource.Resource{
		Name:        reln.Object.ID,
		NamespaceID: reln.Object.NamespaceID,
	}, action.Action{ID: schema.EditPermission})
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
		return nil, status.Errorf(codes.PermissionDenied, errpkg.ErrForbidden.Error())
	}

	err = h.relationService.DeleteV2(ctx, relation.RelationV2{
		Object: relation.Object{
			ID: request.GetObjectId(),
		},
		Subject: relation.Subject{
			ID:     request.GetSubjectId(),
			RoleID: request.GetRole(),
		},
	})
	if err != nil {
		logger.Error(err.Error())
		switch {
		case errors.Is(err, relation.ErrNotExist),
			errors.Is(err, relation.ErrInvalidUUID),
			errors.Is(err, relation.ErrInvalidID):
			return nil, grpcRelationNotFoundErr
		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrMissingEmail):
			return nil, grpcUnauthenticated
		default:
			return nil, grpcInternalServerError
		}
	}

	return &shieldv1beta1.DeleteRelationResponse{
		Message: "Relation deleted",
	}, nil
}

func transformRelationV2ToPB(relation relation.RelationV2) (shieldv1beta1.Relation, error) {
	return shieldv1beta1.Relation{
		Id:              relation.ID,
		ObjectId:        relation.Object.ID,
		ObjectNamespace: relation.Object.NamespaceID,
		Subject:         generateSubject(relation.Subject.ID, relation.Subject.Namespace),
		RoleName:        relation.Subject.RoleID,
		CreatedAt:       nil,
		UpdatedAt:       nil,
	}, nil
}

func extractSubjectFromPrincipal(principal string) (string, string) {
	splits := strings.Split(principal, ":")
	return splits[0], splits[1]
}

func generateSubject(subjectId, principal string) string {
	return fmt.Sprintf("%s:%s", principal, subjectId)
}

func (h Handler) verifyPrincipal(ctx context.Context, principal, subjectID string) error {
	var err error
	switch principal {
	case strings.Split(schema.UserPrincipal, "/")[1]:
		_, err = h.userService.Get(ctx, subjectID)
	case strings.Split(schema.GroupPrincipal, "/")[1]:
		_, err = h.groupService.Get(ctx, subjectID)
	}
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNotExist), errors.Is(err, group.ErrNotExist):
			return grpcBadBodyError
		default:
			return grpcInternalServerError
		}
	}

	return nil
}
