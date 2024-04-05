package namespace

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/user"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	AuditKeyNamespaceCreate = "namespace.create"
	AuditKeyNamespaceUpdate = "namespace.update"
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

	namespaceLogData := newNamespace.ToNameSpaceLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(namespaceLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyNamespaceCreate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
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

	namespaceLogData := updatedNamespace.ToNameSpaceLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(namespaceLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyNamespaceUpdate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}

	return updatedNamespace, nil
}
