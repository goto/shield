package policy

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/user"
	pkgctx "github.com/goto/shield/pkg/context"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	auditKeyPolicyCreate = "policy.create"
	auditKeyPolicyUpdate = "policy.update"
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

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		policyLogData := policy.ToPolicyLogData(policyId)
		var logDataMap map[string]interface{}
		if err := mapstructure.Decode(policyLogData, &logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
		if err := s.activityService.Log(ctx, auditKeyPolicyCreate, currentUser.ID, logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

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

	go func() {
		ctx := pkgctx.WithoutCancel(ctx)
		policyLogData := pol.ToPolicyLogData(policyId)
		var logDataMap map[string]interface{}
		if err := mapstructure.Decode(policyLogData, &logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
		if err := s.activityService.Log(ctx, auditKeyPolicyUpdate, currentUser.ID, logDataMap); err != nil {
			s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
		}
	}()

	return policies, err
}
