package group

import (
	"context"
	"time"

	"github.com/goto/shield/core/relation"
	"github.com/goto/shield/pkg/metadata"
)

const AuditEntity = "group"

type Repository interface {
	Create(ctx context.Context, grp Group) (Group, error)
	GetByID(ctx context.Context, id string) (Group, error)
	GetByIDs(ctx context.Context, groupIDs []string) ([]Group, error)
	GetBySlug(ctx context.Context, slug string) (Group, error)
	List(ctx context.Context, flt Filter) ([]Group, error)
	UpdateByID(ctx context.Context, toUpdate Group) (Group, error)
	UpdateBySlug(ctx context.Context, toUpdate Group) (Group, error)
	ListUserGroups(ctx context.Context, userId string, roleId string) ([]Group, error)
	ListGroupRelations(ctx context.Context, objectId, subjectType, role string) ([]relation.RelationV2, error)
}

type CachedRepository interface {
	GetBySlug(ctx context.Context, slug string) (Group, error)
}

type Group struct {
	ID             string
	Name           string
	Slug           string
	OrganizationID string `json:"orgId"`
	Metadata       metadata.Metadata
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type LogData struct {
	Entity         string `mapstructure:"entity"`
	ID             string `mapstructure:"id"`
	Name           string `mapstructure:"name"`
	Slug           string `mapstructure:"slug"`
	OrganizationID string `mapstructure:"organization_id"`
}

func (group Group) ToLogData() LogData {
	return LogData{
		Entity:         AuditEntity,
		ID:             group.ID,
		Name:           group.Name,
		Slug:           group.Slug,
		OrganizationID: group.OrganizationID,
	}
}
