package organization

import (
	"context"
	"time"

	"github.com/goto/shield/core/user"
	"github.com/goto/shield/pkg/metadata"
)

const auditEntity = "organization"

type Repository interface {
	GetByID(ctx context.Context, id string) (Organization, error)
	GetBySlug(ctx context.Context, slug string) (Organization, error)
	Create(ctx context.Context, org Organization) (Organization, error)
	List(ctx context.Context) ([]Organization, error)
	UpdateByID(ctx context.Context, org Organization) (Organization, error)
	UpdateBySlug(ctx context.Context, org Organization) (Organization, error)
	ListAdminsByOrgID(ctx context.Context, id string) ([]user.User, error)
}

type Organization struct {
	ID        string
	Name      string
	Slug      string
	Metadata  metadata.Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrganizationLogData struct {
	Entity string `mapstructure:"entity"`
	ID     string `mapstructure:"id"`
	Name   string `mapstructure:"name"`
	Slug   string `mapstructure:"slug"`
}

func (organization Organization) ToOrganizationLogData() OrganizationLogData {
	logData := OrganizationLogData{
		Entity: auditEntity,
		ID:     organization.ID,
		Name:   organization.Name,
		Slug:   organization.Slug,
	}

	return logData
}
