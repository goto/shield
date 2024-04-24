package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
	shieldlogger "github.com/goto/shield/pkg/logger"

	"github.com/goto/shield/core/mocks"
	"github.com/goto/shield/pkg/logger"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    user.User
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name: "CreateUserWithUpperCase",
			user: user.User{
				Name:  "John Doe",
				Email: "John.Doe@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(logger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					Create(mock.Anything, user.User{
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com"}).
					Return(user.User{
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com"}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserCreate, activity.Actor{}, user.UserLogData{Entity: "user", Name: "John Doe", Email: "john.doe@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, repository, activityService)
			},
			want: user.User{
				Name:  "John Doe",
				Email: "john.doe@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.TODO(), "jane.doe@gotocompany.com")
			got, err := svc.Create(ctx, tt.user)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_UpdateByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    user.User
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name: "UpdateUserWithUpperCase",
			user: user.User{
				ID:    "1",
				Name:  "John Doe",
				Email: "John.Doe2@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(logger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					UpdateByID(mock.Anything, user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}).
					Return(user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, activity.Actor{}, user.UserLogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, repository, activityService)
			},
			want: user.User{
				ID:    "1",
				Name:  "John Doe",
				Email: "john.doe2@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.TODO(), "jane.doe@gotocompany.com")
			got, err := svc.UpdateByID(ctx, tt.user)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_UpdateByEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		user    user.User
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name: "UpdateUserWithUpperCase",
			user: user.User{
				Name:  "John Doe",
				Email: "John.Doe2@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(logger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					UpdateByEmail(mock.Anything, user.User{
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}).
					Return(user.User{
						Name:  "John Doe",
						Email: "john.doe2@gotocompany.com"}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, activity.Actor{}, user.UserLogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, repository, activityService)
			},
			want: user.User{
				Name:  "John Doe",
				Email: "john.doe2@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.TODO(), "jane.doe@gotocompany.com")
			got, err := svc.UpdateByEmail(ctx, tt.user)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetByEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		email   string
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name:  "GetUserByEmail",
			email: "john.doe@gotocompany.com",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(logger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "john.doe@gotocompany.com").
					Return(user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com"}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, "", user.UserLogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, repository, activityService)
			},
			want: user.User{
				ID:    "1",
				Name:  "John Doe",
				Email: "john.doe@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.GetByEmail(context.Background(), tt.email)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_GetByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		uuid    string
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name: "GetUserByEmail",
			uuid: "qwer-1234-tyui-5678-opas-90",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(logger.Config{})
				repository.EXPECT().
					GetByID(mock.Anything, "qwer-1234-tyui-5678-opas-90").
					Return(user.User{
						ID:    "qwer-1234-tyui-5678-opas-90",
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com"}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, "", user.UserLogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, repository, activityService)
			},
			want: user.User{
				ID:    "qwer-1234-tyui-5678-opas-90",
				Name:  "John Doe",
				Email: "john.doe@gotocompany.com",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.GetByID(context.Background(), tt.uuid)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_FetchCurrentUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		email   string
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name:  "FetchCurrentUser",
			email: "john.doe@gotocompany.com",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(logger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "john.doe@gotocompany.com").
					Return(user.User{
						ID:    "qwer-1234-tyui-5678-opas-90",
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com"}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, "", user.UserLogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, repository, activityService)
			},
			want: user.User{
				ID:    "qwer-1234-tyui-5678-opas-90",
				Name:  "John Doe",
				Email: "john.doe@gotocompany.com",
			},
			wantErr: nil,
		},
		{
			name: "FetchCurrentUserMissingEmail",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(logger.Config{})
				return user.NewService(logger, repository, activityService)
			},
			wantErr: user.ErrMissingEmail,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.Background(), tt.email)

			got, err := svc.FetchCurrentUser(ctx)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
