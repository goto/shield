package group

import (
	"context"
	"time"

	"golang.org/x/exp/maps"

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

type Group struct {
	ID             string
	Name           string
	Slug           string
	OrganizationID string `json:"orgId"`
	Metadata       metadata.Metadata
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (group Group) ToGroupAuditData() (map[string]string, error) {
	logData := map[string]string{
		"entity": AuditEntity,
		"id":     group.ID,
		"name":   group.Name,
		"slug":   group.Slug,
		"orgId":  group.OrganizationID,
	}
	groupMetadata, err := group.Metadata.ToStringValueMap()
	if err != nil {
		return logData, err
	}

	maps.Copy(logData, groupMetadata)
	return logData, nil
}
