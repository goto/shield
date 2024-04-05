package action

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/user"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	AuditKeyActionCreate = "action.create"
	AuditKeyActionUpdate = "action.update"
)

type UserService interface {
	FetchCurrentUser(ctx context.Context) (user.User, error)
}

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data map[string]interface{}) error
}

type Service struct {
	repository      Repository
	userService     UserService
	activityService ActivityService
	logger          *zap.SugaredLogger
}

func NewService(repository Repository, userService UserService, activityService ActivityService, logger *zap.SugaredLogger) *Service {
	return &Service{
		repository:      repository,
		userService:     userService,
		activityService: activityService,
		logger:          logger,
	}
}

func (s Service) Get(ctx context.Context, id string) (Action, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) Create(ctx context.Context, action Action) (Action, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Action{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	newAction, err := s.repository.Create(ctx, action)
	if err != nil {
		return Action{}, err
	}

	actionLogData := newAction.ToActionLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(actionLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyActionCreate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}

	return newAction, nil
}

func (s Service) List(ctx context.Context) ([]Action, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, id string, action Action) (Action, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Action{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	updatedAction, err := s.repository.Update(ctx, Action{
		Name:        action.Name,
		ID:          id,
		NamespaceID: action.NamespaceID,
	})
	if err != nil {
		return Action{}, err
	}

	actionLogData := updatedAction.ToActionLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(actionLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyActionUpdate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}

	return updatedAction, nil
}
