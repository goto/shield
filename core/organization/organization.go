package organization

import (
	"context"
	"time"

	"github.com/goto/shield/core/user"
	"github.com/goto/shield/pkg/metadata"
)

const AuditEntity = "organization"

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
	Entity string
	ID     string
	Name   string
	Slug   string
}

func (organization Organization) ToOrganizationLogData() OrganizationLogData {
	logData := OrganizationLogData{
		Entity: AuditEntity,
		ID:     organization.ID,
		Name:   organization.Name,
		Slug:   organization.Slug,
	}

	return logData
}
