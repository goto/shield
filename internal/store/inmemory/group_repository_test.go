package inmemory

import (
	"context"
	"errors"
	"testing"

	"github.com/goto/shield/core/group"
	"github.com/goto/shield/internal/store/inmemory/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	testCacheConfig = Config{
		NumCounters:  10000000,
		MaxCost:      1073741824,
		BufferItems:  64,
		Metrics:      true,
		TTLInSeconds: 3600,
	}
	testGroupSlug = "test-group-slug"
	testGroup     = group.Group{
		ID:   "test-group-id",
		Slug: testGroupSlug,
		Name: "test group",
	}
)

func TestGetBySlug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		slug        string
		setup       func(t *testing.T) *CachedGroupRepository
		want        group.Group
		wantErr     error
	}{
		{
			description: "should retrieve group from cache",
			slug:        testGroupSlug,
			setup: func(t *testing.T) *CachedGroupRepository {
				t.Helper()
				groupRepository := &mocks.GroupRepository{}
				c, err := NewCache(testCacheConfig)
				if err != nil {
					return nil
				}
				c.Set(getKey(testGroupSlug), testGroup, 0)
				c.Wait()
				return NewCachedGroupRepository(c, groupRepository)
			},
			want: testGroup,
		},
		{
			description: "should retrieve group from repository",
			slug:        testGroupSlug,
			setup: func(t *testing.T) *CachedGroupRepository {
				t.Helper()
				groupRepository := &mocks.GroupRepository{}
				c, err := NewCache(testCacheConfig)
				if err != nil {
					return nil
				}
				groupRepository.EXPECT().GetBySlug(mock.Anything, testGroupSlug).
					Return(testGroup, nil)
				return NewCachedGroupRepository(c, groupRepository)
			},
			want: testGroup,
		},
		{
			description: "should return parse error if cache data invalid",
			slug:        testGroupSlug,
			setup: func(t *testing.T) *CachedGroupRepository {
				t.Helper()
				groupRepository := &mocks.GroupRepository{}
				c, err := NewCache(testCacheConfig)
				if err != nil {
					return nil
				}
				c.Set(getKey(testGroupSlug), "invalid-group-data", 0)
				c.Wait()
				return NewCachedGroupRepository(c, groupRepository)
			},
			wantErr: ErrParsing,
		},
		{
			description: "should return error from repository",
			slug:        testGroupSlug,
			setup: func(t *testing.T) *CachedGroupRepository {
				t.Helper()
				groupRepository := &mocks.GroupRepository{}
				c, err := NewCache(testCacheConfig)
				if err != nil {
					return nil
				}
				groupRepository.EXPECT().GetBySlug(mock.Anything, testGroupSlug).
					Return(group.Group{}, group.ErrInvalidDetail)
				return NewCachedGroupRepository(c, groupRepository)
			},
			wantErr: group.ErrInvalidDetail,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			cacheRepo := tc.setup(t)
			assert.NotNil(t, cacheRepo)

			ctx := context.TODO()
			got, err := cacheRepo.GetBySlug(ctx, tc.slug)
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr))
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, got)
		})
	}
}
