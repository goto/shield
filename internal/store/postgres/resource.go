package postgres

import (
	"database/sql"
	"time"

	"github.com/goto/shield/core/resource"
	"github.com/goto/shield/internal/schema"
)

type Resource struct {
	ID             string         `db:"id"`
	URN            string         `db:"urn"`
	Name           string         `db:"name"`
	ProjectID      string         `db:"project_id"`
	Project        Project        `db:"project"`
	OrganizationID string         `db:"org_id"`
	Organization   Organization   `db:"organization"`
	NamespaceID    string         `db:"namespace_id"`
	Namespace      Namespace      `db:"namespace"`
	User           User           `db:"user"`
	UserID         sql.NullString `db:"user_id"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
}

func (from Resource) transformToResource() resource.Resource {
	// TODO: remove *ID
	return resource.Resource{
		Idxa:           from.ID,
		URN:            from.URN,
		Name:           from.Name,
		ProjectID:      from.ProjectID,
		NamespaceID:    from.NamespaceID,
		OrganizationID: from.OrganizationID,
		UserID:         from.UserID.String,
		CreatedAt:      from.CreatedAt,
		UpdatedAt:      from.UpdatedAt,
	}
}

type ResourceCols struct {
	ID             string         `db:"id"`
	URN            string         `db:"urn"`
	Name           string         `db:"name"`
	ProjectID      string         `db:"project_id"`
	OrganizationID string         `db:"org_id"`
	NamespaceID    string         `db:"namespace_id"`
	UserID         sql.NullString `db:"user_id"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

type ResourceConfig struct {
	ID        uint32    `db:"id"`
	Name      string    `db:"name"`
	Config    string    `db:"config"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (from ResourceConfig) transformToResourceConfig() schema.Config {
	return schema.Config{
		ID:        from.ID,
		Name:      from.Name,
		Config:    from.Config,
		CreatedAt: from.CreatedAt,
		UpdatedAt: from.UpdatedAt,
	}
}
