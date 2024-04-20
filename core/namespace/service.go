package namespace

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/user"
	pkgctx "github.com/goto/shield/pkg/context"
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

func (s Service) Get(ctx context.Context, id string) (Namespace, error) {
	return s.repository.Get(ctx, id)
}

// Create is actually does upsert, this is called every period to sync rules and resources from buckets
// the periodic jobs do not need to logs the activity to avoid spamming activity logs
// we could remove this option once we have on-demand approach of syncing rules and resources
func (s Service) Create(ctx context.Context, ns Namespace, opts ...ServiceOption) (Namespace, error) {
	opt := &serviceOpts{}

	for _, f := range opts {
		f(opt)
	}

	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", user.ErrInvalidEmail.Error(), err.Error()))
	}

	newNamespace, err := s.repository.Create(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	if opt.withActivityLogs {
		go func() {
			ctx := pkgctx.WithoutCancel(ctx)
			namespaceLogData := newNamespace.ToNameSpaceLogData()
			if err := s.activityService.Log(ctx, auditKeyNamespaceCreate, currentUser.ID, namespaceLogData); err != nil {
				s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
			}
		}()
	}

	return newNamespace, nil
}

func (s Service) List(ctx context.Context) ([]Namespace, error) {
	return s.repository.List(ctx)
}

func (s Service) Update(ctx context.Context, ns Namespace) (Namespace, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", user.ErrInvalidEmail.Error(), err.Error()))
	}

	updatedNamespace, err := s.repository.Update(ctx, ns)
	if err != nil {
		return Namespace{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		namespaceLogData := updatedNamespace.ToNameSpaceLogData()
		if err := s.activityService.Log(ctx, auditKeyNamespaceUpdate, currentUser.ID, namespaceLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return updatedNamespace, nil
}
