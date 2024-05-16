package action

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
	pkgctx "github.com/goto/shield/pkg/context"
)

const (
	auditKeyActionCreate = "action.create"
	auditKeyActionUpdate = "action.update"
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

func (s Service) Get(ctx context.Context, id string) (Action, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) Create(ctx context.Context, action Action) (Action, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Action{}, err
	}

	newAction, err := s.repository.Create(ctx, action)
	if err != nil {
		return Action{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		actionLogData := newAction.ToActionLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyActionCreate, actor, actionLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return newAction, nil
}

func (s Service) List(ctx context.Context) ([]Action, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, id string, action Action) (Action, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Action{}, err
	}

	updatedAction, err := s.repository.Update(ctx, Action{
		Name:        action.Name,
		ID:          id,
		NamespaceID: action.NamespaceID,
	})
	if err != nil {
		return Action{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		actionLogData := updatedAction.ToActionLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyActionUpdate, actor, actionLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return updatedAction, nil
}
