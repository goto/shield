package servicedata

import (
	"context"
	"fmt"

	"github.com/goto/shield/core/user"
)

type Repository interface {
	Transactor
	CreateKey(ctx context.Context, key Key) (Key, error)
	Upsert(ctx context.Context, servicedata ServiceData) (ServiceData, error)
	GetKeyByURN(ctx context.Context, URN string) (Key, error)
	Get(ctx context.Context, filter Filter) ([]ServiceData, error)
	ListUsers(ctx context.Context, filter ListUsersFilter, servicedataKeyResourceIds []string) ([]user.User, error)
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context, err error) error
	Commit(ctx context.Context) error
}

type Key struct {
	ID          string
	URN         string
	ProjectID   string
	ProjectSlug string
	Key         string
	Description string
	ResourceID  string
}

type ServiceData struct {
	ID          string
	NamespaceID string
	EntityID    string
	Key         Key
	Value       string
}

type Filter struct {
	ID        string
	Namespace string
	Entities  []string
	EntityIDs [][]string
	Project   string
}

type ListUsersFilter struct {
	ServiceData     map[string]string
	Project         string
	Limit           int32
	Page            int32
	WithServiceData bool
}

func (key Key) CreateURN() string {
	return fmt.Sprintf("%s:servicedata_key:%s", key.ProjectSlug, key.Key)
}
