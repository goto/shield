package servicedata

import (
	"context"
	"fmt"
)

const auditEntityServiceDataKey = "service_data_key"

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
	Name        string
	Description string
	ResourceID  string
}

type ServiceData struct {
	ID          string
	NamespaceID string
	EntityID    string
	Key         Key
	Value       any
}

type KeyLogData struct {
	Entity      string `mapstructure:"entity"`
	URN         string `mapstructure:"urn"`
	ProjectSlug string `mapstructure:"project_slug"`
	Key         string `mapstructure:"key"`
	Description string `mapstructure:"description"`
}

type Filter struct {
	ID        string
	Namespace string
	Entities  []string
	EntityIDs [][]string
	Project   string
}

func CreateURN(projectSlug, keyName string) string {
	return fmt.Sprintf("%s:servicedata_key:%s", projectSlug, keyName)
}

func (key Key) ToKeyLogData() KeyLogData {
	return KeyLogData{
		Entity:      auditEntityServiceDataKey,
		URN:         key.URN,
		ProjectSlug: key.ProjectSlug,
		Key:         key.Name,
		Description: key.Description,
	}
}
