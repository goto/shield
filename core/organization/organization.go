package organization

import (
	"context"
	"time"

	"github.com/goto/shield/core/user"
	"github.com/goto/shield/pkg/metadata"
	"golang.org/x/exp/maps"
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

func (organization Organization) ToOrganizationAuditData() (map[string]string, error) {
	logData := map[string]string{
		"entity": AuditEntity,
		"id":     organization.ID,
		"name":   organization.Name,
		"slug":   organization.Slug,
	}
	orgMetadata, err := organization.Metadata.ToStringValueMap()
	if err != nil {
		return logData, err
	}

	maps.Copy(logData, orgMetadata)
	return logData, nil
}
