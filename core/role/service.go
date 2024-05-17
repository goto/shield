package role

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
	pkgctx "github.com/goto/shield/pkg/context"
)

const (
	auditKeyRoleCreate = "role.create"
	auditKeyRoleUpdate = "role.update"
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

func (s Service) Create(ctx context.Context, toCreate Role) (Role, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Role{}, err
	}

	roleID, err := s.repository.Create(ctx, toCreate)
	if err != nil {
		return Role{}, err
	}

	newRole, err := s.repository.Get(ctx, roleID)
	if err != nil {
		return Role{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		roleLogData := newRole.ToRoleLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyRoleCreate, actor, roleLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return newRole, nil
}

func (s Service) Get(ctx context.Context, id string) (Role, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) List(ctx context.Context) ([]Role, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, toUpdate Role) (Role, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return Role{}, err
	}

	roleID, err := s.repository.Update(ctx, toUpdate)
	if err != nil {
		return Role{}, err
	}

	updatedRole, err := s.repository.Get(ctx, roleID)
	if err != nil {
		return Role{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		roleLogData := updatedRole.ToRoleLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyRoleUpdate, actor, roleLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return updatedRole, nil
}
