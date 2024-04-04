package namespace

import (
	"context"

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
	newNamespace, err := s.repository.Create(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := newNamespace.ToNameSpaceAuditData()
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
	logData := updatedNamespace.ToNameSpaceAuditData()
	if err := s.activityService.Log(ctx, AuditKeyNamespaceUpdate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return updatedNamespace, nil
}
