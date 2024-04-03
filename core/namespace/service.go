package namespace

import (
	"context"

	"github.com/goto/shield/core/user"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	AuditKeyNamespaceCreate = "namespace.create"
	AuditKeyNamespaceUpdate = "namespace.update"

	AuditEntity = "namespace"
)

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data map[string]string) error
}

type Service struct {
	repository      Repository
	userService     UserService
	activityService ActivityService
}

func NewService(repository Repository, userService UserService, activityService ActivityService) *Service {
	return &Service{
		repository:      repository,
		userService:     userService,
		activityService: activityService,
	}
}

func (s Service) Get(ctx context.Context, id string) (Namespace, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) Create(ctx context.Context, ns Namespace) (Namespace, error) {
	newNamespace, err := s.repository.Create(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := map[string]string{
		"entity":       AuditEntity,
		"id":           newNamespace.ID,
		"name":         newNamespace.Name,
		"backend":      newNamespace.Backend,
		"resourceType": newNamespace.ResourceType,
	}
	if err := s.activityService.Log(ctx, AuditKeyNamespaceCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return newNamespace, nil
}

func (s Service) List(ctx context.Context) ([]Namespace, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, ns Namespace) (Namespace, error) {
	updatedNamespace, err := s.repository.Update(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := map[string]string{
		"entity":       AuditEntity,
		"id":           updatedNamespace.ID,
		"name":         updatedNamespace.Name,
		"backend":      updatedNamespace.Backend,
		"resourceType": updatedNamespace.ResourceType,
	}
	if err := s.activityService.Log(ctx, AuditKeyNamespaceUpdate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return updatedNamespace, nil
}
