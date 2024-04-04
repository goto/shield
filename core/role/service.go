package role

import (
	"context"

	"github.com/goto/shield/core/user"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	AuditKeyRoleCreate = "role.create"
	AuditKeyRoleUpdate = "role.update"
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

func (s Service) Create(ctx context.Context, toCreate Role) (Role, error) {
	roleID, err := s.repository.Create(ctx, toCreate)
	if err != nil {
		return Role{}, err
	}

	newRole, err := s.repository.Get(ctx, roleID)
	if err != nil {
		return Role{}, err
	}

	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := newRole.ToRoleAuditData()
	if err := s.activityService.Log(ctx, AuditKeyRoleCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return newRole, nil
}

func (s Service) Get(ctx context.Context, id string) (Role, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) List(ctx context.Context) ([]Role, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, toUpdate Role) (Role, error) {
	roleID, err := s.repository.Update(ctx, toUpdate)
	if err != nil {
		return Role{}, err
	}

	updatedRole, err := s.repository.Get(ctx, roleID)
	if err != nil {
		return Role{}, err
	}

	currentUser, _ := s.userService.FetchCurrentUser(ctx)
	logData := updatedRole.ToRoleAuditData()
	if err := s.activityService.Log(ctx, AuditKeyRoleCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return updatedRole, nil
}
