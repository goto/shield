package namespace

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/user"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	AuditKeyNamespaceCreate = "namespace.create"
	AuditKeyNamespaceUpdate = "namespace.update"
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
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Namespace{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	newNamespace, err := s.repository.Create(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	logData := newNamespace.ToNameSpaceAuditData()
	if err := s.activityService.Log(ctx, AuditKeyNamespaceCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
	}

	return newNamespace, nil
}

func (s Service) List(ctx context.Context) ([]Namespace, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, ns Namespace) (Namespace, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Namespace{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	updatedNamespace, err := s.repository.Update(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	logData := updatedNamespace.ToNameSpaceAuditData()
	if err := s.activityService.Log(ctx, AuditKeyNamespaceUpdate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
	}

	return updatedNamespace, nil
}
