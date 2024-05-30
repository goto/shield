package servicedata

import (
	"context"
	"fmt"
)

type Repository interface {
	Transactor
	CreateKey(ctx context.Context, key Key) (Key, error)
	Upsert(ctx context.Context, servicedata ServiceData) (ServiceData, error)
	GetKeyByURN(ctx context.Context, URN string) (Key, error)
	Get(ctx context.Context, filter Filter) ([]ServiceData, error)
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
	Entity    []string
	EntityIDs [][]string
	Project   string
}

func (key Key) CreateURN() string {
	return fmt.Sprintf("%s:servicedata_key:%s", key.ProjectSlug, key.Key)
}
