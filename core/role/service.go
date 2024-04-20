package role

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
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
	Log(ctx context.Context, action string, actor string, data any) error
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

type serviceOpts struct {
	withActivityLogs bool
}

type ServiceOption func(*serviceOpts)

// WithActivityLogs logs activity in the method
func WithActivityLogs() ServiceOption {
	return func(g *serviceOpts) {
		g.withActivityLogs = true
	}
}

// Create is actually does upsert, this is called every period to sync rules and resources from buckets
// the periodic jobs do not need to logs the activity to avoid spamming activity logs
// we could remove this option once we have on-demand approach of syncing rules and resources
func (s Service) Create(ctx context.Context, toCreate Role, opts ...ServiceOption) (Role, error) {
	opt := &serviceOpts{}

	for _, f := range opts {
		f(opt)
	}

	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", user.ErrInvalidEmail.Error(), err.Error()))
	}

	roleID, err := s.repository.Create(ctx, toCreate)
	if err != nil {
		return Role{}, err
	}

	newRole, err := s.repository.Get(ctx, roleID)
	if err != nil {
		return Role{}, err
	}

	if opt.withActivityLogs {
		go func() {
			ctx := pkgctx.WithoutCancel(ctx)
			roleLogData := newRole.ToRoleLogData()
			if err := s.activityService.Log(ctx, auditKeyRoleCreate, currentUser.ID, roleLogData); err != nil {
				s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
			}
		}()
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
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", user.ErrInvalidEmail.Error(), err.Error()))
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
		if err := s.activityService.Log(ctx, auditKeyRoleUpdate, currentUser.ID, roleLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return updatedRole, nil
}
