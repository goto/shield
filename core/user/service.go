package user

import (
	"context"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/goto/shield/pkg/uuid"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

const (
	AuditKeyUserCreate            = "user.create"
	AuditKeyUserUpdate            = "user.update"
	AuditKeyUserMetadataKeyCreate = "user_metadata_key.create"

	AuditEntityUser         = "user"
	AuditEntityUserMetadata = "user_metadata_key"
)

type ActivityService interface {
	Log(ctx context.Context, action string, actor string, data map[string]string) error
}

type Service struct {
	repository      Repository
	activityService ActivityService
}

func NewService(repository Repository, activityService ActivityService) *Service {
	return &Service{
		repository:      repository,
		activityService: activityService,
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
	newUser, err := s.repository.Create(ctx, User{
		Name:     user.Name,
		Email:    user.Email,
		Metadata: user.Metadata,
	})
	if err != nil {
		return User{}, err
	}

	currentUser, _ := s.FetchCurrentUser(ctx)
	logData := map[string]string{
		"entity": AuditEntityUser,
		"name":   newUser.Name,
		"email":  newUser.Email,
	}
	maps.Copy(logData, newUser.Metadata.ToStringValueMap())
	if err := s.activityService.Log(ctx, AuditKeyUserCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
	}

	return newUser, nil
}

func (s Service) CreateMetadataKey(ctx context.Context, key UserMetadataKey) (UserMetadataKey, error) {
	newUserMetadataKey, err := s.repository.CreateMetadataKey(ctx, UserMetadataKey{
		Key:         key.Key,
		Description: key.Description,
	})
	if err != nil {
		return UserMetadataKey{}, err
	}

	currentUser, _ := s.FetchCurrentUser(ctx)
	logData := map[string]string{
		"entity":      AuditEntityUserMetadata,
		"key":         newUserMetadataKey.Key,
		"description": newUserMetadataKey.Description,
	}
	if err := s.activityService.Log(ctx, AuditKeyUserMetadataKeyCreate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
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
	updatedUser, err := s.repository.UpdateByID(ctx, toUpdate)
	if err != nil {
		return User{}, err
	}

	currentUser, _ := s.FetchCurrentUser(ctx)
	logData := map[string]string{
		"entity": AuditEntityUser,
		"name":   updatedUser.Name,
		"email":  updatedUser.Email,
	}
	maps.Copy(logData, updatedUser.Metadata.ToStringValueMap())
	if err := s.activityService.Log(ctx, AuditKeyUserUpdate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(err.Error())
	}

	return updatedUser, nil
}

func (s Service) UpdateByEmail(ctx context.Context, toUpdate User) (User, error) {
	updatedUser, err := s.repository.UpdateByEmail(ctx, toUpdate)
	if err != nil {
		return User{}, err
	}

	currentUser, _ := s.FetchCurrentUser(ctx)
	logData := map[string]string{
		"entity": AuditEntityUser,
		"name":   updatedUser.Name,
		"email":  updatedUser.Email,
	}
	maps.Copy(logData, updatedUser.Metadata.ToStringValueMap())
	if err := s.activityService.Log(ctx, AuditKeyUserUpdate, currentUser.ID, logData); err != nil {
		logger := grpczap.Extract(ctx)
		logger.Error(ErrLogActivity.Error())
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
