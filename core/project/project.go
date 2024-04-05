package project

import (
	"context"
	"time"

	"github.com/goto/shield/core/organization"
	"github.com/goto/shield/core/user"
	"github.com/goto/shield/pkg/metadata"
)

const AuditEntity = "project"

type Repository interface {
	GetByID(ctx context.Context, id string) (Project, error)
	GetBySlug(ctx context.Context, slug string) (Project, error)
	Create(ctx context.Context, org Project) (Project, error)
	List(ctx context.Context) ([]Project, error)
	UpdateByID(ctx context.Context, toUpdate Project) (Project, error)
	UpdateBySlug(ctx context.Context, toUpdate Project) (Project, error)
	ListAdmins(ctx context.Context, id string) ([]user.User, error)
}

type Project struct {
	ID           string
	Name         string
	Slug         string
	Organization organization.Organization
	Metadata     metadata.Metadata
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ProjectLogData struct {
	Entity string
	ID     string
	Name   string
	Slug   string
	OrgID  string
}

func (project Project) ToProjectLogData() ProjectLogData {
	return ProjectLogData{
		Entity: AuditEntity,
		ID:     project.ID,
		Name:   project.Name,
		Slug:   project.Slug,
		OrgID:  project.Organization.ID,
	}
}
