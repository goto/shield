package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/goto/shield/pkg/uuid"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

const (
	AuditKeyUserCreate            = "user.create"
	AuditKeyUserUpdate            = "user.update"
	AuditKeyUserMetadataKeyCreate = "user_metadata_key.create"
)

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data map[string]interface{}) error
}

type Service struct {
	repository      Repository
	activityService ActivityService
	logger          *zap.SugaredLogger
}

func NewService(repository Repository, activityService ActivityService, logger *zap.SugaredLogger) *Service {
	return &Service{
		repository:      repository,
		activityService: activityService,
		logger:          logger,
	}
}

func (s Service) Get(ctx context.Context, idOrEmail string) (User, error) {
	if uuid.IsValid(idOrEmail) {
		return s.repository.GetByID(ctx, idOrEmail)
	}
	return s.repository.GetByEmail(ctx, idOrEmail)
}

func (s Service) GetByID(ctx context.Context, id string) (User, error) {
	return s.repository.GetByID(ctx, id)
}

func (s Service) GetByIDs(ctx context.Context, userIDs []string) ([]User, error) {
	return s.repository.GetByIDs(ctx, userIDs)
}

func (s Service) GetByEmail(ctx context.Context, email string) (User, error) {
	return s.repository.GetByEmail(ctx, email)
}

func (s Service) Create(ctx context.Context, user User) (User, error) {
	currentUser, err := s.FetchCurrentUser(ctx)
	if err != nil {
		return User{}, fmt.Errorf("%w: %s", ErrInvalidEmail, err.Error())
	}

	newUser, err := s.repository.Create(ctx, User{
		Name:     user.Name,
		Email:    user.Email,
		Metadata: user.Metadata,
	})
	if err != nil {
		return User{}, err
	}

	userLogData := newUser.ToUserLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(userLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyUserCreate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}

	return newUser, nil
}

func (s Service) CreateMetadataKey(ctx context.Context, key UserMetadataKey) (UserMetadataKey, error) {
	currentUser, err := s.FetchCurrentUser(ctx)
	if err != nil {
		return UserMetadataKey{}, fmt.Errorf("%w: %s", ErrInvalidEmail, err.Error())
	}

	newUserMetadataKey, err := s.repository.CreateMetadataKey(ctx, UserMetadataKey{
		Key:         key.Key,
		Description: key.Description,
	})
	if err != nil {
		return UserMetadataKey{}, err
	}

	userMetadataKeyLogData := newUserMetadataKey.ToUserMetadataKeyLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(userMetadataKeyLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyUserMetadataKeyCreate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}

	return newUserMetadataKey, nil
}

func (s Service) List(ctx context.Context, flt Filter) (PagedUsers, error) {
	users, err := s.repository.List(ctx, flt)
	if err != nil {
		return PagedUsers{}, err
	}
	//TODO might better to do this in handler level
	return PagedUsers{
		Count: int32(len(users)),
		Users: users,
	}, nil
}

func (s Service) UpdateByID(ctx context.Context, toUpdate User) (User, error) {
	currentUser, err := s.FetchCurrentUser(ctx)
	if err != nil {
		return User{}, fmt.Errorf("%w: %s", ErrInvalidEmail, err.Error())
	}

	updatedUser, err := s.repository.UpdateByID(ctx, toUpdate)
	if err != nil {
		return User{}, err
	}

	userLogData := updatedUser.ToUserLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(userLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyUserUpdate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}

	return updatedUser, nil
}

func (s Service) UpdateByEmail(ctx context.Context, toUpdate User) (User, error) {
	currentUser, err := s.FetchCurrentUser(ctx)
	if err != nil {
		return User{}, fmt.Errorf("%w: %s", ErrInvalidEmail, err.Error())
	}

	updatedUser, err := s.repository.UpdateByEmail(ctx, toUpdate)
	if err != nil {
		return User{}, err
	}

	userLogData := updatedUser.ToUserLogData()
	var logDataMap map[string]interface{}
	if err := mapstructure.Decode(userLogData, &logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}
	if err := s.activityService.Log(ctx, AuditKeyUserUpdate, currentUser.ID, logDataMap); err != nil {
		s.logger.Errorf("%s: %s", ErrLogActivity.Error(), err.Error())
	}

	return updatedUser, nil
}

func (s Service) FetchCurrentUser(ctx context.Context) (User, error) {
	email, ok := GetEmailFromContext(ctx)
	if !ok {
		return User{}, ErrMissingEmail
	}

	email = strings.TrimSpace(email)
	if email == "" {
		return User{}, ErrMissingEmail
	}

	fetchedUser, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return fetchedUser, nil
}
