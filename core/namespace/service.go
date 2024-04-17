package namespace

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/user"
	pkgctx "github.com/goto/shield/pkg/context"
	"go.uber.org/zap"
)

const (
	auditKeyNamespaceCreate = "namespace.create"
	auditKeyNamespaceUpdate = "namespace.update"
)

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data any) error
}

type Service struct {
	logger          *zap.SugaredLogger
	repository      Repository
	userService     UserService
	activityService ActivityService
}

func NewService(logger *zap.SugaredLogger, repository Repository, userService UserService, activityService ActivityService) *Service {
	return &Service{
		logger:          logger,
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

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		namespaceLogData := newNamespace.ToNameSpaceLogData()
		if err := s.activityService.Log(ctx, auditKeyNamespaceCreate, currentUser.ID, namespaceLogData); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

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

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		namespaceLogData := updatedNamespace.ToNameSpaceLogData()
		if err := s.activityService.Log(ctx, auditKeyNamespaceUpdate, currentUser.ID, namespaceLogData); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return updatedNamespace, nil
}
