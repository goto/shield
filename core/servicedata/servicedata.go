package servicedata

import (
	"context"
	"fmt"
)

type Repository interface {
	Transactor
	CreateKey(ctx context.Context, key Key) (Key, error)
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
	Key         string
	Description string
	ResourceID  string
}

func (key Key) CreateURN() string {
	return fmt.Sprintf("%s:servicedata_key:%s", key.ProjectID, key.Key)
}
