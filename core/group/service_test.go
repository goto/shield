package group_test

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/core/group/mocks"
	"github.com/goto/shield/pkg/logger"
	"github.com/goto/shield/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testLogger = logger.InitLogger(logger.Config{
		Level:  "info",
		Format: "json",
	})
	testOrgID     = uuid.NewString()
	testGroupID   = uuid.NewString()
	testGroupSlug = "test-group-slug"
	testGroup     = group.Group{
		ID:             testGroupID,
		Name:           "Test Group",
		Slug:           testGroupSlug,
		OrganizationID: testOrgID,
	}
)

func TestService_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		idOrSlug string
		setup    func(t *testing.T) *group.Service
		want     group.Group
		wantErr  error
	}{
		{
			name:     "GetGroupById",
			idOrSlug: testGroupID,
			setup: func(t *testing.T) *group.Service {
				t.Helper()
				repository := &mocks.Repository{}
				cachedRepository := &mocks.CachedRepository{}
				relationService := &mocks.RelationService{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByID(mock.Anything, testGroupID).Return(testGroup, nil)
				return group.NewService(testLogger, repository, cachedRepository, relationService, userService, activityService)
			},
			want: testGroup,
		},
		{
			name:     "GetGroupBySlug",
			idOrSlug: testGroupSlug,
			setup: func(t *testing.T) *group.Service {
				t.Helper()
				repository := &mocks.Repository{}
				cachedRepository := &mocks.CachedRepository{}
				relationService := &mocks.RelationService{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetBySlug(mock.Anything, testGroupSlug).Return(testGroup, nil)
				return group.NewService(testLogger, repository, cachedRepository, relationService, userService, activityService)
			},
			want: testGroup,
		},
		{
			name:     "GetGroupByIdErr",
			idOrSlug: testGroupID,
			setup: func(t *testing.T) *group.Service {
				t.Helper()
				repository := &mocks.Repository{}
				cachedRepository := &mocks.CachedRepository{}
				relationService := &mocks.RelationService{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetByID(mock.Anything, testGroupID).Return(group.Group{}, group.ErrNotExist)
				return group.NewService(testLogger, repository, cachedRepository, relationService, userService, activityService)
			},
			wantErr: group.ErrNotExist,
		},
		{
			name:     "GetGroupBySlugErr",
			idOrSlug: testGroupSlug,
			setup: func(t *testing.T) *group.Service {
				t.Helper()
				repository := &mocks.Repository{}
				cachedRepository := &mocks.CachedRepository{}
				relationService := &mocks.RelationService{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				repository.EXPECT().GetBySlug(mock.Anything, testGroupSlug).Return(group.Group{}, group.ErrNotExist)
				return group.NewService(testLogger, repository, cachedRepository, relationService, userService, activityService)
			},
			wantErr: group.ErrNotExist,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.Get(context.TODO(), tt.idOrSlug)
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

func TestService_GetBySlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		idOrSlug string
		setup    func(t *testing.T) *group.Service
		want     group.Group
		wantErr  error
	}{
		{
			name:     "GetBySlug",
			idOrSlug: testGroupSlug,
			setup: func(t *testing.T) *group.Service {
				t.Helper()
				repository := &mocks.Repository{}
				cachedRepository := &mocks.CachedRepository{}
				relationService := &mocks.RelationService{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				cachedRepository.EXPECT().GetBySlug(mock.Anything, testGroupSlug).Return(testGroup, nil)
				return group.NewService(testLogger, repository, cachedRepository, relationService, userService, activityService)
			},
			want: testGroup,
		},
		{
			name:     "GetBySlugErr",
			idOrSlug: testGroupSlug,
			setup: func(t *testing.T) *group.Service {
				t.Helper()
				repository := &mocks.Repository{}
				cachedRepository := &mocks.CachedRepository{}
				relationService := &mocks.RelationService{}
				userService := &mocks.UserService{}
				activityService := &mocks.ActivityService{}
				cachedRepository.EXPECT().GetBySlug(mock.Anything, testGroupSlug).Return(group.Group{}, group.ErrNotExist)
				return group.NewService(testLogger, repository, cachedRepository, relationService, userService, activityService)
			},
			wantErr: group.ErrNotExist,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			svc := tt.setup(t)

			assert.NotNil(t, svc)

			got, err := svc.GetBySlug(context.TODO(), tt.idOrSlug)
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
