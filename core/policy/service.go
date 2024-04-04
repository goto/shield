package policy

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/user"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	AuditKeyPolicyCreate = "policy.create"
	AuditKeyPolicyUpdate = "policy.update"
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

func (s Service) Get(ctx context.Context, id string) (Policy, error) {
	return s.repository.Get(ctx, id)
}

func (s Service) List(ctx context.Context) ([]Policy, error) {
	return s.repository.List(ctx)
}

func (s Service) Create(ctx context.Context, policy Policy) ([]Policy, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return []Policy{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	policyId, err := s.repository.Create(ctx, policy)
	if err != nil {
		return []Policy{}, err
	}
	policies, err := s.repository.List(ctx)
	if err != nil {
		return []Policy{}, err
	}

	logData := policy.ToPolicyLogData(policyId)
	if err := s.activityService.Log(ctx, AuditKeyPolicyCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
	}

	return policies, err
}

func (s Service) Update(ctx context.Context, pol Policy) ([]Policy, error) {
	currentUser, err := s.userService.FetchCurrentUser(ctx)
	if err != nil {
		return []Policy{}, fmt.Errorf("%w: %s", user.ErrInvalidEmail, err.Error())
	}

	policyId, err := s.repository.Update(ctx, pol)
	if err != nil {
		return []Policy{}, err
	}

	policies, err := s.repository.List(ctx)
	if err != nil {
		return []Policy{}, err
	}

	logData := pol.ToPolicyLogData(policyId)
	if err := s.activityService.Log(ctx, AuditKeyPolicyUpdate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(fmt.Sprintf("%s: %s", ErrLogActivity.Error(), err.Error()))
	}

	return policies, err
}
