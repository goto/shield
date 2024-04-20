package policy

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/user"
	pkgctx "github.com/goto/shield/pkg/context"
)

const (
	auditKeyPolicyCreate = "policy.create"
	auditKeyPolicyUpdate = "policy.update"
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

func (s Service) Get(ctx context.Context, id string) (Policy, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) List(ctx context.Context) ([]Policy, error) {
	return s.repository.List(ctx)
}

// Create is actually does upsert, this is called every period to sync rules and resources from buckets
// the periodic jobs do not need to logs the activity to avoid spamming activity logs
// we could remove this option once we have on-demand approach of syncing rules and resources
func (s Service) Create(ctx context.Context, policy Policy, opts ...ServiceOption) ([]Policy, error) {
	opt := &serviceOpts{}

	for _, f := range opts {
		f(opt)
	}

	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", user.ErrInvalidEmail.Error(), err.Error()))
	}

	policyId, err := s.repository.Create(ctx, policy)
	if err != nil {
		return []Policy{}, err
	}
	policies, err := s.repository.List(ctx)
	if err != nil {
		return []Policy{}, err
	}

	if opt.withActivityLogs {
		go func() {
			ctx := pkgctx.WithoutCancel(ctx)
			policyLogData := policy.ToPolicyLogData(policyId)
			if err := s.activityService.Log(ctx, auditKeyPolicyCreate, currentUser.ID, policyLogData); err != nil {
				s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
			}
		}()
	}

	return policies, err
}

func (s Service) Update(ctx context.Context, pol Policy) ([]Policy, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		s.logger.Error(fmt.Sprintf("%s: %s", user.ErrInvalidEmail.Error(), err.Error()))
	}

	policyId, err := s.repository.Update(ctx, pol)
	if err != nil {
		return []Policy{}, err
	}

	policies, err := s.repository.List(ctx)
	if err != nil {
		return []Policy{}, err
	}

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		policyLogData := pol.ToPolicyLogData(policyId)
		if err := s.activityService.Log(ctx, auditKeyPolicyUpdate, currentUser.ID, policyLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return policies, err
}
