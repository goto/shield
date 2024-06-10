package inmemory

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/group"
)

var keyPrefix = "group"

type GroupRepository interface {
	GetBySlug(ctx context.Context, slug string) (group.Group, error)
}

type CachedGroupRepository struct {
	cache      Cache
	repository GroupRepository
}

func NewCachedGroupRepository(cache Cache, repository GroupRepository) *CachedGroupRepository {
	return &CachedGroupRepository{
		cache:      cache,
		repository: repository,
	}
}

func getKey(identifier string) string {
	return fmt.Sprintf("%s:%s", keyPrefix, identifier)
}

func (r CachedGroupRepository) GetBySlug(ctx context.Context, slug string) (group.Group, error) {
	key := getKey(slug)
	grp, found := r.cache.Get(key)
	if !found {
		grp, err := r.repository.GetBySlug(ctx, slug)
		if err != nil {
			return group.Group{}, err
		}

		r.cache.Set(key, grp, 0)
		return grp, nil
	}

	grpParsed, ok := grp.(group.Group)
	if !ok {
		return group.Group{}, ErrParsing
	}

	return grpParsed, nil
}
