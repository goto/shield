package postgres

import (
	"database/sql"
	"time"

	"github.com/goto/shield/core/servicedata"
)

type Key struct {
	ID          string       `db:"id"`
	URN         string       `db:"urn"`
	ProjectID   string       `db:"project_id"`
	Key         string       `db:"key"`
	Description string       `db:"description"`
	ResourceID  string       `db:"resource_id"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

func (from Key) transformToServiceDataKey() servicedata.Key {
	return servicedata.Key{
		ID:          from.ID,
		URN:         from.URN,
		ProjectID:   from.ProjectID,
		Key:         from.Key,
		Description: from.Description,
		ResourceID:  from.ResourceID,
	}
}
