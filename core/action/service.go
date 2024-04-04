package action

import (
	"context"

	"github.com/goto/shield/core/user"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	AuditKeyActionCreate = "action.create"
	AuditKeyActionUpdate = "action.update"
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

func (s Service) Get(ctx context.Context, id string) (Action, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) Create(ctx context.Context, action Action) (Action, error) {
	newAction, err := s.repository.Create(ctx, action)
	if err != nil {
		return Action{}, err
	}

	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := newAction.ToActionLogData()
	if err := s.activityService.Log(ctx, AuditKeyActionCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return newAction, nil
}

func (s Service) List(ctx context.Context) ([]Action, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, id string, action Action) (Action, error) {
	updatedAction, err := s.repository.Update(ctx, Action{
		Name:        action.Name,
		ID:          id,
		NamespaceID: action.NamespaceID,
	})
	if err != nil {
		return Action{}, err
	}

	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := updatedAction.ToActionLogData()
	if err := s.activityService.Log(ctx, AuditKeyActionUpdate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return updatedAction, nil
}
