package action

import (
	"context"
	"time"
)

const auditEntity = "action"

type Repository interface {
	Get(ctx context.Context, id string) (Action, error)
	Upsert(ctx context.Context, action Action) (Action, error)
	List(ctx context.Context) ([]Action, error)
	Update(ctx context.Context, action Action) (Action, error)
}

type Action struct {
	ID          string
	Name        string
	NamespaceID string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ActionLogData struct {
	Entity      string `mapstructure:"entity"`
	ID          string `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	NamespaceID string `mapstructure:"namespace_id"`
}

func (action Action) ToActionLogData() ActionLogData {
	return ActionLogData{
		Entity:      auditEntity,
		ID:          action.ID,
		Name:        action.Name,
		NamespaceID: action.NamespaceID,
	}
}
