package role

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/goto/shield/core/namespace"
	"github.com/goto/shield/pkg/metadata"
)

const AuditEntity = "role"

type Repository interface {
	Get(ctx context.Context, id string) (Role, error)
	List(ctx context.Context) ([]Role, error)
	Create(ctx context.Context, role Role) (string, error)
	Update(ctx context.Context, toUpdate Role) (string, error)
}

type Role struct {
	ID          string
	Name        string
	Types       []string
	NamespaceID string
	Metadata    metadata.Metadata
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoleLogData struct {
	Entity      string
	ID          string
	Name        string
	Types       []string
	NamespaceID string
}

func (role Role) ToRoleLogData() RoleLogData {
	return RoleLogData{
		Entity:      AuditEntity,
		ID:          role.ID,
		Name:        role.Name,
		Types:       role.Types,
		NamespaceID: role.NamespaceID,
	}
}

func GetOwnerRole(ns namespace.Namespace) Role {
	id := fmt.Sprintf("%s_%s", ns.ID, "owner")
	name := fmt.Sprintf("%s_%s", strings.Title(ns.ID), "Owner")
	return Role{
		ID:          id,
		Name:        name,
		Types:       []string{UserType},
		NamespaceID: ns.ID,
	}
}
