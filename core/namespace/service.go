package namespace

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
)

const (
	auditKeyNamespaceUpsert = "namespace.upsert"
	auditKeyNamespaceUpdate = "namespace.update"
)

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor activity.Actor, data any) error
}

type Service struct {
	logger          log.Logger
	repository      Repository
	userService     UserService
	activityService ActivityService
}

func NewService(logger log.Logger, repository Repository, userService UserService, activityService ActivityService) *Service {
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

func (s Service) Upsert(ctx context.Context, ns Namespace) (Namespace, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Namespace{}, err
	}

	newNamespace, err := s.repository.Upsert(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	go func() {
		ctx := context.WithoutCancel(ctx)
		namespaceLogData := newNamespace.ToLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyNamespaceUpsert, actor, namespaceLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
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
		return Namespace{}, err
	}

	updatedNamespace, err := s.repository.Update(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	go func() {
		ctx := context.WithoutCancel(ctx)
		namespaceLogData := updatedNamespace.ToLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyNamespaceUpdate, actor, namespaceLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return updatedNamespace, nil
}
