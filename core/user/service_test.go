package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/goto/shield/core/activity"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/core/user/mocks"
	shieldlogger "github.com/goto/shield/pkg/logger"
	"github.com/goto/shield/pkg/uuid"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		email   string
		user    user.User
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.User
		wantErr error
	}{
		{
			name:  "CreateUserWithUpperCase",
			email: "jane.doe@gotocompany.com",
			user: user.User{
				Name:  "John Doe",
				Email: "John.Doe@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{
						ID:    "test-id",
						Email: "jane.doe@gotocompany.com",
					}, nil)
				repository.EXPECT().
					Create(mock.Anything, user.User{
						Name:  "John Doe",
						Email: "John.Doe@gotocompany.com",
					}).
					Return(user.User{
						Name:  "John Doe",
						Email: "John.Doe@gotocompany.com",
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserCreate, activity.Actor{
						ID:    "test-id",
						Email: "jane.doe@gotocompany.com",
					}, user.LogData{Entity: "user", Name: "John Doe", Email: "John.Doe@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			want: user.User{
				Name:  "John Doe",
				Email: "John.Doe@gotocompany.com",
			},
			wantErr: nil,
		},
		{
			name:  "CreateUserSelfRegister",
			email: "jane.doe@gotocompany.com",
			user: user.User{
				Name:  "Jane Doe",
				Email: "jane.doe@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, user.ErrNotExist)
				repository.EXPECT().
					Create(mock.Anything, user.User{
						Name:  "Jane Doe",
						Email: "jane.doe@gotocompany.com",
					}).
					Return(user.User{
						Name:  "Jane Doe",
						Email: "jane.doe@gotocompany.com",
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserCreate, activity.Actor{
						Email: "jane.doe@gotocompany.com",
					}, user.LogData{Entity: "user", Name: "Jane Doe", Email: "jane.doe@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			want: user.User{
				Name:  "Jane Doe",
				Email: "jane.doe@gotocompany.com",
			},
			wantErr: nil,
		},
		{
			name:  "CreateUserInvalidHeader",
			email: "jane.doe@gotocompany.com",
			user: user.User{
				Name:  "John Doe",
				Email: "john.doe@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, user.ErrNotExist)
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name: "CreateUserMissingHeader",
			user: user.User{
				Name:  "John Doe",
				Email: "john.doe@gotocompany.com",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "").
					Return(user.User{}, user.ErrMissingEmail)
				return user.NewService(logger, user.Config{}, repository, activityService)
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

			ctx := user.SetContextWithEmail(context.TODO(), tt.email)
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

func TestService_CreateMetadataKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		email       string
		metadataKey user.UserMetadataKey
		setup       func(t *testing.T) *user.Service
		want        user.UserMetadataKey
		wantErr     error
	}{
		{
			name:  "CreateMetadataKey",
			email: "jane.doe@gotocompany.com",
			metadataKey: user.UserMetadataKey{
				Key:         "test-key",
				Description: "description for test-key",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					CreateMetadataKey(mock.Anything, user.UserMetadataKey{
						Key:         "test-key",
						Description: "description for test-key",
					}).
					Return(user.UserMetadataKey{
						Key:         "test-key",
						Description: "description for test-key",
					}, nil).Once()
				activityService.EXPECT().
					Log(mock.Anything, "user_metadata_key.create", activity.Actor{}, user.MetadataKeyLogData{
						Entity:      "user_metadata_key",
						Key:         "test-key",
						Description: "description for test-key",
					}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			want: user.UserMetadataKey{
				Key:         "test-key",
				Description: "description for test-key",
			},
		},
		{
			name:  "CreateMetadataKeyError",
			email: "jane.doe@gotocompany.com",
			metadataKey: user.UserMetadataKey{
				Key:         "test-key",
				Description: "description for test-key",
			},
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				activityService := &mocks.ActivityService{}
				repository := &mocks.Repository{}
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					CreateMetadataKey(mock.Anything, user.UserMetadataKey{
						Key:         "test-key",
						Description: "description for test-key",
					}).
					Return(user.UserMetadataKey{}, user.ErrConflict).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			wantErr: user.ErrConflict,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.TODO(), "jane.doe@gotocompany.com")
			got, err := svc.CreateMetadataKey(ctx, tt.metadataKey)
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

func TestService_List(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		email   string
		setup   func(t *testing.T) *user.Service
		urn     string
		want    user.PagedUsers
		wantErr error
	}{
		{
			name:  "ListUser",
			email: "jane.doe@gotocompany.com",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					List(mock.Anything, user.Filter{}).
					Return([]user.User{
						{
							Name:  "John Doe",
							Email: "john.doe@gotocompany.com",
						},
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserCreate, activity.Actor{}, user.LogData{Entity: "user", Name: "John Doe", Email: "john.doe@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			want: user.PagedUsers{
				Users: []user.User{
					{
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com",
					},
				},
				Count: 1,
			},
			wantErr: nil,
		},
		{
			name:  "ListUserError",
			email: "jane.doe@gotocompany.com",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					List(mock.Anything, user.Filter{}).
					Return([]user.User{}, user.ErrInvalidID).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserCreate, activity.Actor{}, user.LogData{Entity: "user", Name: "John Doe", Email: "john.doe@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			// Any error is directly returned by the function
			wantErr: user.ErrInvalidID,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := user.SetContextWithEmail(context.TODO(), tt.email)
			got, err := svc.List(ctx, user.Filter{})
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
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					UpdateByID(mock.Anything, user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "John.Doe2@gotocompany.com",
					}).
					Return(user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "John.Doe2@gotocompany.com",
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, activity.Actor{}, user.LogData{Entity: "user", Name: "John Doe", Email: "John.Doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			want: user.User{
				ID:    "1",
				Name:  "John Doe",
				Email: "John.Doe2@gotocompany.com",
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
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "jane.doe@gotocompany.com").
					Return(user.User{}, nil)
				repository.EXPECT().
					UpdateByEmail(mock.Anything, user.User{
						Name:  "John Doe",
						Email: "John.Doe2@gotocompany.com",
					}).
					Return(user.User{
						Name:  "John Doe",
						Email: "John.Doe2@gotocompany.com",
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, activity.Actor{}, user.LogData{Entity: "user", Name: "John Doe", Email: "John.Doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
			},
			want: user.User{
				Name:  "John Doe",
				Email: "John.Doe2@gotocompany.com",
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
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "john.doe@gotocompany.com").
					Return(user.User{
						ID:    "1",
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com",
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, "", user.LogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
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
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByID(mock.Anything, "qwer-1234-tyui-5678-opas-90").
					Return(user.User{
						ID:    "qwer-1234-tyui-5678-opas-90",
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com",
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, "", user.LogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
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
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository.EXPECT().
					GetByEmail(mock.Anything, "john.doe@gotocompany.com").
					Return(user.User{
						ID:    "qwer-1234-tyui-5678-opas-90",
						Name:  "John Doe",
						Email: "john.doe@gotocompany.com",
					}, nil).Once()

				activityService.EXPECT().
					Log(mock.Anything, user.AuditKeyUserUpdate, "", user.LogData{Entity: "user", Name: "John Doe", Email: "john.doe2@gotocompany.com"}).Return(nil).Once()
				return user.NewService(logger, user.Config{}, repository, activityService)
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
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				return user.NewService(logger, user.Config{}, repository, activityService)
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

func TestService_DeleteUser(t *testing.T) {
	t.Parallel()

	testUserID := uuid.NewString()
	testUserEmail := "john.doe@gotocompany.com"
	emailTag := "inactive"
	tests := []struct {
		name    string
		setup   func(t *testing.T) *user.Service
		id      string
		wantErr error
	}{
		{
			name: "return error from delete by id",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().
					DeleteByID(mock.Anything, testUserID, emailTag).
					Return(user.ErrNotExist)
				return user.NewService(logger, user.Config{InactiveEmailTag: emailTag}, repository, activityService)
			},
			id:      testUserID,
			wantErr: user.ErrNotExist,
		},
		{
			name: "return error from delete by email",
			setup: func(t *testing.T) *user.Service {
				t.Helper()
				logger := shieldlogger.InitLogger(shieldlogger.Config{})
				repository := &mocks.Repository{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().
					DeleteByEmail(mock.Anything, testUserEmail, emailTag).
					Return(user.ErrNotExist)
				return user.NewService(logger, user.Config{InactiveEmailTag: emailTag}, repository, activityService)
			},
			id:      testUserEmail,
			wantErr: user.ErrNotExist,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			ctx := context.Background()

			err := svc.Delete(ctx, tt.id)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tt.wantErr))
		})
	}
}
