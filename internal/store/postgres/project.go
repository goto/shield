package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/project"
)

type Project struct {
	ID        string       `db:"id"`
	Name      string       `db:"name"`
	Slug      string       `db:"slug"`
	OrgID     string       `db:"org_id"`
	Metadata  []byte       `db:"metadata"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

func (from Project) transformToProject() (project.Project, error) {
	var unmarshalledMetadata map[string]any
	if err := json.Unmarshal(from.Metadata, &unmarshalledMetadata); err != nil {
		return project.Project{}, err
	}

	return project.Project{
		ID:           from.ID,
		Name:         from.Name,
		Slug:         from.Slug,
		Organization: organization.Organization{ID: from.OrgID},
		Metadata:     unmarshalledMetadata,
		CreatedAt:    from.CreatedAt,
		UpdatedAt:    from.UpdatedAt,
	}, nil
}
