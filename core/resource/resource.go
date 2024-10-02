package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/internal/schema"
)

const (
	NON_RESOURCE_ID               = "*"
	RESOURCES_CONFIG_STORAGE_PG   = "postgresql"
	RESOURCES_CONFIG_STORAGE_GS   = "gs"
	RESOURCES_CONFIG_STORAGE_FILE = "file"
	RESOURCES_CONFIG_STORAGE_MEM  = "mem"

	AuditEntity = "resource"
)

type Repository interface {
	Transactor
	GetByID(ctx context.Context, id string) (Resource, error)
	GetByURN(ctx context.Context, urn string) (Resource, error)
	Upsert(ctx context.Context, resource Resource) (Resource, error)
	Create(ctx context.Context, resource Resource) (Resource, error)
	List(ctx context.Context, flt Filter) ([]Resource, error)
	Update(ctx context.Context, id string, resource Resource) (Resource, error)
	GetByNamespace(ctx context.Context, name string, ns string) (Resource, error)
}

type Transactor interface {
	WithTransaction(ctx context.Context) context.Context
	Rollback(ctx context.Context, err error) error
	Commit(ctx context.Context) error
}

type SchemaRepository interface {
	UpsertResourceConfigs(ctx context.Context, name string, config schema.NamespaceConfigMapType) (ResourceConfig, error)
}

type Resource struct {
	Idxa           string
	URN            string
	Name           string
	ProjectID      string
	OrganizationID string
	NamespaceID    string
	UserID         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (res Resource) CreateURN() string {
	isSystemNS := namespace.IsSystemNamespaceID(res.NamespaceID)
	if isSystemNS {
		return res.Name
	}
	if res.Name == NON_RESOURCE_ID {
		return fmt.Sprintf("p/%s/%s", res.ProjectID, res.NamespaceID)
	}
	return fmt.Sprintf("r/%s/%s", res.NamespaceID, res.Name)
}

type Filter struct {
	ProjectID      string
	GroupID        string
	OrganizationID string
	NamespaceID    string
	Limit          int32
	Page           int32
}

type YAML struct {
	Name         string              `json:"name" yaml:"name"`
	Backend      string              `json:"backend" yaml:"backend"`
	ResourceType string              `json:"resource_type" yaml:"resource_type"`
	Actions      map[string][]string `json:"actions" yaml:"actions"`
}

type PagedResources struct {
	Count     int32
	Resources []Resource
}

type ResourcePermissions = map[string][]string

type ResourceConfig struct {
	ID        uint32
	Name      string
	Config    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AppConfig struct {
	ConfigStorage string
}

type LogData struct {
	Entity         string `mapstructure:"entity"`
	URN            string `mapstructure:"urn"`
	Name           string `mapstructure:"name"`
	OrganizationID string `mapstructure:"organization_id"`
	ProjectID      string `mapstructure:"project_id"`
	NamespaceID    string `mapstructure:"namespace_id"`
	UserID         string `mapstructure:"user_id"`
}

func (resource Resource) ToLogData() LogData {
	return LogData{
		Entity:         AuditEntity,
		URN:            resource.URN,
		Name:           resource.Name,
		OrganizationID: resource.OrganizationID,
		ProjectID:      resource.ProjectID,
		NamespaceID:    resource.NamespaceID,
		UserID:         resource.UserID,
	}
}
