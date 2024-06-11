package policy

import (
	"context"
	"fmt"

	"github.com/goto/salt/log"
	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
)

const (
	auditKeyPolicyUpsert = "policy.upsert"
	auditKeyPolicyUpdate = "policy.update"
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

func (s Service) Get(ctx context.Context, id string) (Policy, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) List(ctx context.Context) ([]Policy, error) {
	return s.repository.List(ctx)
}

func (s Service) Upsert(ctx context.Context, pol *Policy) ([]Policy, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return []Policy{}, err
	}

	policyID, err := s.repository.Upsert(ctx, pol)
	if err != nil {
		return []Policy{}, err
	}
	pol.ID = policyID

	policies, err := s.repository.List(ctx)
	if err != nil {
		return []Policy{}, err
	}

	go func() {
		ctx := context.WithoutCancel(ctx)
		policyLogData := pol.ToLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyPolicyUpsert, actor, policyLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return policies, err
}

func (s Service) Update(ctx context.Context, pol *Policy) ([]Policy, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return []Policy{}, err
	}

	policyID, err := s.repository.Update(ctx, pol)
	if err != nil {
		return []Policy{}, err
	}
	pol.ID = policyID

	policies, err := s.repository.List(ctx)
	if err != nil {
		return []Policy{}, err
	}

	go func() {
		ctx := context.WithoutCancel(ctx)
		policyLogData := pol.ToLogData()
		actor := activity.Actor{ID: currentUser.ID, Email: currentUser.Email}
		if err := s.activityService.Log(ctx, auditKeyPolicyUpdate, actor, policyLogData); err != nil {
			s.logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
		}
	}()

	return policies, err
}
