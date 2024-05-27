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

type UpsertServiceData struct {
	ID          string    `db:"data.id"`
	NamespaceID string    `db:"data.namespace_id"`
	EntityID    string    `db:"data.entity_id"`
	KeyID       string    `db:"data.key_id"`
	KeyURN      string    `db:"key.urn"`
	Value       string    `db:"data.value"`
	CreatedAt   time.Time `db:"data.created_at"`
	UpdatedAt   time.Time `db:"data.updated_at"`
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

func (from UpsertServiceData) transformToServiceData() servicedata.ServiceData {
	data := servicedata.ServiceData{
		ID:          from.ID,
		NamespaceID: from.NamespaceID,
		EntityID:    from.EntityID,
		Key: servicedata.Key{
			ID:  from.KeyID,
			URN: from.KeyURN,
		},
		Value: from.Value,
	}
	return data
}
